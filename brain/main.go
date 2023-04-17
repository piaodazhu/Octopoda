package main

import (
	"brain/config"
	"brain/heartbeat"
	"brain/logger"
	"brain/model"
	"brain/network"
	"brain/rdb"
	"brain/ticker"
)

func main() {
	config.InitConfig()
	logger.InitLogger()
	ticker.InitTicker()
	rdb.InitRedis()
	model.InitNodeMap()
	model.InitScenarioMap()

	heartbeat.InitHeartbeat()

	network.InitTentacleFace()
	network.InitBrainFace()

	network.Run()
}
