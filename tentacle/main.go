package main

import (
	"tentacle/app"
	"tentacle/config"
	"tentacle/heartbeat"
	"tentacle/logger"
	"tentacle/network"
	"tentacle/service"
)

func main() {
	config.InitConfig()
	logger.InitLogger()
	app.InitAppModel()
	network.InitListener()
	heartbeat.InitHeartbeat()
	service.InitService()

	network.Run()
}
