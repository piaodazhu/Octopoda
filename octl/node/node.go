package node

import (
	"fmt"
	"io"
	"net/http"
	"octl/config"

	"github.com/hokaccha/go-prettyjson"
)

func NodeInfo(name string) {
	url := fmt.Sprintf("http://%s:%d/%s%s?name=%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.NodeInfo,
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

func NodesInfo() {
	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.NodesInfo,
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

func NodePrune() {
	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.NodePrune,
	)
	res, err := http.Get(url)
	if err != nil {
		panic("Get")
	}
	res.Body.Close()
}
