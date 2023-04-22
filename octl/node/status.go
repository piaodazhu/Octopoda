package node

import (
	"fmt"
	"io"
	"net/http"
	"octl/config"
	"octl/output"
)

func NodeStatus(name string) {
	url := fmt.Sprintf("http://%s:%d/%s%s?name=%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.NodeStatus,
		name,
	)
	res, err := http.Get(url)
	if err != nil {
		output.PrintFatal("Get")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	output.PrintJSON(raw)
}

func NodesStatus() {
	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.NodesStatus,
	)
	res, err := http.Get(url)
	if err != nil {
		output.PrintFatal("Get")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	output.PrintJSON(raw)
}
