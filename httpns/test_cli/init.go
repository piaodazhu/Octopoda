package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"net/http"
	"os"
)

var port int
var ip string
var caCertFile string
var cliCertFile string
var cliKeyFile string

var httpsClient *http.Client

func NewTransPort() *http.Transport {
	ca, err := os.ReadFile(caCertFile)
	if err != nil {
		log.Fatalln(err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(ca)

	clientCrt, err := tls.LoadX509KeyPair(cliCertFile, cliKeyFile)
	if err != nil {
		log.Fatalln(err)
	}

	return &http.Transport{
		TLSHandshakeTimeout: 0,
		TLSClientConfig: &tls.Config{
			RootCAs:            certPool,
			InsecureSkipVerify: false,
			ClientAuth:         tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{
				clientCrt,
			},
		},
	}
}

func NewHttpsClient() *http.Client {
	return &http.Client{
		Timeout:   0,
		Transport: NewTransPort(),
	}
}

func InitHttpsClient() {
	httpsClient = NewHttpsClient()
}

func PingServer() error {
	res, err := httpsClient.Get(host + "/ping")
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New("ping Error")
	}
	return nil
}
