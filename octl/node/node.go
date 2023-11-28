package node

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

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
	// output.PrintJSON(info.ToText())

	return &info, nil
}

func NodesInfoWithFilter(names []string, stateFilter string) (*protocols.NodesInfo, *errs.OctlError) {
	infos, err := NodesInfo(names)
	if err != nil {
		return nil, err
	}

	infos_filtered := &protocols.NodesInfo{
		BrainName:    infos.BrainAddr,
		BrainVersion: infos.BrainVersion,
		BrainAddr:    infos.BrainAddr,
		InfoList:     nil,
	}

	for _, info := range infos.InfoList {
		infoText := info.ToText()
		if !filterMatch(infoText.State, stateFilter) {
			continue
		}
		infos_filtered.InfoList = append(infos_filtered.InfoList, info)
	}

	output.PrintJSON(infos_filtered.ToText())
	return infos_filtered, nil
}

func filterMatch(value, target string) bool {
	if target == "" {
		return true
	}
	v := strings.ToLower(value)
	t := strings.ToLower(target)
	return strings.Contains(v, t)
}

func NodesPrune(names []string) *errs.OctlError {
	nodes, err := NodesParseNoCheck(names)
	if err != nil {
		emsg := "node parse." + err.Error()
		output.PrintFatalln(emsg)
		return errs.New(errs.OctlNodeParseError, emsg)
	}
	nodes_serialized, _ := config.Jsoner.Marshal(&nodes)
	url := fmt.Sprintf("https://%s/%s%s?targetNodes=%s",
		httpclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodePrune,
		nodes_serialized,
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

// func NodePrune() *errs.OctlError {
// 	url := fmt.Sprintf("https://%s/%s%s",
// 		httpclient.BrainAddr,
// 		config.GlobalConfig.Brain.ApiPrefix,
// 		config.API_NodePrune,
// 	)
// 	res, err := httpclient.BrainClient.Get(url)
// 	if err != nil {
// 		emsg := "http get error: " + err.Error()
// 		output.PrintFatalln(emsg)
// 		return errs.New(errs.OctlHttpRequestError, emsg)
// 	}
// 	res.Body.Close()
// 	return nil
// }

func NodesParse(names []string) ([]string, error) {
	result, err := nodesParse(names)
	if err != nil {
		return nil, err
	}
	if len(result.InvalidNames) != 0 {
		return nil, fmt.Errorf("node parse return invalid names: %v", result.InvalidNames)
	}
	if len(result.UnhealthyNodes) != 0 {
		return nil, fmt.Errorf("node parse return unhealthy nodes: %v", result.UnhealthyNodes)
	}
	return result.OutputNames, nil
}

func NodesParseNoCheck(names []string) ([]string, error) {
	result, err := nodesParse(names)
	if err != nil {
		return nil, err
	}
	return result.OutputNames, nil
}

func nodesParse(names []string) (protocols.NodeParseResult, error) {
	body, _ := json.Marshal(names)
	result := protocols.NodeParseResult{}

	url := fmt.Sprintf("https://%s/%s%s",
		httpclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodesParse,
	)
	res, err := httpclient.BrainClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		output.PrintFatalln("nodeparse post")
		return result, err
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode == http.StatusOK {
		err := json.Unmarshal(raw, &result)
		if err != nil {
			output.PrintFatalln(err)
			return result, err
		}
		return result, nil
	} else {
		return result, errors.New(string(raw))
	}
}
