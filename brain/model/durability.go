package model

import (
	"io"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/protocols"
)

var Busy int32

func SetBusy()         { atomic.StoreInt32(&Busy, 0) }
func SetReady()        { atomic.StoreInt32(&Busy, 1) }
func CheckReady() bool { return atomic.LoadInt32(&Busy) == 1 }

func InitScenarioMap() {
	SetBusy()
	defer SetReady()
	var file strings.Builder
	file.WriteString(config.GlobalConfig.Workspace.Root)
	file.WriteString(diskFileName)
	f, err := os.Open(file.String())
	if err != nil {
		logger.SysInfo.Printf("Scenario file not found. New...")
		ScenLock.Lock()
		ScenarioMap = make(map[string]*ScenarioModel)
		ScenLock.Unlock()
	} else {
		logger.SysInfo.Printf("Scenario file found. Loading...")
		defer f.Close()
		content, _ := io.ReadAll(f)
		ScenLock.Lock()
		if err := config.Jsoner.Unmarshal(content, &ScenarioMap); err != nil {
			logger.Exceptions.Fatal("Invalid scenarios file!")
		}
		ScenLock.Unlock()

		go AutoFix()
	}
}

// bad code...
type AppBasic struct {
	Name        string
	Scenario    string
	Description string
	Message     string
}

type Version struct {
	Time int64
	Hash string
	Msg  string
}

// Define some error type
type ErrorScenarioNotFound struct{}

func (ErrorScenarioNotFound) Error() string { return "ErrorScenarioNotFound" }

type ErrorScenarioDirty struct{}

func (ErrorScenarioDirty) Error() string { return "ErrorScenarioDirty" }

type ErrorNodeOffline struct{}

func (ErrorNodeOffline) Error() string { return "ErrorNodeOffline" }

type ErrorNodeDisconnect struct{}

func (ErrorNodeDisconnect) Error() string { return "ErrorNodeDisconnect" }

type ErrorNodeAppError struct{}

func (ErrorNodeAppError) Error() string { return "ErrorNodeAppError" }

type ErrorAddScenNodeApp struct{}

func (ErrorAddScenNodeApp) Error() string { return "ErrorAddScenNodeApp" }

type ErrorUpdateScenario struct{}

func (ErrorUpdateScenario) Error() string { return "ErrorUpdateScenario" }

func Fix(name string) error {
	var scen *ScenarioModel
	var found bool

	// get and copy current nodeapp versions
	ScenLock.RLock()
	if scen, found = ScenarioMap[name]; !found {
		ScenLock.RUnlock()
		return ErrorScenarioNotFound{}
	}
	// there are uncommitted change. Don't fix
	if scen.modified {
		ScenLock.RUnlock()
		return ErrorScenarioDirty{}
	}

	curNodeApps := []NodeAppItem{}
	if len(scen.Versions) > 0 {
		for _, app := range scen.Versions[len(scen.Versions)-1].Apps {
			for _, nodeapp := range app.NodeApp {
				curNodeApps = append(curNodeApps, NodeAppItem{
					AppName:  app.Name,
					ScenName: scen.Name,
					NodeName: nodeapp.Name,
					Version:  nodeapp.Version,
				})
			}
		}
	}
	ScenLock.RUnlock()

	// check real nodeapp versions
	curNodeAppLatestVersion := make([]Version, len(curNodeApps))
	for i := range curNodeApps {
		aParams := AppBasic{
			Name:     curNodeApps[i].AppName,
			Scenario: curNodeApps[i].ScenName,
		}

		payload, _ := config.Jsoner.Marshal(&aParams)
		raw, err := Request(curNodeApps[i].NodeName, protocols.TypeAppLatestVersion, payload)
		if err != nil {
			return ErrorNodeAppError{}
		}

		var latest Version
		err = config.Jsoner.Unmarshal(raw, &latest)
		if err != nil || len(latest.Hash) == 0 {
			return ErrorNodeAppError{}
		}

		curNodeAppLatestVersion[i] = latest
	}

	for i := range curNodeApps {
		// check the latest version
		if curNodeApps[i].Version != curNodeAppLatestVersion[i].Hash {
			if !AddScenNodeApp(curNodeApps[i].ScenName, curNodeApps[i].AppName, "", curNodeApps[i].NodeName, curNodeAppLatestVersion[i].Hash, true) {
				return ErrorAddScenNodeApp{}
			}
		}
	}

	_, ok := UpdateScenario(name, "System Data Fix")
	if !ok {
		// should not happen
		return ErrorUpdateScenario{}
	}
	return nil
}

func AutoFix() {
	// 10s to wait active nodes connect to brain
	time.Sleep(10 * time.Second)

	scenList := []string{}

	ScenLock.RLock()
	for name := range ScenarioMap {
		if err := Fix(name); err != nil {
			scenList = append(scenList, name)
		}
	}
	ScenLock.RUnlock()

	for len(scenList) > 0 {
		unfixed := []string{}
		logger.SysInfo.Printf("Scenarios to be fixed: %v", scenList)
		for _, name := range scenList {
			// each 30s, check a scenario
			time.Sleep(30 * time.Second)
			// fix failed is allowed. Because scenarios are dynamically changing.
			if err := Fix(name); err != nil {
				unfixed = append(unfixed, name)
			}
		}
		scenList = unfixed
		// sleep 1min then start next roll fix
		time.Sleep(60 * time.Second)
	}
	logger.SysInfo.Println("All scenarios are fixed.")
}

func saveNoLock() {
	var file strings.Builder
	file.WriteString(config.GlobalConfig.Workspace.Root)
	file.WriteString(diskFileName)
	_, err := os.Stat(file.String())
	if os.IsExist(err) {
		os.Rename(file.String(), file.String()+".bk")
	}
	serialized, _ := config.Jsoner.Marshal(&ScenarioMap)
	err = os.WriteFile(file.String(), serialized, os.ModePerm)
	if err != nil {
		logger.Exceptions.Print("cannot WriteFile")
	}
}
