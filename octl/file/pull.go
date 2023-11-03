package file

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/task"
	"github.com/piaodazhu/Octopoda/protocols"
)

func PullFile(pathtype string, node string, fileOrDir string, targetdir string) (*protocols.ExecutionResults, error) {
	if len(node) == 0 || node[0] == '@' {
		emsg := "command pull not support node group"
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}
	url := fmt.Sprintf("http://%s/%s%s?pathtype=%s&name=%s&fileOrDir=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_FilePull,
		pathtype,
		node,
		fileOrDir,
	)
	res, err := http.Get(url)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}
	defer res.Body.Close()
	msg, err := io.ReadAll(res.Body)
	if err != nil {
		emsg := "http read body: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}

	if res.StatusCode != http.StatusAccepted {
		emsg := fmt.Sprintf("http request error msg=%s, status=%d.", msg, res.StatusCode)
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}

	// have to wait
	results, err := task.WaitTask("PULLING...", string(msg))
	if err != nil {
		emsg := "Task processing error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}

	if len(results) != 1 || len(results[0].Result) == 0 {
		emsg := fmt.Sprintf("number of results of this command should be only 1 but get %d", len(results))
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	} 
	
	// marshal the file info
	finfo := protocols.FilePullParams{}
	err = config.Jsoner.Unmarshal([]byte(results[0].Result), &finfo)
	if err != nil {
		emsg := "config.Jsoner.Unmarshal([]byte(fileResults[0].Result), &finfo) error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}
	// unpack result.Output
	err = unpackFiles(finfo.FileBuf, finfo.PackName, targetdir)
	if err != nil {
		emsg := "unpackFiles(finfo.FileBuf, finfo.PackName, targetdir) error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}
	output.PrintInfoln("Success")

	results[0].Result = ""
	return &results[0], nil
}
