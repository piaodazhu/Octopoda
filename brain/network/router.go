package network

import (
	"brain/api"

	"github.com/gin-gonic/gin"
)

func initRouter(engine *gin.Engine) {
	group := engine.Group("/api/v1")
	{
		group.GET("/node/info", api.NodeInfo)
		group.GET("/node/status", api.NodeState)
		group.GET("/node/apps", NotImpl)
		group.GET("/node/log", api.NodeLog)
		group.GET("/node/reboot", api.NodeReboot)
		group.GET("/node/prune", api.NodePrune)

		group.GET("/nodes/info", api.NodesInfo)
		group.GET("/nodes/status", api.NodesState)

		group.GET("/scenario/info", NotImpl)
		group.GET("/scenario/versions", NotImpl)
		group.GET("/scenario/log", NotImpl)

		group.POST("/file/upload", api.FileUpload)
		group.POST("/file/spread", api.FileSpread)
		group.GET("/file/tree", api.FileTree)

		group.GET("/sshinfo", api.SSHInfo)
	}
}

func NotImpl(ctx *gin.Context) {
	ctx.JSON(501, struct{}{})
}