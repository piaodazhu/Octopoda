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
		entry, err = GetNameConfigDao().List(pattern)
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
	var params ConfigUploadParam
	err := ctx.ShouldBind(&params)
	if err != nil {
		log.Println("ctx.ShouldBind():", err.Error())
		ctx.JSON(400, Response{Message: err.Error()})
		return
	}

	entry := ConfigEntry{
		ConfigUploadParam: params,
		TimeStamp:     time.Now().UnixMilli(),
	}

	switch entry.Method {
	case "clear": 
		err = GetNameConfigDao().Del(params.Name)
		if err != nil {
			log.Println("NameConfigDao.Del():", err.Error())
			ctx.JSON(400, Response{Message: err.Error()})
			return
		}
	case "reset":
		err = GetNameConfigDao().Del(params.Name)
		if err != nil {
			log.Println("NameConfigDao.Del():", err.Error())
			ctx.JSON(400, Response{Message: err.Error()})
			return
		}
		fallthrough // reset need to clear then append
	case "append":
		entry.Method = ""
		err = GetNameConfigDao().Append(params.Name, entry)
		if err != nil {
			log.Println("NameConfigDao.Append():", err.Error())
			ctx.JSON(400, Response{Message: err.Error()})
			return
		}
	default:
		log.Println("NameConfigDao.???():", err.Error())
		ctx.JSON(400, Response{Message: err.Error()})
		return
	}
	ctx.JSON(200, Response{Message: "OK"})
}

func DownloadConfig(ctx *gin.Context) {
	var params ConfigQueryParam
	err := ctx.ShouldBindQuery(&params)
	if err != nil {
		log.Println("ctx.ShouldBindQuery():", err.Error())
		ctx.JSON(400, Response{Message: err.Error()})
		return
	}
	entry, err := GetNameConfigDao().GetRange(params.Name, params.Index, params.Amount)
	if err != nil {
		log.Println("dao.DaoGet():", err.Error())
		ctx.JSON(400, Response{Message: err.Error()})
		return
	}
	ctx.JSON(200, Response{Message: "OK", RawConfig: entry})
}

func UploadSshInfo(ctx *gin.Context) {
	var params SshInfoUploadParam
	err := ctx.ShouldBind(&params)
	if err != nil {
		log.Println("ctx.ShouldBind():", err.Error())
		ctx.JSON(400, Response{Message: err.Error()})
		return
	}

	entry := SshInfo{
		SshInfoUploadParam: params,
		TimeStamp:     time.Now().UnixMilli(),
	}
	err = GetSshInfoDao().Set(params.Name, entry, 0)
	if err != nil {
		log.Println("SshInfoDao.Set():", err.Error())
		ctx.JSON(400, Response{Message: err.Error()})
		return
	}
	ctx.JSON(200, Response{Message: "OK"})
}

func DownloadSshInfo(ctx *gin.Context) {
	key, ok := ctx.GetQuery("name")
	if !ok {
		log.Println("ctx.GetQuery(): no name")
		ctx.JSON(400, Response{Message: "no name"})
		return
	}
	entry, err := GetSshInfoDao().Get(key)
	if err != nil {
		log.Println("dao.DaoGet():", err.Error())
		ctx.JSON(400, Response{Message: err.Error()})
		return
	}
	ctx.JSON(200, Response{Message: "OK", SshInfo: entry})
}
