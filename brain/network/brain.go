package network

import (
	"brain/config"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var engine *gin.Engine
var listenaddr string

func InitBrainFace() {
	var sb strings.Builder
	sb.WriteString(config.GlobalConfig.BrainFace.Ip)
	sb.WriteByte(':')
	sb.WriteString(strconv.Itoa(int(config.GlobalConfig.BrainFace.Port)))
	listenaddr = sb.String()

	// gin.SetMode(gin.DebugMode)
	// engine = gin.Default()

	gin.SetMode(gin.ReleaseMode)
	engine = gin.New()
	engine.Use(gin.Recovery())

	initRouter(engine)
}

func ListenCommand() {
	engine.Run(listenaddr)
}
