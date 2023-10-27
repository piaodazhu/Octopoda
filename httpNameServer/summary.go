package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/piaodazhu/Octopoda/protocols"
)

var summary *protocols.Summary

func ServiceInit() {
	summary = &protocols.Summary{
		TotalRequests: 0,
		Since:         time.Now().UnixMilli(),
		ApiStats:      map[string]*protocols.ApiStat{},
	}
}

func ServiceSummary(ctx *gin.Context) {
	ctx.JSON(200, summary)
}
