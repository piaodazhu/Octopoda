package file

import (
	"bytes"
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
	"github.com/piaodazhu/Octopoda/octl/httpclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/task"
	"github.com/piaodazhu/Octopoda/octl/workgroup"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/errs"

	"github.com/mholt/archiver/v3"
)

func Upload(localFileOrDir string, remoteTargetPath string, names []string, isForce bool) ([]protocols.ExecutionResult, *errs.OctlError) {
	nodes, err := workgroup.NodesParse(names)
	if err != nil {
		emsg := "node parse error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlNodeParseError, emsg)
	}

	if remoteTargetPath == "." {
		remoteTargetPath = ""
	} else if remoteTargetPath[len(remoteTargetPath)-1] != '/' {
		remoteTargetPath = remoteTargetPath + "/"
	}

	pwd, _ := os.Getwd()
	srcPath := pathFixing(localFileOrDir, pwd+string(filepath.Separator))

	// wrap the files first
	wrapName := fmt.Sprintf("%d.wrap", time.Now().Nanosecond())
	os.Mkdir(wrapName, os.ModePerm)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell.exe", "/C", fmt.Sprintf("cp -Force -r %s %s", srcPath, wrapName))
	} else {
		cmd = exec.Command("/bin/bash", "-c", fmt.Sprintf("cp -r %s %s", srcPath, wrapName))
	}
	err = cmd.Run()
	if err != nil {
		emsg := fmt.Sprintf("wrap files %s to %s error: %s", srcPath, wrapName, err.Error())
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlFileOperationError, emsg)
	}
	defer os.RemoveAll(wrapName)

	packName := fmt.Sprintf("%s.zip", wrapName)

	archiver.DefaultZip.OverwriteExisting = true
	err = archiver.DefaultZip.Archive([]string{wrapName}, packName)
	if err != nil {
		emsg := fmt.Sprintf("archiver.DefaultZip.Archive([]string{%s}, %s) error: %s", wrapName, packName, err.Error())
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlFileOperationError, emsg)
	}
	defer os.Remove(packName)

	f, err := os.OpenFile(packName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		emsg := fmt.Sprintf("os.OpenFile(%s, os.O_RDONLY, os.ModePerm) errors: %s", packName, err.Error())
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlFileOperationError, emsg)
	}
	defer f.Close()

	nodes_serialized, _ := config.Jsoner.Marshal(&nodes)

	bodyBuffer := bytes.Buffer{}
	bodyWriter := multipart.NewWriter(&bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("packfiles", packName)
	io.Copy(fileWriter, f)
	bodyWriter.WriteField("targetPath", remoteTargetPath)
	bodyWriter.WriteField("isForce", fmt.Sprint(isForce))
	bodyWriter.WriteField("targetNodes", string(nodes_serialized))

	contentType := bodyWriter.FormDataContentType()

	bodyWriter.Close()

	url := fmt.Sprintf("https://%s/%s%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_FileDistrib,
	)

	req, _ := http.NewRequest("POST", url, &bodyBuffer)
	workgroup.SetHeader(req)
	req.Header.Set("Content-Type", contentType)
	res, err := httpclient.BrainClient.Do(req)
	if err != nil {
		emsg := "http post error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpRequestError, emsg)
	}
	msg, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		emsg := "http read body: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpRequestError, emsg)
	}

	if res.StatusCode != http.StatusAccepted {
		emsg := fmt.Sprintf("http request error msg=%s, status=%d.", msg, res.StatusCode)
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpStatusError, emsg)
	}
	results, err := task.WaitTask("DISTRIBUTING...", string(msg))
	if err != nil {
		emsg := "Task processing error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlTaskWaitingError, emsg)
	}
	output.PrintJSON(results)
	return results, nil
}
