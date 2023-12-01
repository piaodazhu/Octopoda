package scenario

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/httpclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/workgroup"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

func ScenariosInfo() ([][]byte, *errs.OctlError) {
	url := fmt.Sprintf("https://%s/%s%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_ScenariosInfo,
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

func ScenarioInfo(name string) ([]byte, *errs.OctlError) {
	url := fmt.Sprintf("https://%s/%s%s?name=%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_ScenarioInfo,
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

	output.PrintJSON(raw)
	return raw, nil
}

func ScenarioFix(name string) (string, *errs.OctlError) {
	url := fmt.Sprintf("https://%s/%s%s?name=%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_ScenarioFix,
		name,
	)
	req, _ := http.NewRequest("GET", url, nil)
	workgroup.SetHeader(req)
	res, err := httpclient.BrainClient.Do(req)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return emsg, errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	output.PrintJSON(raw)
	return string(raw), nil
}

func ScenarioVersion(name string, offset, limit int) ([]byte, *errs.OctlError) {
	url := fmt.Sprintf("https://%s/%s%s?name=%s&offset=%d&limit=%d",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_ScenarioVersion,
		name,
		offset,
		limit,
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

	output.PrintJSON(raw)
	return raw, nil
}

func ScenarioReset(name string, version string, message string) (string, *errs.OctlError) {
	url := fmt.Sprintf("https://%s/%s%s?name=%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_ScenarioVersion,
		name,
	)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	contentType := writer.FormDataContentType()
	writer.WriteField("name", name)
	writer.WriteField("version", version)
	writer.WriteField("message", message)
	writer.Close()

	req, _ := http.NewRequest("POST", url, body)
	workgroup.SetHeader(req)
	req.Header.Set("Content-Type", contentType)
	res, err := httpclient.BrainClient.Do(req)
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
