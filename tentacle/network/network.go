package network

import (
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

func dialWithDevice(addr, dev string) (net.Conn, error) {
	remote, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		logger.Network.Println("cannot resolve remote address: ", err)
		return nil, err
	}

	localip, err := getIpByDevice(dev)
	if err != nil {
		logger.Network.Println("cannot get local ip: ", err)
		return nil, err
	}

	local, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", localip))
	if err != nil {
		logger.Network.Println("cannot resolve local address: ", err)
		return nil, err
	}

	tcpConn, err := net.DialTCP("tcp", local, remote)
	if err != nil {
		// fall back to dail
		return net.Dial("tcp", addr)
	}
	return tcpConn, nil
}

func Dial(addr string) (net.Conn, error) {
	dev := config.GlobalConfig.NetDevice
	if len(dev) == 0 {
		return net.Dial("tcp", addr)
	}
	return dialWithDevice(addr, dev)
}

func Run() {
	wg.Add(1)
	joinwg.Add(1)
	KeepAlive()
	ReadAndServe()
	wg.Wait()
}
