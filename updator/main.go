package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/state", GetStateHandler)
	r.GET("/history", GetHistoryHandler)
	r.POST("/update", UpdateHandler)
	r.POST("/redo", RedoHandler)
	r.POST("/tryupdate", TryUpdateHandler)
	r.POST("/confirm", ConfirmHandler)
	
	r.Run(":3458")
}
