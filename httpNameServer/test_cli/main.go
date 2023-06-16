package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var httpsClient *http.Client
var host = "https://127.0.0.1:3455"

func InitHttpsClient(caCert, cliCert, cliKey string) {
	ca, err := os.ReadFile(caCert)
	if err != nil {
		log.Fatalln(err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(ca)

	clientCrt, err := tls.LoadX509KeyPair(cliCert, cliKey)
	if err != nil {
		log.Fatalln(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            certPool,
			InsecureSkipVerify: false,
			ClientAuth:         tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{
				clientCrt,
			},
		},
	}
	httpsClient = &http.Client{
		Transport: tr,
	}
}

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

func main() {
	var port int
	var caCertFile string
	var cliCertFile string
	var cliKeyFile string
	flag.IntVar(&port, "p", 9931, "listening port")
	flag.StringVar(&caCertFile, "ca", "ca.pem", "ca certificate")
	flag.StringVar(&cliCertFile, "crt", "client.pem", "client certificate")
	flag.StringVar(&cliKeyFile, "key", "client.key", "client private key")
	flag.Parse()

	InitHttpsClient(caCertFile, cliCertFile, cliKeyFile)
	
	entry1 := RegisterParam{
		Type:        "brain",
		Name:        "net1.Brain.01",
		Ip:          "192.168.3.181",
		Port:        3456,
		Description: "for test",
		TTL: 0,
	}
	entry2 := RegisterParam{
		Type:        "brain",
		Name:        "net1.Brain.02",
		Ip:          "192.168.3.182",
		Port:        3456,
		Description: "hello world",
		TTL: 10,
	}
	entry3 := RegisterParam{
		Type:        "tentacle",
		Name:        "net1.Tentacle.01",
		Ip:          "192.168.3.183",
		Port:        3456,
		Description: "3333",
		TTL: 10,
	}
	entry4 := RegisterParam{
		Type:        "tentacle",
		Name:        "net1.Tentacle.02",
		Ip:          "192.168.3.184",
		Port:        3456,
		Description: "4444",
		TTL: 10,
	}
	fakeconf1, _ := json.Marshal(entry1)
	fakeconf2, _ := json.Marshal(entry2)
	fakeconf3, _ := json.Marshal(entry3)
	fakeconf4, _ := json.Marshal(entry4)
	conf1 := ConfigUploadParam{
		Type: "brain",
		Name: "net1.Brain.01",
		Method: "append",
		RawConfig: string(fakeconf1),
	}
	conf2 := ConfigUploadParam{
		Type: "brain",
		Name: "net1.Brain.01",
		Method: "append",
		RawConfig: string(fakeconf2),
	}
	conf3 := ConfigUploadParam{
		Type: "brain",
		Name: "net1.Brain.01",
		Method: "append",
		RawConfig: string(fakeconf3),
	}
	conf4 := ConfigUploadParam{
		Type: "brain",
		Name: "net1.Brain.01",
		Method: "reset",
		RawConfig: string(fakeconf4),
	}
	conf5 := ConfigUploadParam{
		Type: "brain",
		Name: "net1.Brain.01",
		Method: "clear",
		RawConfig: "{should be ignored}",
	}
	sshinfo1 := SshInfoUploadParam{
		Type: "brain",
		Name: "net1.Brain.05",
		Username: "pi",
		Ip: "192.168.3.181",
		Port: 22,
		Password: "sshinfo1",
	}
	sshinfo2 := SshInfoUploadParam{
		Type: "brain",
		Name: "net1.Brain.06",
		Username: "pi",
		Ip: "192.168.3.182",
		Port: 33,
		Password: "sshinfo2",
	}

	fmt.Println("------- register name -------")
	err := NameRegister(&entry1)
	if err != nil {
		fmt.Println(err)
	}
	err = NameRegister(&entry2)
	if err != nil {
		fmt.Println(err)
	}
	err = NameRegister(&entry3)
	if err != nil {
		fmt.Println(err)
	}
	err = NameRegister(&entry4)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------- query name -------")
	err = NameQuery(entry1.Name)
	if err != nil {
		fmt.Println(err)
	}
	err = NameQuery(entry2.Name)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------- register configs -------")
	err = ConfigRegister(&conf1)
	if err != nil {
		fmt.Println(err)
	}
	err = ConfigRegister(&conf2)
	if err != nil {
		fmt.Println(err)
	}
	err = ConfigRegister(&conf3)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------- query configs -------")
	q := ConfigQueryParam{
		Name: conf1.Name,
		Index: 0,
		Amount: 2,
	}
	err = ConfigQuery(&q)
	if err != nil {
		fmt.Println(err)
	}
	q.Amount = 1
	err = ConfigQuery(&q)
	if err != nil {
		fmt.Println(err)
	}
	q.Index = 1
	q.Amount = 1
	err = ConfigQuery(&q)
	if err != nil {
		fmt.Println(err)
	}
	
	fmt.Println("------- register sshinfo -------")
	err = SshinfoRegister(&sshinfo1)
	if err != nil {
		fmt.Println(err)
	}
	err = SshinfoRegister(&sshinfo2)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------- query sshinfo -------")
	err = SshinfoQuery(sshinfo1.Name)
	if err != nil {
		fmt.Println(err)
	}
	err = SshinfoQuery(sshinfo2.Name)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------- query list -------")
	lqp := ListQueryParam{
		Scope: "name",
		Method: "all",
	}
	err = NameList(&lqp)
	if err != nil {
		fmt.Println(err)
	}
	lqp.Scope = "config"
	err = NameList(&lqp)
	if err != nil {
		fmt.Println(err)
	}
	lqp.Scope = "ssh"
	err = NameList(&lqp)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------- register and query confs again -------")
	err = ConfigRegister(&conf4)
	if err != nil {
		fmt.Println(err)
	}
	q.Index = 0
	q.Amount = 10
	err = ConfigQuery(&q)
	if err != nil {
		fmt.Println(err)
	}
	err = ConfigRegister(&conf5)
	if err != nil {
		fmt.Println(err)
	}
	q.Index = 0
	q.Amount = 10
	err = ConfigQuery(&q)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------- get summary -------")
	err = GetSummary()
	if err != nil {
		fmt.Println(err)
	}
}
