package api

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"brain/model"
	"net"

	"github.com/gin-gonic/gin"
)

func NodeAppsInfo(ctx *gin.Context) {
	var name string
	var ok bool
	var conn *net.Conn
	rmsg := message.Result{
		Rmsg: "OK",
	}

	if name, ok = ctx.GetQuery("name"); !ok {
		rmsg.Rmsg = "Lack node name"
		ctx.JSON(404, rmsg)
		return
	}

	if conn, ok = model.GetNodeMsgConn(name); !ok {
		rmsg.Rmsg = "Invalid node"
		ctx.JSON(404, rmsg)
		return
	}
	raw, err := message.Request(conn, message.TypeAppsInfo, []byte{})
	if err != nil {
		logger.Comm.Println("NodeAppsInfo", err)
		rmsg.Rmsg = "NodeAppsInfo"
		ctx.JSON(404, rmsg)
		return
	}
	ctx.Data(200, "application/json", raw)
}

func NodeAppVersion(ctx *gin.Context) {
	var name, app, scen string
	var ok bool
	var conn *net.Conn
	rmsg := message.Result{
		Rmsg: "OK",
	}

	name = ctx.Query("name")
	app = ctx.Query("app")
	scen = ctx.Query("scenario")
	if len(name) == 0 || len(app) == 0 || len(scen) == 0 {
		rmsg.Rmsg = "ERROR: Wrong Args"
		ctx.JSON(400, rmsg)
		return
	}

	if conn, ok = model.GetNodeMsgConn(name); !ok {
		rmsg.Rmsg = "Invalid node"
		ctx.JSON(404, rmsg)
		return
	}

	aParams := &AppBasic{
		Name:     app,
		Scenario: scen,
	}
	payload, _ := config.Jsoner.Marshal(aParams)
	raw, err := message.Request(conn, message.TypeAppVersion, payload)
	if err != nil {
		logger.Comm.Println("NodeAppVersion", err)
		rmsg.Rmsg = "NodeAppVersion"
		ctx.JSON(404, rmsg)
		return
	}
	ctx.Data(200, "application/json", raw)
}

func NodeAppReset(ctx *gin.Context) {
	var name, app, scen, version, msg string
	var ok bool
	rmsg := message.Result{
		Rmsg: "OK",
	}

	name = ctx.PostForm("name")
	app = ctx.PostForm("app")
	scen = ctx.PostForm("scenario")
	version = ctx.PostForm("version")
	msg = ctx.PostForm("message")
	if len(name) == 0 || len(app) == 0 || len(scen) == 0 || len(version) == 0 || len(msg) == 0 {
		rmsg.Rmsg = "ERROR: Wrong Args"
		ctx.JSON(400, rmsg)
		return
	}

	var conn *net.Conn
	if conn, ok = model.GetNodeMsgConn(name); !ok {
		rmsg.Rmsg = "Invalid node"
		ctx.JSON(404, rmsg)
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
	raw, err := message.Request(conn, message.TypeAppReset, payload)
	if err != nil {
		logger.Comm.Println("NodeAppsInfo", err)
		rmsg.Rmsg = "NodeAppsInfo"
		ctx.JSON(404, rmsg)
		return
	}

	var result message.Result
	err = config.Jsoner.Unmarshal(raw, &result)
	if err != nil {
		logger.Exceptions.Println("NodeAppReset Unmarshal", err)
		rmsg.Rmsg = "NodeApp Result:" + err.Error()
		ctx.JSON(500, rmsg)
		return
	}
	if result.Rmsg != "OK" {
		logger.Exceptions.Println("NodeAppReset", err)
		ctx.JSON(404, result)
		return
	}

	// update scenario version
	success := model.AddScenNodeApp(arParams.Scenario, arParams.Name, arParams.Description, name, result.Version, result.Modified)
	if !success {
		logger.Exceptions.Print("Failed: AddScenNodeApp")
		rmsg.Rmsg = "AddScenNodeApp"
		ctx.JSON(500, rmsg)
	}
	success = model.UpdateScenario(scen, msg)
	if !success {
		logger.Exceptions.Print("Failed: UpdateScenario")
		rmsg.Rmsg = "UpdateScenario"
		ctx.JSON(500, rmsg)
	}
	ctx.JSON(200, result)
}
