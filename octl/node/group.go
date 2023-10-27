package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"octl/config"
	"octl/nameclient"
	"octl/output"
	"protocols"
)

func GroupGetAll() (string, error) {
	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_Groups,
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

func GroupGet(name string) (string, error) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_Group,
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

func GroupDel(name string) (string, error) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_Group,
		name,
	)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		emsg := "http delete error."
		output.PrintFatalln(emsg, err)
		return emsg, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		output.PrintFatalln("DELETE, status=", res.StatusCode)
	}
	defer res.Body.Close()
	return "OK", nil
}

func GroupSet(name string, nocheck bool, names []string) (string, error) {
	nodes, err := NodesParse(names)
	if err != nil {
		msg := "node parse."
		output.PrintFatalln(msg, err)
		return msg, err
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
		emsg := "http get error."
		output.PrintFatalln(emsg, err)
		return emsg, err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode == 200 {
		output.PrintInfoln(string(raw))
	} else {
		output.PrintFatalln(string(raw))
	}
	return string(raw), nil
}
