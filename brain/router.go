package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/piaodazhu/Octopoda/brain/api"
	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/brain/workgroup"

	"github.com/piaodazhu/Octopoda/protocols"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func Run() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	initRouter(engine)
	listenTLS(engine)
}

func initRouter(engine *gin.Engine) {
	engine.Use(gin.Recovery())
	engine.Use(BusyBlocker())
	engine.Use(WorkgroupAuth())
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
		group.GET("/node/app/info", api.NodeAppInfo)
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
		group.POST("/scenario/app/commit", api.AppCommit)

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

		group.GET("/workgroup/info", api.WorkgroupInfo)
		group.POST("/workgroup/info", api.WorkgroupGrant)
		group.GET("/workgroup/children", api.WorkgroupChildren)
		group.GET("/workgroup/members", api.WorkgroupMembers)
		group.POST("/workgroup/members", api.WorkgroupMembersOperation)
	}
}

func NotImpl(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, struct{}{})
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
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, protocols.Result{
				Rcode: -1,
				Rmsg:  "Server Busy",
			})
			return 
		}
		c.Next()
	}
}

func listenTLS(engine *gin.Engine) {
	// config TLS server
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(config.GlobalConfig.Sslinfo.CaCert)
	if err != nil {
		log.Panic(err)
	}

	ok := certPool.AppendCertsFromPEM(ca)
	if !ok {
		log.Panic(ok)
	}
	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		ClientCAs:  certPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	// run HTTP server with TLS
	s := http.Server{
		Addr:      fmt.Sprintf("%s:%d", config.GlobalConfig.OctlFace.Ip, config.GlobalConfig.OctlFace.Port),
		Handler:   engine,
		TLSConfig: tlsConfig,
	}
	logger.SysInfo.Fatal(s.ListenAndServeTLS(config.GlobalConfig.Sslinfo.ServerCert, config.GlobalConfig.Sslinfo.ServerKey))
}

func RateLimiter(r float64, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(r), burst)
	return func(ctx *gin.Context) {
		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusForbidden, fmt.Sprintf("rate limit(r=%.1f,b=%d)", r, burst))
			return 
		}
		ctx.Next()
	}
}

func WorkgroupAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rootpath := ctx.GetHeader("rootpath")
		rootpath = strings.TrimSuffix(rootpath, "/")
		currentPath := ctx.GetHeader("currentpath")
		currentPath = strings.TrimSuffix(currentPath, "/")
		password := ctx.GetHeader("password")
		if len(password) == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !workgroup.IsSameOrSubPath(currentPath, rootpath) {
			fmt.Println("!workgroup.IsSameOrSubPath(currentPath, rootpath)")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		} 
		fmt.Printf("%s, %s\n", rootpath, currentPath)

		info, err := workgroup.Info(rootpath)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return 
		}
		if info == nil { // rootgroup not found
			fmt.Println("info == nil")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if info.Password != password { // rootgroup unauth
			fmt.Println("info.Password != password: ", info.Password, password)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		members, err := workgroup.Members(currentPath)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return 
		}
		if len(currentPath) > 0 && len(members) == 0 {
			fmt.Println("len(members) == 0")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("octopoda_scope", workgroup.MakeScope(members))
		ctx.Next()
	}
}
