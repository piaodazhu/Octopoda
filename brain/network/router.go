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
		group.GET("/node/log", api.NodeLog)
		group.GET("/node/reboot", api.NodeReboot)
		group.GET("/node/prune", api.NodePrune)

		group.GET("/node/apps", api.NodeAppsInfo)
		group.GET("/node/app/version", api.NodeAppVersion)
		group.POST("/node/app/version", api.NodeAppReset)

		group.GET("/nodes/info", api.NodesInfo)
		group.GET("/nodes/status", api.NodesState)

		group.GET("/scenario/info", api.ScenarioInfo)
		group.POST("/scenario/info", api.ScenarioCreate)
		group.DELETE("/scenario/info", api.ScenarioDelete)
		group.POST("/scenario/update", api.ScenarioUpdate)

		group.GET("/scenario/version", api.ScenarioVersion)
		group.POST("/scenario/version", api.ScenarioReset)

		group.GET("/scenario/fix", api.ScenarioFix)

		group.GET("/scenarios/info", api.ScenariosInfo)

		group.POST("/scenario/app/prepare", api.AppPrepare)
		group.POST("/scenario/app/deployment", api.AppDeploy)

		group.POST("/file/upload", api.FileUpload)
		group.POST("/file/spread", api.FileSpread)
		group.POST("/file/distrib", api.FileDistrib)
		group.GET("/file/tree", api.FileTree)

		group.GET("/sshinfo", api.SSHInfo)
		group.POST("/run/script", api.RunScript)
		group.POST("/run/cmd", api.RunCmd)
	}
}

func NotImpl(ctx *gin.Context) {
	ctx.JSON(501, struct{}{})
}
