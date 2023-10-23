package file

import (
	"fmt"
	"io"
	"net/http"
	"octl/config"
	"octl/nameclient"
	"octl/output"
	"octl/task"
	"protocols"
)

func PullFile(pathtype string, node string, fileOrDir string, targetdir string) {
	if len(node) == 0 || node[0] == '@' {
		output.PrintFatalln("command pull not support node group")
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
		output.PrintFatalln("Get")
	}
	defer res.Body.Close()
	msg, err := io.ReadAll(res.Body)
	if err != nil {
		output.PrintFatalln("ReadAll")
	}

	result := protocols.Result{}
	finfo := protocols.FilePullParams{}
	if res.StatusCode == 200 {
		// get the file info structure from brain
		err = config.Jsoner.Unmarshal(msg, &result)
		if err != nil {
			output.PrintFatalln(err.Error())
		}

		// marshal the file info
		err = config.Jsoner.Unmarshal([]byte(result.Output), &finfo)
		if err != nil {
			output.PrintFatalln(err.Error())
		}

		// unpack result.Output
		err = unpackFiles(finfo.FileBuf, finfo.PackName, targetdir)
		if err != nil {
			output.PrintFatalln(err.Error())
		}
		output.PrintInfoln("Success")
	} else if res.StatusCode == 202 {
		// have to wait
		resultmsg, err := task.WaitTask("PULLING...", string(msg))
		if err != nil {
			output.PrintFatalln("Task processing error: " + err.Error())
			return
		}
		config.Jsoner.Unmarshal(resultmsg, &result)
		if len(result.Output) == 0 {
			output.PrintJSON(resultmsg)
		} else {
			// marshal the file info
			err = config.Jsoner.Unmarshal([]byte(result.Output), &finfo)
			if err != nil {
				output.PrintFatalln(err.Error())
			}
			// unpack result.Output
			err = unpackFiles(finfo.FileBuf, finfo.PackName, targetdir)
			if err != nil {
				output.PrintFatalln(err.Error())
			}
			output.PrintInfoln("Success")
		}
	} else {
		// some error
		output.PrintJSON(msg)
	}
}
