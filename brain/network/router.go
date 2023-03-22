package network

import "github.com/gin-gonic/gin"

func initRouter(engine *gin.Engine) {
	engine.GET("/node/info/", foobar)
	engine.GET("/node/status/", foobar)
	engine.GET("/node/apps/", foobar)
	engine.GET("/node/log/", foobar)

	engine.GET("/scenario/info/", foobar)
	engine.GET("/scenario/versions/", foobar)
	engine.GET("/scenario/log/", foobar)

}

func foobar(ctx *gin.Context) {
	ctx.JSON(200, struct{}{})
}