package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DowngradeHandler(ctx *gin.Context) {
	Busy.Lock()
	defer Busy.Unlock()
	res := Response{Msg: "OK"}
	if State.StateType != STABLE {
		res.Msg = "Not in stable state, confirm or install first"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	go doDownGrade()
	ctx.JSON(200, res)
}

func CancelHandler(ctx *gin.Context) {
	Busy.Lock()
	defer Busy.Unlock()
	res := Response{Msg: "OK"}
	if State.StateType != PREVIEW {
		res.Msg = "Not in preview state, not need cancel"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	go doCancel()
	ctx.JSON(200, res)
}

func UpgradeHandler(ctx *gin.Context) {
	Busy.Lock()
	defer Busy.Unlock()
	res := Response{Msg: "OK"}
	if State.StateType != STABLE {
		res.Msg = "Not in stable state, confirm or install first"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	version := ctx.PostForm("version")
	if version == "" {
		res.Msg = "Invalid version number"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	if version == State.Version2 {
		res.Msg = "required version is already installed"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	go doUpgrade(version)
	ctx.JSON(200, res)
}

func ConfirmHandler(ctx *gin.Context) {
	Busy.Lock()
	defer Busy.Unlock()
	res := Response{Msg: "OK"}
	if State.StateType != PREVIEW {
		res.Msg = "Not in preview state, not need confirm"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	go doConfirm()
	ctx.JSON(200, res)
}

func InstallHandler(ctx *gin.Context) {
	Busy.Lock()
	defer Busy.Unlock()
	res := Response{Msg: "OK"}
	if State.StateType != EMPTY {
		res.Msg = "Use upgrade instead"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	version := ctx.PostForm("version")
	if version == "" {
		res.Msg = "Invalid version number"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	go doInstall(version)
	ctx.JSON(200, res)
}

func GetStateHandler(ctx *gin.Context) {
	Busy.Lock()
	defer Busy.Unlock()
	statusCode := http.StatusOK
	res := Response{
		Msg:       "OK",
		StateType: State.StateType,
		StateMsg:  StateMsg[State.StateType],
		Version1:  State.Version1,
		Version2:  State.Version2,
		Version3:  State.Version3,
	}

	if PakmaError != nil {
		res.Msg = PakmaError.Error()
		statusCode = http.StatusBadRequest
		PakmaError = nil
	}
	ctx.JSON(statusCode, res)
}

func CleanHandler(ctx *gin.Context) {
	Busy.Lock()
	defer Busy.Unlock()
	res := Response{Msg: "OK"}
	doClean()
	ctx.JSON(200, res)
}
