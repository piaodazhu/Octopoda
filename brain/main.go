package main

import (
	"brain/config"
	"brain/heartbeat"
	"brain/logger"
	"brain/model"
	"brain/network"
	"brain/rdb"
	"brain/ticker"
	"flag"
	"fmt"
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
	flag.BoolVar(&stdout, "p", false, "print log to stdout, default is false")
	flag.BoolVar(&askver, "version", false, "tell version number")
	flag.Parse()

	if askver {
		fmt.Printf("Octopoda Brain\nbuild name:\t%s\nbuild ver:\t%s\nbuild time:\t%s\nCommit ID:\t%s\n", BuildName, BuildVersion, BuildTime, CommitID)
		return
	}

	config.InitConfig()
	logger.InitLogger(stdout)
	ticker.InitTicker()
	rdb.InitRedis()
	model.InitNodeMap()
	model.InitScenarioMap()

	heartbeat.InitHeartbeat()

	network.InitTentacleFace()
	network.InitBrainFace()

	network.Run()
}
