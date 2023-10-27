package file

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"octl/config"
	"octl/nameclient"
	"octl/node"
	"octl/output"
	"octl/task"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mholt/archiver/v3"
)

func DistribFile(localFileOrDir string, targetPath string, names []string) (string, error) {
	nodes, err := node.NodesParse(names)
	if err != nil {
		msg := "node parse."
		output.PrintFatalln(msg, err)
		return msg, err
	}

	if targetPath == "." {
		targetPath = ""
	} else if targetPath[len(targetPath)-1] != '/' {
		targetPath = targetPath + "/"
	}

	pwd, _ := os.Getwd()
	srcPath := pathFixing(localFileOrDir, pwd+string(filepath.Separator))

	// wrap the files first
	wrapName := fmt.Sprintf("%d.wrap", time.Now().Nanosecond())
	os.Mkdir(wrapName, os.ModePerm)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell.exe", "/C", fmt.Sprintf("cp -Force %s %s", srcPath, wrapName))
	} else {
		cmd = exec.Command("/bin/bash", "-c", fmt.Sprintf("cp -r %s %s", srcPath, wrapName))
	}
	err = cmd.Run()
	if err != nil {
		msg := "Wrap files: " + srcPath + "-->" + wrapName + " | " + cmd.String()
		output.PrintFatalln(msg, err)
		return msg, err
	}
	defer os.RemoveAll(wrapName)

	packName := fmt.Sprintf("%s.zip", wrapName)

	archiver.DefaultZip.OverwriteExisting = true
	err = archiver.DefaultZip.Archive([]string{wrapName}, packName)
	if err != nil {
		msg := fmt.Sprintf("archiver.DefaultZip.Archive([]string{%s}, %s).", wrapName, packName)
		output.PrintFatalln(msg, err)
		return msg, err
	}
	defer os.Remove(packName)

	f, err := os.OpenFile(packName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		msg := fmt.Sprintf("os.OpenFile(%s, os.O_RDONLY, os.ModePerm).", packName)
		output.PrintFatalln(msg, err)
		return msg, err
	}
	defer f.Close()

	nodes_serialized, _ := config.Jsoner.Marshal(&nodes)

	bodyBuffer := bytes.Buffer{}
	bodyWriter := multipart.NewWriter(&bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("packfiles", packName)
	io.Copy(fileWriter, f)
	bodyWriter.WriteField("targetPath", targetPath)
	bodyWriter.WriteField("targetNodes", string(nodes_serialized))

	contentType := bodyWriter.FormDataContentType()

	bodyWriter.Close()

	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_FileDistrib,
	)

	res, err := http.Post(url, contentType, &bodyBuffer)
	if err != nil {
		emsg := "http post error."
		output.PrintFatalln(emsg, err)
		return emsg, err
	}
	msg, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		emsg := "http read body."
		output.PrintFatalln(emsg, err)
		return emsg, err
	}

	if res.StatusCode != 202 {
		emsg := fmt.Sprintf("http request error msg=%s, status=%d. ", msg, res.StatusCode)
		output.PrintFatalln(emsg)
		return emsg, errors.New(emsg)
	}
	results, err := task.WaitTask("DISTRIBUTING...", string(msg))
	if err != nil {
		emsg := "Task processing error: " + err.Error()
		output.PrintFatalln(emsg, err)
		return emsg, err
	}
	output.PrintJSON(results)
	return string(results), nil
}
