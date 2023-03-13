package main

import (
	"nworkerd/config"
	"nworkerd/heartbeat"
	"nworkerd/logger"
	"nworkerd/network"
	"nworkerd/service"
)

func main() {
	config.InitConfig()
	logger.InitLogger()
	network.InitListener()
	heartbeat.InitHeartbeat()
	service.InitService()

	network.Run()
}
