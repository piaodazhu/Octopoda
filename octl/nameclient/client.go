package nameclient

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/output"
	"os"
	"github.com/piaodazhu/Octopoda/protocols"
)

var nsAddr string
var HttpsClient *http.Client
var BrainAddr string
var BrainIp string

func InitClient() error {
	BrainIp = config.GlobalConfig.Brain.Ip
	defaultBrainAddr := fmt.Sprintf("%s:%d", config.GlobalConfig.Brain.Ip, config.GlobalConfig.Brain.Port)
	if !config.GlobalConfig.HttpsNameServer.Enabled {
		output.PrintWarningf("NameService client is disabled")
		BrainAddr = defaultBrainAddr
		return nil
	}

	BrainAddr = ""
	nsAddr = fmt.Sprintf("%s:%d", config.GlobalConfig.HttpsNameServer.Host, config.GlobalConfig.HttpsNameServer.Port)
	output.PrintInfof("NameService client is enabled. nsAddr=%s", nsAddr)
	// init https client
	err := initHttpsClient(config.GlobalConfig.Sslinfo.CaCert, config.GlobalConfig.Sslinfo.ClientCert, config.GlobalConfig.Sslinfo.ClientKey)
	if err != nil {
		emsg := "InitHttpsClient:" + err.Error()
		output.PrintFatalln(emsg)
		return err
	}
	err = pingNameServer()
	if err != nil {
		output.PrintFatalf("Could not ping NameServer: %s (%s)", nsAddr, err.Error())
		return err
	}

	entry, err := NameQuery(config.GlobalConfig.Brain.Name + ".octlFace")
	if err != nil {
		output.PrintWarningf("Could not resolve name %s (%s)", config.GlobalConfig.Brain.Name+".octlFace", err.Error())
		return err
	}
	BrainIp = entry.Ip
	BrainAddr = fmt.Sprintf("%s:%d", BrainIp, entry.Port)
	return nil
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

func NameQuery(name string) (*protocols.NameEntry, error) {
	res, err := HttpsClient.Get(fmt.Sprintf("https://%s/query?name=%s", nsAddr, name))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("NameQuery status code = %d", res.StatusCode)
	}
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response protocols.Response
	err = json.Unmarshal(buf, &response)
	if err != nil {
		return nil, err
	}
	return response.NameEntry, nil
}
