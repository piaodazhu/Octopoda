package service

import (
	"encoding/base64"
	"fmt"
	"io"
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
	PathType   string
	TargetPath string
	FileBuf    string
}

func FilePush(conn net.Conn, raw []byte) {
	var file strings.Builder
	rmsg := message.Result{
		Rmsg: "OK",
	}

	fileinfo := FileParams{}
	err := config.Jsoner.Unmarshal(raw, &fileinfo)
	if err != nil {
		logger.Exceptions.Println("FilePush")
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
	payload, _ := config.Jsoner.Marshal(&rmsg)
	err = message.SendMessage(conn, message.TypeFilePushResponse, payload)
	if err != nil {
		logger.Comm.Println("FilePush send error")
	}
}

func FilePull(conn net.Conn, raw []byte) {
	var pathsb strings.Builder
	var err error
	var packName string
	payload := []byte{}
	fileinfo := FileParams{}
	err = config.Jsoner.Unmarshal(raw, &fileinfo)
	if err != nil {
		logger.Exceptions.Println("FilePull")
		goto errorout
	}
	switch fileinfo.PathType {
	case "store":
		pathsb.WriteString(config.GlobalConfig.Workspace.Store)
	case "log":
		pathsb.WriteString(config.GlobalConfig.Logger.Path)
	case "nodeapp":
		pathsb.WriteString(config.GlobalConfig.Workspace.Root)
	default:
		goto errorout
	}
	pathsb.WriteString(fileinfo.TargetPath)
	_, err = os.Stat(pathsb.String())
	if err != nil {
		goto errorout
	}
	// pack the file or dir
	packName = packFile(pathsb.String())
	if packName == "" {
		logger.Exceptions.Println("packFile")
		goto errorout
	}
	defer os.Remove(packName)
	payload = []byte(loadFile(packName))
errorout:
	err = message.SendMessage(conn, message.TypeFilePullResponse, payload)
	if err != nil {
		logger.Comm.Println("FilePull send error")
	}
}

func FileTree(conn net.Conn, raw []byte) {
	var pathsb strings.Builder
	res := []byte{}
	pathinfo := FileParams{}
	err := config.Jsoner.Unmarshal(raw, &pathinfo)
	if err != nil {
		goto errorout
	}

	switch pathinfo.PathType {
	case "store":
		pathsb.WriteString(config.GlobalConfig.Workspace.Store)
	case "log":
		pathsb.WriteString(config.GlobalConfig.Logger.Path)
	case "nodeapp":
		pathsb.WriteString(config.GlobalConfig.Workspace.Root)
	default:
		goto errorout
	}
	pathsb.WriteString(pathinfo.TargetPath)

	res = allFiles(pathsb.String())
errorout:
	err = message.SendMessage(conn, message.TypeFileTreeResponse, res)
	if err != nil {
		logger.Comm.Println("FileTree send error")
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
	serialized, _ := config.Jsoner.Marshal(&finfos)
	return serialized
}

func walkDir(path string, files *[]FileInfo) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return
	}

	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.Name()[0] == '.' {
			continue
		} else if fi.IsDir() {
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
		logger.Exceptions.Println("OpenFile:", filename)
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
			logger.Exceptions.Println("FileDecode:", err.Error())
			return err
		}
		_, err = f.Write(content)
		if err != nil {
			logger.Exceptions.Println("Write:", err.Error())
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

	logger.SysInfo.Print("want save: ", file.String(), path)

	archiver.DefaultZip.OverwriteExisting = true
	err = archiver.DefaultZip.Unarchive(file.String(), path)
	if err != nil {
		logger.Exceptions.Print(err)
		return ErrUnpack{}
	}

	os.Remove(file.String())
	return nil
}

func packFile(fileOrDir string) string {
	packName := fmt.Sprintf("%d.zip", time.Now().Nanosecond())
	archiver.DefaultZip.OverwriteExisting = true
	err := archiver.DefaultZip.Archive([]string{fileOrDir}, packName)
	if err != nil {
		return ""
	}
	return packName
}

func loadFile(packName string) string {
	var filebufb64 strings.Builder
	f, err := os.Open(packName)
	if err != nil {
		return ""
	}
	defer f.Close()

	// prepare enough buffer capacity
	info, _ := f.Stat()
	filebufb64.Grow(base64.RawStdEncoding.EncodedLen(int(info.Size())))

	// read and encode to base64
	ChunkSize := 4096 * 4
	ChunkBuf := make([]byte, ChunkSize)
	for {
		n, err := f.Read(ChunkBuf)
		if err == io.EOF {
			break
		}
		filebufb64.WriteString(base64.RawStdEncoding.EncodeToString(ChunkBuf[:n]))
	}
	return filebufb64.String()
}
