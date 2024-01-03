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
)

func PakmaCmd(ctx *gin.Context) {
	params := PakmaParams{}
	params.Command = ctx.PostForm("command")
	params.Version = ctx.PostForm("version")
	params.Time = ctx.PostForm("time")
	params.Limit, _ = strconv.Atoi(ctx.PostForm("limit"))

	targetNodes := ctx.PostForm("targetNodes")

	rmsg := protocols.Result{
		Rmsg: "OK",
	}

	if params.Command == "" || len(targetNodes) == 0 {
		logger.Request.Println("PakmaCmd Args Error")
		rmsg.Rmsg = "ERORR: arguments"
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	nodes := []string{}
	err := config.Jsoner.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Rmsg = "ERROR: targetNodes"
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	payload, _ := config.Jsoner.Marshal(&params)
	results := make([]protocols.ExecutionResult, len(nodes))
	var wg sync.WaitGroup

	for i := range nodes {
		name := nodes[i]
		results[i].Name = name
		wg.Add(1)
		if name == "brain" {
			go pakmaLocal(params, &wg, &results[i])
		} else {
			go pakmaRemote(name, payload, &wg, &results[i])
		}
	}
	wg.Wait()
	ctx.JSON(http.StatusOK, results)
}

func pakmaRemote(name string, payload []byte, wg *sync.WaitGroup, result *protocols.ExecutionResult) {
	defer wg.Done()
	result.Code = protocols.ExecOK

	raw, err := model.Request(name, protocols.TypePakmaCommand, payload)
	if err != nil {
		emsg := fmt.Sprintf("Request error %v", err)
		logger.Comm.Println(emsg)
		result.Code = protocols.ExecCommunicationError
		result.CommunicationErrorMsg = emsg
		return
	}
	result.Result = string(raw)
}

func pakmaLocal(params PakmaParams, wg *sync.WaitGroup, result *protocols.ExecutionResult) {
	defer wg.Done()
	result.Code = protocols.ExecOK

	var response []byte
	var err error
	var valid bool = true

	switch params.Command {
	case "install":
		if !checkVersion(params.Version) {
			valid = false
			emsg := "PakmaCommand invalid version"
			logger.Exceptions.Println(emsg)
			result.Code = protocols.ExecCommunicationError
			result.CommunicationErrorMsg = emsg
			break
		}
		response, err = pakmaInstall(params.Version)
	case "upgrade":
		if !checkVersion(params.Version) {
			valid = false
			emsg := "PakmaCommand invalid version"
			logger.Exceptions.Println(emsg)
			result.Code = protocols.ExecCommunicationError
			result.CommunicationErrorMsg = emsg
			break
		}
		response, err = pakmaUpgrade(params.Version)
	case "state":
		response, err = pakmaState()
	case "confirm":
		response, err = pakmaConfirm()
	case "cancel":
		response, err = pakmaCancel()
	case "clean":
		response, err = pakmaClean()
	case "downgrade":
		response, err = pakmaDowngrade()
	case "history":
		response, err = pakmaHistory(params.Time, params.Limit)
	default:
		valid = false
		emsg := "PakmaCommand unsupport command"
		logger.Exceptions.Println(emsg)
		result.Code = protocols.ExecCommunicationError
		result.CommunicationErrorMsg = emsg
	}
	if valid {
		// fmt.Println(string(response), err)
		if err != nil {
			emsg := "Pakma request error"
			logger.Exceptions.Println(emsg)
			result.Code = protocols.ExecCommunicationError
			result.CommunicationErrorMsg = emsg
		} else {
			result.Result = string(response)
		}
	}
}
