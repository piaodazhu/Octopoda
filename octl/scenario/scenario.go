package scenario

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

func ScenariosInfo() (string, *errs.OctlError) {
	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_ScenariosInfo,
	)
	res, err := http.Get(url)
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

func ScenarioInfo(name string) (string, *errs.OctlError) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_ScenarioInfo,
		name,
	)
	res, err := http.Get(url)
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

func ScenarioFix(name string) (string, *errs.OctlError) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_ScenarioFix,
		name,
	)
	res, err := http.Get(url)
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

func ScenarioVersion(name string) (string, *errs.OctlError) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_ScenarioVersion,
		name,
	)
	res, err := http.Get(url)
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

func ScenarioReset(name string, version string, message string) (string, *errs.OctlError) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
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

	res, err := http.Post(url, contentType, body)
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
