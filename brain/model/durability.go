package model

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"encoding/json"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

var FirstFixed chan struct{}

func InitScenarioMap() {
	FirstFixed = make(chan struct{}, 1)
	var file strings.Builder
	file.WriteString(config.GlobalConfig.Workspace.Root)
	file.WriteString(diskFileName)
	f, err := os.Open(file.String())
	if err != nil {
		logger.Brain.Printf("Scenario file not found. New...")
		ScenLock.Lock()
		ScenarioMap = make(map[string]*ScenarioModel)
		ScenLock.Unlock()
		// new storage need not fix
		FirstFixed <- struct{}{}
	} else {
		logger.Brain.Printf("Scenario file found. Loading...")
		defer f.Close()
		content, _ := io.ReadAll(f)
		ScenLock.Lock()
		if err := json.Unmarshal(content, &ScenarioMap); err != nil {
			logger.Brain.Fatal("Invalid scenarios file!")
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
	ScenLock.RUnlock()

	// check real nodeapp versions
	for i := range curNodeApps {
		if addr, ok := GetNodeAddress(curNodeApps[i].NodeName); !ok {
			return ErrorNodeOffline{}
		} else {
			if conn, err := net.Dial("tcp", addr); err != nil {
				return ErrorNodeDisconnect{}
			} else {
				defer conn.Close()
				aParams := AppBasic{
					Name:     curNodeApps[i].AppName,
					Scenario: curNodeApps[i].ScenName,
				}

				payload, _ := json.Marshal(&aParams)
				message.SendMessage(conn, message.TypeAppLatestVersion, payload)
				mtype, raw, err := message.RecvMessage(conn)
				if err != nil || mtype != message.TypeAppLatestVersionResponse {
					return ErrorNodeAppError{}
				}

				var latest Version
				err = json.Unmarshal(raw, &latest)
				if err != nil || len(latest.Hash) == 0 {
					return ErrorNodeAppError{}
				}

				// check the latest version
				if curNodeApps[i].Version != latest.Hash {
					if !AddScenNodeApp(curNodeApps[i].ScenName, curNodeApps[i].AppName, "", curNodeApps[i].NodeName, latest.Hash, true) {
						return ErrorAddScenNodeApp{}
					}
				}
			}
		}
	}
	if !UpdateScenario(name, "System Data Fix") {
		return ErrorUpdateScenario{}
	}
	return nil
}

func AutoFix() {
	// 3s to wait active nodes connect to brain
	time.Sleep(10 * time.Second)

	for name := range ScenarioMap {
		if err := Fix(name); err != nil {
			logger.Brain.Printf("Warning: scenario < %s >: %s Try to manually fix this scenario later.", name, err.Error())
		}
	}
	// first roll fix done. Then http request can be accepted
	logger.Brain.Println("First Fix Done")
	FirstFixed <- struct{}{}

	for {
		namelist := []string{}
		ScenLock.RLock()
		for name := range ScenarioMap {
			namelist = append(namelist, name)
		}
		ScenLock.RUnlock()

		for _, name := range namelist {
			// each 30s, check a scenario
			time.Sleep(30 * time.Second)

			// fix failed is allowed. Because scenarios are dynamically changing.
			if err := Fix(name); err != nil {
				logger.Brain.Printf("Warning: scenario < %s >: %s Try to manually fix this scenario later.", name, err.Error())
			}
		}

		// sleep 1min then start next roll fix
		time.Sleep(60 * time.Second)
	}
}

func saveNoLock() {
	var file strings.Builder
	file.WriteString(config.GlobalConfig.Workspace.Root)
	file.WriteString(diskFileName)
	_, err := os.Stat(file.String())
	if os.IsExist(err) {
		os.Rename(file.String(), file.String()+".bk")
	}
	serialized, _ := json.Marshal(&ScenarioMap)
	err = os.WriteFile(file.String(), serialized, os.ModePerm)
	if err != nil {
		logger.Brain.Print("cannot WriteFile")
	}
}
