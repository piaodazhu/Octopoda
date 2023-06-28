package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/message"
	"tentacle/snp"
	"time"

	"github.com/mholt/archiver/v3"
)

type FileParams struct {
	PackName   string
	PathType   string
	TargetPath string
	FileBuf    string
}

func pathFixing(path string, base string) string {
	// dstPath: the unpacked files will be moved under this path
	var result strings.Builder
	// find ~
	homePos := -1
	for i, c := range path {
		if c == '~' {
			homePos = i
		}
	}
	if homePos != -1 {
		result.WriteString(path[homePos:])
	} else {
		if len(path) == 0 || path[0] != '/' {
			result.WriteString(base)
		}
		result.WriteString(path)
	}
	return result.String()
}

func FilePush(conn net.Conn, raw []byte) {
	rmsg := message.Result{
		Rmsg: "OK",
	}

	var file string
	fileinfo := FileParams{}
	err := config.Jsoner.Unmarshal(raw, &fileinfo)
	if err != nil {
		logger.Exceptions.Println("FilePush")
		rmsg.Rmsg = "FilePush"
		goto errorout
	}

	file = pathFixing(fileinfo.TargetPath, config.GlobalConfig.Workspace.Store)
	err = unpackFiles(fileinfo.FileBuf, fileinfo.PackName, file)
	if err != nil {
		rmsg.Rmsg = "unpack Files"
		goto errorout
	}

errorout:
	payload, _ := config.Jsoner.Marshal(&rmsg)
	err = message.SendMessageUnique(conn, message.TypeFilePushResponse, snp.GenSerial(), payload)
	if err != nil {
		logger.Comm.Println("FilePush send error")
	}
}

func FilePull(conn net.Conn, raw []byte) {
	var pathsb strings.Builder
	var err error
	var packName, wrapName, srcPath, pwd string
	var cmd *exec.Cmd

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

	pwd, _ = os.Getwd()
	srcPath = pathFixing(pathsb.String(), pwd+string(filepath.Separator))
	// _, err = os.Stat(srcPath)
	// if err != nil {
	// 	logger.Exceptions.Println(srcPath, " not exist")
	// 	goto errorout
	// }

	// wrap the files first
	wrapName = fmt.Sprintf("%d.wrap", time.Now().Nanosecond())
	os.MkdirAll(wrapName, os.ModePerm)

	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell.exe", "/C", fmt.Sprintf("cp -Force %s %s", srcPath, wrapName))
	} else {
		cmd = exec.Command("/bin/bash", "-c", fmt.Sprintf("cp -r %s %s", srcPath, wrapName))
	}
	err = cmd.Run()
	if err != nil {
		logger.Exceptions.Println("Wrap files: " + srcPath + "-->" + wrapName + " | " + cmd.String())
		goto errorout
	}
	defer os.RemoveAll(wrapName)

	packName = fmt.Sprintf("%s.zip", wrapName)

	archiver.DefaultZip.OverwriteExisting = true
	err = archiver.DefaultZip.Archive([]string{wrapName}, packName)
	if err != nil {
		logger.Exceptions.Println("packFile")
	}
	defer os.Remove(packName)

	fileinfo.FileBuf = loadFile(packName)
	fileinfo.PackName = packName

	payload, _ = json.Marshal(&fileinfo)
errorout:
	err = message.SendMessageUnique(conn, message.TypeFilePullResponse, snp.GenSerial(), payload)
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
	err = message.SendMessageUnique(conn, message.TypeFileTreeResponse, snp.GenSerial(), res)
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

func unpackFiles(packb64 string, packName string, targetDir string) error {
	pname := filepath.Base(packName)
	extPos := -1
	for i := len(pname) - 1; i > 0; i-- {
		if pname[i] == '.' {
			extPos = i
			break
		}
	}
	if extPos == -1 {
		return fmt.Errorf("extPos == -1")
	}
	wname := pname[:extPos]
	tmpDir := config.GlobalConfig.Workspace.Store + ".octopoda_tmp/"
	os.Mkdir(tmpDir, os.ModePerm)
	defer os.RemoveAll(tmpDir)

	ppath := tmpDir + pname
	wpath := tmpDir + wname

	err := saveFile(packb64, ppath)
	if err != nil {
		return err
	}

	archiver.DefaultZip.OverwriteExisting = true
	err = archiver.DefaultZip.Unarchive(ppath, tmpDir)
	if err != nil {
		logger.Exceptions.Print(err)
		return ErrUnpack{}
	}

	// os.Mkdir(targetDir, os.ModePerm)
	// fmt.Println(";;;;", fmt.Sprintf("cp -r %s/* %s", wpath, targetDir))
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("mkdir -p %s && cp -r %s/* %s", targetDir, wpath, targetDir))
	err = cmd.Run()
	if err != nil {
		logger.Exceptions.Print(err)
		return fmt.Errorf("cp -r")
	}

	return nil
}

func unpackFilesNoWrap(packb64 string, dir string) error {
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

	archiver.DefaultZip.OverwriteExisting = true
	err = archiver.DefaultZip.Unarchive(file.String(), path)
	if err != nil {
		return ErrUnpack{}
	}

	os.Remove(file.String())
	return nil
}

// func packFile(fileOrDir string) string {
// 	packName := fmt.Sprintf("%d.zip", time.Now().Nanosecond())
// 	archiver.DefaultZip.OverwriteExisting = true
// 	err := archiver.DefaultZip.Archive([]string{fileOrDir}, packName)
// 	if err != nil {
// 		return ""
// 	}
// 	return packName
// }

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
	ChunkSize := 4096 * 3
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
