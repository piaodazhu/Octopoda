package shell

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"octl/config"
	"octl/nameclient"
	"octl/node"
	"octl/output"
	"os"
	"path/filepath"
	"strconv"
)

func XRun(task string, params []string) {
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
				return
			}
			delay = x
		default:
			names = append(names, params[i])
		}
	}
	runTask(task, names, delay)
}

func Run(task string, names []string) {
	runTask(task, names, -1)
}

func runTask(task string, names []string, delay int) {
	nodes, err := node.NodesParse(names)
	if err != nil {
		output.PrintFatalln(err)
	}

	isScript := true
	isBackground := true
	if task[0] == '{' {
		if len(task) > 2 && task[len(task)-1] == '}' {
			isScript = false
			isBackground = false
			task = task[1 : len(task)-1]
		} else {
			return
		}
	} else if task[0] == '(' {
		if len(task) > 2 && task[len(task)-1] == ')' {
			isScript = false
			task = task[1 : len(task)-1]
		} else {
			return
		}
	}
	if isScript {
		runScript(task, nodes, delay)
	} else {
		runCmd(task, nodes, isBackground, delay)
	}
}

func runScript(task string, names []string, delay int) {
	nodes, err := node.NodesParse(names)
	if err != nil {
		output.PrintFatalln(err)
	}

	f, err := os.OpenFile(task, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(task, " is not a script")
		return
	}
	defer f.Close()
	fname := filepath.Base(task)

	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.GlobalConfig.Api.RunScript,
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
		output.PrintFatalln("Post")
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	output.PrintJSON(raw)
}

func runCmd(task string, names []string, bg bool, delay int) {
	nodes, err := node.NodesParse(names)
	if err != nil {
		output.PrintFatalln(err)
	}

	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.GlobalConfig.Api.RunCmd,
	)
	nodes_serialized, _ := config.Jsoner.Marshal(&nodes)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("command", task)
	if bg {
		writer.WriteField("background", "true")
	}
	writer.WriteField("delayTime", fmt.Sprint(delay))
	writer.WriteField("targetNodes", string(nodes_serialized))
	writer.Close()

	res, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		output.PrintFatalln("Post")
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	output.PrintJSON(raw)
}
