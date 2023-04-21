package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/message"
	"time"

	"github.com/mholt/archiver/v3"
)

type FileParams struct {
	PackName   string
	TargetPath string
	FileBuf    string
}

func FilePush(conn net.Conn, raw []byte) {
	var file strings.Builder
	rmsg := message.Result{
		Rmsg: "OK",
	}

	fileinfo := FileParams{}
	err := json.Unmarshal(raw, &fileinfo)
	if err != nil {
		logger.Server.Println("FilePush")
		rmsg.Rmsg = "FilePush"
		goto errorout
	}
	// logger.Server.Println(fileinfo)
	file.WriteString(config.GlobalConfig.Workspace.Store)
	file.WriteString(fileinfo.TargetPath)

	err = unpackFiles(fileinfo.FileBuf, file.String())
	if err != nil {
		rmsg.Rmsg = "unpack Files"
		goto errorout
	}

	// os.Mkdir(file.String(), os.ModePerm)
	// if file.String()[file.Len()-1] != '/' {
	// 	file.WriteByte('/')
	// }

	// file.WriteString(fileinfo.TarName)

	// err = saveFile(fileinfo.FileBuf, file.String())
	// if err != nil {
	// 	rmsg.Rmsg = "FileNotSave"
	// 	goto errorout
	// }
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
	Offset := 0
	Len := len(filebufb64)
	ChunkSize := 4096 * 4
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		logger.Server.Println("OpenFile:", filename)
		return err
	}
	defer f.Close()

	for Offset < Len {
		end := Offset + ChunkSize
		if Offset+ChunkSize > Len {
			end = Len
		}
		content, err := base64.RawStdEncoding.DecodeString(filebufb64[Offset:end])
		if err != nil {
			logger.Server.Println("FileDecode:", err.Error())
			return err
		}
		_, err = f.Write(content)
		if err != nil {
			logger.Server.Println("Write:", err.Error())
			return err
		}
		Offset += ChunkSize
	}
	// content, err := base64.RawStdEncoding.DecodeString(filebufb64)
	// if err != nil {
	// 	logger.Server.Println("FileDecode")
	// 	return err
	// }
	// err = os.WriteFile(filename, content, os.ModePerm)
	// if err != nil {
	// 	logger.Server.Println("WriteFile:", filename)
	// 	return err
	// }
	return nil
}

type ErrUnpack struct{}

func (ErrUnpack) Error() string { return "ErrUnpack" }

func unpackFiles(packb64 string, dir string) error {
	var file strings.Builder
	// tmp tar file name
	fname := fmt.Sprintf("%d.zip", time.Now().Nanosecond())

	// make complete path and filename
	// file.WriteString(config.GlobalConfig.Workspace.Root)
	file.WriteString(dir)
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

	// cmd := exec.Command("tar", "-xf", file.String(), "-C", path)
	// err = cmd.Run()
	// if err != nil {
	// 	return err
	// }
	// // err = cmd.Wait()
	// // if err != nil {
	// // 	return err
	// // }
	// if !cmd.ProcessState.Success() {
	// 	return ErrUnpack{}
	// }

	err = archiver.Unarchive(file.String(), path)
	if err != nil {
		return ErrUnpack{}
	}

	os.Remove(file.String())
	return nil
}
