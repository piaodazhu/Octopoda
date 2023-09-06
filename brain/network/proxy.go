package network

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"fmt"
	"time"

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
		entry := &message.RegisterParam{
			Type:        "other",
			Name:        ctx.ServiceInfo().Name,
			Ip:          octlFaceIp,
			Port:        int(ctx.ServiceInfo().Port),
			Description: ctx.ServiceInfo().Message,
			TTL:         0,
		}
		err := nameRegister(entry)
		if err != nil {
			logger.Exceptions.Printf("add proxylite service %s failed: %v", entry.Name, err)
		}
	})
	proxyliteServer.OnTunnelDestroyed(func(ctx *proxylite.Context) {
		err := nameUnregister(ctx.ServiceInfo().Name)
		if err != nil {
			logger.Exceptions.Printf("delete proxylite service %s failed: %v", ctx.ServiceInfo().Name, err)
		}
	})
	go func() {
		err := proxyliteServer.Run(fmt.Sprintf("0.0.0.0:%d", config.GlobalConfig.ProxyliteServer.Port))
		if err != nil {
			panic(err)
		}
	}()

	// register self
	nameEntry := &message.RegisterParam{
		Type:        "other",
		Name:        config.GlobalConfig.Name + ".proxyliteFace",
		Ip:          tentacleFaceIp,
		Port:        int(config.GlobalConfig.ProxyliteServer.Port),
		Description: "proxylite serve port",
		TTL:         1000 * (config.GlobalConfig.ProxyliteServer.FreshTime + 10),
	}

	err = nameRegister(nameEntry)
	fmt.Println("REGISTER: ", nameEntry.Name)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(config.GlobalConfig.ProxyliteServer.FreshTime))
			err := nameRegister(nameEntry)
			if err != nil {
				logger.Exceptions.Print("fresh proxylite server name register: ", err)
			}
		}
	}()
}
