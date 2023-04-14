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

func ScenarioUpdate(ctx *gin.Context) {
	name := ctx.PostForm("name")
	message := ctx.PostForm("message")
	rmsg := RMSG{"OK"}
	if len(name) == 0 || len(message) == 0 {
		rmsg.Msg = "ERROR: Wrong Args"
		ctx.JSON(400, rmsg)
		return
	}
	if !model.UpdateScenario(name, message) {
		rmsg.Msg = "ERROR: UpdateScenario"
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
	name := ctx.Query("name")
	rmsg := RMSG{"OK"}
	if len(name) == 0 {
		rmsg.Msg = "ERROR: Wrong Args"
		ctx.JSON(400, rmsg)
		return
	}

	// must exists
	if _, exists := model.GetScenarioInfoByName(name); !exists {
		rmsg.Msg = "ERROR: Scenario Not Exists"
		ctx.JSON(404, rmsg)
		return
	}

	// delete all app in nodes
	dlist := model.GetNodeApps(name, "")

	// check target nodes
	// delete app
	results := make([]BasicNodeResults, len(dlist))
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
		if rmsg.Msg != "OK" {
			*result = "NodeError"
		}
	}
}

func ScenarioInfo(ctx *gin.Context) {
	var name string
	var ok bool
	var scen *model.ScenarioInfo
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	if scen, ok = model.GetScenarioInfoByName(name); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	ctx.JSON(200, scen)
}

func ScenariosInfo(ctx *gin.Context) {
	var scens []model.ScenarioDigest
	var ok bool

	if scens, ok = model.GetScenariosDigestAll(); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	ctx.JSON(200, scens)
}

func ScenarioVersion(ctx *gin.Context) {
	var versions []model.BasicVersionModel
	var name string
	var ok bool
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	versions = model.GetScenarioVersionByName(name)
	ctx.JSON(200, versions)
}

type AppResetParams struct {
	AppBasic
	VersionHash string
	Mode        string
}

func ScenarioReset(ctx *gin.Context) {
	name := ctx.PostForm("name")
	prefix := ctx.PostForm("version")
	message := ctx.PostForm("message")
	rmsg := RMSG{"OK"}
	if len(name) == 0 || len(message) == 0 || len(prefix) < 2 || len(prefix) > 40 {
		rmsg.Msg = "ERROR: Wrong Args. (should specific name and prefix. prefix length should be in [2, 40])"
		ctx.JSON(400, rmsg)
		return
	}

	// must exists
	if _, exists := model.GetScenarioInfoByName(name); !exists {
		rmsg.Msg = "ERROR: Scenario Not Exists"
		ctx.JSON(404, rmsg)
		return
	}

	vlist := model.GetScenarioVersionByName(name)
	version := ""
	ambiguity := false
	for i := range vlist {
		if vlist[i].Version[:len(prefix)] == prefix {
			// prefix matched
			if version != "" {
				// ambiguity arises
				ambiguity = true
				break
			}
			// assign complete version string to version
			version = vlist[i].Version
		}
	}
	if len(version) == 0 {
		rmsg.Msg = "ERROR: Version Not Exists"
		ctx.JSON(404, rmsg)
		return
	}
	if ambiguity {
		rmsg.Msg = "ERROR: Version Ambiguity"
		ctx.JSON(404, rmsg)
		return
	}

	// get reset nodeapp list
	// then reset all nodeapp
	rlist := model.GetNodeApps(name, version)

	// check target nodes
	results := make([]BasicNodeResults, len(rlist))
	var wg sync.WaitGroup

	for i := range rlist {
		//payload
		arg := &AppResetParams{
			AppBasic: AppBasic{
				Name: rlist[i].AppName,
				Scenario: rlist[i].ScenName,
				Message: message,
			},
			VersionHash: rlist[i].Version,
			Mode: "undef",
		}
		payload, _ := json.Marshal(arg)

		// item name
		item := rlist[i].NodeName + ":" + rlist[i].AppName + "@" + rlist[i].ScenName
		results[i].Name = item

		// node name
		nodename := rlist[i].NodeName
		if addr, exists := model.GetNodeAddress(nodename); exists {
			wg.Add(1)
			// payload?
			go resetApp(addr, payload, &wg, &results[i].Result)
		} else {
			results[i].Result = "NodeNotExists"
		}
	}
	wg.Wait()

	// finally reset this scenario locallly
	model.ResetScenario(name, version, message)
	ctx.JSON(200, results)
}

func resetApp(addr string, payload []byte, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "OK"

	if conn, err := net.Dial("tcp", addr); err != nil {
		return
	} else {
		defer conn.Close()

		message.SendMessage(conn, message.TypeAppReset, payload)
		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeAppResetResponse {
			logger.Tentacle.Println("TypeAppResetResponse", err)
			*result = "NetError"
			return
		}

		var rmsg RDATA
		err = json.Unmarshal(raw, &rmsg)
		if err != nil {
			logger.Tentacle.Println("Unmarshal", err)
			*result = "MasterError"
			return
		}
		*result = rmsg.Msg
	}
}

func ScenarioFix(ctx *gin.Context) {
	var name string
	var ok bool
	var rmsg RMSG
	rmsg.Msg = "OK"

	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	err := model.Fix(name)
	if err != nil {
		rmsg.Msg = err.Error()
		ctx.JSON(400, rmsg)
	}
	ctx.JSON(200, rmsg)
}