package network

import (
	"fmt"
	"sync"
	"time"

	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
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
	})
	proxyliteServer.OnTunnelDestroyed(func(ctx *proxylite.Context) {
		DelSshInfo(ctx.ServiceInfo().Name)
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
}

func ProxyServices() ([]proxylite.ServiceInfo, error) {
	services, err := proxylite.DiscoverServices(fmt.Sprintf(":%d", config.GlobalConfig.ProxyliteServer.Port))
	if err != nil {
		return nil, err
	}

	// clean sshinfo
	set := map[string]struct{}{}
	for _, s := range services {
		set[s.Name] = struct{}{}
	}
	itemToBeDel := []string{}
	sshInfos.Range(func(key, value any) bool {
		name := key.(string)
		if _, found := set[name]; !found {
			itemToBeDel = append(itemToBeDel, name)
		}
		return true
	})
	for _, name := range itemToBeDel {
		DelSshInfo(name)
	}

	return services, err
}

type SSHInfo struct {
	Ip       string
	Port     uint32
	Username string
	Password string
}

var sshInfos sync.Map

func init() {
	sshInfos = sync.Map{}
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

