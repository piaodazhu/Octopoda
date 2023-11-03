package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/protocols"
)

func NodeAppsInfo(ctx *gin.Context) {
	var name string
	var ok bool
	rmsg := protocols.Result{
		Rmsg: "OK",
	}

	if name, ok = ctx.GetQuery("name"); !ok {
		rmsg.Rmsg = "Lack node name"
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}

	raw, err := model.Request(name, protocols.TypeAppsInfo, []byte{})
	if err != nil {
		logger.Comm.Println("NodeAppsInfo", err)
		rmsg.Rmsg = "NodeAppsInfo"
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}
	ctx.Data(200, "application/json", raw)
}

func NodeAppVersion(ctx *gin.Context) {
	var name, app, scen string
	rmsg := protocols.Result{
		Rmsg: "OK",
	}

	name = ctx.Query("name")
	app = ctx.Query("app")
	scen = ctx.Query("scenario")
	if len(name) == 0 || len(app) == 0 || len(scen) == 0 {
		rmsg.Rmsg = "ERROR: Wrong Args"
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	aParams := &AppBasic{
		Name:     app,
		Scenario: scen,
	}
	payload, _ := config.Jsoner.Marshal(aParams)
	raw, err := model.Request(name, protocols.TypeAppVersion, payload)
	if err != nil {
		logger.Comm.Println("NodeAppVersion", err)
		rmsg.Rmsg = "NodeAppVersion"
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}
	ctx.Data(200, "application/json", raw)
}

func NodeAppReset(ctx *gin.Context) {
	var name, app, scen, version, msg string
	rmsg := protocols.Result{
		Rmsg: "OK",
	}

	name = ctx.PostForm("name")
	app = ctx.PostForm("app")
	scen = ctx.PostForm("scenario")
	version = ctx.PostForm("version")
	msg = ctx.PostForm("message")
	if len(name) == 0 || len(app) == 0 || len(scen) == 0 || len(version) == 0 || len(msg) == 0 {
		rmsg.Rmsg = "ERROR: Wrong Args"
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	arParams := &AppResetParams{
		AppBasic: AppBasic{
			Name:     app,
			Scenario: scen,
			Message:  msg,
		},
		VersionHash: version,
	}
	payload, _ := config.Jsoner.Marshal(arParams)
	raw, err := model.Request(name, protocols.TypeAppReset, payload)
	if err != nil {
		logger.Comm.Println("NodeAppsInfo", err)
		rmsg.Rmsg = "NodeAppsInfo"
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}

	var result protocols.Result
	err = config.Jsoner.Unmarshal(raw, &result)
	if err != nil {
		logger.Exceptions.Println("NodeAppReset Unmarshal", err)
		rmsg.Rmsg = "NodeApp Result:" + err.Error()
		ctx.JSON(http.StatusInternalServerError, rmsg)
		return
	}
	if result.Rmsg != "OK" {
		logger.Exceptions.Println("NodeAppReset", err)
		ctx.JSON(http.StatusNotFound, result)
		return
	}

	// update scenario version
	success := model.AddScenNodeApp(arParams.Scenario, arParams.Name, arParams.Description, name, result.Version, result.Modified)
	if !success {
		logger.Exceptions.Print("Failed: AddScenNodeApp")
		rmsg.Rmsg = "AddScenNodeApp"
		ctx.JSON(http.StatusInternalServerError, rmsg)
	}
	modified, success := model.UpdateScenario(scen, msg)
	if !success {
		logger.Exceptions.Print("Failed: UpdateScenario")
		rmsg.Rmsg = "UpdateScenario"
		ctx.JSON(http.StatusInternalServerError, rmsg)
	}
	rmsg.Modified = modified
	ctx.JSON(200, result)
}
