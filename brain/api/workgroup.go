package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piaodazhu/Octopoda/brain/workgroup"
	"github.com/piaodazhu/Octopoda/protocols"
)

func WorkgroupInfo(ctx *gin.Context) {
	path := ctx.Query("path")
	if !workgroup.IsSameOrSubPath(path, ctx.GetHeader("rootpath")) {
		ctx.String(http.StatusUnauthorized, "target path is not under rootpath")
		return
	}
	info, err := workgroup.Info(path)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "internal error")
		return
	}
	if info == nil {
		ctx.String(http.StatusNotFound, "path not exist")
		return
	}
	ctx.JSON(http.StatusOK, info)
}

func WorkgroupGrant(ctx *gin.Context) {
	info := protocols.WorkgroupInfo{}
	err := ctx.BindJSON(&info)
	if err != nil {
		ctx.String(http.StatusBadRequest, "invalid grant info")
		return
	}
	if len(info.Password) < 6 {
		ctx.String(http.StatusBadRequest, "password must longer than 6")
		return
	}
	if !workgroup.IsDirectSubPath(info.Path, ctx.GetHeader("rootpath")) {
		ctx.String(http.StatusUnauthorized, "target path is not under rootpath")
		return
	}

	err = workgroup.Grant(info.Path, info.Password)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "internal error: " + err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func WorkgroupChildren(ctx *gin.Context) {
	path := ctx.Query("path")
	if !workgroup.IsSameOrSubPath(path, ctx.GetHeader("rootpath")) {
		ctx.String(http.StatusUnauthorized, "target path is not under rootpath")
		return
	}

	info, err := workgroup.Info(path)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "internal error: " + err.Error())
		return
	}
	if info == nil {
		ctx.String(http.StatusNotFound, "path not exist")
		return
	}

	children, err := workgroup.Children(path)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "internal error")
		return
	}
	ctx.JSON(http.StatusOK, children)
}

func WorkgroupMembers(ctx *gin.Context) {
	path := ctx.Query("path")
	if !workgroup.IsSameOrSubPath(path, ctx.GetHeader("rootpath")) {
		ctx.String(http.StatusUnauthorized, "target path is not under rootpath")
		return
	}

	info, err := workgroup.Info(path)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "internal error: " + err.Error())
		return
	}
	if info == nil {
		ctx.String(http.StatusNotFound, "path not exist")
		return
	}

	members, err := workgroup.Members(path)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "internal error: " + err.Error())
		return
	}
	ctx.JSON(http.StatusOK, members)
}

func WorkgroupMembersOperation(ctx *gin.Context) {
	params := protocols.WorkgroupMembersPostParams{}
	err := ctx.BindJSON(&params)
	if err != nil {
		ctx.String(http.StatusBadRequest, "invalid operation params")
		return
	}
	
	if !workgroup.IsDirectSubPath(params.Path, ctx.GetHeader("rootpath")) {
		ctx.String(http.StatusUnauthorized, "target path is not under rootpath")
		return
	}

	if params.IsAdd {
		// must not be empty
		if len(params.Members) == 0 {
			ctx.String(http.StatusBadRequest, "add empty members is not allows")
			return
		}
		// must be valid names
		scope := ctx.GetStringMapString("octopoda_scope")
		if !workgroup.IsInScope(scope, params.Members...) {
			ctx.String(http.StatusBadRequest, "names is valid or beyond scope")
			return
		}

		// it will add child and create info if params.Path not exist
		err = workgroup.AddMembers(ctx.GetHeader("rootpath"), params.Path, params.Members)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "add members error: " + err.Error())
			return
		}
	} else {
		// it will delete the path if params.Members is empty
		members, err := workgroup.Members(params.Path)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "get current members error: " + err.Error())
			return
		}

		if !workgroup.IsSubSet(params.Members, members) {
			ctx.String(http.StatusBadRequest, "names is valid or beyond scope")
			return
		}

		// it will delete the members along all children
		err = workgroup.RemoveMembers(ctx.GetHeader("rootpath"), params.Path, params.Members)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "remove members error: " + err.Error())
			return
		}
	}
	ctx.Status(http.StatusOK)
}
