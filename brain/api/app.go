package api

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"brain/model"
	"brain/rdb"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/gin-gonic/gin"
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
	rmsg := message.Result{
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

	// fast return
	taskid := rdb.TaskIdGen()
	if !rdb.TaskNew(taskid, 3600) {
		logger.Exceptions.Print("TaskNew")
	}
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
		// payload, _ := config.Jsoner.Marshal(&acParams)

		// check target nodes
		// spread tar file
		results := make([]CommitResults, len(nodes))
		var wg sync.WaitGroup

		for i := range nodes {
			name := nodes[i]
			results[i].Name = name
			if addr, exists := model.GetNodeAddress(name); exists {
				wg.Add(1)
				go createApp(name, addr, &acParams, &wg, &results[i].Result)
			} else {
				results[i].Result = "NodeNotExists"
			}
		}
		wg.Wait()
		// logger.Brain.Println(results)
		if !rdb.TaskMarkDone(taskid, results, 3600) {
			logger.Exceptions.Print("TaskMarkDone")
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
	rmsg := message.Result{
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

	// fast return
	taskid := rdb.TaskIdGen()
	if !rdb.TaskNew(taskid, 3600) {
		logger.Exceptions.Print("TaskNew")
	}
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
		// payload, _ := config.Jsoner.Marshal(&adParams)

		// check target nodes
		// run scripts
		results := make([]CommitResults, len(nodes))
		var wg sync.WaitGroup

		for i := range nodes {
			name := nodes[i]
			results[i].Name = name
			if addr, exists := model.GetNodeAddress(name); exists {
				wg.Add(1)
				go deployApp(name, addr, &adParams, &wg, &results[i].Result)
			} else {
				results[i].Result = "NodeNotExists"
			}
		}
		wg.Wait()
		// logger.Brain.Println(results)
		if !rdb.TaskMarkDone(taskid, results, 3600) {
			logger.Exceptions.Print("TaskMarkDone Error")
		}
	}()
}

func createApp(node string, addr string, acParams *AppCreateParams, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "UnknownError"

	if conn, err := net.Dial("tcp", addr); err != nil {
		return
	} else {
		defer conn.Close()

		payload, _ := config.Jsoner.Marshal(acParams)
		message.SendMessage(conn, message.TypeAppCreate, payload)
		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeAppCreateResponse {
			logger.Comm.Println("TypeAppCreateResponse", err)
			*result = "NetError"
			return
		}

		var rmsg message.Result
		err = config.Jsoner.Unmarshal(raw, &rmsg)
		if err != nil {
			logger.Exceptions.Println("UnmarshalNodeState", err)
			*result = "MasterError"
			return
		}
		// logger.Tentacle.Print(rmsg.Rmsg)
		// *result = rmsg.Rmsg
		*result = fmt.Sprintf("[%s]: %s", rmsg.Rmsg, rmsg.Output)

		// update scenario version
		success := model.AddScenNodeApp(acParams.Scenario, acParams.Name, acParams.Description, node, rmsg.Version, rmsg.Modified)
		if success {
			logger.Request.Print("Success: AddScenNodeApp")
		} else {
			logger.Exceptions.Print("Failed: AddScenNodeApp")
			*result = "Failed: AddScenNodeApp"
		}
	}
}

func deployApp(node string, addr string, adParams *AppDeployParams, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "UnknownError"

	if conn, err := net.Dial("tcp", addr); err != nil {
		return
	} else {
		defer conn.Close()

		payload, _ := config.Jsoner.Marshal(adParams)
		message.SendMessage(conn, message.TypeAppDeploy, payload)
		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeAppDeployResponse {
			logger.Comm.Println("TypeAppDeployResponse", err)
			*result = "NetError"
			return
		}

		var rmsg message.Result
		err = config.Jsoner.Unmarshal(raw, &rmsg)
		if err != nil {
			logger.Exceptions.Println("UnmarshalNodeState", err)
			*result = "MasterError"
			return
		}
		// *result = string(rmsg.Output)
		// *result = rmsg.Output
		*result = fmt.Sprintf("[%s]: %s", rmsg.Rmsg, rmsg.Output)

		// update scenario version
		success := model.AddScenNodeApp(adParams.Scenario, adParams.Name, adParams.Description, node, rmsg.Version, rmsg.Modified)
		if success {
			logger.Request.Print("Success: AddScenNodeApp")
		} else {
			logger.Exceptions.Print("Failed: AddScenNodeApp")
			*result = "Failed: AddScenNodeApp"
		}
	}
}
