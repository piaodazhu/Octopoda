package file

import (
	"fmt"
	"io"
	"net/http"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/httpclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/task"
	"github.com/piaodazhu/Octopoda/octl/workgroup"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

func Download(remoteFileOrDir string, localTargetPath string, node string) (*protocols.ExecutionResult, *errs.OctlError) {
	if len(node) == 0 || node[0] == '@' {
		emsg := "command pull not support node group"
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlArgumentError, emsg)
	}
	url := fmt.Sprintf("https://%s/%s%s?name=%s&fileOrDir=%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_FilePull,
		node,
		remoteFileOrDir,
	)
	req, _ := http.NewRequest("GET", url, nil)
	workgroup.SetHeader(req)
	res, err := httpclient.BrainClient.Do(req)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	msg, err := io.ReadAll(res.Body)
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

	// have to wait
	results, err := task.WaitTask("PULLING...", string(msg))
	if err != nil {
		emsg := "Task processing error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlTaskWaitingError, emsg)
	}

	if len(results) != 1 {
		emsg := fmt.Sprintf("number of results of this command should be only 1 but get %d", len(results))
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlMessageParseError, emsg)
	}

	if len(results[0].Result) == 0 {
		emsg := "file is empty or not found."
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlMessageParseError, emsg)
	}

	// marshal the file info
	finfo := protocols.FilePullParams{}
	err = config.Jsoner.Unmarshal([]byte(results[0].Result), &finfo)
	if err != nil {
		emsg := "config.Jsoner.Unmarshal([]byte(fileResults[0].Result), &finfo) error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlMessageParseError, emsg)
	}
	// unpack result.Output
	err = unpackFiles(finfo.FileBuf, finfo.PackName, localTargetPath)
	if err != nil {
		emsg := "unpackFiles(finfo.FileBuf, finfo.PackName, targetdir) error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlFileOperationError, emsg)
	}
	output.PrintInfoln("Success")

	results[0].Result = ""
	return &results[0], nil
}
