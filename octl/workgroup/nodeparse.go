package workgroup

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
)

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
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodesParse,
	)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	SetHeader(req)
	req.Header.Set("Content-Type", "application/json")
	res, err := httpclient.BrainClient.Do(req)
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
