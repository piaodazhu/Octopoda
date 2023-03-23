package service

import (
	"encoding/base64"
	"encoding/json"
	"net"
	"os"
	"strings"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/message"
)

type FileInfo struct {
	FileName   string
	TargetPath string
	FileBuf    string
}

type RMSG struct {
	Msg string
}

func FilePush(conn net.Conn, raw []byte) {
	var file strings.Builder
	var content []byte
	var f *os.File

	fileinfo := FileInfo{}
	err := json.Unmarshal(raw, &fileinfo)
	rmsg := RMSG{"OK"}
	if err != nil {
		logger.Server.Println("FilePush")
		rmsg.Msg = "FilePush"
		goto errorout
	}

	content, err = base64.RawStdEncoding.DecodeString(fileinfo.FileBuf)
	if err != nil {
		logger.Server.Println("FileDecode")
		rmsg.Msg = "FileDecode"
		goto errorout
	}

	
	file.WriteString(config.GlobalConfig.Workspace.Root)
	file.WriteString(fileinfo.TargetPath)
	
	os.Mkdir(file.String(), os.ModePerm)

	file.WriteString(fileinfo.FileName)
	f, err = os.Create(file.String())
	if err != nil {
		logger.Server.Println("FileCreate")
		rmsg.Msg = "FileCreate"
		goto errorout
	}
	defer f.Close()
	f.Write(content)
errorout:
	payload, _ := json.Marshal(&rmsg)
	err = message.SendMessage(conn, message.TypeFilePushResponse, payload)
	if err != nil {
		logger.Server.Println("FilePush send error")
	}
}
