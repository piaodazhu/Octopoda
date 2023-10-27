package file

import (
	"fmt"
	"io"
	"net/http"
	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/output"
)

func ListAllFile(pathtype string, node string, subdir string) (string, error) {
	url := fmt.Sprintf("http://%s/%s%s?pathtype=%s&name=%s&subdir=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_FileTree,
		pathtype,
		node,
		subdir,
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
	output.PrintJSON(msg)
	return string(msg), nil
}
