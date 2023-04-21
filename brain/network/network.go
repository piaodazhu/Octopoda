package network

import (
	"brain/model"
)

func Run() {
	ListenNodeJoin()

	// must finish first fix.
	// logger.SysInfo.Println("starting...")
	<-model.FirstFixed
	// logger.SysInfo.Println("start")
	ListenCommand()
}
