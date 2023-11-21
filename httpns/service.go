package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piaodazhu/Octopoda/protocols"
)

func NameRegister(ctx *gin.Context) {
	var entries []*protocols.NameServiceEntry
	err := ctx.ShouldBind(&entries)
	if err != nil {
		log.Println("ctx.ShouldBind():", err.Error())
		ctx.JSON(http.StatusBadRequest, protocols.Response{Message: err.Error()})
		return
	}

	for _, entry := range entries {
		err = GetNameEntryDao().Set(entry.Key, *entry, entry.TTL)
		if err != nil {
			// log.Println("dao.DaoSet():", err.Error())
			ctx.JSON(http.StatusBadRequest, protocols.Response{Message: err.Error()})
			return
		}
	}
	ctx.JSON(200, protocols.Response{Message: "OK"})
}

func NameDelete(ctx *gin.Context) {
	key, ok := ctx.GetPostForm("name")
	if !ok {
		log.Println("ctx.GetQuery(): no name")
		ctx.JSON(http.StatusBadRequest, protocols.Response{Message: "no name"})
		return
	}
	scope, ok := ctx.GetPostForm("scope")
	if !ok {
		scope = "name"
	}
	var err error
	switch scope {
	case "name":
		err = GetNameEntryDao().Del(key)
	default:
		log.Println("ctx.GetQuery(): invalid scope")
		ctx.JSON(http.StatusNotFound, protocols.Response{Message: "invalid scope"})
		return
	}
	if err != nil {
		// log.Println("dao.DaoDel():", err.Error())
		ctx.JSON(http.StatusBadRequest, protocols.Response{Message: err.Error()})
		return
	}
	ctx.JSON(200, protocols.Response{Message: "OK"})
}

func NameQuery(ctx *gin.Context) {
	key, ok := ctx.GetQuery("name")
	if !ok {
		log.Println("ctx.GetQuery(): no name")
		ctx.JSON(http.StatusBadRequest, protocols.Response{Message: "no name"})
		return
	}
	entry, err := GetNameEntryDao().Get(key)
	if err != nil {
		// log.Println("dao.DaoGet():", err.Error())
		ctx.JSON(http.StatusBadRequest, protocols.Response{Message: err.Error()})
		return
	}
	ctx.JSON(200, protocols.Response{Message: "OK", NameEntry: entry})
}

func NameList(ctx *gin.Context) {
	var params protocols.ListQueryParam
	err := ctx.ShouldBindQuery(&params)
	if err != nil {
		log.Println("ctx.ShouldBindQuery():", err.Error())
		ctx.JSON(http.StatusBadRequest, protocols.Response{Message: err.Error()})
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

	names, err := GetNameEntryDao().List(pattern)
	if err != nil {
		// log.Println("dao.DaoGet():", err.Error())
		ctx.JSON(http.StatusBadRequest, protocols.Response{Message: err.Error()})
		return
	}
	ctx.JSON(200, protocols.Response{Message: "OK", NameList: names})
}
