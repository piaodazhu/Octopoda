package file

import (
	"fmt"
	"io"
	"net/http"
	"octl/config"
	"octl/nameclient"
	"octl/output"
)

func ListAllFile(pathtype string, node string, subdir string) {
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
		output.PrintFatalln("Get")
	}
	defer res.Body.Close()
	msg, err := io.ReadAll(res.Body)
	if err != nil {
		output.PrintFatalln("ReadAll")
	}
	output.PrintJSON(msg)
}
