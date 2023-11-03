package file

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/task"
	"github.com/piaodazhu/Octopoda/protocols"

	"github.com/mholt/archiver/v3"
)

func UpLoadFile(localFileOrDir string, targetPath string) (*protocols.ExecutionResults, error) {
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
	err := cmd.Run()
	if err != nil {
		emsg := fmt.Sprintf("warp files %s to %s error: %s", srcPath, wrapName, err.Error())
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}
	defer os.RemoveAll(wrapName)

	packName := fmt.Sprintf("%s.zip", wrapName)
	archiver.DefaultZip.OverwriteExisting = true
	err = archiver.DefaultZip.Archive([]string{wrapName}, packName)
	if err != nil {
		emsg := fmt.Sprintf("archiver.DefaultZip.Archive([]string{%s}, %s) error: %s", wrapName, packName, err.Error())
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}
	defer os.Remove(packName)

	f, err := os.OpenFile(packName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		emsg := fmt.Sprintf("os.OpenFile(%s, os.O_RDONLY, os.ModePerm) errors: %s", packName, err.Error())
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}
	defer f.Close()

	bodyBuffer := bytes.Buffer{}
	bodyWriter := multipart.NewWriter(&bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("tarfile", packName)
	io.Copy(fileWriter, f)

	bodyWriter.WriteField("targetPath", targetPath)

	contentType := bodyWriter.FormDataContentType()

	bodyWriter.Close()

	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_FileUpload,
	)
	res, err := http.Post(url, contentType, &bodyBuffer)
	if err != nil {
		emsg := "http post error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}

	msg, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		emsg := "http read body: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}

	if res.StatusCode != http.StatusAccepted {
		emsg := fmt.Sprintf("http request error msg=%s, status=%d. ", msg, res.StatusCode)
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}
	results, err := task.WaitTask("UPLOADING...", string(msg))
	if err != nil {
		emsg := "Task processing error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}
	
	if len(results) != 1 {
		emsg := fmt.Sprintf("number of results of this command should be only 1 but get %d", len(results))
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}
	output.PrintJSON(results[0])
	return &results[0], nil
}
