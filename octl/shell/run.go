package shell

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/node"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/task"
	"github.com/piaodazhu/Octopoda/protocols"
)

func XRun(runtask string, params []string) ([]protocols.ExecutionResults, error) {
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
				return nil, errors.New(emsg)
			}
			delay = x
		default:
			names = append(names, params[i])
		}
	}
	return runTask(runtask, names, delay)
}

func Run(runtask string, names []string) ([]protocols.ExecutionResults, error) {
	return runTask(runtask, names, -1)
}

func runTask(runtask string, names []string, delay int) ([]protocols.ExecutionResults, error) {
	nodes, err := node.NodesParse(names)
	if err != nil {
		emsg := "node parse error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
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
			return nil, errors.New(emsg)
		}
	} else if runtask[0] == '(' {
		if len(runtask) > 2 && runtask[len(runtask)-1] == ')' {
			isScript = false
			runtask = runtask[1 : len(runtask)-1]
		} else {
			emsg := fmt.Sprintf("runtask=%s invalid.", runtask)
			output.PrintFatalln(emsg)
			return nil, errors.New(emsg)
		}
	}
	if isScript {
		return runScript(runtask, nodes, delay)
	} else {
		return runCmd(runtask, nodes, isBackground, delay)
	}
}

func runScript(runtask string, names []string, delay int) ([]protocols.ExecutionResults, error) {
	nodes, err := node.NodesParse(names)
	if err != nil {
		emsg := "node parse error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}

	f, err := os.OpenFile(runtask, os.O_RDONLY, os.ModePerm)
	if err != nil {
		emsg := "open script file error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}
	defer f.Close()
	fname := filepath.Base(runtask)

	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
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

	res, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		emsg := "http post error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
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
		return nil, errors.New(emsg)
	}
	output.PrintJSON(results)
	return results, nil
}

func runCmd(runtask string, names []string, bg bool, delay int) ([]protocols.ExecutionResults, error) {
	nodes, err := node.NodesParse(names)
	if err != nil {
		emsg := "node parse error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}

	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
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

	res, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		emsg := "http post error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
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
		return nil, errors.New(emsg)
	}
	output.PrintJSON(results)
	return results, nil
}

func RunCancel(taskid string) {
	URL := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_RunCancel,
	)
	values := url.Values{}
	values.Add("taskid", taskid)
	res, err := http.PostForm(URL, values)
	if err != nil {
		output.PrintFatalln("runCancel http post error.")
	}
	res.Body.Close()
}
