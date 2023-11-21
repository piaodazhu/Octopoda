package file

import (
	"fmt"
	"io"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/httpclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

func ListAllFile(pathtype string, node string, subdir string) (string, *errs.OctlError) {
	url := fmt.Sprintf("https://%s/%s%s?pathtype=%s&name=%s&subdir=%s",
		httpclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_FileTree,
		pathtype,
		node,
		subdir,
	)
	res, err := httpclient.BrainClient.Get(url)
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
