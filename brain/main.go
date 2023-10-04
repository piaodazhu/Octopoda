package main

import (
	"brain/alert"
	"brain/buildinfo"
	"brain/config"
	"brain/logger"
	"brain/model"
	"brain/sys"

	"brain/network"
	"brain/rdb"
	"flag"
	"fmt"
)

var (
	BuildVersion string = "dev"
	BuildTime    string
	BuildName    string = "brain"
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
		fmt.Printf("Octopoda Brain\nbuild name:\t%s\nbuild ver:\t%s\nbuild time:\t%s\nCommit ID:\t%s\n", BuildName, BuildVersion, BuildTime, CommitID)
		return
	}

	buildinfo.SetBuildInfo(buildinfo.BuildInfo{
		BuildVersion: BuildVersion,
		BuildTime:    BuildTime,
		BuildName:    BuildName,
		CommitID:     CommitID,
	})

	config.InitConfig(conf)
	logger.InitLogger(stdout)

	rdb.InitRedis()
	alert.InitAlert()

	model.InitNodeMap()
	model.InitScenarioMap()
	sys.InitNodeStatus()

	network.InitNameClient()
	network.InitTentacleFace()
	network.WaitNodeJoin()
	network.InitProxyServer()

	if err := Run(); err != nil {
		panic(err)
	}
}
