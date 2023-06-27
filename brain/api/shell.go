package api

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"brain/model"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
)

type ScriptParams struct {
	FileName   string
	TargetPath string
	FileBuf    string
}

func RunScript(ctx *gin.Context) {
	script, _ := ctx.FormFile("script")
	targetNodes := ctx.PostForm("targetNodes")
	rmsg := message.Result{
		Rmsg: "OK",
	}

	if script.Size == 0 || len(targetNodes) == 0 {
		logger.Request.Println("RunScript Args Error")
		rmsg.Rmsg = "ERORR: arguments"
		ctx.JSON(400, rmsg)
		return
	}

	f, err := script.Open()
	if err != nil {
		rmsg.Rmsg = "Open:" + err.Error()
		ctx.JSON(400, rmsg)
		return
	}
	defer f.Close()

	raw, _ := io.ReadAll(f)
	content := base64.RawStdEncoding.EncodeToString(raw)
	sparams := ScriptParams{
		FileName:   script.Filename,
		TargetPath: "scripts/",
		FileBuf:    content,
	}
	payload, _ := config.Jsoner.Marshal(&sparams)

	nodes := []string{}
	err = config.Jsoner.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Rmsg = "ERROR: targetNodes"
		ctx.JSON(400, rmsg)
		return
	}

	results := make([]BasicNodeResults, len(nodes))
	var wg sync.WaitGroup

	for i := range nodes {
		name := nodes[i]
		results[i].Name = name
		// if addr, exists := model.GetNodeAddress(name); exists {
		// 	wg.Add(1)
		// 	go runScript(addr, payload, &wg, &results[i].Result)
		// } else {
		// 	results[i].Result = "NodeNotExists"
		// }
		wg.Add(1)
		go runScript(name, payload, &wg, &results[i].Result)
	}
	wg.Wait()
	ctx.JSON(200, results)
}

type CommandParams struct {
	Command    string
	Background bool
}

func RunCmd(ctx *gin.Context) {
	cmd := ctx.PostForm("command")
	bg := ctx.PostForm("background")
	var isbg bool = false
	if len(bg) != 0 {
		isbg = true
	}
	targetNodes := ctx.PostForm("targetNodes")
	rmsg := message.Result{
		Rmsg: "OK",
	}

	if len(cmd) == 0 || len(targetNodes) == 0 {
		logger.Request.Println("RunCmd Args Error")
		rmsg.Rmsg = "ERORR: arguments"
		ctx.JSON(400, rmsg)
		return
	}
	cParams := CommandParams{
		Command:    cmd,
		Background: isbg,
	}
	payload, _ := json.Marshal(cParams)

	nodes := []string{}
	err := config.Jsoner.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Rmsg = "Unmarshal:" + err.Error()
		ctx.JSON(400, rmsg)
		return
	}

	results := make([]BasicNodeResults, len(nodes))
	var wg sync.WaitGroup

	for i := range nodes {
		name := nodes[i]
		results[i].Name = name
		// if addr, exists := model.GetNodeAddress(name); exists {
		// 	wg.Add(1)
		// 	go runCmd(addr, cmd, &wg, &results[i].Result)
		// } else {
		// 	results[i].Result = "NodeNotExists"
		// }
		wg.Add(1)
		go runCmd(name, payload, &wg, &results[i].Result)
	}
	wg.Wait()
	ctx.JSON(200, results)
}

func runCmd(name string, payload []byte, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "UnknownError"

	raw, err := model.Request(name, message.TypeRunCommand, payload)
	if err != nil {
		logger.Comm.Println("Request", err)
		*result = "Request error"
		return
	}
	var rmsg message.Result
	err = config.Jsoner.Unmarshal(raw, &rmsg)
	if err != nil {
		logger.Exceptions.Println("UnmarshalRuncmd", err)
		*result = "BrainError"
		return
	}
	// logger.SysInfo.Print(rmsg.Rmsg)
	// *result = rmsg.Output
	*result = fmt.Sprintf("[%s]\n%s", rmsg.Rmsg, rmsg.Output)
}

func runScript(name string, payload []byte, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "UnknownError"

	raw, err := model.Request(name, message.TypeRunScript, payload)
	if err != nil {
		logger.Comm.Println("Request", err)
		*result = "Request error"
		return
	}
	var rmsg message.Result
	err = config.Jsoner.Unmarshal(raw, &rmsg)
	if err != nil {
		logger.Exceptions.Println("UnmarshalRunscript", err)
		*result = "Brain Error"
		return
	}
	// logger.SysInfo.Print(rmsg.Rmsg)
	*result = fmt.Sprintf("[%s]\n%s", rmsg.Rmsg, rmsg.Output)
}
