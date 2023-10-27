package api

import (
	"encoding/base64"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/protocols"
)

type AppBasic struct {
	Name        string
	Scenario    string
	Description string
	Message     string
}

type AppCreateParams struct {
	AppBasic
	FilePack string
}

type AppDeployParams struct {
	AppBasic
	Script string
}

type CommitResults BasicNodeResults

func AppPrepare(ctx *gin.Context) {
	appName := ctx.PostForm("appName")
	scenario := ctx.PostForm("scenario")
	description := ctx.PostForm("description")
	messages := ctx.PostForm("message")
	targetNodes := ctx.PostForm("targetNodes")
	files, err := ctx.FormFile("files")
	rmsg := protocols.Result{
		Rmsg: "OK",
	}

	if len(appName) == 0 || len(scenario) == 0 || len(description) == 0 || len(messages) == 0 || len(targetNodes) == 0 || err != nil {
		rmsg.Rmsg = "ERROR: Wrong Args"
		ctx.JSON(400, rmsg)
		return
	}

	if _, exists := model.GetScenarioInfoByName(scenario); !exists {
		rmsg.Rmsg = "ERROR: Scenario Not Exists"
		ctx.JSON(404, rmsg)
		return
	}

	nodes := []string{}
	err = config.Jsoner.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Rmsg = "targetNodes:" + err.Error()
		ctx.JSON(400, rmsg)
		return
	}

	multipart, err := files.Open()
	if err != nil {
		rmsg.Rmsg = "Open:" + err.Error()
		ctx.JSON(400, rmsg)
		return
	}

	raw, err := io.ReadAll(multipart)
	multipart.Close()
	if err != nil {
		rmsg.Rmsg = "ReadAll:" + err.Error()
		ctx.JSON(400, rmsg)
		return
	}

	taskid := model.BrainTaskManager.CreateTask(len(nodes))
	ctx.String(202, taskid)

	// async processing
	go func() {
		content := base64.RawStdEncoding.EncodeToString(raw)
		acParams := AppCreateParams{
			AppBasic{
				Name:        appName,
				Scenario:    scenario,
				Description: description,
				Message:     messages,
			},
			content,
		}

		for i := range nodes {
			go createApp(taskid, nodes[i], &acParams)
		}
	}()
}

func AppDeploy(ctx *gin.Context) {
	appName := ctx.PostForm("appName")
	scenario := ctx.PostForm("scenario")
	description := ctx.PostForm("description")
	messages := ctx.PostForm("message")
	targetNodes := ctx.PostForm("targetNodes")
	file, err := ctx.FormFile("script")
	rmsg := protocols.Result{
		Rmsg: "OK",
	}

	if len(appName) == 0 || len(scenario) == 0 || len(description) == 0 || len(messages) == 0 || len(targetNodes) == 0 || err != nil {
		rmsg.Rmsg = "ERROR: Wrong Args"
		ctx.JSON(400, rmsg)
		return
	}

	if _, exists := model.GetScenarioInfoByName(scenario); !exists {
		rmsg.Rmsg = "ERROR: Scenario Not Exists"
		ctx.JSON(404, rmsg)
		return
	}

	nodes := []string{}
	err = config.Jsoner.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Rmsg = "targetNodes:" + err.Error()
		ctx.JSON(400, rmsg)
		return
	}

	multipart, err := file.Open()
	if err != nil {
		rmsg.Rmsg = "Open:" + err.Error()
		ctx.JSON(400, rmsg)
		return
	}
	defer multipart.Close()

	raw, err := io.ReadAll(multipart)
	if err != nil {
		rmsg.Rmsg = "ReadAll:" + err.Error()
		ctx.JSON(400, rmsg)
		return
	}

	taskid := model.BrainTaskManager.CreateTask(len(nodes))
	ctx.String(202, taskid)

	// async processing
	go func() {
		content := base64.RawStdEncoding.EncodeToString(raw)
		adParams := AppDeployParams{
			AppBasic{
				Name:        appName,
				Scenario:    scenario,
				Description: description,
				Message:     messages,
			},
			content,
		}
		for i := range nodes {
			go deployApp(taskid, nodes[i], &adParams)
		}
	}()
}

func createApp(taskid string, name string, acParams *AppCreateParams) {
	payload, _ := config.Jsoner.Marshal(acParams)
	rmsg, err := runAndWait(taskid, name, payload, protocols.TypeAppCreate)
	if err != nil {
		return
	}
	// update scenario version
	success := model.AddScenNodeApp(acParams.Scenario, acParams.Name, acParams.Description, name, rmsg.Version, rmsg.Modified)
	if success {
		logger.Request.Print("Success: AddScenNodeApp")
	} else {
		logger.Exceptions.Printf("failed to add nodeapp to scen: name=%s, %s@%s", name, acParams.Name, acParams.Scenario)
		// *result = "Failed: AddScenNodeApp"
	}
}

func deployApp(taskid string, name string, adParams *AppDeployParams) {
	payload, _ := config.Jsoner.Marshal(adParams)
	rmsg, err := runAndWait(taskid, name, payload, protocols.TypeAppDeploy)
	if err != nil {
		return
	}

	// update scenario version
	success := model.AddScenNodeApp(adParams.Scenario, adParams.Name, adParams.Description, name, rmsg.Version, rmsg.Modified)
	if success {
		logger.Request.Print("Success: AddScenNodeApp")
	} else {
		logger.Exceptions.Printf("failed to add nodeapp to scen: name=%s, %s@%s", name, adParams.Name, adParams.Scenario)
		// *result = "Failed: AddScenNodeApp"
	}
}
