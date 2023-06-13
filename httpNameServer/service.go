package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func NameRegister(ctx *gin.Context) {
	var params RegisterParam
	err := ctx.ShouldBind(&params)
	if err != nil {
		log.Println("ctx.ShouldBind():", err.Error())
		ctx.JSON(400, Response{Message: err.Error()})
		return
	}

	entry := NameEntry{
		RegisterParam: params,
		TimeStamp:     time.Now().UnixMilli(),
	}
	err = GetNameEntryDao().Set(params.Name, entry, params.TTL)
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
	err := GetNameEntryDao().Del(key)
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
	entry, err := GetNameEntryDao().Get(key)
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
	default:
		pattern = ""
	}

	var entry []string
	switch params.Scope {
	case "name":
		entry, err = GetNameEntryDao().List(pattern)
	case "config":
		entry, err = GetNamConfigDao().List(pattern)
	case "ssh":
		entry, err = GetSshInfoDao().List(pattern)
	default:
	}

	if err != nil {
		log.Println("dao.DaoGet():", err.Error())
		ctx.JSON(400, Response{Message: err.Error()})
		return
	}
	ctx.JSON(200, Response{Message: "OK", NameList: entry})
}

func UploadConfig(ctx *gin.Context) {

}

func DownloadConfig(ctx *gin.Context) {

}

func UploadSshInfo(ctx *gin.Context) {

}

func DownloadSshInfo(ctx *gin.Context) {

}
