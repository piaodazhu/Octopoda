package node

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/output"
)

func NodeAppsInfo(node string) (string, error) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodeApps,
		node,
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

func NodeAppInfo(node, app, scenario string) (string, error) {
	url := fmt.Sprintf("http://%s/%s%s?name=%s&app=%s&scenario=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodeAppVersion,
		node,
		app,
		scenario,
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

func NodeAppReset(node, app, scenario, version, message string) (string, error) {
	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
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

	res, err := http.Post(url, contentType, body)
	if err != nil {
		emsg := "http post error."
		output.PrintFatalln(emsg, err)
		return emsg, err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	output.PrintJSON(raw)
	return string(raw), nil
}
