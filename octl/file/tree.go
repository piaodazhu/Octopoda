package file

import (
	"fmt"
	"io"
	"net/http"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/httpclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/workgroup"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

func ListAllFile(subdir string, node string) (string, *errs.OctlError) {
	url := fmt.Sprintf("https://%s/%s%s?name=%s&subdir=%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_FileTree,
		node,
		subdir,
	)
	req, _ := http.NewRequest("GET", url, nil)
	workgroup.SetHeader(req)
	res, err := httpclient.BrainClient.Do(req)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return emsg, errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	msg, err := io.ReadAll(res.Body)
	if err != nil {
		emsg := "http read body: " + err.Error()
		output.PrintFatalln(emsg)
		return emsg, errs.New(errs.OctlHttpRequestError, emsg)
	}
	output.PrintJSON(msg)
	return string(msg), nil
}
