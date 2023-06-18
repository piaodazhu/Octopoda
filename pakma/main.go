package main

import (
	"flag"
	"pakma/config"
	"pakma/httpsclient"

	"github.com/gin-gonic/gin"
)

func main() {
	var conf string
	flag.StringVar(&conf, "c", "", "config file of target app")
	flag.Parse()

	config.InitConfig(conf)
	httpsclient.InitClient()
	InitFiniteStateMachine()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/state", GetStateHandler)       // http://127.0.0.1/state
	r.GET("/history", GetHistoryHandler)   // http://127.0.0.1/history?time=2023-06-18-15:04:05&limit=10
	r.POST("/downgrade", DowngradeHandler) // http://127.0.0.1/downgrade
	r.POST("/cancel", CancelHandler)       // http://127.0.0.1/cancel
	r.POST("/upgrade", UpgradeHandler)    // http://127.0.0.1/upgrade form: version=1.3.5
	r.POST("/confirm", ConfirmHandler)     // http://127.0.0.1/confirm
	r.POST("/install", InstallHandler)     // http://127.0.0.1/install

	r.Run(":3458")
}
