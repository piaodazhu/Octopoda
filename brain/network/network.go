package network

import (
	"brain/logger"
	"brain/model"
)

func Run() {
	ListenNodeJoin()

	// must finish first fix.
	logger.Brain.Println("starting...")
	<-model.FirstFixed
	logger.Brain.Println("start")
	ListenCommand()
}
