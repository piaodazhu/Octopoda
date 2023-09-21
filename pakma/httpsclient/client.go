package httpsclient

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"pakma/config"
)

var nsAddr string
var HttpsClient *http.Client

func InitClient() {
	if !config.GlobalConfig.HttpsNameServer.Enabled {
		panic("httpsNameServer is disabled, so Pakma will not work.")
	}
	nsAddr = fmt.Sprintf("%s:%d", config.GlobalConfig.HttpsNameServer.Host, config.GlobalConfig.HttpsNameServer.Port)

	// init https client
	err := initHttpsClient(config.GlobalConfig.Sslinfo.CaCert, config.GlobalConfig.Sslinfo.ClientCert, config.GlobalConfig.Sslinfo.ClientKey)
	if err != nil {
		panic(fmt.Sprint("InitHttpsClient:", err.Error()))
	}
	err = pingNameServer()
	if err != nil {
		fmt.Printf("Could not ping NameServer: %s (%s)", nsAddr, err.Error())
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
		Timeout: 0,
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

func FetchReleasePackage(packname string, targetpath string) error {
	res, err := HttpsClient.Get(fmt.Sprintf("https://%s/release/%s.tar.xz", nsAddr, packname))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("fetchReleasePackage status code = %d", res.StatusCode)
	}

	f, err := os.Create(targetpath + packname + ".tar.xz")
	if err != nil {
		return fmt.Errorf("cannot create target file: %s", err.Error())
	}

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return fmt.Errorf("error in io Copy: %s", err.Error())
	}
	return nil
}
