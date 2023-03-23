package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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

	url := "http://127.0.0.1:12345/file/upload"
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

	url := "http://127.0.0.1:12345/file/spread"

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
