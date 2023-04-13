package scenario

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"octl/config"
)

func ScenariosInfo() {
	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.ScenariosInfo,
	)
	res, err := http.Get(url)
	if err != nil {
		panic("Get")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	fmt.Println(string(raw))
}

func ScenarioInfo(name string) {
	url := fmt.Sprintf("http://%s:%d/%s%s?name=%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.ScenarioInfo,
		name,
	)
	res, err := http.Get(url)
	if err != nil {
		panic("Get")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	fmt.Println(string(raw))
}

func ScenarioVersion(name string) {
	url := fmt.Sprintf("http://%s:%d/%s%s?name=%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.ScenarioVersion,
		name,
	)
	res, err := http.Get(url)
	if err != nil {
		panic("Get")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	fmt.Println(string(raw))
}

func ScenarioReset(name string, version string, message string) {
	url := fmt.Sprintf("http://%s:%d/%s%s?name=%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.ScenarioVersion,
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
		panic("Post")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	fmt.Println(string(raw))
}
