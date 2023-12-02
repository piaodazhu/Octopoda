package node

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/httpclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/workgroup"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

func NodeInfo(name string) (*protocols.NodeInfo, *errs.OctlError) {
	url := fmt.Sprintf("https://%s/%s%s?name=%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodeInfo,
		name,
	)
	req, _ := http.NewRequest("GET", url, nil)
	workgroup.SetHeader(req)
	res, err := httpclient.BrainClient.Do(req)
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
	nodes, err := workgroup.NodesParse(names)
	if err != nil {
		emsg := "node parse." + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlNodeParseError, emsg)
	}
	nodes_serialized, _ := config.Jsoner.Marshal(&nodes)
	url := fmt.Sprintf("https://%s/%s%s?targetNodes=%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodesInfo,
		string(nodes_serialized),
	)
	req, _ := http.NewRequest("GET", url, nil)
	workgroup.SetHeader(req)
	res, err := httpclient.BrainClient.Do(req)
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
	nodes, err := workgroup.NodesParseNoCheck(names)
	if err != nil {
		emsg := "node parse." + err.Error()
		output.PrintFatalln(emsg)
		return errs.New(errs.OctlNodeParseError, emsg)
	}
	nodes_serialized, _ := config.Jsoner.Marshal(&nodes)
	url := fmt.Sprintf("https://%s/%s%s?targetNodes=%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodePrune,
		nodes_serialized,
	)
	req, _ := http.NewRequest("GET", url, nil)
	workgroup.SetHeader(req)
	res, err := httpclient.BrainClient.Do(req)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return errs.New(errs.OctlHttpRequestError, emsg)
	}
	res.Body.Close()
	return nil
}
