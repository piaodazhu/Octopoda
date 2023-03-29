package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"octl/config"
	"os"
	"path/filepath"
)

func UpLoadFile(localFile string, targetPath string) {
	f, err := os.OpenFile(localFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic("err")
	}
	defer f.Close()
	fname := filepath.Base(localFile)

	bodyBuffer := bytes.Buffer{}
	bodyWriter := multipart.NewWriter(&bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("file", fname)
	io.Copy(fileWriter, f)

	bodyWriter.WriteField("targetPath", targetPath)

	contentType := bodyWriter.FormDataContentType()

	bodyWriter.Close()

	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.FileUpload,
	)
	res, err := http.Post(url, contentType, &bodyBuffer)
	if err != nil {
		panic("post")
	}
	defer res.Body.Close()
	msg, err := io.ReadAll(res.Body)
	if err != nil {
		panic("ReadAll")
	}
	fmt.Println(msg)
}

type FileSpreadParams struct {
	SourcePath  string
	TargetPath  string
	FileName    string
	TargetNodes []string
}

func SpreadFile(fileName string, sourcePath string, targetPath string, nodes []string) {
	fsParams := &FileSpreadParams{
		SourcePath:  sourcePath,
		TargetPath:  targetPath,
		FileName:    fileName,
		TargetNodes: nodes,
	}
	buf, _ := json.Marshal(fsParams)

	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.FileUpload,
	)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		panic("Post")
	}
	defer res.Body.Close()
	msg, err := io.ReadAll(res.Body)
	if err != nil {
		panic("ReadAll")
	}
	fmt.Println(msg)
}

type FileDistribParams struct {
	LocalFile  string
	TargetPath  string
	TargetNodes []string
}

func DistribFile(localFile string, targetPath string, nodes []string) {
	if targetPath[len(targetPath) - 1] != '/' {
		targetPath = targetPath + "/"
	}

	f, err := os.OpenFile(localFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic("err")
	}
	defer f.Close()
	fname := filepath.Base(localFile)
	nodes_serialized, _ := json.Marshal(&nodes)

	bodyBuffer := bytes.Buffer{}
	bodyWriter := multipart.NewWriter(&bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("file", fname)
	io.Copy(fileWriter, f)
	bodyWriter.WriteField("targetPath", targetPath)
	bodyWriter.WriteField("targetNodes", string(nodes_serialized))

	contentType := bodyWriter.FormDataContentType()

	bodyWriter.Close()

	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.FileDistrib,
	)

	res, err := http.Post(url, contentType, &bodyBuffer)
	if err != nil {
		panic("post")
	}
	defer res.Body.Close()
	msg, err := io.ReadAll(res.Body)
	if err != nil {
		panic("ReadAll")
	}
	fmt.Println(string(msg))
}

func ListAllFile(node string, subdir string) {
	url := fmt.Sprintf("http://%s:%d/%s%s?name=%s&subdir=%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.FileTree,
		node,
		subdir,
	)
	res, err := http.Get(url)
	if err != nil {
		panic("Get")
	}
	defer res.Body.Close()
	msg, err := io.ReadAll(res.Body)
	if err != nil {
		panic("ReadAll")
	}
	fmt.Println(string(msg))
}