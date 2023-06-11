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
var host = "https://192.168.3.181:3455"

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

func NameRegister(entry *NameEntry) error {
	form := url.Values{}
	form.Set("name", entry.Name)
	form.Set("ip", entry.Ip)
	form.Set("port", strconv.Itoa(entry.Port))
	form.Set("type", entry.Type)
	form.Set("description", entry.Description)
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

func NameList(params *ListQueryParam) error {
	res, err := httpsClient.Get(fmt.Sprintf("%s/list?match=%s&method=%s", host, params.Match, params.Method))
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
	flag.StringVar(&caCertFile, "ca", "../ca/ca.pem", "ca certificate")
	flag.StringVar(&cliCertFile, "crt", "../ca/net1.Brain.01/client.pem", "client certificate")
	flag.StringVar(&cliKeyFile, "key", "../ca/net1.Brain.01/client.key", "client private key")
	flag.Parse()

	InitHttpsClient(caCertFile, cliCertFile, cliKeyFile)

	fmt.Println("------- 1 -------")
	entry1 := NameEntry{
		Type:        "brain",
		Name:        "net1.Brain.01",
		Ip:          "192.168.3.181",
		Port:        3456,
		Description: "for test",
	}
	entry2 := NameEntry{
		Type:        "brain",
		Name:        "net1.Brain.02",
		Ip:          "192.168.3.182",
		Port:        3456,
		Description: "hello world",
	}
	err := NameRegister(&entry1)
	if err != nil {
		fmt.Println(err)
	}
	err = NameRegister(&entry2)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------- 2 -------")
	err = NameQuery(entry1.Name)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------- 3 -------")
	params := ListQueryParam{
		Match:  "Brain",
		Method: "contain",
	}
	err = NameList(&params)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------- 4 -------")
	err = NameDelete(entry1.Name)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------- 5 -------")
	err = NameQuery(entry1.Name)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------- 6 -------")
	err = GetSummary()
	if err != nil {
		fmt.Println(err)
	}
}
