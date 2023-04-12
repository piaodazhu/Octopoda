package api

import (
	"brain/logger"
	"brain/message"
	"brain/model"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

func NodeAppsInfo(ctx *gin.Context) {
	var name, addr string
	var ok bool
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	if addr, ok = model.GetNodeAddress(name); !ok {
		ctx.JSON(404, struct{}{})
		return
	}

	if conn, err := net.Dial("tcp", addr); err != nil {
		ctx.JSON(404, struct{}{})
	} else {
		defer conn.Close()
		message.SendMessage(conn, message.TypeAppsInfo, []byte{})
		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeAppsInfoResponse {
			logger.Tentacle.Println("NodeAppsInfo", err)
			ctx.JSON(404, struct{}{})
			return
		}
		ctx.Data(200, "application/json", raw)
	}
}

func NodeAppVersion(ctx *gin.Context) {
	var name, app, scen, addr string
	var ok bool

	name = ctx.Query("name")
	app = ctx.Query("app")
	scen = ctx.Query("scenario")
	if len(name) == 0 || len(app) == 0 || len(scen) == 0 {
		ctx.JSON(400, struct{}{})
		return
	}

	if addr, ok = model.GetNodeAddress(name); !ok {
		ctx.JSON(404, struct{}{})
		return
	}

	if conn, err := net.Dial("tcp", addr); err != nil {
		ctx.JSON(404, struct{}{})
	} else {
		defer conn.Close()

		aParams := &AppBasic{
			Name:     app,
			Scenario: scen,
		}
		payload, _ := json.Marshal(aParams)
		message.SendMessage(conn, message.TypeAppVersion, payload)
		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeAppVersionResponse {
			logger.Tentacle.Println("NodeAppVersion", err)
			ctx.JSON(404, struct{}{})
			return
		}
		ctx.Data(200, "application/json", raw)
	}
}

func NodeAppReset(ctx *gin.Context) {
	var name, app, scen, version, addr string
	var ok bool

	name = ctx.Query("name")
	app = ctx.Query("app")
	scen = ctx.Query("scenario")
	version = ctx.Query("version")
	if len(name) == 0 || len(app) == 0 || len(scen) == 0 || len(version) == 0 {
		ctx.JSON(400, struct{}{})
		return
	}

	if addr, ok = model.GetNodeAddress(name); !ok {
		ctx.JSON(404, struct{}{})
		return
	}

	if conn, err := net.Dial("tcp", addr); err != nil {
		ctx.JSON(404, struct{}{})
	} else {
		defer conn.Close()

		arParams := &AppResetParams{
			AppBasic: AppBasic{
				Name: app,
				Scenario: scen,
			},
			VersionHash: version,
		}
		payload, _ := json.Marshal(arParams)
		message.SendMessage(conn, message.TypeAppVersion, payload)
		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeAppVersionResponse {
			logger.Tentacle.Println("NodeAppReset", err)
			ctx.JSON(500, struct{}{})
			return
		}

		var rmsg RDATA
		err = json.Unmarshal(raw, &rmsg)
		if err != nil {
			logger.Tentacle.Println("NodeAppReset Unmarshal", err)
			ctx.JSON(500, struct{}{})
			return
		}
		if rmsg.Msg != "OK" {
			logger.Tentacle.Println("NodeAppReset MSG", err)
			ctx.JSON(404, rmsg)
			return
		}

		// update scenario version
		success := model.AddScenNodeApp(arParams.Scenario, arParams.Name, arParams.Description, name, rmsg.Version, rmsg.Modified)
		if success {
			logger.Tentacle.Print("Success: AddScenNodeApp")
			ctx.JSON(200, rmsg)
		} else {
			logger.Tentacle.Print("Failed: AddScenNodeApp")
			ctx.JSON(500, struct{}{})
		}
	}
}
