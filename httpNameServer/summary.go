package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

var summary *Summary

func ServiceInit() {
	summary = &Summary{
		TotalRequests: 0,
		Since:         time.Now().UnixMilli(),
		ApiStats:      map[string]*ApiStat{},
	}
}

func ServiceSummary(ctx *gin.Context) {
	ctx.JSON(200, summary)
}
