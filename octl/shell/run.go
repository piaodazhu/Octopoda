package shell

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

	"github.com/hokaccha/go-prettyjson"
)

func RunTask(task string, nodes []string) {
	isScript := true
	if task[0] == '{' {
		if len(task) > 2 && task[len(task)-1] == '}' {
			isScript = false
			task = task[1 : len(task)-1]
		} else {
			return
		}
	}
	if isScript {
		runScript(task, nodes)
	} else {
		runCmd(task, nodes)
	}
}

func runScript(task string, nodes []string) {
	f, err := os.OpenFile(task, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(task, " is not a script")
		return
	}
	defer f.Close()
	fname := filepath.Base(task)

	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.RunScript,
	)
	nodes_serialized, _ := json.Marshal(&nodes)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("targetNodes", string(nodes_serialized))
	fileWriter, _ := writer.CreateFormFile("script", fname)
	io.Copy(fileWriter, f)
	writer.Close()

	res, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		panic("Post")
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	s, _ := prettyjson.Format(raw)
	fmt.Println(string(s))
}

func runCmd(task string, nodes []string) {
	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.RunCmd,
	)
	nodes_serialized, _ := json.Marshal(&nodes)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("command", task)
	writer.WriteField("targetNodes", string(nodes_serialized))
	writer.Close()

	res, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		panic("Post")
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	s, _ := prettyjson.Format(raw)
	fmt.Println(string(s))
}
