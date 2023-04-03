package api

import (
	"brain/logger"
	"brain/message"
	"brain/model"
	"encoding/base64"
	"encoding/json"
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

func AppPrepare(ctx *gin.Context) {
	appName := ctx.PostForm("appName")
	scenario := ctx.PostForm("scenario")
	description := ctx.PostForm("description")
	messages := ctx.PostForm("message")
	targetNodes := ctx.PostForm("targetNodes")
	files, err := ctx.FormFile("files")
	rmsg := RMSG{}

	if len(appName) == 0 || len(scenario) == 0 || len(description) == 0 || len(messages) == 0 || len(targetNodes) == 0 || err != nil {
		rmsg.Msg = "ERROR: Wrong Args"
		ctx.JSON(400, rmsg)
		return
	}

	nodes := []string{}
	err = json.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Msg = "ERROR: targetNodes"
		ctx.JSON(400, rmsg)
		return
	}

	multipart, err := files.Open()
	if err != nil {
		rmsg.Msg = "ERROR: Open File"
		ctx.JSON(400, rmsg)
		return
	}
	defer multipart.Close()

	raw, err := io.ReadAll(multipart)
	if err != nil {
		rmsg.Msg = "ERROR: Read File"
		ctx.JSON(400, rmsg)
		return
	}

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
	payload, _ := json.Marshal(&acParams)

	// check target nodes
	// spread tar file
	results := make([]UploadResults, len(nodes))
	var wg sync.WaitGroup

	for i := range nodes {
		name := nodes[i]
		results[i].Name = name
		if addr, exists := model.GetNodeAddress(name); exists {
			wg.Add(1)
			go createApp(addr, payload, &wg, &results[i].Result)
		} else {
			results[i].Result = "NodeNotExists"
		}
	}
	wg.Wait()
	ctx.JSON(200, results)
}

func AppDeploy(ctx *gin.Context) {
	appName := ctx.PostForm("appName")
	scenario := ctx.PostForm("scenario")
	description := ctx.PostForm("description")
	messages := ctx.PostForm("message")
	targetNodes := ctx.PostForm("targetNodes")
	file, err := ctx.FormFile("script")
	rmsg := RMSG{}

	if len(appName) == 0 || len(scenario) == 0 || len(description) == 0 || len(messages) == 0 || len(targetNodes) == 0 || err != nil {
		rmsg.Msg = "ERROR: Wrong Args"
		ctx.JSON(400, rmsg)
		return
	}

	nodes := []string{}
	err = json.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Msg = "ERROR: targetNodes"
		ctx.JSON(400, rmsg)
		return
	}

	multipart, err := file.Open()
	if err != nil {
		rmsg.Msg = "ERROR: Open File"
		ctx.JSON(400, rmsg)
		return
	}
	defer multipart.Close()

	raw, err := io.ReadAll(multipart)
	if err != nil {
		rmsg.Msg = "ERROR: Read File"
		ctx.JSON(400, rmsg)
		return
	}

	content := base64.RawStdEncoding.EncodeToString(raw)
	adParams := AppDeployParams{
		AppBasic{
			Name: appName,
			Scenario: scenario,
			Description: description,
			Message: messages,
		},
		content,
	}
	payload, _ := json.Marshal(&adParams)

	// check target nodes
	// run scripts
	results := make([]UploadResults, len(nodes))
	var wg sync.WaitGroup

	for i := range nodes {
		name := nodes[i]
		results[i].Name = name
		if addr, exists := model.GetNodeAddress(name); exists {
			wg.Add(1)
			go deployApp(addr, payload, &wg, &results[i].Result)
		} else {
			results[i].Result = "NodeNotExists"
		}
	}
	wg.Wait()
	ctx.JSON(200, results)
}

func createApp(addr string, payload []byte, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "OK"

	if conn, err := net.Dial("tcp", addr); err != nil {
		return
	} else {
		defer conn.Close()

		message.SendMessage(conn, message.TypeAppCreate, payload)
		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeAppCreateResponse {
			logger.Tentacle.Println("TypeAppCreateResponse", err)
			*result = "NetError"
			return
		}

		var rmsg RMSG
		err = json.Unmarshal(raw, &rmsg)
		if err != nil {
			logger.Tentacle.Println("UnmarshalNodeState", err)
			*result = "MasterError"
			return
		}
		logger.Tentacle.Print(rmsg.Msg)
		if rmsg.Msg != "OK" {
			*result = "NodeError"
		}
	}
}

func deployApp(addr string, payload []byte, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "OK"

	if conn, err := net.Dial("tcp", addr); err != nil {
		return
	} else {
		defer conn.Close()

		message.SendMessage(conn, message.TypeAppDeploy, payload)
		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeAppDeployResponse {
			logger.Tentacle.Println("TypeAppDeployResponse", err)
			*result = "NetError"
			return
		}

		var rmsg RMSG
		err = json.Unmarshal(raw, &rmsg)
		if err != nil {
			logger.Tentacle.Println("UnmarshalNodeState", err)
			*result = "MasterError"
			return
		}
		logger.Tentacle.Print(rmsg.Msg)
		*result = rmsg.Msg
	}
}