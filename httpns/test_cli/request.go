package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"
)

func NameRegister(entry *RegisterParam) error {
	form := url.Values{}
	form.Set("name", entry.Name)
	form.Set("ip", entry.Ip)
	form.Set("port", strconv.Itoa(entry.Port))
	form.Set("type", entry.Type)
	form.Set("description", entry.Description)
	form.Set("ttl", strconv.Itoa(entry.TTL))
	res, err := httpsClient.PostForm(host+"/register", form)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response Response
	json.Unmarshal(buf, &response)
	fmt.Println(response.Message)
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
	var response Response
	json.Unmarshal(buf, &response)
	fmt.Println(response.NameEntry)
	return nil
}

func ConfigRegister(conf *ConfigUploadParam) error {
	form := url.Values{}
	form.Set("name", conf.Name)
	form.Set("method", conf.Method)
	form.Set("type", conf.Type)
	form.Set("conf", conf.RawConfig)
	res, err := httpsClient.PostForm(host+"/conf", form)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response Response
	json.Unmarshal(buf, &response)
	fmt.Println(response.Message)
	return nil
}

func ConfigQuery(query *ConfigQueryParam) error {
	res, err := httpsClient.Get(fmt.Sprintf("%s/conf?name=%s&index=%d&amount=%d", host, query.Name, query.Index, query.Amount))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response Response
	json.Unmarshal(buf, &response)
	fmt.Println("configlist:")
	for _, c := range response.RawConfig {
		fmt.Println("  ", *c)
	}
	return nil
}

func SshinfoRegister(sshinfo *SshInfoUploadParam) error {
	form := url.Values{}
	form.Set("name", sshinfo.Name)
	form.Set("ip", sshinfo.Ip)
	form.Set("port", strconv.Itoa(sshinfo.Port))
	form.Set("type", sshinfo.Type)
	form.Set("username", sshinfo.Username)
	form.Set("password", sshinfo.Password)
	res, err := httpsClient.PostForm(host+"/sshinfo", form)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response Response
	json.Unmarshal(buf, &response)
	fmt.Println(response.Message)
	return nil
}

func SshinfoQuery(name string) error {
	res, err := httpsClient.Get(fmt.Sprintf("%s/sshinfo?name=%s", host, name))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response Response
	json.Unmarshal(buf, &response)
	fmt.Println(response.SshInfo)
	return nil
}

func NameList(params *ListQueryParam) error {
	res, err := httpsClient.Get(fmt.Sprintf("%s/list?match=%s&method=%s&scope=%s", host, params.Match, params.Method, params.Scope))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response Response
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
	var response Response
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
	var summary Summary
	json.Unmarshal(buf, &summary)

	fmt.Println(summary.TotalRequests, time.UnixMilli(summary.Since).Format("2006-01-02 15:04:05"))
	for url, stats := range summary.ApiStats {
		fmt.Printf("url=%s, total request=%d, success request=%d\n", url, stats.Requests, stats.Success)
	}
	return nil
}
