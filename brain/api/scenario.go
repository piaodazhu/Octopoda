package api

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"brain/model"
	"sync"

	"github.com/gin-gonic/gin"
)

func ScenarioCreate(ctx *gin.Context) {
	name := ctx.PostForm("name")
	description := ctx.PostForm("description")
	rmsg := message.Result{
		Rmsg: "OK",
	}
	if len(name) == 0 || len(description) == 0 {
		rmsg.Rmsg = "ERROR: Wrong Args"
		ctx.JSON(400, rmsg)
		return
	}
	if !model.AddScenario(name, description) {
		rmsg.Rmsg = "ERROR: Scenario Exists"
		ctx.JSON(404, rmsg)
		return
	}
	ctx.JSON(200, rmsg)
}

func ScenarioUpdate(ctx *gin.Context) {
	name := ctx.PostForm("name")
	msg := ctx.PostForm("message")
	rmsg := message.Result{
		Rmsg: "OK",
	}
	if len(name) == 0 || len(msg) == 0 {
		rmsg.Rmsg = "Lack scenario name or message"
		ctx.JSON(400, rmsg)
		return
	}
	modified, ok := model.UpdateScenario(name, msg)
	if !ok {
		rmsg.Rmsg = "ERROR: UpdateScenario"
		ctx.JSON(404, rmsg)
		return
	}
	rmsg.Modified = modified
	ctx.JSON(200, rmsg)
}

type AppDeleteParams struct {
	Name     string
	Scenario string
}

func ScenarioDelete(ctx *gin.Context) {
	// This is not so simple. Should delete all apps of a scenario in this function
	name := ctx.Query("name")
	rmsg := message.Result{
		Rmsg: "OK",
	}
	if len(name) == 0 {
		rmsg.Rmsg = "Lack scenario name"
		ctx.JSON(400, rmsg)
		return
	}

	// must exists
	if _, exists := model.GetScenarioInfoByName(name); !exists {
		rmsg.Rmsg = "ERROR: Scenario Not Exists"
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
		payload, _ := config.Jsoner.Marshal(&AppDeleteParams{dlist[i].AppName, dlist[i].ScenName})

		// item name
		item := dlist[i].NodeName + ":" + dlist[i].AppName + "@" + dlist[i].ScenName
		results[i].Name = item

		// node name
		nodename := dlist[i].NodeName

		wg.Add(1)
		go deleteApp(nodename, payload, &wg, &results[i].Result)
	}
	wg.Wait()

	// finally delete this scenario locallly
	model.DelScenario(name)
	ctx.JSON(200, results)
}

func deleteApp(name string, payload []byte, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "UnknownError"

	raw, err := model.Request(name, message.TypeAppDelete, payload)
	if err != nil {
		logger.Comm.Println("TypeAppDeleteResponse", err)
		*result = "TypeAppDeleteResponse"
		return
	}
	var rmsg message.Result
	err = config.Jsoner.Unmarshal(raw, &rmsg)
	if err != nil {
		logger.Exceptions.Println("Unmarshal", err)
		*result = "BrainError"
		return
	}
	if rmsg.Rmsg != "OK" {
		*result = "NodeError:" + rmsg.Rmsg
	} else {
		*result = "OK"
	}
}

func ScenarioInfo(ctx *gin.Context) {
	var name string
	var ok bool
	var scen *model.ScenarioInfo
	rmsg := message.Result{
		Rmsg: "OK",
	}
	if name, ok = ctx.GetQuery("name"); !ok {
		rmsg.Rmsg = "Lack scenario name"
		ctx.JSON(404, rmsg)
		return
	}
	if scen, ok = model.GetScenarioInfoByName(name); !ok {
		rmsg.Rmsg = "Error: GetScenarioInfoByName"
		ctx.JSON(404, rmsg)
		return
	}
	ctx.JSON(200, scen)
}

func ScenariosInfo(ctx *gin.Context) {
	var scens []model.ScenarioDigest
	var ok bool
	rmsg := message.Result{
		Rmsg: "OK",
	}

	if scens, ok = model.GetScenariosDigestAll(); !ok {
		rmsg.Rmsg = "Error: GetScenarioInfoByName"
		ctx.JSON(404, rmsg)
		return
	}
	ctx.JSON(200, scens)
}

func ScenarioVersion(ctx *gin.Context) {
	var versions []model.BasicVersionModel
	var name string
	var ok bool
	rmsg := message.Result{
		Rmsg: "OK",
	}

	if name, ok = ctx.GetQuery("name"); !ok {
		rmsg.Rmsg = "Lack Name"
		ctx.JSON(404, rmsg)
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
	msg := ctx.PostForm("message")
	rmsg := message.Result{
		Rmsg: "OK",
	}

	if len(name) == 0 || len(msg) == 0 || len(prefix) < 2 || len(prefix) > 40 {
		rmsg.Rmsg = "ERROR: Wrong Args. (should specific name and prefix. prefix length should be in [2, 40])"
		ctx.JSON(400, rmsg)
		return
	}

	// must exists
	if _, exists := model.GetScenarioInfoByName(name); !exists {
		rmsg.Rmsg = "ERROR: Scenario Not Exists"
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
		rmsg.Rmsg = "ERROR: Version Not Exists"
		ctx.JSON(404, rmsg)
		return
	}
	if ambiguity {
		rmsg.Rmsg = "ERROR: Version Ambiguity"
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
				Name:     rlist[i].AppName,
				Scenario: rlist[i].ScenName,
				Message:  msg,
			},
			VersionHash: rlist[i].Version,
			Mode:        "undef",
		}
		payload, _ := config.Jsoner.Marshal(arg)

		// item name
		item := rlist[i].NodeName + ":" + rlist[i].AppName + "@" + rlist[i].ScenName
		results[i].Name = item

		// node name
		nodename := rlist[i].NodeName
		wg.Add(1)
		go resetApp(nodename, payload, &wg, &results[i].Result)
	}
	wg.Wait()

	// finally reset this scenario locallly
	model.ResetScenario(name, version, msg)
	ctx.JSON(200, results)
}

func resetApp(name string, payload []byte, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "UnknownError"

	raw, err := model.Request(name, message.TypeAppReset, payload)
	if err != nil {
		logger.Comm.Println("TypeAppResetResponse", err)
		*result = "TypeAppResetResponse"
		return
	}
	var rmsg message.Result
	err = config.Jsoner.Unmarshal(raw, &rmsg)
	if err != nil {
		logger.Exceptions.Println("Unmarshal", err)
		*result = "BrainError"
		return
	}
	*result = rmsg.Rmsg
}

func ScenarioFix(ctx *gin.Context) {
	var name string
	var ok bool
	rmsg := message.Result{
		Rmsg: "OK",
	}

	if name, ok = ctx.GetQuery("name"); !ok {
		rmsg.Rmsg = "Lack scenario name"
		ctx.JSON(404, rmsg)
		return
	}
	err := model.Fix(name)
	if err != nil {
		rmsg.Rmsg = "Fix:" + err.Error()
		ctx.JSON(400, rmsg)
	}
	ctx.JSON(200, rmsg)
}
