package file

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/httpclient"
	"github.com/piaodazhu/Octopoda/octl/node"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/task"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

func SpreadFile(FileOrDir string, targetPath string, names []string) ([]protocols.ExecutionResults, *errs.OctlError) {
	nodes, err := node.NodesParse(names)
	if err != nil {
		emsg := "node parse error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlNodeParseError, emsg)
	}

	fsParams := &protocols.FileSpreadParams{
		TargetPath:  targetPath,
		FileOrDir:   FileOrDir,
		TargetNodes: nodes,
	}
	buf, _ := config.Jsoner.Marshal(fsParams)

	url := fmt.Sprintf("https://%s/%s%s",
		httpclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_FileSpread,
	)

	res, err := httpclient.BrainClient.Post(url, "application/json", bytes.NewBuffer(buf))
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
		emsg := fmt.Sprintf("http request error msg=%s, status=%d. ", msg, res.StatusCode)
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpStatusError, emsg)
	}
	results, err := task.WaitTask("SPREADING...", string(msg))
	if err != nil {
		emsg := "Task processing error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlTaskWaitingError, emsg)
	}
	output.PrintJSON(results)
	return results, nil
}
