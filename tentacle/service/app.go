package service

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"tentacle/app"
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

func AppVersion(conn net.Conn, raw []byte) {
	err := message.SendMessage(conn, message.TypeAppVersionResponse, app.Versions(string(raw)))
	if err != nil {
		logger.Server.Println("AppVersion send error")
	}
}

func AppsInfo(conn net.Conn, raw []byte) {
	err := message.SendMessage(conn, message.TypeAppsInfoResponse, app.Digest())
	if err != nil {
		logger.Server.Println("AppsInfo send error")
	}
}

type AppResetParams struct {
	Name string
	Hash string
	Mode string
}

func AppReset(conn net.Conn, raw []byte) {
	arParams := &AppResetParams{}
	rmsg := RMSG{"OK"}
	var payload []byte
	var longhash string
	var ok bool

	err := json.Unmarshal(raw, arParams)
	if err != nil {
		logger.Client.Println(err)
		rmsg.Msg = "Invalid Params"
		goto errorout
	}
	if longhash, ok = app.ConvertHash(arParams.Name, arParams.Hash); !ok {
		rmsg.Msg = "No app or version"
		goto errorout
	}

	if err = app.GitReset(arParams.Name, longhash, arParams.Mode); err != nil {
		rmsg.Msg = "Failed to GitRevert: " + err.Error()
		goto errorout
	}
	if ok = app.Reset(arParams.Name, arParams.Hash); !ok {
		rmsg.Msg = "Failed to Revert: " + err.Error()
		goto errorout
	}
	app.Save()
errorout:
	payload, _ = json.Marshal(&rmsg)
	err = message.SendMessage(conn, message.TypeAppResetResponse, payload)
	if err != nil {
		logger.Server.Println("AppReset send error")
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
	rmsg := RMSG{"OK"}
	var output, payload []byte
	var ok bool
	var dirName string
	var version app.Version

	err := json.Unmarshal(raw, acParams)
	if err != nil {
		logger.Client.Println(err)
		rmsg.Msg = "Invalid Params"
		goto errorout
	}
	dirName = acParams.Name + "@" + acParams.Scenario

	// is exists?
	if ok = app.Exists(acParams.Name, acParams.Scenario); !ok {
		if !app.Create(acParams.Name, acParams.Scenario, acParams.Description) {
			logger.Client.Println("app.Create")
			rmsg.Msg = "Failed Create App"
			goto errorout
		}
		if !app.GitCreate(dirName) {
			logger.Client.Println("app.GitCreate")
			rmsg.Msg = "Failed Init the Repo"
			goto errorout
		}
	}
	// unpack files
	err = unpackFiles(acParams.FilePack, dirName)
	if err != nil {
		logger.Client.Println("unpack Files")
		rmsg.Msg = err.Error()
		goto errorout
	} else {
		rmsg.Msg = string(output)
	}
	// commit
	version, err = app.GitCommit(dirName, acParams.Message)
	if err != nil {
		logger.Client.Println("app.GitCommit")
		rmsg.Msg = err.Error()
		goto errorout
	}

	// update nodeApps
	ok = app.Update(dirName, version)
	if !ok {
		logger.Client.Println("app.Update")
		rmsg.Msg = "Faild to update app version"
	}
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
	rmsg := RMSG{"OK"}
	var output, payload []byte
	var ok bool
	var sparams ScriptParams
	var dirName string
	var version app.Version

	err := json.Unmarshal(raw, adParams)
	if err != nil {
		logger.Client.Println(err)
		rmsg.Msg = "Invalid Params"
		goto errorout
	}
	dirName = adParams.Name + "@" + adParams.Scenario

	// is new?
	if ok = app.Exists(adParams.Name, adParams.Scenario); !ok {
		if !app.Create(adParams.Name, adParams.Scenario, adParams.Description) {
			logger.Client.Println("app.Create")
			rmsg.Msg = "Failed Create App"
			goto errorout
		}
		if !app.GitCreate(dirName) {
			logger.Client.Println("app.GitCreate")
			rmsg.Msg = "Failed Init the Repo"
			goto errorout
		}
	}
	// run script
	sparams = ScriptParams{
		FileName:   fmt.Sprintf("script_%s.sh", time.Now().Format("2006_01_02_15_04")),
		TargetPath: dirName,
		FileBuf:    adParams.Script,
	}
	output, err = execScript(&sparams)
	if err != nil {
		logger.Client.Println("execScript")
		rmsg.Msg = err.Error()
		goto errorout
	} else {
		rmsg.Msg = string(output)
	}
	// commit
	version, err = app.GitCommit(dirName, adParams.Message)
	if err != nil {
		logger.Client.Println("app.GitCommit")
		rmsg.Msg = err.Error()
		goto errorout
	}

	// update nodeApps
	ok = app.Update(dirName, version)
	if !ok {
		logger.Client.Println("app.Update")
		rmsg.Msg = "Faild to update app version"
	}
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
	var dirName string

	err := json.Unmarshal(raw, adParams)
	if err != nil {
		logger.Client.Println(err)
		rmsg.Msg = "Invalid Params"
		goto errorout
	}
	dirName = adParams.Name + "@" + adParams.Scenario

	// is new?
	if ok = app.Exists(adParams.Name, adParams.Scenario); ok {
		if !app.Delete(adParams.Name, adParams.Scenario) {
			logger.Client.Println("app.Delete")
			rmsg.Msg = "Failed Delete App"
			goto errorout
		}
	}
	os.Remove(dirName)
	app.Save()
errorout:
	payload, _ = json.Marshal(&rmsg)
	err = message.SendMessage(conn, message.TypeAppDeployResponse, payload)
	if err != nil {
		logger.Server.Println("AppDelete send error")
	}
}
