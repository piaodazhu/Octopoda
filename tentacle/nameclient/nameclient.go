package nameclient

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/security"
	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
)

var nsAddr string
var httpsClient *http.Client
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

	nsAddr = fmt.Sprintf("%s:%d", config.GlobalConfig.HttpsNameServer.Host, config.GlobalConfig.HttpsNameServer.Port)
	logger.Network.Printf("NameService client is enabled. nsAddr=%s\n", nsAddr)
	// init https client
	err := InitHttpsClient(config.GlobalConfig.Sslinfo.CaCert, config.GlobalConfig.Sslinfo.ClientCert, config.GlobalConfig.Sslinfo.ClientKey)
	if err != nil {
		logger.Network.Fatal("InitHttpsClient:", err.Error())
		return
	}
	err = pingNameServer()
	if err != nil {
		logger.Network.Fatal("pingNameServer:", err.Error())
		return
	}
	ResolveBrain()
	fetchTokens()
}

func ResolveBrain() error {
	if !config.GlobalConfig.HttpsNameServer.Enabled {
		return errors.New("name resolution not enabled")
	}
	entry, err := NameQuery(config.GlobalConfig.Brain.Name + ".tentacleFace")
	if err != nil {
		logger.Network.Println("NameQuery error:", err.Error())
		return err
	}
	BrainHeartAddr = fmt.Sprintf("%s:%d", entry.Ip, entry.Port)
	BrainMsgAddr = fmt.Sprintf("%s:%d", entry.Ip, entry.Port2)
	return nil
}

func InitHttpsClient(caCert, cliCert, cliKey string) error {
	ca, err := os.ReadFile(caCert)
	if err != nil {
		return err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(ca)

	clientCrt, err := tls.LoadX509KeyPair(cliCert, cliKey)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            certPool,
			InsecureSkipVerify: false,
			ClientAuth:         tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{
				clientCrt,
			},
		},
	}
	httpsClient = &http.Client{
		Transport: tr,
		Timeout:   0,
	}
	return nil
}

func pingNameServer() error {
	res, err := httpsClient.Get(fmt.Sprintf("https://%s/ping", nsAddr))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("cannot Ping https Nameserver")
	}
	return nil
}

func NameQuery(name string) (*protocols.NameEntry, error) {
	res, err := httpsClient.Get(fmt.Sprintf("https://%s/query?name=%s", nsAddr, name))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("NameQuery status code = %d", res.StatusCode)
	}
	var response protocols.Response
	err = json.Unmarshal(buf, &response)
	if err != nil {
		return nil, err
	}
	return response.NameEntry, nil
}

func GetToken() (*protocols.Tokens, error) {
	res, err := httpsClient.Get(fmt.Sprintf("https://%s/tokens", nsAddr))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("cannot get token from Nameserver")
	}
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	tks := &protocols.Tokens{}
	err = json.Unmarshal(raw, tks)
	if err != nil {
		return nil, err
	}
	return tks, nil
}

func fetchTokens() {
	ticker := time.NewTicker(security.Fetchinterval)
	go func() {
		fetchAndUpdate()
		for range ticker.C {
			if err := fetchAndUpdate(); err != nil {
				logger.Exceptions.Println("can not get token:", err)
				continue
			}
		}
	}()
}

func fetchAndUpdate() error {
	tks, err := GetToken()
	if err != nil {
		return err
	}
	cur := security.Token{
		Raw:    []byte(tks.CurToken),
		Serial: tks.CurSerial,
		Age:    tks.CurAge,
	}
	prev := security.Token{
		Raw:    []byte(tks.PrevToken),
		Serial: tks.PrevSerial,
		Age:    tks.PrevAge,
	}
	security.UpdateTokens(cur, prev)
	return nil
}
