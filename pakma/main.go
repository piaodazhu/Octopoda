package main

import (
	"flag"
	"fmt"
	"pakma/config"
	"pakma/httpsclient"

	"github.com/gin-gonic/gin"
)

var (
	BuildVersion string = "dev"
	BuildTime    string
	BuildName    string = "pakma"
	CommitID     string = "snapshot"
)

func main() {
	var conf string
	var askver bool

	flag.BoolVar(&askver, "version", false, "tell version number")
	flag.StringVar(&conf, "c", "", "config file of target app")
	flag.Parse()

	if askver {
		fmt.Printf("Octopoda Octl\nbuild name:\t%s\nbuild ver:\t%s\nbuild time:\t%s\nCommit ID:\t%s\n", BuildName, BuildVersion, BuildTime, CommitID)
		return
	}

	config.InitConfig(conf)
	httpsclient.InitClient()
	InitFiniteStateMachine()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/state", GetStateHandler)       // http://127.0.0.1/state
	r.GET("/history", GetHistoryHandler)   // http://127.0.0.1/history?time=2023-06-18@15:04:05&limit=10
	r.POST("/downgrade", DowngradeHandler) // http://127.0.0.1/downgrade
	r.POST("/cancel", CancelHandler)       // http://127.0.0.1/cancel
	r.POST("/upgrade", UpgradeHandler)     // http://127.0.0.1/upgrade form: version=1.3.5
	r.POST("/confirm", ConfirmHandler)     // http://127.0.0.1/confirm
	r.POST("/install", InstallHandler)     // http://127.0.0.1/install
	r.POST("/clean", CleanHandler)         // http://127.0.0.1/clean

	err := r.Run(fmt.Sprintf(":%d", config.GlobalConfig.Packma.ServePort))
	if err != nil {
		panic(err)
	}
}
