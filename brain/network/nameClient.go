package network

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/protocols"
)

var nsAddr string
var httpsClient *http.Client

func InitNameClient() {
	if !config.GlobalConfig.HttpsNameServer.Enabled {
		logger.Network.Println("NameService client is disabled")
		return
	}

	nsAddr = fmt.Sprintf("%s:%d", config.GlobalConfig.HttpsNameServer.Host, config.GlobalConfig.HttpsNameServer.Port)
	logger.Network.Printf("NameService client is enabled. nsAddr=%s\n", nsAddr)

	// init https client
	err := InitHttpsClient(config.GlobalConfig.Sslinfo.CaCert, config.GlobalConfig.Sslinfo.ServerCert, config.GlobalConfig.Sslinfo.ServerKey)
	if err != nil {
		logger.Network.Fatal("InitHttpsClient:", err.Error())
		return
	}

	err = pingNameServer()
	if err != nil {
		logger.Network.Fatal("pingNameServer:", err.Error())
		return
	}

	retry := 6
	success := false
	var entries []*protocols.NameServiceEntry
	for retry > 0 {
		entries, err = getRegisterEntries()
		if err != nil {
			logger.Network.Println("getRegisterEntries: ", err)
			time.Sleep(time.Second)
			retry--
			continue
		}
		err = entriesRegister(entries)
		if err != nil {
			logger.Network.Println("entriesRegister: ", err)
			time.Sleep(time.Second)
			retry--
			continue
		}
		success = true
		retry = 0
	}
	if !success {
		logger.Network.Fatal("Exit because cannot get IPv4 address of netDevice")
		return
	}
	// periodical report to ns
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(config.GlobalConfig.HttpsNameServer.RequestInterval))
			entries, err = getRegisterEntries()
			if err != nil {
				logger.Network.Println("getRegisterEntries: ", err)
				continue
			}
			err = entriesRegister(entries)
			if err != nil {
				logger.Network.Println("entriesRegister: ", err)
			}
		}
	}()
}

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

func getTentacleFaceIp() (string, error) {
	if config.GlobalConfig.TentacleFace.Ip != "" {
		return config.GlobalConfig.TentacleFace.Ip, nil
	}
	return getIpByDevice(config.GlobalConfig.TentacleFace.NetDevice)
}

func getOctlFaceIp() (string, error) {
	if config.GlobalConfig.OctlFace.Ip != "" {
		return config.GlobalConfig.OctlFace.Ip, nil
	}
	return getIpByDevice(config.GlobalConfig.OctlFace.NetDevice)
}

func GetOctlFaceIp() (string, error) {
	return getOctlFaceIp()
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
			ServerName:         "octopoda",
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

func getRegisterEntries() ([]*protocols.NameServiceEntry, error) {
	tentacleFaceIp, err := getTentacleFaceIp()
	if err != nil {
		logger.Network.Println("getTentacleFaceIp", err)
		return nil, err
	}
	octlFaceIp, err := getOctlFaceIp()
	if err != nil {
		logger.Network.Println("getTentacleFaceIp", err)
		return nil, err
	}

	var entries [3]*protocols.NameServiceEntry
	// first report to ns
	entries[0] = &protocols.NameServiceEntry{
		Key:         config.GlobalConfig.Name + ".tentacleFace.heartbeat",
		Type:        "addr",
		Value:       fmt.Sprintf("%s:%d", tentacleFaceIp, config.GlobalConfig.TentacleFace.HeartbeatPort),
		Description: "first update heartbeat address since running",
		TTL:         0,
	}
	entries[1] = &protocols.NameServiceEntry{
		Key:         config.GlobalConfig.Name + ".tentacleFace.message",
		Type:        "addr",
		Value:       fmt.Sprintf("%s:%d", tentacleFaceIp, config.GlobalConfig.TentacleFace.MessagePort),
		Description: "first update message address since running",
		TTL:         0,
	}
	entries[2] = &protocols.NameServiceEntry{
		Key:         config.GlobalConfig.Name + ".octlFace.request",
		Type:        "addr",
		Value:       fmt.Sprintf("%s:%d", octlFaceIp, config.GlobalConfig.OctlFace.Port),
		Description: "first update octl request address since running",
		TTL:         0,
	}
	return entries[:], nil
}

func entriesRegister(entries []*protocols.NameServiceEntry) error {
	body, _ := json.Marshal(entries)
	res, err := httpsClient.Post(fmt.Sprintf("https://%s/register", nsAddr), "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response protocols.Response
	json.Unmarshal(buf, &response)
	if response.Message != "OK" || res.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Message)
	}
	return nil
}

func pingNameServer() error {
	res, err := httpsClient.Get(fmt.Sprintf("https://%s/ping", nsAddr))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot Ping https Nameserver")
	}
	return nil
}
