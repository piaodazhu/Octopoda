package node

import (
	"fmt"
	"io"
	"net/http"
	"octl/config"
	"octl/nameclient"
	"octl/output"
)

func NodeStatus(name string) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.GlobalConfig.Api.NodeStatus,
		name,
	)
	res, err := http.Get(url)
	if err != nil {
		output.PrintFatalln("Get")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	output.PrintJSON(raw)
}

func NodesStatus() {
	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.GlobalConfig.Api.NodesStatus,
	)
	res, err := http.Get(url)
	if err != nil {
		output.PrintFatalln("Get")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	output.PrintJSON(raw)
}
