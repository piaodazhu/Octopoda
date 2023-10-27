package main

import (
	"fmt"
	"time"

	"github.com/piaodazhu/Octopoda/brain/api"
	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/brain/model"

	"github.com/piaodazhu/Octopoda/protocols"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
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
		group.POST("/nodes/parse", api.NodesParse)

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
		group.POST("/run/cancel", api.CancelRun)

		group.GET("/ssh", api.SshLoginInfo)
		group.POST("/ssh", api.SshRegister)
		group.DELETE("/ssh", api.SshUnregister)

		group.POST("/pakma", api.PakmaCmd)

		group.POST("/group", api.GroupSetGroup)
		group.GET("/group", api.GroupGetGroup)
		group.DELETE("/group", api.GroupDeleteGroup)

		group.GET("/groups", RateLimiter(1, 1), api.GroupGetAll)
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
			c.AbortWithStatusJSON(503, protocols.Result{
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

func RateLimiter(r float64, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(r), burst)
	return func(ctx *gin.Context) {
		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(403, fmt.Sprintf("rate limit(r=%.1f,b=%d)", r, burst))
		}
		ctx.Next()
	}
}
