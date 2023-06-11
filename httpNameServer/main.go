package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// read arg
	var port int
	var caCertFile string
	var svrCertFile string
	var svrKeyFile string
	flag.IntVar(&port, "p", 9931, "listening port")
	flag.StringVar(&caCertFile, "ca", "ca/ca.pem", "ca certificate")
	flag.StringVar(&svrCertFile, "crt", "ca/nameserver/server.pem", "server certificate")
	flag.StringVar(&svrKeyFile, "key", "ca/nameserver/server.key", "server private key")
	flag.Parse()

	// init dao and service
	DaoInit()
	ServiceInit()

	// config GIN handler
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(StatsMiddleWare())
	
	r.GET("/query", NameQuery)
	r.GET("/list", NameList)
	r.POST("/register", NameRegister)
	r.POST("/delete", NameDelete)
	r.GET("/summary", ServiceSummary)

	// config TLS server
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(caCertFile)
	if err != nil {
		log.Panic(err)
	}

	ok := certPool.AppendCertsFromPEM(ca)
	if !ok {
		log.Panic(ok)
	}
	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		ClientCAs:  certPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	// run HTTP server with TLS
	s := http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   r,
		TLSConfig: tlsConfig,
	}
	log.Fatal(s.ListenAndServeTLS(svrCertFile, svrKeyFile))
}
