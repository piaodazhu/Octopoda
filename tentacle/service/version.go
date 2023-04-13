package service

import (
	"encoding/json"
	"net"
	"tentacle/app"
	"tentacle/logger"
	"tentacle/message"
)

func AppLatestVersion(conn net.Conn, raw []byte) {
	aParams := AppBasic{}
	err := json.Unmarshal(raw, &aParams)
	var payload []byte  
	if err != nil {
		logger.Client.Println(err)
		payload = []byte{}
	} else {
		v := app.CurVersion(aParams.Name, aParams.Scenario)
		payload, _ = json.Marshal(&v)
	}
	
	err = message.SendMessage(conn, message.TypeAppVersionResponse, payload)
	if err != nil {
		logger.Server.Println("AppLatestVersion send error")
	}
}

func AppVersions(conn net.Conn, raw []byte) {
	aParams := AppBasic{}
	err := json.Unmarshal(raw, &aParams)
	var payload []byte  
	if err != nil {
		logger.Client.Println(err)
		payload = []byte{}
	} else {
		payload = app.Versions(aParams.Name, aParams.Scenario)
	}
	err = message.SendMessage(conn, message.TypeAppVersionResponse, payload)
	if err != nil {
		logger.Server.Println("AppVersion send error")
	}
}

type AppResetParams struct {
	AppBasic
	VersionHash string
	Mode        string
}

func AppReset(conn net.Conn, raw []byte) {
	arParams := &AppResetParams{}
	rmsg := RDATA{}
	rmsg.Modified = false
	var payload []byte
	var longhash, fullname string
	var ok bool

	err := json.Unmarshal(raw, arParams)
	if err != nil {
		logger.Client.Println(err)
		rmsg.Msg = "Invalid Params"
		goto errorout
	}
	if longhash, ok = app.ConvertHash(arParams.Name, arParams.Scenario, arParams.VersionHash); !ok {
		rmsg.Msg = "No app or version"
		goto errorout
	}
	fullname = arParams.Name + "@" + arParams.Scenario
	if err = app.GitReset(fullname, longhash, arParams.Mode); err != nil {
		rmsg.Msg = "Failed to GitRevert: " + err.Error()
		goto errorout
	}
	if ok = app.Reset(arParams.Name, arParams.Scenario, arParams.VersionHash, arParams.Message); !ok {
		rmsg.Msg = "Failed to Revert: " + err.Error()
		goto errorout
	}
	rmsg.Msg = "OK"
	rmsg.Version = longhash
	rmsg.Modified = true
	app.Save()
errorout:
	payload, _ = json.Marshal(&rmsg)
	logger.Client.Println(rmsg)
	err = message.SendMessage(conn, message.TypeAppResetResponse, payload)
	if err != nil {
		logger.Server.Println("AppReset send error")
	}
}
