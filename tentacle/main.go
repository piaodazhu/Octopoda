package main

import (
	"flag"
	"fmt"
	"tentacle/app"
	"tentacle/config"
	"tentacle/heartbeat"
	"tentacle/logger"
	"tentacle/network"
	"tentacle/service"
)

var (
	BuildVersion string
	BuildTime    string
	BuildName    string
	CommitID     string
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
	network.InitListener()
	heartbeat.InitHeartbeat()
	service.InitService()

	network.Run()
}
