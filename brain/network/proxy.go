package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/brain/rdb"
	"github.com/piaodazhu/Octopoda/protocols"

	"github.com/piaodazhu/proxylite"
)

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
			info := protocols.SSHInfoDump{
				Name:     ctx.ServiceInfo().Name,
				Username: sshinfo.Username,
				Password: sshinfo.Password,
				Port:     sshinfo.Port,
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

const sshInfoDumpKey = "sshInfoDumpKey"

var sshInfos sync.Map

func init() {
	sshInfos = sync.Map{}
}

func dumpSshInfos() error {
	infos := []protocols.SSHInfoDump{}
	sshInfos.Range(func(key, value any) bool {
		sshInfo := value.(protocols.SSHInfo)
		if len(sshInfo.Ip) == 0 { // haven't complete. continue
			return true
		}
		infos = append(infos, protocols.SSHInfoDump{
			Name:     key.(string),
			Username: value.(protocols.SSHInfo).Username,
			Password: value.(protocols.SSHInfo).Password,
			Ip:       value.(protocols.SSHInfo).Ip,
			Port:     value.(protocols.SSHInfo).Port,
		})
		return true
	})
	serialized, _ := json.Marshal(infos)
	return rdb.SetString(sshInfoDumpKey, string(serialized))
}

func CreateSshInfo(name string, username, password string) {
	sshInfos.Store(name, protocols.SSHInfo{
		Username: username,
		Password: password,
	})
}

func CompleteSshInfo(name string, ip string, port uint32) {
	if v, found := sshInfos.Load(name); found {
		info := v.(protocols.SSHInfo)
		info.Ip = ip
		info.Port = port
		sshInfos.Store(name, info)
	}
}

func DelSshInfo(name string) {
	sshInfos.Delete(name)
	dumpSshInfos()
}

func GetSshInfo(name string) (protocols.SSHInfo, bool) {
	if v, found := sshInfos.Load(name); found {
		info := v.(protocols.SSHInfo)
		if len(info.Username) == 0 || len(info.Ip) == 0 {
			fmt.Println("call del 1")
			DelSshInfo(name)
			return protocols.SSHInfo{}, false
		}
		return v.(protocols.SSHInfo), true
	}
	return protocols.SSHInfo{}, false
}

func askRegister(name string, port uint32) error {
	if name == "brain" {
		ip, _ := GetOctlFaceIp()
		CompleteSshInfo(name, ip, uint32(config.GlobalConfig.OctlFace.SshPort))
		return nil
	}
	if state, ok := model.GetNodeState(name); !ok || state != protocols.NodeStateReady {
		return fmt.Errorf("invalid node %s", name)
	}

	raw, err := model.Request(name, protocols.TypeSshRegister, []byte(strconv.Itoa(int(port))))
	if err != nil {
		return fmt.Errorf("cannot request node %s: %s", name, err.Error())
	}

	pmsg := protocols.ProxyMsg{}
	err = json.Unmarshal(raw, &pmsg)
	if err != nil {
		return fmt.Errorf("cannot marshal response from node %s: %s", name, err.Error())
	}
	if pmsg.Code != 0 {
		return fmt.Errorf("askRegister client failed %s: %s", name, pmsg.Msg)
	}
	return nil
}

func innerRestore(info protocols.SSHInfoDump) error {
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
	if err := askRegister(info.Name, info.Port); err != nil { // 不成功就删除?
		// DelSshInfo(info.Name)
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

	infos := []protocols.SSHInfoDump{}
	if err := json.Unmarshal([]byte(res), &infos); err != nil {
		logger.Exceptions.Println("restoreSessions: ", err)
		return
	}

	time.Sleep(time.Second * 20)

	retryCnt = 10
retryRestore:
	failed := []protocols.SSHInfoDump{}
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

	for _, info := range failed {
		fmt.Println("call del 2")
		DelSshInfo(info.Name)
	}
}

func autoRestore(info protocols.SSHInfoDump) {
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
	fmt.Println("call del 3")
	DelSshInfo(info.Name)
}
