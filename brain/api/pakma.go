package api

import (
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
		ctx.JSON(400, rmsg)
		return
	}

	nodes := []string{}
	err := config.Jsoner.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Rmsg = "ERROR: targetNodes"
		ctx.JSON(400, rmsg)
		return
	}

	payload, _ := config.Jsoner.Marshal(&params)
	results := make([]BasicNodeResults, len(nodes))
	var wg sync.WaitGroup

	for i := range nodes {
		name := nodes[i]
		results[i].Name = name
		wg.Add(1)
		if name == "brain" {
			go pakmaLocal(params, &wg, &results[i].Result)
		} else {
			go pakmaRemote(name, payload, &wg, &results[i].Result)
		}
	}
	wg.Wait()
	ctx.JSON(200, results)
}

func pakmaRemote(name string, payload []byte, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "UnknownError"

	raw, err := model.Request(name, protocols.TypePakmaCommand, payload)
	if err != nil {
		logger.Comm.Println("Request", err)
		*result = "Request error"
		return
	}
	*result = string(raw)
}

func pakmaLocal(params PakmaParams, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "UnknownError"

	var response []byte
	var err error
	var valid bool = true

	switch params.Command {
	case "install":
		if !checkVersion(params.Version) {
			valid = false
			logger.Exceptions.Println("PakmaCommand invalid version")
			*result = "PakmaCommand invalid version"
			break
		}
		response, err = pakmaInstall(params.Version)
	case "upgrade":
		if !checkVersion(params.Version) {
			valid = false
			logger.Exceptions.Println("PakmaCommand invalid version")
			*result = "PakmaCommand invalid version"
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
		logger.Exceptions.Println("PakmaCommand unsupport command")
		*result = "PakmaCommand unsupport command"
	}
	if valid {
		// fmt.Println(string(response), err)
		if err != nil {
			logger.Exceptions.Println("Pakma request error")
			*result = "Pakma request error"
		} else {
			*result = string(response)
		}
	}
}
