package file

import (
	"bytes"
	"errors"
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

func SpreadFile(FileOrDir string, targetPath string, names []string) (string, error) {
	nodes, err := node.NodesParse(names)
	if err != nil {
		msg := "node parse."
		output.PrintFatalln(msg, err)
		return msg, err
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
	results, err := task.WaitTask("SPREADING...", string(msg))
	if err != nil {
		emsg := "Task processing error: " + err.Error()
		output.PrintFatalln(emsg, err)
		return emsg, err
	}
	output.PrintJSON(results)
	return string(results), nil
}
