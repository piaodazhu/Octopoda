package shell

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/httpclient"
	"github.com/piaodazhu/Octopoda/octl/node"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/task"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

func XRun(runtask string, params []string) ([]protocols.ExecutionResults, *errs.OctlError) {
	delay := 0
	names := []string{}
	for i := range params {
		if len(params[i]) < 3 {
			names = append(names, params[i])
			continue
		}
		switch params[i][:2] {
		case "-d":
			x, err := strconv.Atoi(params[i][2:])
			if err != nil {
				emsg := "invalid args: " + err.Error()
				return nil, errs.New(errs.OctlArgumentError, emsg)
			}
			delay = x
		default:
			names = append(names, params[i])
		}
	}
	return runTask(runtask, names, delay)
}

func Run(runtask string, names []string) ([]protocols.ExecutionResults, *errs.OctlError) {
	return runTask(runtask, names, -1)
}

func runTask(runtask string, names []string, delay int) ([]protocols.ExecutionResults, *errs.OctlError) {
	nodes, err := node.NodesParse(names)
	if err != nil {
		emsg := "node parse error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlNodeParseError, emsg)
	}

	isScript := true
	isBackground := true
	if runtask[0] == '{' {
		if len(runtask) > 2 && runtask[len(runtask)-1] == '}' {
			isScript = false
			isBackground = false
			runtask = runtask[1 : len(runtask)-1]
		} else {
			emsg := fmt.Sprintf("runtask=%s invalid.", runtask)
			output.PrintFatalln(emsg)
			return nil, errs.New(errs.OctlArgumentError, emsg)
		}
	} else if runtask[0] == '(' {
		if len(runtask) > 2 && runtask[len(runtask)-1] == ')' {
			isScript = false
			runtask = runtask[1 : len(runtask)-1]
		} else {
			emsg := fmt.Sprintf("runtask=%s invalid.", runtask)
			output.PrintFatalln(emsg)
			return nil, errs.New(errs.OctlArgumentError, emsg)
		}
	}
	if isScript {
		return runScript(runtask, nodes, delay)
	} else {
		return runCmd(runtask, nodes, isBackground, delay)
	}
}

func runScript(runtask string, names []string, delay int) ([]protocols.ExecutionResults, *errs.OctlError) {
	nodes, err := node.NodesParse(names)
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
		httpclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_RunScript,
	)
	nodes_serialized, _ := config.Jsoner.Marshal(&nodes)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("delayTime", fmt.Sprint(delay))
	writer.WriteField("targetNodes", string(nodes_serialized))
	fileWriter, _ := writer.CreateFormFile("script", fname)
	io.Copy(fileWriter, f)
	writer.Close()

	res, err := httpclient.BrainClient.Post(url, writer.FormDataContentType(), body)
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
		RunCancel(tid)
	}(string(taskid))

	results, err := task.WaitTask("", string(taskid))
	if err != nil {
		emsg := "Task processing error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlTaskWaitingError, emsg)
	}
	output.PrintJSON(results)
	return results, nil
}

func runCmd(runtask string, names []string, bg bool, delay int) ([]protocols.ExecutionResults, *errs.OctlError) {
	nodes, err := node.NodesParse(names)
	if err != nil {
		emsg := "node parse error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlNodeParseError, emsg)
	}

	url := fmt.Sprintf("https://%s/%s%s",
		httpclient.BrainAddr,
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
	writer.WriteField("targetNodes", string(nodes_serialized))
	writer.Close()

	res, err := httpclient.BrainClient.Post(url, writer.FormDataContentType(), body)
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
		RunCancel(tid)
	}(string(taskid))

	results, err := task.WaitTask("", string(taskid))
	if err != nil {
		emsg := "Task processing error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlTaskWaitingError, emsg)
	}
	output.PrintJSON(results)
	return results, nil
}

func RunCancel(taskid string) {
	URL := fmt.Sprintf("https://%s/%s%s",
		httpclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_RunCancel,
	)
	values := url.Values{}
	values.Add("taskid", taskid)
	res, err := httpclient.BrainClient.PostForm(URL, values)
	if err != nil {
		output.PrintFatalln("runCancel http post error.")
	}
	res.Body.Close()
}
