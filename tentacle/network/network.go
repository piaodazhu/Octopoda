package network

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"

	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
)

var wg, joinwg sync.WaitGroup

func getIpByDevice(device string) (string, error) {
	iface, err := net.InterfaceByName(device)
	if err != nil {
		return "", err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("IPv4 address not found with device %s", device)
}

func Dial(addr string) (net.Conn, error) {
	dailer := net.Dialer{KeepAlive: -1}
	dev := config.GlobalConfig.NetDevice
	if len(dev) != 0 {
		localip, err := getIpByDevice(dev)
		if err != nil {
			logger.Network.Println("cannot get local ip: ", err)
			// fall back to dail
		} else {
			local, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", localip))
			if err != nil {
				logger.Network.Println("cannot resolve local address: ", err)
				// fall back to dail
			}
			dailer.LocalAddr = local
		}
	}
	conn, err := tls.DialWithDialer(&dailer, "tcp", addr, config.TLSConfig)
	if err != nil {
		logger.Network.Println("cannot dail with dialer: ", err)
		// fall back to dail
		return tls.Dial("tcp", addr, config.TLSConfig)
	}
	return conn, nil
}

func Run() {
	wg.Add(1)
	joinwg.Add(1)
	KeepAlive()
	ReadAndServe()
	wg.Wait()
}
