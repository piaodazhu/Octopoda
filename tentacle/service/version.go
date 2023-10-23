package service

import (
	"net"
	"protocols"
	"tentacle/app"
	"tentacle/config"
	"tentacle/logger"
)

func AppLatestVersion(conn net.Conn, serialNum uint32, raw []byte) {
	aParams := AppBasic{}
	err := config.Jsoner.Unmarshal(raw, &aParams)
	var payload []byte
	if err != nil {
		logger.Exceptions.Println(err)
		payload = []byte{}
	} else {
		v := app.CurVersion(aParams.Name, aParams.Scenario)
		payload, _ = config.Jsoner.Marshal(&v)
	}

	err = protocols.SendMessageUnique(conn, protocols.TypeAppLatestVersionResponse, serialNum, payload)
	if err != nil {
		logger.Comm.Println("AppLatestVersion send error")
	}
}

func AppVersions(conn net.Conn, serialNum uint32, raw []byte) {
	aParams := AppBasic{}
	err := config.Jsoner.Unmarshal(raw, &aParams)
	var payload []byte
	if err != nil {
		logger.Exceptions.Println(err)
		payload = []byte{}
	} else {
		payload = app.Versions(aParams.Name, aParams.Scenario)
	}
	err = protocols.SendMessageUnique(conn, protocols.TypeAppVersionResponse, serialNum, payload)
	if err != nil {
		logger.Comm.Println("AppVersion send error")
	}
}

type AppResetParams struct {
	AppBasic
	VersionHash string
	Mode        string
}

func AppReset(conn net.Conn, serialNum uint32, raw []byte) {
	arParams := &AppResetParams{}
	rmsg := protocols.Result{
		Rmsg: "OK",
	}
	rmsg.Modified = false
	var payload []byte
	var longhash, fullname string
	var ok bool

	err := config.Jsoner.Unmarshal(raw, arParams)
	if err != nil {
		logger.Exceptions.Println(err)
		rmsg.Rmsg = "Invalid Params"
		goto errorout
	}
	if longhash, ok = app.ConvertHash(arParams.Name, arParams.Scenario, arParams.VersionHash); !ok {
		rmsg.Rmsg = "No app or version"
		goto errorout
	}
	fullname = arParams.Name + "@" + arParams.Scenario
	if err = app.GitReset(fullname, longhash, arParams.Mode); err != nil {
		rmsg.Rmsg = "Failed to GitRevert: " + err.Error()
		goto errorout
	}
	if ok = app.Reset(arParams.Name, arParams.Scenario, arParams.VersionHash, arParams.Message); !ok {
		rmsg.Rmsg = "Failed to Revert: " + err.Error()
		goto errorout
	}
	rmsg.Rmsg = "OK"
	rmsg.Version = longhash
	rmsg.Modified = true
	app.Save()
errorout:
	payload, _ = config.Jsoner.Marshal(&rmsg)
	err = protocols.SendMessageUnique(conn, protocols.TypeAppResetResponse, serialNum, payload)
	if err != nil {
		logger.Comm.Println("AppReset send error")
	}
}
