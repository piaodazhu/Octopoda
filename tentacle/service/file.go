package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/message"
	"time"
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
	rmsg := RMSG{"OK"}

	fileinfo := FileParams{}
	err := json.Unmarshal(raw, &fileinfo)
	if err != nil {
		logger.Server.Println("FilePush")
		rmsg.Msg = "FilePush"
		goto errorout
	}

	file.WriteString(config.GlobalConfig.Workspace.Store)
	file.WriteString(fileinfo.TargetPath)

	os.Mkdir(file.String(), os.ModePerm)
	if file.String()[file.Len()-1] != '/' {
		file.WriteByte('/')
	}

	file.WriteString(fileinfo.FileName)

	err = saveFile(fileinfo.FileBuf, file.String())
	if err != nil {
		rmsg.Msg = "FileNotSave"
		goto errorout
	}
errorout:
	payload, _ := json.Marshal(&rmsg)
	err = message.SendMessage(conn, message.TypeFilePushResponse, payload)
	if err != nil {
		logger.Server.Println("FilePush send error")
	}
}

func FileTree(conn net.Conn, raw []byte) {
	var pathsb strings.Builder
	pathsb.WriteString(config.GlobalConfig.Workspace.Store)
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
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	finfos := []FileInfo{}
	walkDir(path, &finfos)
	hideRoot(path, &finfos)
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

func hideRoot(root string, files *[]FileInfo) {
	for i := range *files {
		(*files)[i].Name = (*files)[i].Name[len(root):]
	}
}

func saveFile(filebufb64, filename string) error {
	content, err := base64.RawStdEncoding.DecodeString(filebufb64)
	if err != nil {
		logger.Server.Println("FileDecode")
		return err
	}
	err = os.WriteFile(filename, content, os.ModePerm)
	if err != nil {
		logger.Server.Println("WriteFile")
		return err
	}
	return nil
}

type ErrUnpack struct{}
func (ErrUnpack) Error() string { return "ErrUnpack" }

func unpackFiles(packb64 string, subdir string) error {
	var file strings.Builder
	// tmp tar file name
	fname := fmt.Sprintf("%d.tar", time.Now().Nanosecond())

	// make complete path and filename
	file.WriteString(config.GlobalConfig.Workspace.Root)
	file.WriteString(subdir)
	path := file.String()
	os.Mkdir(path, os.ModePerm)
	if path[len(path)-1] != '/' {
		file.WriteByte('/')
	}
	file.WriteString(fname)
	
	err := saveFile(packb64, file.String())
	if err != nil {
		return err
	}

	cmd := exec.Command("tar", "-xf", file.String(), "-C", path)
	err = cmd.Run()
	if err != nil {
		return err
	}
	// err = cmd.Wait()
	// if err != nil {
	// 	return err
	// }
	if !cmd.ProcessState.Success() {
		return ErrUnpack{}
	}

	os.Remove(file.String())
	return nil
}
