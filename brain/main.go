package main

import (
	"brain/config"
	"brain/heartbeat"
	"brain/logger"
	"brain/model"
	"brain/network"
	"brain/ticker"
)

func main() {
	config.InitConfig()
	logger.InitLogger()
	ticker.InitTicker()
	model.InitNodeMap()
	model.InitScenarioMap()

	heartbeat.InitHeartbeat()

	network.InitTentacleFace()
	network.InitBrainFace()

	network.Run()
}
