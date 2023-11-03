package node

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/protocols"
)

func NodeInfo(name string) (*protocols.NodeInfo, error) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodeInfo,
		name,
	)
	res, err := http.Get(url)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		emsg := fmt.Sprintf("[%d]msg=%s\n", res.StatusCode, string(raw))
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}

	info := protocols.NodeInfo{}
	err = json.Unmarshal(raw, &info)
	if err != nil {
		emsg := "res unmarshal error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}
	output.PrintJSON(info.ToText())
	return &info, nil
}

func NodesInfo(names []string) (*protocols.NodesInfo, error) {
	nodes, err := NodesParse(names)
	if err != nil {
		emsg := "node parse." + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
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
		emsg := "http get error." + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		emsg := fmt.Sprintf("[%d]msg=%s\n", res.StatusCode, string(raw))
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}

	info := protocols.NodesInfo{}
	err = json.Unmarshal(raw, &info)
	if err != nil {
		emsg := "res unmarshal error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errors.New(emsg)
	}
	output.PrintJSON(info.ToText())

	return &info, nil
}

func NodePrune() error {
	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodePrune,
	)
	res, err := http.Get(url)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return errors.New(emsg)
	}
	res.Body.Close()
	return nil
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
	if res.StatusCode == http.StatusOK {
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
