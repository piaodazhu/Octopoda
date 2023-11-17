package node

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/httpclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

func NodeInfo(name string) (*protocols.NodeInfo, *errs.OctlError) {
	url := fmt.Sprintf("https://%s/%s%s?name=%s",
		httpclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodeInfo,
		name,
	)
	res, err := httpclient.BrainClient.Get(url)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		emsg := fmt.Sprintf("[%d]msg=%s\n", res.StatusCode, string(raw))
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpStatusError, emsg)
	}

	info := protocols.NodeInfo{}
	err = json.Unmarshal(raw, &info)
	if err != nil {
		emsg := "res unmarshal error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlMessageParseError, emsg)
	}
	output.PrintJSON(info.ToText())
	return &info, nil
}

func NodesInfo(names []string) (*protocols.NodesInfo, *errs.OctlError) {
	nodes, err := NodesParse(names)
	if err != nil {
		emsg := "node parse." + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlNodeParseError, emsg)
	}
	nodes_serialized, _ := config.Jsoner.Marshal(&nodes)
	url := fmt.Sprintf("https://%s/%s%s?targetNodes=%s",
		httpclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodesInfo,
		string(nodes_serialized),
	)
	res, err := httpclient.BrainClient.Get(url)
	if err != nil {
		emsg := "http get error." + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		emsg := fmt.Sprintf("[%d]msg=%s\n", res.StatusCode, string(raw))
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpStatusError, emsg)
	}

	info := protocols.NodesInfo{}
	err = json.Unmarshal(raw, &info)
	if err != nil {
		emsg := "res unmarshal error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlMessageParseError, emsg)
	}
	output.PrintJSON(info.ToText())

	return &info, nil
}

func NodePrune() *errs.OctlError {
	url := fmt.Sprintf("https://%s/%s%s",
		httpclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodePrune,
	)
	res, err := httpclient.BrainClient.Get(url)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return errs.New(errs.OctlHttpRequestError, emsg)
	}
	res.Body.Close()
	return nil
}

func NodesParse(names []string) ([]string, error) {
	// parse nodes
	body, _ := json.Marshal(names)

	url := fmt.Sprintf("https://%s/%s%s",
		httpclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodesParse,
	)
	res, err := httpclient.BrainClient.Post(url, "application/json", bytes.NewReader(body))
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
