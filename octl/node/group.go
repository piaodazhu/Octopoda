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

func GroupGetAll() {
	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_Groups,
	)
	res, err := http.Get(url)
	if err != nil {
		output.PrintFatalln("Get")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	output.PrintJSON(raw)
}

func GroupGet(name string) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_Group,
		name,
	)
	res, err := http.Get(url)
	if err != nil {
		output.PrintFatalln("Get")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	output.PrintJSON(raw)
}

func GroupDel(name string) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_Group,
		name,
	)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		output.PrintFatalln(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		output.PrintFatalln("DELETE, status=", res.StatusCode)
	}
	defer res.Body.Close()
}

func GroupSet(name string, nocheck bool, names []string) {
	nodes, err := NodesParse(names)
	if err != nil {
		output.PrintFatalln(err)
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
		output.PrintFatalln("POST")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode == 200 {
		output.PrintInfoln(string(raw))
	} else {
		output.PrintFatalln(string(raw))
	}
}
