package main

import (
	"fmt"

	"github.com/piaodazhu/Octopoda/protocols"
)

func basicTest() {
	var err error
	entry1 := protocols.NameServiceEntry{
		Key:         "test1",
		Type:        "addr",
		Value:       "value1",
		Description: "desc",
		TTL:         0,
	}
	entry2 := protocols.NameServiceEntry{
		Key:         "test2",
		Type:        "string",
		Value:       "value2",
		Description: "desc",
		TTL:         0,
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

	fmt.Println("------- query name -------")
	err = NameQuery(entry1.Key)
	if err != nil {
		fmt.Println(err)
	}
	err = NameQuery(entry2.Key)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------- query list -------")
	lQuery := protocols.ListQueryParam{
		Match:  "test",
		Method: "prefix",
	}

	err = NameList(&lQuery)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------- get summary -------")
	err = GetSummary()
	if err != nil {
		fmt.Println(err)
	}
}
