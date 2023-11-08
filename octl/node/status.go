package node

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

func NodeStatus(name string) (*protocols.Status, *errs.OctlError) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodeStatus,
		name,
	)
	res, err := http.Get(url)
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

	info := protocols.Status{}
	err = json.Unmarshal(raw, &info)
	if err != nil {
		emsg := "res unmarshal error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlMessageParseError, emsg)
	}
	output.PrintJSON(info.ToText())

	return &info, nil
}

func NodesStatus(names []string) (*protocols.NodesStatus, *errs.OctlError) {
	nodes, err := NodesParse(names)
	if err != nil {
		emsg := "node parse." + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlNodeParseError, emsg)
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

	info := protocols.NodesStatus{}
	err = json.Unmarshal(raw, &info)
	if err != nil {
		emsg := "res unmarshal error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlMessageParseError, emsg)
	}
	output.PrintJSON(info.ToText())

	return &info, nil
}
