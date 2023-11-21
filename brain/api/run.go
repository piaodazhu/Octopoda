package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/ostp"
)

func RunScript(ctx *gin.Context) {
	script, _ := ctx.FormFile("script")
	delayStr := ctx.PostForm("delayTime")
	targetNodes := ctx.PostForm("targetNodes")
	rmsg := protocols.Result{
		Rmsg: "OK",
	}

	if script.Size == 0 || len(targetNodes) == 0 {
		logger.Request.Println("RunScript Args Error")
		rmsg.Rmsg = "ERORR: arguments"
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	var delay int
	var err error
	if delay, err = strconv.Atoi(delayStr); err != nil {
		logger.Request.Println("RunScript Delay Arg Error: ", err.Error())
		rmsg.Rmsg = "ERORR: arguments: " + err.Error()
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	f, err := script.Open()
	if err != nil {
		rmsg.Rmsg = "Open:" + err.Error()
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}
	defer f.Close()

	raw, _ := io.ReadAll(f)
	content := base64.RawStdEncoding.EncodeToString(raw)

	nodes := []string{}
	err = config.Jsoner.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Rmsg = "ERROR: targetNodes"
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	sparams := protocols.ScriptParams{
		FileName:   script.Filename,
		TargetPath: "scripts/",
		FileBuf:    content,
		DelayTime:  delay,
		ExecTs: ostp.ExtimateExecTs(model.GetNodesMaxDelay(nodes)),
	}
	payload, _ := config.Jsoner.Marshal(&sparams)

	taskid := model.BrainTaskManager.CreateTask(len(nodes))
	ctx.String(http.StatusAccepted, taskid)

	// async processing
	go func() {
		for i := range nodes {
			go runScript(taskid, nodes[i], payload)
		}
	}()
}

func runScript(taskid string, name string, payload []byte) {
	runAndWait(taskid, name, payload, protocols.TypeRunScript)
}

func RunCmd(ctx *gin.Context) {
	cmd := ctx.PostForm("command")
	bg := ctx.PostForm("background")
	delayStr := ctx.PostForm("delayTime")
	var isbg bool = false
	if len(bg) != 0 {
		isbg = true
	}
	targetNodes := ctx.PostForm("targetNodes")
	rmsg := protocols.Result{
		Rmsg: "OK",
	}

	if len(cmd) == 0 || len(targetNodes) == 0 {
		logger.Request.Println("RunCmd Args Error")
		rmsg.Rmsg = "ERORR: arguments"
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	var delay int
	var err error
	if delay, err = strconv.Atoi(delayStr); err != nil {
		logger.Request.Println("RunScript Delay Arg Error: ", err.Error())
		rmsg.Rmsg = "ERORR: arguments: " + err.Error()
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	nodes := []string{}
	err = config.Jsoner.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Rmsg = "Unmarshal:" + err.Error()
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	cParams := protocols.CommandParams{
		Command:    cmd,
		Background: isbg,
		DelayTime:  delay,
		ExecTs:     ostp.ExtimateExecTs(model.GetNodesMaxDelay(nodes)),
	}
	payload, _ := json.Marshal(cParams)

	taskid := model.BrainTaskManager.CreateTask(len(nodes))
	ctx.String(http.StatusAccepted, taskid)

	// async processing
	go func() {
		for i := range nodes {
			go runCmd(taskid, nodes[i], payload)
		}
	}()
}

func runCmd(taskid string, name string, payload []byte) {
	runAndWait(taskid, name, payload, protocols.TypeRunCommand)
}

func runAndWait(taskid string, name string, payload []byte, rtype int) (*protocols.Result, error) {
	var rstr string
	result := protocols.ExecutionResults{
		Name: name,
		Code: protocols.ExecOK,
	}

	if rtype == protocols.TypeRunScript {
		rstr = "script"
	} else if rtype == protocols.TypeRunCommand {
		rstr = "command"
	} else if rtype == protocols.TypeAppCreate {
		rstr = "createApp"
	} else if rtype == protocols.TypeAppDeploy {
		rstr = "deployApp"
	} else {
		logger.Comm.Println("unsupported rtype in runAndWait: ", rtype)
		result.Code = protocols.ExecCommunicationError
		result.CommunicationErrorMsg = "unsupported rtype"
		model.BrainTaskManager.AddFailedSubTask(taskid, model.TaskIdGen(), &result)
		return nil, errors.New(result.CommunicationErrorMsg)
	}

	subid, err := model.Request(name, rtype, payload)
	if err != nil || len(subid) == 0 {
		emsg := fmt.Sprintf("Send %s request: %v", rstr, err)
		logger.Comm.Println(emsg)
		result.Code = protocols.ExecCommunicationError
		result.CommunicationErrorMsg = emsg
		model.BrainTaskManager.AddFailedSubTask(taskid, model.TaskIdGen(), &result)
		return nil, errors.New(emsg)
	}
	subtask_id := string(subid)
	model.BrainTaskManager.AddSubTask(taskid, subtask_id, &result)
	defer model.BrainTaskManager.DoneSubTask(taskid, subtask_id, &result)

	raw, err := model.Request(name, protocols.TypeWaitTask, subid)
	if err != nil || len(raw) == 0 {
		emsg := fmt.Sprintf("Send %s wait task request error: %v", rstr, err)
		logger.Comm.Println(emsg)
		result.Code = protocols.ExecCommunicationError
		result.CommunicationErrorMsg = emsg
		return nil, errors.New(emsg)
	}

	var rmsg protocols.Result
	err = config.Jsoner.Unmarshal(raw, &rmsg)
	if err != nil {
		emsg := fmt.Sprintf("unmarshal %s response error: %v", rstr, err)
		logger.Comm.Println(emsg)
		result.Code = protocols.ExecCommunicationError
		result.CommunicationErrorMsg = emsg
		return nil, errors.New(emsg)
	}

	if rmsg.Rmsg != "OK" {
		result.Code = protocols.ExecProcessError
		result.ProcessErrorMsg = rmsg.Rmsg
	}
	result.Result = rmsg.Output
	return &rmsg, nil
}

func CancelRun(ctx *gin.Context) {
	taskid := ctx.PostForm("taskid")
	ctx.Status(200)

	plist := model.BrainTaskManager.PendingList(taskid)
	for i := range plist {
		subtaskid := plist[i].SubTaskId
		name := plist[i].Result.(*protocols.ExecutionResults).Name
		go cancelNodeRun(subtaskid, name)
	}
}

func cancelNodeRun(subtaskid, name string) {
	raw, err := model.Request(name, protocols.TypeCancelTask, []byte(subtaskid))
	if err != nil || len(raw) == 0 {
		emsg := fmt.Sprintf("Send cancel task to %s:%s request error: %v", name, subtaskid, err)
		logger.Comm.Println(emsg)
		return
	}
}
