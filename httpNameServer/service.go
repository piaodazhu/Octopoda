package main

import (
	"log"
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

func NameRegister(ctx *gin.Context) {
	var params RegisterParam
	err := ctx.ShouldBind(&params)
	if err != nil {
		log.Println("ctx.ShouldBind():", err.Error())
		ctx.JSON(400, Response{Message: err.Error()})
		return
	}

	entry := NameEntry{
		Type:        params.Type,
		Name:        params.Name,
		Ip:          params.Ip,
		Port:        params.Port,
		Description: params.Description,
		TimeStamp:   time.Now().UnixMilli(),
	}
	err = GetNameEntryDao().DaoSet(params.Name, entry, params.TTL)
	if err != nil {
		log.Println("dao.DaoSet():", err.Error())
		ctx.JSON(400, Response{Message: err.Error()})
		return
	}
	ctx.JSON(200, Response{Message: "OK"})
}

func NameDelete(ctx *gin.Context) {
	key, ok := ctx.GetPostForm("name")
	if !ok {
		log.Println("ctx.GetQuery(): no name")
		ctx.JSON(400, Response{Message: "no name"})
		return
	}
	err := GetNameEntryDao().DaoDel(key)
	if err != nil {
		log.Println("dao.DaoDel():", err.Error())
		ctx.JSON(400, Response{Message: err.Error()})
		return
	}
	ctx.JSON(200, Response{Message: "OK"})
}

func NameQuery(ctx *gin.Context) {
	key, ok := ctx.GetQuery("name")
	if !ok {
		log.Println("ctx.GetQuery(): no name")
		ctx.JSON(400, Response{Message: "no name"})
		return
	}
	entry, err := GetNameEntryDao().DaoGet(key)
	if err != nil {
		log.Println("dao.DaoGet():", err.Error())
		ctx.JSON(400, Response{Message: err.Error()})
		return
	}
	ctx.JSON(200, Response{Message: "OK", NameEntry: entry})
}

func NameList(ctx *gin.Context) {
	var params ListQueryParam
	err := ctx.ShouldBindQuery(&params)
	if err != nil {
		log.Println("ctx.ShouldBindQuery():", err.Error())
		ctx.JSON(400, Response{Message: err.Error()})
		return
	}

	var pattern string
	switch params.Method {
	case "prefix":
		pattern = params.Match + "*"
	case "suffix":
		pattern = "*" + params.Match
	case "contain":
		pattern = "*" + params.Match + "*"
	case "equal":
		pattern = params.Match
	case "all":
		pattern = "*"
	}

	entry, err := GetNameEntryDao().DaoList(pattern)
	if err != nil {
		log.Println("dao.DaoGet():", err.Error())
		ctx.JSON(400, Response{Message: err.Error()})
		return
	}
	ctx.JSON(200, Response{Message: "OK", NameList: entry})
}

func ServiceSummary(ctx *gin.Context) {
	ctx.JSON(200, summary)
}
