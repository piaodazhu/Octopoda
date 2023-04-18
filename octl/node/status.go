package node

import (
	"fmt"
	"io"
	"net/http"
	"octl/config"

	"github.com/hokaccha/go-prettyjson"
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
		panic("Get")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	s, _ := prettyjson.Format(raw)
	fmt.Println(string(s))
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
		panic("Get")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	s, _ := prettyjson.Format(raw)
	fmt.Println(string(s))
}
