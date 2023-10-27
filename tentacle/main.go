package main

import (
	"flag"
	"fmt"
	"github.com/piaodazhu/Octopoda/tentacle/app"
	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
	"github.com/piaodazhu/Octopoda/tentacle/nameclient"
	"github.com/piaodazhu/Octopoda/tentacle/network"
	"github.com/piaodazhu/Octopoda/tentacle/service"
	"github.com/piaodazhu/Octopoda/protocols/buildinfo"
)

var (
	BuildVersion string = "dev"
	BuildTime    string
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

	buildinfo.SetBuildInfo(buildinfo.BuildInfo{
		BuildVersion: BuildVersion,
		BuildTime:    BuildTime,
		BuildName:    BuildName,
		CommitID:     CommitID,
	})

	config.InitConfig(conf)
	logger.InitLogger(stdout)

	app.InitAppModel()

	nameclient.InitNameClient()
	service.InitService()

	network.Run()
}
