package api

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/san"
)

func ScenarioCreate(ctx *gin.Context) {
	name := ctx.PostForm("name")
	description := ctx.PostForm("description")
	rmsg := protocols.Result{
		Rmsg: "OK",
	}
	if len(name) == 0 || len(description) == 0 {
		rmsg.Rmsg = "ERROR: Wrong Args"
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	if !model.AddScenario(name, description) {
		rmsg.Rmsg = "ERROR: Scenario Exists"
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}
	ctx.JSON(http.StatusOK, rmsg)
}

func ScenarioUpdate(ctx *gin.Context) {
	name := ctx.PostForm("name")
	msg := ctx.PostForm("message")
	rmsg := protocols.Result{
		Rmsg: "OK",
	}
	if len(name) == 0 || len(msg) == 0 {
		rmsg.Rmsg = "Lack scenario name or message"
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	modified, ok := model.UpdateScenario(name, msg)
	if !ok {
		rmsg.Rmsg = "ERROR: UpdateScenario"
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}
	rmsg.Modified = modified
	ctx.JSON(http.StatusOK, rmsg)
}

func ScenarioDelete(ctx *gin.Context) {
	// This is not so simple. Should delete all apps of a scenario in this function
	name := ctx.Query("name")
	rmsg := protocols.Result{
		Rmsg: "OK",
	}
	if len(name) == 0 {
		rmsg.Rmsg = "Lack scenario name"
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	// must exists
	if _, exists := model.GetScenarioInfoByName(name); !exists {
		rmsg.Rmsg = "ERROR: Scenario Not Exists"
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}

	// delete all app in nodes
	dlist := model.GetNodeApps(name, "")

	// check target nodes
	// delete app
	results := make([]protocols.ExecutionResults, len(dlist))
	var wg sync.WaitGroup

	for i := range dlist {
		//payload
		payload, _ := config.Jsoner.Marshal(&san.AppDeleteParams{
			Name: dlist[i].AppName, 
			Scenario: dlist[i].ScenName,
		})

		// item name
		item := dlist[i].NodeName + ":" + dlist[i].AppName + "@" + dlist[i].ScenName
		results[i].Name = item

		// node name
		nodename := dlist[i].NodeName

		wg.Add(1)
		go deleteApp(nodename, payload, &wg, &results[i])
	}
	wg.Wait()

	// finally delete this scenario locallly
	model.DelScenario(name)
	ctx.JSON(http.StatusOK, results)
}

func deleteApp(name string, payload []byte, wg *sync.WaitGroup, result *protocols.ExecutionResults) {
	defer wg.Done()
	result.Code = protocols.ExecOK

	raw, err := model.Request(name, protocols.TypeAppDelete, payload)
	if err != nil {
		emsg := fmt.Sprintf("Send TypeAppDeleteResponse request error: %v", err)
		logger.Comm.Println(emsg)
		result.Code = protocols.ExecCommunicationError
		result.CommunicationErrorMsg = emsg
		return
	}
	var rmsg protocols.Result
	err = config.Jsoner.Unmarshal(raw, &rmsg)
	if err != nil {
		emsg := fmt.Sprintf("unmarshal response error: %v", err)
		logger.Comm.Println(emsg)
		result.Code = protocols.ExecCommunicationError
		result.CommunicationErrorMsg = emsg
		return
	}
	if rmsg.Rmsg != "OK" {
		result.Code = protocols.ExecProcessError
		result.ProcessErrorMsg = rmsg.Rmsg
	}
	result.Result = rmsg.Output
}

func ScenarioInfo(ctx *gin.Context) {
	var name string
	var ok bool
	var scen *san.ScenarioInfo
	rmsg := protocols.Result{
		Rmsg: "OK",
	}
	if name, ok = ctx.GetQuery("name"); !ok {
		rmsg.Rmsg = "Lack scenario name"
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}

	if scen, ok = model.GetScenarioInfoByName(name); !ok {
		rmsg.Rmsg = "Error: GetScenarioInfoByName"
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}
	ctx.JSON(http.StatusOK, scen)
}

func ScenariosInfo(ctx *gin.Context) {
	var scens []san.ScenarioDigest
	var ok bool
	rmsg := protocols.Result{
		Rmsg: "OK",
	}

	if scens, ok = model.GetScenariosDigestAll(); !ok {
		rmsg.Rmsg = "Error: GetScenarioInfoByName"
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}
	ctx.JSON(http.StatusOK, scens)
}

func ScenarioVersion(ctx *gin.Context) {
	var vlist []san.Version
	var name, offsetStr, limitStr string
	var ok bool
	var err error
	rmsg := protocols.Result{
		Rmsg: "OK",
	}

	if name, ok = ctx.GetQuery("name"); !ok {
		rmsg.Rmsg = "Lack Name"
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}

	offset := 0
	limit := 5

	offsetStr = ctx.Query("offset")
	if len(offsetStr) != 0 {
		if offset, err = strconv.Atoi(offsetStr); err != nil {
			rmsg.Rmsg = "ERROR: Wrong Args: offset=" + offsetStr
			ctx.JSON(http.StatusBadRequest, rmsg)
			return
		}
	}

	limitStr = ctx.Query("limit")
	if len(limitStr) != 0 {
		if limit, err = strconv.Atoi(limitStr); err != nil {
			rmsg.Rmsg = "ERROR: Wrong Args: limit=" + limitStr
			ctx.JSON(http.StatusBadRequest, rmsg)
			return
		}
	}

	allVersions := model.GetScenarioVersionByName(name)
	if len(allVersions) <= offset {
		ctx.JSON(http.StatusOK, vlist)
		return
	}
	end := offset + limit
	if len(allVersions) < end {
		end = len(allVersions)
	}

	// reverse list
	offset = len(allVersions) - 1 - offset
	end = len(allVersions) - 1 - end
	for i := offset; i > end; i-- {
		vlist = append(vlist, allVersions[i])
	}

	ctx.JSON(http.StatusOK, vlist)
}

func ScenarioReset(ctx *gin.Context) {
	name := ctx.PostForm("name")
	prefix := ctx.PostForm("version")
	msg := ctx.PostForm("message")
	rmsg := protocols.Result{
		Rmsg: "OK",
	}

	if len(name) == 0 || len(msg) == 0 || len(prefix) < 2 || len(prefix) > 40 {
		rmsg.Rmsg = "ERROR: Wrong Args. (should specific name and prefix. prefix length should be in [2, 40])"
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	// must exists
	if _, exists := model.GetScenarioInfoByName(name); !exists {
		rmsg.Rmsg = "ERROR: Scenario Not Exists"
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}

	vlist := model.GetScenarioVersionByName(name)
	version := ""
	ambiguity := false
	for i := range vlist {
		if vlist[i].Hash[:len(prefix)] == prefix {
			// prefix matched
			if version != "" {
				// ambiguity arises
				ambiguity = true
				break
			}
			// assign complete version string to version
			version = vlist[i].Hash
		}
	}
	if len(version) == 0 {
		rmsg.Rmsg = "ERROR: Version Not Exists"
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}
	if ambiguity {
		rmsg.Rmsg = "ERROR: Version Ambiguity"
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}

	// get reset nodeapp list
	// then reset all nodeapp
	rlist := model.GetNodeApps(name, version)

	// check target nodes
	results := make([]protocols.ExecutionResults, len(rlist))
	var wg sync.WaitGroup

	for i := range rlist {
		//payload
		arg := &san.AppResetParams{
			AppBasic: san.AppBasic{
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
		go resetApp(nodename, payload, &wg, &results[i])
	}
	wg.Wait()

	// finally reset this scenario locallly
	model.ResetScenario(name, version, msg)
	ctx.JSON(http.StatusOK, results)
}

func resetApp(name string, payload []byte, wg *sync.WaitGroup, result *protocols.ExecutionResults) {
	defer wg.Done()
	result.Code = protocols.ExecOK

	raw, err := model.Request(name, protocols.TypeAppReset, payload)
	if err != nil {
		emsg := fmt.Sprintf("Request error %v", err)
		logger.Comm.Println(emsg)
		result.Code = protocols.ExecCommunicationError
		result.CommunicationErrorMsg = emsg
		return
	}
	var rmsg protocols.Result
	err = config.Jsoner.Unmarshal(raw, &rmsg)
	if err != nil {
		emsg := fmt.Sprintf("Brain unmarshal error %v", err)
		logger.Comm.Println(emsg)
		result.Code = protocols.ExecCommunicationError
		result.CommunicationErrorMsg = emsg
		return
	}
	result.Result = rmsg.Rmsg
}

func ScenarioFix(ctx *gin.Context) {
	var name string
	var ok bool
	rmsg := protocols.Result{
		Rmsg: "OK",
	}

	if name, ok = ctx.GetQuery("name"); !ok {
		rmsg.Rmsg = "Lack scenario name"
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}

	err := model.Fix(name)
	if err != nil {
		rmsg.Rmsg = "Fix:" + err.Error()
		ctx.JSON(http.StatusBadRequest, rmsg)
	}
	ctx.JSON(http.StatusOK, rmsg)
}
