package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"octl/config"
	"octl/output"
	"octl/task"
	"os"
	"time"

	"github.com/mholt/archiver"
)

func UpLoadFile(localFileOrDir string, targetPath string) {
	if targetPath == "." {
		targetPath = ""
	} else if targetPath[len(targetPath)-1] != '/' {
		targetPath = targetPath + "/"
	}

	if localFileOrDir[len(localFileOrDir)-1] == '/' {
		localFileOrDir = localFileOrDir[:len(localFileOrDir)-1]
	}

	packName := fmt.Sprintf("%d.zip", time.Now().Nanosecond())
	// err := exec.Command("tar", "-cf", packName, "-C", filepath.Dir(localFileOrDir), filepath.Base(localFileOrDir)).Run()

	// if err != nil {
	// 	panic("cmd.Run")
	// }
	archiver.DefaultZip.OverwriteExisting = true
	err := archiver.DefaultZip.Archive([]string{localFileOrDir}, packName)
	if err != nil {
		panic("Archive")
	}
	defer os.Remove(packName)

	f, err := os.OpenFile(packName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic("err")
	}
	defer f.Close()

	bodyBuffer := bytes.Buffer{}
	bodyWriter := multipart.NewWriter(&bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("tarfile", packName)
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

	msg, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic("ReadAll")
	}

	if res.StatusCode != 202 {
		fmt.Println("Request submit error: ", string(msg))
		return
	}
	results, err := task.WaitTask("UPLOADING...", string(msg))
	if err != nil {
		fmt.Println("Task processing error: ", err.Error())
		return
	}
	output.PrintJSON(results)
}

type FileSpreadParams struct {
	TargetPath  string
	FileOrDir   string
	TargetNodes []string
}

func SpreadFile(FileOrDir string, targetPath string, nodes []string) {
	fsParams := &FileSpreadParams{
		TargetPath:  targetPath,
		FileOrDir:   FileOrDir,
		TargetNodes: nodes,
	}
	buf, _ := json.Marshal(fsParams)

	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.FileSpread,
	)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		panic("Post")
	}
	msg, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic("ReadAll")
	}

	if res.StatusCode != 202 {
		fmt.Println("Request submit error: ", string(msg))
		return
	}
	results, err := task.WaitTask("SPREADING...", string(msg))
	if err != nil {
		fmt.Println("Task processing error: ", err.Error())
		return
	}
	output.PrintJSON(results)
}

type FileDistribParams struct {
	LocalFile   string
	TargetPath  string
	TargetNodes []string
}

func DistribFile(localFileOrDir string, targetPath string, nodes []string) {
	if targetPath == "." {
		targetPath = ""
	} else if targetPath[len(targetPath)-1] != '/' {
		targetPath = targetPath + "/"
	}

	if localFileOrDir[len(localFileOrDir)-1] == '/' {
		localFileOrDir = localFileOrDir[:len(localFileOrDir)-1]
	}

	packName := fmt.Sprintf("%d.zip", time.Now().Nanosecond())
	// err := exec.Command("tar", "-cf", tarName, "-C", filepath.Dir(localFileOrDir), filepath.Base(localFileOrDir)).Run()
	// if err != nil {
	// 	panic("cmd.Run")
	// }
	archiver.DefaultZip.OverwriteExisting = true
	err := archiver.DefaultZip.Archive([]string{localFileOrDir}, packName)
	if err != nil {
		panic("Archive")
	}
	defer os.Remove(packName)

	f, err := os.OpenFile(packName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic("err")
	}
	defer f.Close()

	nodes_serialized, _ := json.Marshal(&nodes)

	bodyBuffer := bytes.Buffer{}
	bodyWriter := multipart.NewWriter(&bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("packfiles", packName)
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
	msg, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic("ReadAll")
	}

	if res.StatusCode != 202 {
		fmt.Println("Request submit error: ", string(msg))
		return
	}
	results, err := task.WaitTask("DISTRIBUTING...", string(msg))
	if err != nil {
		fmt.Println("Task processing error: ", err.Error())
		return
	}
	output.PrintJSON(results)
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
	output.PrintJSON(msg)
}
