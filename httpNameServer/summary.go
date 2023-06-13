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

func StatsMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var stats *ApiStat
		var exists bool
		if stats, exists = summary.ApiStats[c.Request.URL.Path]; exists {
			stats.Requests++
		} else {
			stats = &ApiStat{Requests: 1}
			summary.ApiStats[c.Request.URL.Path] = stats
		}
		summary.TotalRequests++
		c.Next()
		if c.Writer.Status() == 200 {
			stats.Success++
		}
	}
}

func ServiceSummary(ctx *gin.Context) {
	ctx.JSON(200, summary)
}