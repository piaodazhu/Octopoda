package nameclient

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"octl/config"
	"octl/output"
	"os"
	"strconv"
)

var nsAddr string
var HttpsClient *http.Client
var BrainAddr string
var BrainIp string

func InitClient() {
	BrainIp = config.GlobalConfig.Brain.Ip
	defaultBrainAddr := fmt.Sprintf("%s:%d", config.GlobalConfig.Brain.Ip, config.GlobalConfig.Brain.Port)
	if !config.GlobalConfig.HttpsNameServer.Enabled {
		output.PrintWarningf("NameService client is disabled")
		BrainAddr = defaultBrainAddr
		return
	}

	BrainAddr = ""
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

	entry, err := NameQuery(config.GlobalConfig.Brain.Name + ".octlFace")
	if err != nil {
		output.PrintWarningf("Could not resolve name %s (%s)", config.GlobalConfig.Brain.Name+".octlFace", err.Error())
		return
	}
	BrainIp = entry.Ip
	BrainAddr = fmt.Sprintf("%s:%d", BrainIp, entry.Port)
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

func NameQuery(name string) (*NameEntry, error) {
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
	var response Response
	err = json.Unmarshal(buf, &response)
	if err != nil {
		return nil, err
	}
	return response.NameEntry, nil
}

func SshinfoRegister(sshinfo *SshInfoUploadParam) error {
	form := url.Values{}
	form.Set("name", sshinfo.Name)
	form.Set("ip", sshinfo.Ip)
	form.Set("port", strconv.Itoa(sshinfo.Port))
	form.Set("type", sshinfo.Type)
	form.Set("username", sshinfo.Username)
	form.Set("password", sshinfo.Password)
	res, err := HttpsClient.PostForm(fmt.Sprintf("https://%s/sshinfo", nsAddr), form)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("ssh info register rejected by server")
	}
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response Response
	json.Unmarshal(buf, &response)
	return nil
}

func SshinfoQuery(name string) (*SshInfo, error) {
	res, err := HttpsClient.Get(fmt.Sprintf("https://%s/sshinfo?name=%s", nsAddr, name))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("ssh info not found")
	}
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response Response
	json.Unmarshal(buf, &response)
	return response.SshInfo, nil
}

func NameDelete(name string, scope string) error {
	form := url.Values{}
	form.Set("name", name)
	form.Set("scope", scope)
	res, err := HttpsClient.PostForm(fmt.Sprintf("https://%s/delete", nsAddr), form)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("ssh info not found")
	}
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response Response
	json.Unmarshal(buf, &response)
	return nil
}
