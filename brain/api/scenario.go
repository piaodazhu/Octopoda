package api

import (
	"brain/logger"
	"brain/message"
	"brain/model"
	"encoding/json"
	"net"
	"sync"

	"github.com/gin-gonic/gin"
)

func ScenarioCreate(ctx *gin.Context) {
	name := ctx.PostForm("name")
	description := ctx.PostForm("description")
	rmsg := RMSG{"OK"}
	if len(name) == 0 || len(description) == 0 {
		rmsg.Msg = "ERROR: Wrong Args"
		ctx.JSON(400, rmsg)
		return
	}
	if !model.AddScenario(name, description) {
		rmsg.Msg = "ERROR: Scenario Exists"
		ctx.JSON(404, rmsg)
		return
	}
	ctx.JSON(200, rmsg)
}

type AppDeleteParams struct {
	Name     string
	Scenario string
}

func ScenarioDelete(ctx *gin.Context) {
	// This is not so simple. Should delete all apps of a scenario in this function
	name := ctx.PostForm("name")
	rmsg := RMSG{"OK"}
	if len(name) == 0 {
		rmsg.Msg = "ERROR: Wrong Args"
		ctx.JSON(400, rmsg)
		return
	}

	// delete all app in nodes
	dlist := model.GetNodeApps(name)

	// check target nodes
	// delete app
	results := make([]UploadResults, len(dlist))
	var wg sync.WaitGroup

	for i := range dlist {
		//payload
		payload, _ := json.Marshal(&AppDeleteParams{dlist[i].AppName, dlist[i].ScenName})

		// item name
		item := dlist[i].NodeName + ":" + dlist[i].AppName + "@" + dlist[i].ScenName
		results[i].Name = item

		// node name
		nodename := dlist[i].NodeName
		if addr, exists := model.GetNodeAddress(nodename); exists {
			wg.Add(1)
			// payload?
			go deleteApp(addr, payload, &wg, &results[i].Result)
		} else {
			results[i].Result = "NodeNotExists"
		}
	}
	wg.Wait()

	// finally delete this scenario locallly
	model.DelScenario(name)
	ctx.JSON(200, results)
}

func deleteApp(addr string, payload []byte, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "OK"

	if conn, err := net.Dial("tcp", addr); err != nil {
		return
	} else {
		defer conn.Close()

		message.SendMessage(conn, message.TypeAppDelete, payload)
		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeAppDeleteResponse {
			logger.Tentacle.Println("TypeAppDeleteResponse", err)
			*result = "NetError"
			return
		}

		var rmsg RMSG
		err = json.Unmarshal(raw, &rmsg)
		if err != nil {
			logger.Tentacle.Println("Unmarshal", err)
			*result = "MasterError"
			return
		}
		logger.Tentacle.Print(rmsg.Msg)
		if rmsg.Msg != "OK" {
			*result = "NodeError"
		}
	}
}

func ScenarioInfo(ctx *gin.Context) {

}

func ScenariosInfo(ctx *gin.Context) {

}

func ScenarioReset(ctx *gin.Context) {

}
