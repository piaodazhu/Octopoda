package file

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"octl/config"
	"octl/nameclient"
	"octl/node"
	"octl/output"
	"octl/task"
	"protocols"
)

func SpreadFile(FileOrDir string, targetPath string, names []string) {
	nodes, err := node.NodesParse(names)
	if err != nil {
		output.PrintFatalln(err)
	}

	fsParams := &protocols.FileSpreadParams{
		TargetPath:  targetPath,
		FileOrDir:   FileOrDir,
		TargetNodes: nodes,
	}
	buf, _ := config.Jsoner.Marshal(fsParams)

	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_FileSpread,
	)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		output.PrintFatalln("Post")
	}
	msg, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		output.PrintFatalln("ReadAll")
	}

	if res.StatusCode != 202 {
		output.PrintFatalln("Request submit error: " + string(msg))
		return
	}
	results, err := task.WaitTask("SPREADING...", string(msg))
	if err != nil {
		output.PrintFatalln("Task processing error: " + err.Error())
		return
	}
	output.PrintJSON(results)
}
