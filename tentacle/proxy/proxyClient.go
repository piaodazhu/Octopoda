package proxy

import (
	"errors"
	"fmt"

	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/nameclient"
	"github.com/piaodazhu/proxylite"
)

var proxyClient *proxylite.ProxyLiteClient = nil
var cancelFunc func()
var done chan struct{}

func RegisterSshService() (string, error) {
	var err error
	var proxyServerIp string
	var proxyServerPort int
	if proxyClient == nil {
		proxyServerIp, proxyServerPort, err = getProxyServerAddr()
		if err != nil {
			return "", err
		}
		proxyClient = proxylite.NewProxyLiteClient(fmt.Sprintf("%s:%d", proxyServerIp, proxyServerPort))
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
	return fmt.Sprintf("%s:%d", proxyServerIp, port), nil
}

func UnregisterSshService() {
	if proxyClient != nil {
		cancelFunc()
		<-done
		proxyClient = nil
	}
}

func getProxyServerAddr() (string, int, error) {
	entry, err := nameclient.NameQuery(config.GlobalConfig.Brain.Name + ".proxyliteFace")
	if err != nil {
		fmt.Println("QUERY: ", config.GlobalConfig.Brain.Name+".proxyliteFace")
		return "", 0, err
	}
	return entry.Ip, entry.Port, nil
}
