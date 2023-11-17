package main

import (
	"flag"
	"fmt"

	"github.com/piaodazhu/Octopoda/brain/alert"
	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/brain/network"
	"github.com/piaodazhu/Octopoda/brain/rdb"
	"github.com/piaodazhu/Octopoda/brain/sys"
	"github.com/piaodazhu/Octopoda/protocols/buildinfo"
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

	Run()
}
