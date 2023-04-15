package service

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"tentacle/app"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/message"
	"time"
)

// show apps
// show versions [name]
// reset version [name] [version]
//
// exists app [name]
// create app [name] [scenario] [description]
// commit app [name] [message]

type RDATA struct {
	Msg      string
	Version  string
	Modified bool
}

func AppsInfo(conn net.Conn, raw []byte) {
	err := message.SendMessage(conn, message.TypeAppsInfoResponse, app.Digest())
	if err != nil {
		logger.Server.Println("AppsInfo send error")
	}
}

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

func AppCreate(conn net.Conn, raw []byte) {
	acParams := &AppCreateParams{}
	rmsg := RDATA{}
	rmsg.Modified = false
	var output, payload []byte
	var ok bool
	var fullname string
	var version app.Version
	// os.WriteFile("./dump", raw, os.ModePerm)
	// logger.Client.Println(len(raw))

	err := json.Unmarshal(raw, acParams)
	if err != nil {
		logger.Client.Println(err)
		rmsg.Msg = "Invalid Params"
		goto errorout
	}
	fullname = acParams.Name + "@" + acParams.Scenario

	// is exists?
	if ok = app.Exists(acParams.Name, acParams.Scenario); !ok {
		if !app.Create(acParams.Name, acParams.Scenario, acParams.Description) {
			logger.Client.Println("app.Create")
			rmsg.Msg = "Failed Create App"
			goto errorout
		}
		if !app.GitCreate(fullname) {
			logger.Client.Println("app.GitCreate")
			rmsg.Msg = "Failed Init the Repo"
			goto errorout
		}
	}
	// unpack files
	err = unpackFiles(acParams.FilePack, fullname)
	if err != nil {
		logger.Client.Println("unpack Files")
		rmsg.Msg = err.Error()
		goto errorout
	} else {
		rmsg.Msg = string(output)
	}
	// commit
	version, err = app.GitCommit(fullname, acParams.Message)
	if err != nil {
		logger.Client.Println("app.GitCommit")
		if _, ok := err.(app.EmptyCommitError); ok {
			rmsg.Msg = "OK: No Change"
			rmsg.Version = app.CurVersion(acParams.Name, acParams.Scenario).Hash
		} else {
			rmsg.Msg = err.Error()
		}
		goto errorout
	}

	// update nodeApps
	ok = app.Update(acParams.Name, acParams.Scenario, version)
	if !ok {
		logger.Client.Println("app.Update")
		rmsg.Msg = "Faild to update app version"
	}
	rmsg.Msg = "OK"
	rmsg.Version = version.Hash
	rmsg.Modified = true
	app.Save()
errorout:
	payload, _ = json.Marshal(&rmsg)
	err = message.SendMessage(conn, message.TypeAppCreateResponse, payload)
	if err != nil {
		logger.Server.Println("AppCreate send error")
	}
}

func AppDeploy(conn net.Conn, raw []byte) {
	adParams := &AppDeployParams{}
	rmsg := RDATA{}
	rmsg.Modified = false
	var output, payload []byte
	var ok bool
	var sparams ScriptParams
	var fullname string
	var version app.Version

	err := json.Unmarshal(raw, adParams)
	if err != nil {
		logger.Client.Println(err)
		rmsg.Msg = "Invalid Params"
		goto errorout
	}
	fullname = adParams.Name + "@" + adParams.Scenario

	// is new?
	if ok = app.Exists(adParams.Name, adParams.Scenario); !ok {
		if !app.Create(adParams.Name, adParams.Scenario, adParams.Description) {
			logger.Client.Println("app.Create")
			rmsg.Msg = "Failed Create App"
			goto errorout
		}
		if !app.GitCreate(fullname) {
			logger.Client.Println("app.GitCreate")
			rmsg.Msg = "Failed Init the Repo"
			goto errorout
		}
	}
	// run script
	sparams = ScriptParams{
		FileName:   fmt.Sprintf("script_%s.sh", time.Now().Format("2006_01_02_15_04")),
		TargetPath: "scripts/",
		FileBuf:    adParams.Script,
	}
	output, err = execScript(&sparams, config.GlobalConfig.Workspace.Root + fullname)
	if err != nil {
		logger.Client.Println("execScript")
		rmsg.Msg = err.Error()
		goto errorout
	} else {
		rmsg.Msg = string(output)
	}
	// commit
	version, err = app.GitCommit(fullname, adParams.Message)
	if err != nil {
		if _, ok := err.(app.EmptyCommitError); ok {
			rmsg.Msg = "OK: No Change"
			rmsg.Version = app.CurVersion(adParams.Name, adParams.Scenario).Hash
		} else {
			rmsg.Msg = err.Error()
			logger.Client.Println("app.GitCommit")
		}
		goto errorout
	}

	// update nodeApps
	ok = app.Update(adParams.Name, adParams.Scenario, version)
	if !ok {
		logger.Client.Println("app.Update")
		rmsg.Msg = "Faild to update app version"
	}
	rmsg.Msg = "OK"
	rmsg.Version = version.Hash
	rmsg.Modified = true
	app.Save()
errorout:
	payload, _ = json.Marshal(&rmsg)
	err = message.SendMessage(conn, message.TypeAppDeployResponse, payload)
	if err != nil {
		logger.Server.Println("AppDeploy send error")
	}
}

type AppDeleteParams struct {
	Name     string
	Scenario string
}

func AppDelete(conn net.Conn, raw []byte) {
	adParams := &AppDeleteParams{}
	rmsg := RMSG{"OK"}
	var payload []byte
	var ok bool
	var subdirName string

	err := json.Unmarshal(raw, adParams)
	if err != nil {
		logger.Client.Println(err)
		rmsg.Msg = "Invalid Params"
		goto errorout
	}
	subdirName = adParams.Name + "@" + adParams.Scenario

	// is new?
	if ok = app.Exists(adParams.Name, adParams.Scenario); ok {
		if !app.Delete(adParams.Name, adParams.Scenario) {
			logger.Client.Println("app.Delete")
			rmsg.Msg = "Failed Delete App"
			goto errorout
		}
	}
	os.RemoveAll(config.GlobalConfig.Workspace.Root + subdirName)
	app.Save()
errorout:
	payload, _ = json.Marshal(&rmsg)
	err = message.SendMessage(conn, message.TypeAppDeleteResponse, payload)
	if err != nil {
		logger.Server.Println("AppDelete send error")
	}
}
