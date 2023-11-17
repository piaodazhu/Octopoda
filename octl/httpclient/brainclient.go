package httpclient

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

var BrainClient *http.Client

func initBrainClient() *errs.OctlError {
	err := initHttpsClient(config.GlobalConfig.Sslinfo.CaCert, config.GlobalConfig.Sslinfo.ClientCert, config.GlobalConfig.Sslinfo.ClientKey)
	if err != nil {
		emsg := "InitHttpsClient for brain:" + err.Error()
		output.PrintFatalln(emsg)
		return errs.New(errs.OctlInitClientError, emsg)
	}
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
	BrainClient = &http.Client{
		Transport: tr,
	}
	return nil
}
