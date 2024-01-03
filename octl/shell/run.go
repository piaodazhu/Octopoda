package shell

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/httpclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/task"
	"github.com/piaodazhu/Octopoda/octl/workgroup"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

func RunCommand(cmd string, isBackgroud bool, shouldAlign bool, names []string) ([]protocols.ExecutionResult, *errs.OctlError) {
	return runCmd(cmd, names, isBackgroud, -1, shouldAlign)
}

func RunScript(cmd string, shouldAlign bool, names []string) ([]protocols.ExecutionResult, *errs.OctlError) {
	return runScript(cmd, names, -1, shouldAlign)
}

func XRunCommand(cmd string, isBackgroud bool, shouldAlign bool, delayExec int, names []string) ([]protocols.ExecutionResult, *errs.OctlError) {
	return runCmd(cmd, names, isBackgroud, delayExec, shouldAlign)
}

func XRunScript(cmd string, shouldAlign bool, delayExec int, names []string) ([]protocols.ExecutionResult, *errs.OctlError) {
	return runScript(cmd, names, delayExec, shouldAlign)
}

func runScript(runtask string, names []string, delay int, align bool) ([]protocols.ExecutionResult, *errs.OctlError) {
	nodes, err := workgroup.NodesParse(names)
	if err != nil {
		emsg := "node parse error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlNodeParseError, emsg)
	}

	f, err := os.OpenFile(runtask, os.O_RDONLY, os.ModePerm)
	if err != nil {
		emsg := "open script file error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlFileOperationError, emsg)
	}
	defer f.Close()
	fname := filepath.Base(runtask)

	url := fmt.Sprintf("https://%s/%s%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_RunScript,
	)
	nodes_serialized, _ := config.Jsoner.Marshal(&nodes)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("delayTime", fmt.Sprint(delay))
	writer.WriteField("needAlign", fmt.Sprint(align))
	writer.WriteField("targetNodes", string(nodes_serialized))
	fileWriter, _ := writer.CreateFormFile("script", fname)
	io.Copy(fileWriter, f)
	writer.Close()

	req, _ := http.NewRequest("POST", url, body)
	workgroup.SetHeader(req)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := httpclient.BrainClient.Do(req)
	if err != nil {
		emsg := "http post error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()

	taskid, _ := io.ReadAll(res.Body)
	go func(tid string) {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		signal.Stop(sigChan)
		RunCancel(tid)
	}(string(taskid))

	results, err := task.WaitTask("", string(taskid))
	if err != nil {
		emsg := "Task processing error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlTaskWaitingError, emsg)
	}
	output.PrintJSON(protocols.ExecutionResults(results))
	return results, nil
}

func runCmd(runtask string, names []string, bg bool, delay int, align bool) ([]protocols.ExecutionResult, *errs.OctlError) {
	nodes, err := workgroup.NodesParse(names)
	if err != nil {
		emsg := "node parse error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlNodeParseError, emsg)
	}

	url := fmt.Sprintf("https://%s/%s%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_RunCmd,
	)
	nodes_serialized, _ := config.Jsoner.Marshal(&nodes)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("command", runtask)
	if bg {
		writer.WriteField("background", "true")
	}
	writer.WriteField("delayTime", fmt.Sprint(delay))
	writer.WriteField("needAlign", fmt.Sprint(align))
	writer.WriteField("targetNodes", string(nodes_serialized))
	writer.Close()

	req, _ := http.NewRequest("POST", url, body)
	workgroup.SetHeader(req)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := httpclient.BrainClient.Do(req)
	if err != nil {
		emsg := "http post error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()

	taskid, _ := io.ReadAll(res.Body)
	go func(tid string) {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		signal.Stop(sigChan)
		RunCancel(tid)
	}(string(taskid))

	results, err := task.WaitTask("", string(taskid))
	if err != nil {
		emsg := "Task processing error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlTaskWaitingError, emsg)
	}
	output.PrintJSON(protocols.ExecutionResults(results))
	return results, nil
}

func RunCancel(taskid string) {
	URL := fmt.Sprintf("https://%s/%s%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_RunCancel,
	)
	values := url.Values{}
	values.Add("taskid", taskid)

	req, _ := http.NewRequest("POST", URL, bytes.NewBufferString(values.Encode()))
	workgroup.SetHeader(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := httpclient.BrainClient.Do(req)
	if err != nil {
		output.PrintFatalln("runCancel http post error.")
	}
	res.Body.Close()
}
