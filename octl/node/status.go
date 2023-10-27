package node

import (
	"fmt"
	"io"
	"net/http"
	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/output"
)

func NodeStatus(name string) (string, error) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodeStatus,
		name,
	)
	res, err := http.Get(url)
	if err != nil {
		emsg := "http get error."
		output.PrintFatalln(emsg, err)
		return emsg, err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	output.PrintJSON(raw)
	return string(raw), nil
}

func NodesStatus(names []string) (string, error) {
	nodes, err := NodesParse(names)
	if err != nil {
		msg := "node parse."
		output.PrintFatalln(msg, err)
		return msg, err
	}
	nodes_serialized, _ := config.Jsoner.Marshal(&nodes)
	url := fmt.Sprintf("http://%s/%s%s?targetNodes=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodesStatus,
		string(nodes_serialized),
	)
	res, err := http.Get(url)
	if err != nil {
		emsg := "http get error."
		output.PrintFatalln(emsg, err)
		return emsg, err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	output.PrintJSON(raw)
	return string(raw), nil
}
