package httpnc

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"octl/config"
	"octl/output"
	"os"
	"time"
)

var nsAddr string
var HttpsClient *http.Client
var BrainAddr string

func InitClient() {
	defaultBrainAddr := fmt.Sprintf("%s:%d", config.GlobalConfig.Brain.Ip, config.GlobalConfig.Brain.Port)
	if !config.GlobalConfig.HttpsNameServer.Enabled {
		output.PrintWarningf("NameService client is disabled")
		BrainAddr = defaultBrainAddr
		return
	}

	nsAddr = fmt.Sprintf("%s:%d", config.GlobalConfig.HttpsNameServer.Host, config.GlobalConfig.HttpsNameServer.Port)
	output.PrintInfof("NameService client is enabled. nsAddr=%s", nsAddr)
	// init https client
	err := initHttpsClient(config.GlobalConfig.Sslinfo.CaCert, config.GlobalConfig.Sslinfo.ClientCert, config.GlobalConfig.Sslinfo.ClientKey)
	if err != nil {
		output.PrintFatalln("InitHttpsClient:", err.Error())
		return
	}
	err = pingNameServer()
	if err != nil {
		output.PrintFatalf("Could not ping NameServer: %s (%s)", nsAddr, err.Error())
		return
	}

	retry := 6
	success := false
	for retry > 0 {
		entry, err := nameQuery(config.GlobalConfig.Brain.Name + ".octlFace")
		if err != nil {
			output.PrintFatalf("Could not resolve name %s (%s)", config.GlobalConfig.Brain.Name + ".octlFace", err.Error())
			time.Sleep(time.Second * time.Duration(config.GlobalConfig.HttpsNameServer.RequestInterval) * 3)
			retry--
			continue
		}
		BrainAddr = fmt.Sprintf("%s:%d", entry.Ip, entry.Port)
		success = true
		retry = 0
	}
	if !success {
		output.PrintFatalln("Exit because cannot get IPv4 address of netDevice")
		return
	}
}

func initHttpsClient(caCert, cliCert, cliKey string) error {
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
	HttpsClient = &http.Client{
		Transport: tr,
	}
	return nil
}

func pingNameServer() error {
	res, err := HttpsClient.Get(fmt.Sprintf("https://%s/ping", nsAddr))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("cannot Ping https Nameserver")
	}
	return nil
}

func nameQuery(name string) (*NameEntry, error) {
	res, err := HttpsClient.Get(fmt.Sprintf("https://%s/query?name=%s", nsAddr, name))
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
	var response Response
	err = json.Unmarshal(buf, &response)
	if err != nil {
		return nil, err
	}
	return response.NameEntry, nil
}
