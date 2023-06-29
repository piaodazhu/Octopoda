package main

import (
	"encoding/json"
	"fmt"
	"os"
	"pakma/config"
	"time"
)

// path: brain_v1.3.0_linux_amd64
// version: 1.3.0

const (
	EMPTY   = iota + 1 // ([ newpack ])
	STABLE             // ([ old ]) [ this ] ([ newpack ])
	PREVIEW            // [ old ] [ prev ] [ this ]
)

type PakmaState struct {
	StateType int
	Version1  string
	Version2  string
	Version3  string
}

var StateMsg map[int]string

func getPathFromVersion(version string) string {
	return fmt.Sprintf("%s_v%s_%s_%s", config.GlobalConfig.AppName, version, config.GlobalConfig.AppOS, config.GlobalConfig.AppArch)
}

func dumpState() {
	serialized, _ := json.Marshal(State)
	if fileExists(config.GlobalConfig.Packma.Root + "pakma.json") {
		os.Rename(config.GlobalConfig.Packma.Root+"pakma.json", config.GlobalConfig.Packma.Root+"pakma.json.bk")
	}
	os.WriteFile(config.GlobalConfig.Packma.Root+"pakma.json", serialized, os.ModePerm)
}

func loadState() error {
	if !fileExists(config.GlobalConfig.Packma.Root + "pakma.json") {
		return fmt.Errorf("pakma.json not found")
	}
	serialized, err := os.ReadFile(config.GlobalConfig.Packma.Root + "pakma.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(serialized, &State)
	if err != nil {
		return err
	}
	if State.StateType == PREVIEW {
		time.AfterFunc(time.Second*time.Duration(config.GlobalConfig.Packma.PreviewDuration), doCancel)
	}
	return nil
}

func initState() {
	State = PakmaState{
		StateType: EMPTY,
		Version1:  "",
		Version2:  "",
		Version3:  "",
	}
	dumpState()
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func dirExists(dirname string) bool {
	s, err := os.Stat(dirname)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return s.IsDir()
}
