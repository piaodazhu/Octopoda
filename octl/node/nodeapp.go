package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/httpclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

func NodeAppsInfo(node string) ([][]byte, *errs.OctlError) {
	url := fmt.Sprintf("https://%s/%s%s?name=%s",
		httpclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodeApps,
		node,
	)
	res, err := httpclient.BrainClient.Get(url)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	itemList := []interface{}{}
	err = json.Unmarshal(raw, &itemList)
	if err != nil {
		emsg := "unmarshal list error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlMessageParseError, emsg)
	}
	rawList := [][]byte{}
	for _, item := range itemList {
		rawItem, _ := json.Marshal(item)
		rawList = append(rawList, rawItem)
	}
	output.PrintJSON(raw)
	return rawList, nil
}

func NodeAppInfo(node, app, scenario string) ([]byte, *errs.OctlError) {
	url := fmt.Sprintf("https://%s/%s%s?name=%s&app=%s&scenario=%s",
		httpclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodeAppVersion,
		node,
		app,
		scenario,
	)
	res, err := httpclient.BrainClient.Get(url)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	output.PrintJSON(raw)
	return raw, nil
}

func NodeAppReset(node, app, scenario, version, message string) (string, *errs.OctlError) {
	url := fmt.Sprintf("https://%s/%s%s",
		httpclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodeAppVersion,
	)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	contentType := writer.FormDataContentType()
	writer.WriteField("name", node)
	writer.WriteField("app", app)
	writer.WriteField("scenario", scenario)
	writer.WriteField("version", version)
	writer.WriteField("message", message)
	writer.Close()

	res, err := httpclient.BrainClient.Post(url, contentType, body)
	if err != nil {
		emsg := "http post error: " + err.Error()
		output.PrintFatalln(emsg)
		return emsg, errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	output.PrintJSON(raw)
	return string(raw), nil
}
