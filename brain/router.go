package main

import (
	"brain/api"
	"brain/config"
	"brain/logger"
	"brain/message"
	"brain/model"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

var engine *gin.Engine

func Run() error {
	// gin.SetMode(gin.DebugMode)
	// engine = gin.Default()

	gin.SetMode(gin.ReleaseMode)
	engine = gin.New()
	engine.Use(gin.Recovery())

	initRouter(engine)
	return ListenCommand()
}

func initRouter(engine *gin.Engine) {
	engine.Use(BusyBlocker())
	group := engine.Group("/api/v1")
	{
		group.GET("/taskstate", api.TaskState)
		group.OPTIONS("/")

		group.Use(OctopodaLogger())
		group.GET("/node/info", api.NodeInfo)
		group.GET("/node/status", api.NodeStatus)
		group.GET("/node/log", api.NodeLog)
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
		group.GET("/file/pull", api.FilePull)

		group.POST("/run/script", api.RunScript)
		group.POST("/run/cmd", api.RunCmd)

		group.POST("/pakma", api.PakmaCmd)
	}
}

func NotImpl(ctx *gin.Context) {
	ctx.JSON(501, struct{}{})
}

func OctopodaLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()
		latency := time.Since(start)
		if latency > time.Minute {
			latency = latency.Truncate(time.Second)
		}
		logger.Request.Printf("[HTTP]| %3d | %13v | %-7s %#v\n",
			c.Writer.Status(),
			latency,
			c.Request.Method,
			path,
		)
	}
}

func BusyBlocker() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !model.CheckReady() {
			c.AbortWithStatusJSON(503, message.Result{
				Rcode: -1,
				Rmsg:  "Server Busy",
			})
		}
		c.Next()
	}
}

func ListenCommand() error {
	listenaddr := fmt.Sprintf("%s:%d", config.GlobalConfig.OctlFace.Ip, config.GlobalConfig.OctlFace.Port)
	return engine.Run(listenaddr)
}
