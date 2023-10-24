package main

import (
	"fmt"
	"httpns/logger"
	"protocols"
	"time"

	"github.com/gin-gonic/gin"
)

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

func StatsMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var stats *protocols.ApiStat
		var exists bool
		key := fmt.Sprintf("%s[%s]", c.Request.URL.Path, c.Request.Method)
		if stats, exists = summary.ApiStats[key]; exists {
			stats.Requests++
		} else {
			stats = &protocols.ApiStat{Requests: 1}
			summary.ApiStats[key] = stats
		}
		summary.TotalRequests++
		c.Next()
		if c.Writer.Status() == 200 {
			stats.Success++
		}
	}
}
