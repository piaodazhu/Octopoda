package network

import (
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

	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/security"
)

var nsAddr string
var httpsClient *http.Client

func InitNameClient() {
	security.TokenEnabled = config.GlobalConfig.HttpsNameServer.Enabled
	if !config.GlobalConfig.HttpsNameServer.Enabled {
		logger.Network.Println("NameService client is disabled")
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

	retry := 6
	success := false
	var nameEntry1, nameEntry2 *protocols.RegisterParam
	for retry > 0 {
		tentacleFaceIp, err := getTentacleFaceIp()
		if err != nil {
			logger.Network.Println("getTentacleFaceIp", err)
			time.Sleep(time.Second * 10)
			retry--
			continue
		}
		octlFaceIp, err := getOctlFaceIp()
		if err != nil {
			logger.Network.Println("getTentacleFaceIp", err)
			time.Sleep(time.Second * 10)
			retry--
			continue
		}

		// first report to ns
		nameEntry1 = &protocols.RegisterParam{
			Type:        "brain",
			Name:        config.GlobalConfig.Name + ".tentacleFace",
			Ip:          tentacleFaceIp,
			Port:        int(config.GlobalConfig.TentacleFace.HeartbeatPort),
			Port2:       int(config.GlobalConfig.TentacleFace.MessagePort),
			Description: "first update since running",
			TTL:         0,
		}
		nameEntry2 = &protocols.RegisterParam{
			Type:        "brain",
			Name:        config.GlobalConfig.Name + ".octlFace",
			Ip:          octlFaceIp,
			Port:        int(config.GlobalConfig.OctlFace.Port),
			Description: "first update since running",
			TTL:         0,
		}

		err = nameRegister(nameEntry1)
		if err != nil {
			logger.Network.Fatal("NameRegister1:", err)
		}
		err = nameRegister(nameEntry2)
		if err != nil {
			logger.Network.Fatal("NameRegister2:", err)
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
			tentacleFaceIp, err := getTentacleFaceIp()
			if err != nil {
				logger.Network.Println("getTentacleFaceIp", err)
				// time.Sleep(time.Second * 10)
				continue
			}
			octlFaceIp, err := getOctlFaceIp()
			if err != nil {
				logger.Network.Println("getTentacleFaceIp", err)
				// time.Sleep(time.Second * 10)
				continue
			}
			nameEntry1.Description = "periodical"
			nameEntry1.Ip = tentacleFaceIp
			nameEntry2.Description = "periodical"
			nameEntry2.Ip = octlFaceIp
			err = nameRegister(nameEntry1)
			if err != nil {
				logger.Network.Println("NameRegister1:", err)
			}
			err = nameRegister(nameEntry2)
			if err != nil {
				logger.Network.Println("NameRegister2:", err)
			}
		}
	}()
	fetchTokens()
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

func nameRegister(entry *protocols.RegisterParam) error {
	form := url.Values{}
	form.Set("name", entry.Name)
	form.Set("ip", entry.Ip)
	form.Set("port", strconv.Itoa(entry.Port))
	form.Set("port2", strconv.Itoa(entry.Port2))
	form.Set("type", entry.Type)
	form.Set("description", entry.Description)
	form.Set("ttl", strconv.Itoa(entry.TTL))
	res, err := httpsClient.PostForm(fmt.Sprintf("https://%s/register", nsAddr), form)
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
	if response.Message != "OK" || res.StatusCode != 200 {
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
	if res.StatusCode != 200 {
		return fmt.Errorf("cannot Ping https Nameserver")
	}
	return nil
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
