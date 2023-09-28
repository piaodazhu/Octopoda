package main

import (
	"flag"
	"fmt"
	"tentacle/app"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/nameclient"
	"tentacle/network"
	"tentacle/service"
	"time"
)

var (
	BuildVersion string = "dev"
	BuildTime    string = time.Now().UTC().String()
	BuildName    string = "tentacle"
	CommitID     string = "snapshot"
)

func main() {
	var stdout bool
	var askver bool
	var conf string
	flag.BoolVar(&stdout, "p", false, "print log to stdout, default is false")
	flag.BoolVar(&askver, "version", false, "tell version number")
	flag.StringVar(&conf, "c", "", "specify a configuration file")
	flag.Parse()

	if askver {
		fmt.Printf("Octopoda Tentacle\nbuild name:\t%s\nbuild ver:\t%s\nbuild time:\t%s\nCommit ID:\t%s\n", BuildName, BuildVersion, BuildTime, CommitID)
		return
	}

	config.InitConfig(conf)
	logger.InitLogger(stdout)

	app.InitAppModel()

	nameclient.InitNameClient()
	service.InitService()

	network.Run()
}
