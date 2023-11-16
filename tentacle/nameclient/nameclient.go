package nameclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/security"
	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
)

var nsAddr string
var BrainHeartAddr, BrainMsgAddr string

func InitNameClient() {
	defaultBrainHeartAddr := fmt.Sprintf("%s:%d", config.GlobalConfig.Brain.Ip, config.GlobalConfig.Brain.HeartbeatPort)
	defaultBrainMsgAddr := fmt.Sprintf("%s:%d", config.GlobalConfig.Brain.Ip, config.GlobalConfig.Brain.MessagePort)
	security.TokenEnabled = config.GlobalConfig.HttpsNameServer.Enabled
	if !config.GlobalConfig.HttpsNameServer.Enabled {
		logger.Network.Println("NameService client is disabled")
		BrainHeartAddr = defaultBrainHeartAddr
		BrainMsgAddr = defaultBrainMsgAddr
		return
	}

	nsAddr = fmt.Sprintf("%s:%d", config.GlobalConfig.HttpsNameServer.Host, config.GlobalConfig.HttpsNameServer.Port+1) // port: http = https + 1
	logger.Network.Printf("NameService client is enabled. nsAddr=%s\n", nsAddr)
	err := pingNameServer()
	if err != nil {
		logger.Network.Fatal("pingNameServer:", err.Error())
		return
	}
	ResolveBrain()
}

func ResolveBrain() error {
	if !config.GlobalConfig.HttpsNameServer.Enabled {
		return errors.New("name resolution not enabled")
	}
	entry, err := NameQuery(config.GlobalConfig.Brain.Name + ".tentacleFace.heartbeat")
	if err != nil {
		logger.Network.Println("NameQuery heartbeat address error:", err.Error())
		return err
	}
	BrainHeartAddr = entry.Value

	entry, err = NameQuery(config.GlobalConfig.Brain.Name + ".tentacleFace.message")
	if err != nil {
		logger.Network.Println("NameQuery message address error:", err.Error())
		return err
	}
	BrainMsgAddr = entry.Value
	return nil
}

func pingNameServer() error {
	res, err := http.Get(fmt.Sprintf("http://%s/ping", nsAddr))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot Ping https Nameserver")
	}
	return nil
}

func NameQuery(name string) (*protocols.NameServiceEntry, error) {
	res, err := http.Get(fmt.Sprintf("http://%s/query?name=%s", nsAddr, name))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NameQuery status code = %d", res.StatusCode)
	}
	var response protocols.Response
	err = json.Unmarshal(buf, &response)
	if err != nil {
		return nil, err
	}
	return response.NameEntry, nil
}
