package main

import (
	"tentacle/config"
	"tentacle/heartbeat"
	"tentacle/logger"
	"tentacle/network"
	"tentacle/service"
)

func main() {
	config.InitConfig()
	logger.InitLogger()
	network.InitListener()
	heartbeat.InitHeartbeat()
	service.InitService()

	network.Run()
}
