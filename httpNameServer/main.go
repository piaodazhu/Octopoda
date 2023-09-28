package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"httpns/config"
	"httpns/logger"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	BuildVersion string = "dev"
	BuildTime    string = time.Now().UTC().String()
	BuildName    string = "httpns"
	CommitID     string = "snapshot"
)

func main() {
	// read arg
	var stdout bool
	var conf string
	var ip string
	var port int
	var caCertFile string
	var svrCertFile string
	var svrKeyFile string
	var askver bool

	flag.BoolVar(&askver, "version", false, "tell version number")
	flag.StringVar(&conf, "c", "", "specify config file")
	flag.BoolVar(&stdout, "p", false, "print log to stdout, default is false")
	flag.StringVar(&ip, "ip", "", "listening ip")
	flag.IntVar(&port, "port", 0, "listening port")
	flag.StringVar(&caCertFile, "ca", "", "ca certificate")
	flag.StringVar(&svrCertFile, "crt", "", "server certificate")
	flag.StringVar(&svrKeyFile, "key", "", "server private key")
	flag.Parse()

	if askver {
		fmt.Printf("Octopoda Octl\nbuild name:\t%s\nbuild ver:\t%s\nbuild time:\t%s\nCommit ID:\t%s\n", BuildName, BuildVersion, BuildTime, CommitID)
		return
	}

	config.InitConfig(conf)
	logger.InitLogger(stdout)
	if port != 0 {
		config.GlobalConfig.ServePort = uint16(port)
	}
	if caCertFile != "" {
		config.GlobalConfig.Sslinfo.CaCert = caCertFile
	}
	if svrCertFile != "" {
		config.GlobalConfig.Sslinfo.ServerCert = svrCertFile
	}
	if svrKeyFile != "" {
		config.GlobalConfig.Sslinfo.ServerKey = svrKeyFile
	}

	if ip != "" {
		config.GlobalConfig.ServeIp = ip
	} else if config.GlobalConfig.NetDevice == "" {
		config.GlobalConfig.ServeIp = "0.0.0.0"
	} else {
		devIp, err := getIpByDevice(config.GlobalConfig.NetDevice)
		if err != nil {
			panic("serve IP has not been configured or detected!")
		}
		config.GlobalConfig.ServeIp = devIp
	}

	// start rolling token
	startRollingToken()

	// init dao and service
	err := DaoInit()
	if err != nil {
		log.Panic(err)
	}
	ServiceInit()

	// config GIN handler
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(OctopodaLogger())
	r.Use(StatsMiddleWare())

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(200)
	})
	r.GET("/query", NameQuery)
	r.GET("/list", NameList)
	r.GET("/conf", DownloadConfig)
	r.GET("/sshinfo", DownloadSshInfo)
	r.GET("/tokens", DistribToken)

	r.POST("/register", NameRegister)
	r.POST("/delete", NameDelete)
	r.POST("/conf", UploadConfig)
	r.POST("/sshinfo", UploadSshInfo)

	r.GET("/summary", ServiceSummary)

	mountStatic(r)

	// config TLS server
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(config.GlobalConfig.Sslinfo.CaCert)
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
		Addr:      fmt.Sprintf("%s:%d", config.GlobalConfig.ServeIp, config.GlobalConfig.ServePort),
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	log.Fatal(s.ListenAndServeTLS(config.GlobalConfig.Sslinfo.ServerCert, config.GlobalConfig.Sslinfo.ServerKey))
}

func mountStatic(engine *gin.Engine) {
	for _, c := range config.GlobalConfig.StaticDirs {
		engine.Static(c.Url, c.Dir)
	}
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
