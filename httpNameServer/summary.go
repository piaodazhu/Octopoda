package main

import (
	"fmt"
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
		key := fmt.Sprintf("%s[%s]", c.Request.URL.Path, c.Request.Method)
		if stats, exists = summary.ApiStats[key]; exists {
			stats.Requests++
		} else {
			stats = &ApiStat{Requests: 1}
			summary.ApiStats[key] = stats
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
