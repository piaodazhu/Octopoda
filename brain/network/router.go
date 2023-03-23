package network

import (
	"brain/api"

	"github.com/gin-gonic/gin"
)

func initRouter(engine *gin.Engine) {
	engine.Group("/api/v1") 
	{
		engine.GET("/node/info", api.NodeInfo)
		engine.GET("/node/status", api.NodeState)
		engine.GET("/node/apps", NotImpl)
		engine.GET("/node/log", NotImpl)
		engine.GET("/node/reboot", api.NodeReboot)

		engine.GET("/nodes/info", api.NodesInfo)
		engine.GET("/nodes/status", api.NodesState)

		engine.GET("/scenario/info", NotImpl)
		engine.GET("/scenario/versions", NotImpl)
		engine.GET("/scenario/log", NotImpl)

		engine.POST("/file/upload", api.FileUpload)
		engine.POST("/file/spread", api.FileSpread)

		engine.GET("/sshinfo", api.SSHInfo)
	}
}

func NotImpl(ctx *gin.Context) {
	ctx.JSON(200, struct{}{})
}