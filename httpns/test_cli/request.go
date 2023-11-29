package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/piaodazhu/Octopoda/protocols"
)

func NameRegister(entries ...*protocols.NameServiceEntry) error {
	body, _ := json.Marshal(entries)
	res, err := httpsClient.Post(host+"/register", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response protocols.Response
	json.Unmarshal(buf, &response)
	if response.Message != "OK" || res.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Message)
	}
	return nil
}

func NameQuery(name string) error {
	res, err := httpsClient.Get(fmt.Sprintf("%s/query?name=%s", host, name))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("NameQuery status code = %d", res.StatusCode)
	}
	var response protocols.Response
	err = json.Unmarshal(buf, &response)
	if err != nil {
		return err
	}
	fmt.Println(response.NameEntry)
	return nil
}

func NameList(params *protocols.ListQueryParam) error {
	res, err := httpsClient.Get(fmt.Sprintf("%s/list?match=%s&method=%s", host, params.Match, params.Method))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response protocols.Response
	json.Unmarshal(buf, &response)
	fmt.Println(response.NameList)
	return nil
}

func NameDelete(name string) error {
	form := url.Values{}
	form.Set("name", name)
	res, err := httpsClient.PostForm(host+"/delete", form)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response protocols.Response
	json.Unmarshal(buf, &response)
	fmt.Println(response.Message)
	return nil
}

func GetSummary() error {
	res, err := httpsClient.Get(host + "/summary")
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var summary protocols.Summary
	json.Unmarshal(buf, &summary)

	fmt.Println(summary.TotalRequests, time.UnixMilli(summary.Since).Format("2006-01-02 15:04:05"))
	for url, stats := range summary.ApiStats {
		fmt.Printf("url=%s, total request=%d, success request=%d\n", url, stats.Requests, stats.Success)
	}
	return nil
}
