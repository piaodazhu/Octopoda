package main

import (
	"encoding/json"
	"fmt"
)

func basicTest() {
	var err error
	entry1 := RegisterParam{
		Type:        "brain",
		Name:        "net1.Brain.01",
		Ip:          "192.168.3.181",
		Port:        3456,
		Description: "for test",
		TTL:         0,
	}
	entry2 := RegisterParam{
		Type:        "brain",
		Name:        "net1.Brain.02",
		Ip:          "192.168.3.182",
		Port:        3456,
		Description: "hello world",
		TTL:         10,
	}
	entry3 := RegisterParam{
		Type:        "tentacle",
		Name:        "net1.Tentacle.01",
		Ip:          "192.168.3.183",
		Port:        3456,
		Description: "3333",
		TTL:         10,
	}
	entry4 := RegisterParam{
		Type:        "tentacle",
		Name:        "net1.Tentacle.02",
		Ip:          "192.168.3.184",
		Port:        3456,
		Description: "4444",
		TTL:         10,
	}
	fakeconf1, _ := json.Marshal(entry1)
	fakeconf2, _ := json.Marshal(entry2)
	fakeconf3, _ := json.Marshal(entry3)
	fakeconf4, _ := json.Marshal(entry4)
	conf1 := ConfigUploadParam{
		Type:      "brain",
		Name:      "net1.Brain.01",
		Method:    "append",
		RawConfig: string(fakeconf1),
	}
	conf2 := ConfigUploadParam{
		Type:      "brain",
		Name:      "net1.Brain.01",
		Method:    "append",
		RawConfig: string(fakeconf2),
	}
	conf3 := ConfigUploadParam{
		Type:      "brain",
		Name:      "net1.Brain.01",
		Method:    "append",
		RawConfig: string(fakeconf3),
	}
	conf4 := ConfigUploadParam{
		Type:      "brain",
		Name:      "net1.Brain.01",
		Method:    "reset",
		RawConfig: string(fakeconf4),
	}
	conf5 := ConfigUploadParam{
		Type:      "brain",
		Name:      "net1.Brain.01",
		Method:    "clear",
		RawConfig: "{should be ignored}",
	}
	sshinfo1 := SshInfoUploadParam{
		Type:     "brain",
		Name:     "net1.Brain.05",
		Username: "pi",
		Ip:       "192.168.3.181",
		Port:     22,
		Password: "sshinfo1",
	}
	sshinfo2 := SshInfoUploadParam{
		Type:     "brain",
		Name:     "net1.Brain.06",
		Username: "pi",
		Ip:       "192.168.3.182",
		Port:     33,
		Password: "sshinfo2",
	}

	fmt.Println("------- register name -------")
	err = NameRegister(&entry1)
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
		Name:   conf1.Name,
		Index:  0,
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
		Scope:  "name",
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
