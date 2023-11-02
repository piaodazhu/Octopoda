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

func PullFile(pathtype string, node string, fileOrDir string, targetdir string) (string, error) {
	if len(node) == 0 || node[0] == '@' {
		emsg := "command pull not support node group"
		output.PrintFatalln(emsg)
		return emsg, errors.New(emsg)
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
		emsg := "http get error."
		output.PrintFatalln(emsg, err)
		return emsg, err
	}
	defer res.Body.Close()
	msg, err := io.ReadAll(res.Body)
	if err != nil {
		emsg := "http read body."
		output.PrintFatalln(emsg, err)
		return emsg, err
	}

	fileResults := []protocols.ExecutionResults{}
	finfo := protocols.FilePullParams{}

	if res.StatusCode == 202 {
		// have to wait
		resultmsg, err := task.WaitTask("PULLING...", string(msg))
		if err != nil {
			emsg := "Task processing error: " + err.Error()
			output.PrintFatalln(emsg, err)
			return emsg, err
		}
		err = config.Jsoner.Unmarshal(resultmsg, &fileResults)
		if err != nil {
			emsg := "result unmarshal error: " + err.Error()
			output.PrintFatalln(emsg, err)
			return emsg, err
		}
		if len(fileResults) != 1 || len(fileResults[0].Result) == 0 {
			output.PrintJSON(resultmsg)
		} else {
			// marshal the file info
			err = config.Jsoner.Unmarshal([]byte(fileResults[0].Result), &finfo)
			if err != nil {
				emsg := "config.Jsoner.Unmarshal([]byte(fileResults[0].Result), &finfo)."
				output.PrintFatalln(emsg, err)
				return emsg, err
			}
			// unpack result.Output
			err = unpackFiles(finfo.FileBuf, finfo.PackName, targetdir)
			if err != nil {
				emsg := "unpackFiles(finfo.FileBuf, finfo.PackName, targetdir)."
				output.PrintFatalln(emsg, err)
				return emsg, err
			}
			output.PrintInfoln("Success")
		}
	} else {
		// some error
		output.PrintJSON(msg)
	}
	
	return "OK", nil
}
