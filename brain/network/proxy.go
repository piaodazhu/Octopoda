package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/brain/rdb"
	"github.com/piaodazhu/Octopoda/protocols"

	"github.com/piaodazhu/proxylite"
)

type ProxyMsg struct {
	Code int
	Msg  string
	Data string
}

var proxyliteServer *proxylite.ProxyLiteServer

func InitProxyServer() {
	proxyliteServer = proxylite.NewProxyLiteServer()
	proxyliteServer.AddPort(int(config.GlobalConfig.ProxyliteServer.MinMapPort),
		int(config.GlobalConfig.ProxyliteServer.MaxMapPort))
	proxyliteServer.SetLogger(nil)

	tentacleFaceIp, err := getTentacleFaceIp()
	if err != nil {
		panic(err)
	}
	octlFaceIp, err := getOctlFaceIp()
	if err != nil {
		panic(err)
	}

	proxyliteServer.OnTunnelCreated(func(ctx *proxylite.Context) {
		CompleteSshInfo(ctx.ServiceInfo().Name, octlFaceIp, ctx.ServiceInfo().Port)
		time.AfterFunc(time.Second, func() { dumpSshInfos() })
	})
	proxyliteServer.OnTunnelDestroyed(func(ctx *proxylite.Context) {
		if sshinfo, found := GetSshInfo(ctx.ServiceInfo().Name); found {
			info := SSHInfoDump{
				Name:     ctx.ServiceInfo().Name,
				Username: sshinfo.Username,
				Password: sshinfo.Password,
			}
			go autoRestore(info)
		}
	})
	go func() {
		err := proxyliteServer.Run(fmt.Sprintf("0.0.0.0:%d", config.GlobalConfig.ProxyliteServer.Port))
		if err != nil {
			panic(err)
		}
	}()

	// register self
	nameEntry := &protocols.NameServiceEntry{
		Key:         config.GlobalConfig.Name + ".proxyliteFace",
		Type:        "addr",
		Value:       fmt.Sprintf("%s:%d", tentacleFaceIp, config.GlobalConfig.ProxyliteServer.Port),
		Description: "proxylite serve port",
		TTL:         1000 * (config.GlobalConfig.ProxyliteServer.FreshTime + 10),
	}

	err = entriesRegister([]*protocols.NameServiceEntry{nameEntry})
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(config.GlobalConfig.ProxyliteServer.FreshTime))
			err := entriesRegister([]*protocols.NameServiceEntry{nameEntry})
			if err != nil {
				logger.Exceptions.Print("fresh proxylite server name register: ", err)
			}
		}
	}()

	go restoreSessions()
}

func ProxyServices() ([]proxylite.ServiceInfo, error) {
	services, err := proxylite.DiscoverServices(fmt.Sprintf(":%d", config.GlobalConfig.ProxyliteServer.Port))
	if err != nil {
		return nil, err
	}
	return services, err
}

type SSHInfo struct {
	Ip       string
	Port     uint32
	Username string
	Password string
}

type SSHInfoDump struct {
	Name     string
	Username string
	Password string
}

const sshInfoDumpKey = "sshInfoDumpKey"

var sshInfos sync.Map

func init() {
	sshInfos = sync.Map{}
}

func dumpSshInfos() error {
	infos := []SSHInfoDump{}
	sshInfos.Range(func(key, value any) bool {
		sshInfo := value.(SSHInfo)
		if len(sshInfo.Ip) == 0 { // haven't complete. continue
			return true
		}
		infos = append(infos, SSHInfoDump{
			Name:     key.(string),
			Username: value.(SSHInfo).Username,
			Password: value.(SSHInfo).Password,
		})
		return true
	})
	serialized, _ := json.Marshal(infos)
	return rdb.SetString(sshInfoDumpKey, string(serialized))
}

func CreateSshInfo(name string, username, password string) {
	sshInfos.Store(name, SSHInfo{
		Username: username,
		Password: password,
	})
}

func CompleteSshInfo(name string, ip string, port uint32) {
	if v, found := sshInfos.Load(name); found {
		info := v.(SSHInfo)
		info.Ip = ip
		info.Port = port
		sshInfos.Store(name, info)
	}
}

func DelSshInfo(name string) {
	fmt.Println("[X] call DelSshInfo: " + name)
	sshInfos.Delete(name)
}

func GetSshInfo(name string) (SSHInfo, bool) {
	if v, found := sshInfos.Load(name); found {
		info := v.(SSHInfo)
		if len(info.Username) == 0 || len(info.Ip) == 0 {
			DelSshInfo(name)
			return SSHInfo{}, false
		}
		return v.(SSHInfo), true
	}
	return SSHInfo{}, false
}

func askRegister(name string) error {
	if name == "brain" {
		ip, _ := GetOctlFaceIp()
		CompleteSshInfo(name, ip, uint32(config.GlobalConfig.OctlFace.SshPort))
		return nil
	}
	if state, ok := model.GetNodeState(name); !ok || state != protocols.NodeStateReady {
		return fmt.Errorf("invalid node %s", name)
	}
	raw, err := model.Request(name, protocols.TypeSshRegister, []byte{})
	if err != nil {
		return fmt.Errorf("cannot request node %s: %s", name, err.Error())
	}

	pmsg := ProxyMsg{}
	err = json.Unmarshal(raw, &pmsg)
	if err != nil {
		return fmt.Errorf("cannot marshal response from node %s: %s", name, err.Error())
	}
	if pmsg.Code != 0 {
		return fmt.Errorf("askRegister client failed %s: %s", name, pmsg.Msg)
	}
	return nil
}

func innerRestore(info SSHInfoDump) error {
	services, err := ProxyServices()
	if err != nil {
		return err
	}

	for _, s := range services {
		if info.Name == s.Name {
			return nil
		}
	}
	CreateSshInfo(info.Name, info.Username, info.Password)
	if err := askRegister(info.Name); err != nil { // 不成功就删除
		DelSshInfo(info.Name)
		return errors.New("innerRestore: " + err.Error())
	}
	return nil
}

func restoreSessions() {
	fmt.Println("restoreSessions start")
	retryCnt := 3
	found := false
retryDb:
	res, found, err := rdb.GetString(sshInfoDumpKey)
	if err != nil && retryCnt > 0 {
		time.Sleep(time.Second)
		retryCnt--
		goto retryDb
	}

	if !found {
		logger.Exceptions.Println("restoreSessions: info not found")
		return
	}

	infos := []SSHInfoDump{}
	if err := json.Unmarshal([]byte(res), &infos); err != nil {
		logger.Exceptions.Println("restoreSessions: ", err)
		return
	}

	time.Sleep(time.Second * 20)

	retryCnt = 10
retryRestore:
	failed := []SSHInfoDump{}
	for _, info := range infos {
		if err := innerRestore(info); err != nil {
			logger.Exceptions.Printf("restoreSessions restore failed: %s", err.Error())
			failed = append(failed, info)
		}
	}
	if len(failed) != 0 && retryCnt > 0 {
		infos = failed
		time.Sleep(time.Second * 30)
		retryCnt--
		goto retryRestore
	}
}

func autoRestore(info SSHInfoDump) {
	fmt.Println("call autoRestore for " + info.Name)
	retryCnt := 5
	for retryCnt > 0 {
		if _, exists := GetSshInfo(info.Name); !exists { // ssh info have been deleted
			fmt.Println("autoRestore: ssh info have been deleted: " + info.Name)
			return
		}
		if state, ok := model.GetNodeState(info.Name); !ok { // already deleted from cluster
			fmt.Println("autoRestore: already deleted from cluster: " + info.Name)
			return
		} else if state != protocols.NodeStateReady { // just offline
			fmt.Println("autoRestore: node offline: " + info.Name)
			time.Sleep(time.Minute)
			retryCnt = 5
			continue
		}

		if err := innerRestore(info); err != nil {
			logger.Exceptions.Printf("autoRestore failed: %s", err.Error())
			retryCnt--
			time.Sleep(time.Second * 10)
			continue
		}
		fmt.Println("autoRestore: success: " + info.Name)
		return
	}
	DelSshInfo(info.Name)
}
