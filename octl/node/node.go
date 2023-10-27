package node

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"octl/config"
	"octl/nameclient"
	"octl/output"
)

func NodeInfo(name string) (string, error) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodeInfo,
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

func NodesInfo(names []string) (string, error) {
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
		config.API_NodesInfo,
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

func NodePrune() (string, error) {
	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodePrune,
	)
	res, err := http.Get(url)
	if err != nil {
		emsg := "http get error."
		output.PrintFatalln(emsg, err)
		return emsg, err
	}
	res.Body.Close()
	return "OK", nil
}

func NodesParse(names []string) ([]string, error) {
	// parse nodes
	body, _ := json.Marshal(names)

	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodesParse,
	)
	res, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		output.PrintFatalln("nodeparse post")
		return nil, err
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode == 200 {
		nodes := []string{}
		err := json.Unmarshal(raw, &nodes)
		if err != nil {
			output.PrintFatalln(err)
			return nil, err
		}
		return nodes, nil
	} else {
		return nil, errors.New(string(raw))
	}
}
