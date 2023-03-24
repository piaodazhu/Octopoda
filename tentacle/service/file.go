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

type FileParams struct {
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
	rmsg := RMSG{"OK"}

	fileinfo := FileParams{}
	err := json.Unmarshal(raw, &fileinfo)
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

func FileTree(conn net.Conn, raw []byte) {
	var pathsb strings.Builder
	pathsb.WriteString(config.GlobalConfig.Workspace.Root)
	pathsb.Write(raw)
	res := allFiles(pathsb.String())
	err := message.SendMessage(conn, message.TypeFileTreeResponse, res)
	if err != nil {
		logger.Server.Println("FileTree send error")
	}
}

type FileInfo struct {
	Name       string
	Size       int64
	ModifyTime string
}

func allFiles(path string) []byte {
	if path[len(path) - 1] == '/' {
		path = path[:len(path) - 1]
	}
	finfos := []FileInfo{}
	walkDir(path, &finfos)
	serialized, _ := json.Marshal(&finfos)
	return serialized
}

func walkDir(path string, files *[]FileInfo) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return
	}

	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() {
			walkDir(path+PthSep+fi.Name(), files)
		} else {
			detail, _ := fi.Info()
			modtimestr := detail.ModTime().Format("01月02日 15:04")
			finfo := FileInfo{
				Name:       path + PthSep + fi.Name(),
				Size:       detail.Size(),
				ModifyTime: modtimestr,
			}

			*files = append(*files, finfo)
		}
	}
}
