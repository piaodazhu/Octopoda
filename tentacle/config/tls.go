package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
)

var TLSConfig *tls.Config = nil

func InitTLSConfig() error {
	cert, err := tls.LoadX509KeyPair(GlobalConfig.Sslinfo.ClientCert, GlobalConfig.Sslinfo.ClientKey)
	if err != nil {
		return err
	}

	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(GlobalConfig.Sslinfo.CaCert)
	if err != nil {
		return err
	}

	ok := certPool.AppendCertsFromPEM(ca)
	if !ok {
		return errors.New("certPool.AppendCertsFromPEM failed")
	}
	TLSConfig = &tls.Config{
		ServerName:         "octopoda",
		Certificates:       []tls.Certificate{cert},
		RootCAs:            certPool,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		InsecureSkipVerify: false,
	}
	return nil
}
