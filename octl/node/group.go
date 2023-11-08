package node

import (
	"bytes"
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

func GroupGetAll() ([]string, *errs.OctlError) {
	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_Groups,
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

	groups := []string{}
	err = json.Unmarshal(raw, &groups)
	if err != nil {
		emsg := "res unmarshal error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlMessageParseError, emsg)
	}

	output.PrintJSON(groups)
	return groups, nil
}

func GroupGet(name string) (*protocols.GroupInfo, *errs.OctlError) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_Group,
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

	group := protocols.GroupInfo{}
	err = json.Unmarshal(raw, &group)
	if err != nil {
		emsg := "res unmarshal error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlMessageParseError, emsg)
	}

	output.PrintJSON(group)
	return &group, nil
}

func GroupDel(name string) *errs.OctlError {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_Group,
		name,
	)
	req, _ := http.NewRequest("DELETE", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		emsg := "http delete error: " + err.Error()
		output.PrintFatalln(emsg)
		return errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		emsg := fmt.Sprintf("[%d]msg=%s\n", res.StatusCode, string(raw))
		output.PrintFatalln(emsg)
		return errs.New(errs.OctlHttpStatusError, emsg)
	}

	return nil
}

func GroupSet(name string, nocheck bool, names []string) *errs.OctlError {
	nodes, err := NodesParse(names)
	if err != nil {
		emsg := "node parse: " + err.Error()
		output.PrintFatalln(emsg)
		return errs.New(errs.OctlNodeParseError, emsg)
	}
	ginfo := protocols.GroupInfo{
		Name:    name,
		Nodes:   nodes,
		NoCheck: nocheck,
	}
	body, _ := json.Marshal(ginfo)

	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_Group,
	)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		emsg := fmt.Sprintf("[%d]msg=%s\n", res.StatusCode, string(raw))
		output.PrintFatalln(emsg)
		return errs.New(errs.OctlHttpStatusError, emsg)
	}
	output.PrintInfoln(string(raw))
	return nil
}
