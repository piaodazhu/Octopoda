package proxy

import (
	"errors"
	"fmt"
	"strings"

	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/nameclient"
	"github.com/piaodazhu/proxylite"
)

var proxyClient *proxylite.ProxyLiteClient = nil
var cancelFunc func()
var done chan struct{}

func RegisterSshService() (string, error) {
	var err error
	var proxyAddr string
	proxyAddr, err = getProxyServerAddr()
	if err != nil {
		return "", err
	}

	if proxyClient == nil {
		proxyClient = proxylite.NewProxyLiteClient(proxyAddr)
		proxyClient.SetLogger(nil)
	}

	port, ok := proxyClient.AnyPort()
	if !ok {
		return "", errors.New("no any port avaliable on proxy server")
	}
	name := config.GlobalConfig.Name
	cancelFunc, done, err = proxyClient.RegisterInnerService(proxylite.RegisterInfo{
		OuterPort: port,
		InnerAddr: ":22",
		Name:      name,
		Message:   "ssh proxy of " + name,
	}, proxylite.ControlInfo{
		MaxServeConn: 20,
	})
	if err != nil {
		proxyClient = nil
		return "", err
	}

	proxyIpPort := strings.Split(proxyAddr, ":")
	if len(proxyIpPort) != 2 {
		return "", errors.New("invalid proxy service address entry value: " + proxyAddr)
	}
	return fmt.Sprintf("%s:%d", proxyIpPort[0], port), nil
}

func UnregisterSshService() {
	if proxyClient != nil {
		cancelFunc()
		<-done
		proxyClient = nil
	}
}

func getProxyServerAddr() (string, error) {
	entry, err := nameclient.NameQuery(config.GlobalConfig.Brain.Name + ".proxyliteFace")
	if err != nil {
		fmt.Println("QUERY: ", config.GlobalConfig.Brain.Name+".proxyliteFace")
		return "", err
	}
	return entry.Value, nil
}
