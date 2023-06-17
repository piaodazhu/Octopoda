package network

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var nsAddr string
var httpsClient *http.Client

func InitNameClient() {
	if !config.GlobalConfig.HttpsNameServer.Enabled {
		logger.Network.Println("NameService client is disabled")
		return
	}

	nsAddr := fmt.Sprintf("%s:%d", config.GlobalConfig.HttpsNameServer.Host, config.GlobalConfig.HttpsNameServer.Port)
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

	retry := 6
	success := false
	var nameEntry *message.RegisterParam
	for retry > 0 {
		myip, err := getIpByDevice(config.GlobalConfig.NetDevice)
		if err != nil {
			logger.Network.Println("getIpByDevice:", err)
			time.Sleep(time.Second * 10)
			retry--
			continue
		}

		// first report to ns
		nameEntry = &message.RegisterParam{
			Type:        "brain",
			Name:        config.GlobalConfig.Name,
			Ip:          myip,
			Port:        int(config.GlobalConfig.TentacleFace.HeartbeatPort),
			Port2:       int(config.GlobalConfig.TentacleFace.MessagePort),
			Description: "first update since running",
			TTL:         0,
		}

		err = nameRegister(nameEntry)
		if err != nil {
			logger.Network.Fatal("NameRegister:", err)
		}
		success = true
		retry = 0
	}
	if !success {
		logger.Network.Fatal("Exit because cannot get IPv4 address of netDevice ", config.GlobalConfig.NetDevice)
		return
	}
	// periodical report to ns
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(config.GlobalConfig.HttpsNameServer.RequestInterval))
			myip, err := getIpByDevice(config.GlobalConfig.NetDevice)
			if err != nil {
				logger.Network.Println("getIpByDevice:", err)
				time.Sleep(time.Second * 10)
				continue
			}

			nameEntry.Description = "periodical"
			nameEntry.Ip = myip
			err = nameRegister(nameEntry)
			if err != nil {
				logger.Network.Println("NameRegister:", err)
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
	}
	return nil
}

func nameRegister(entry *message.RegisterParam) error {
	form := url.Values{}
	form.Set("name", entry.Name)
	form.Set("ip", entry.Ip)
	form.Set("port", strconv.Itoa(entry.Port))
	form.Set("type", entry.Type)
	form.Set("description", entry.Description)
	form.Set("ttl", strconv.Itoa(entry.TTL))
	res, err := httpsClient.PostForm(fmt.Sprintf("https://%s:%d/register", nsAddr, config.GlobalConfig.HttpsNameServer.Port), form)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response message.Response
	json.Unmarshal(buf, &response)
	if response.Message != "OK" || res.StatusCode != 200 {
		return fmt.Errorf(response.Message)
	}
	return nil
}

func pingNameServer() error {
	res, err := httpsClient.Get(fmt.Sprintf("https://%s:%d/ping", nsAddr, config.GlobalConfig.HttpsNameServer.Port))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("cannot Ping https Nameserver")
	}
	return nil
}
