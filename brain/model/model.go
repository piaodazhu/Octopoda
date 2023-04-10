package model

import (
	"brain/logger"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ScenarioMap map[string]*ScenarioModel
var ScenLock sync.RWMutex

// Scenario: a set of applications
type ScenarioModel struct {
	Id            uint32
	Name          string
	Description   string
	Versions      []*ScenarioVersionModel
	newversionbuf []*AppModel
	modified      bool
}

// ScenarioVersion: a snapshot of a scenario
type ScenarioVersionModel struct {
	AppVersionModel
	Apps []*AppModel
}

// App: a kind of application in a scenario
type AppModel struct {
	Id          uint32
	Name        string
	Description string
	NodeApp     []*NodeAppModel
	// Scenario    *ScenarioModel
}

// NodeApp: a application instence on the node
type NodeAppModel struct {
	Name    string
	Version string
	// Versions []*AppVersionModel
}

// appversion: a snapshot of a NodeApp
type AppVersionModel struct {
	Version   string
	Message   string
	Timestamp int64
}

// -------------------------
// digest of a scenario
type ScenarioDigest struct {
	Name        string
	Description string
	Version     string
	Timestamp   int64
	Message     string
}

// detail info of a scenario
type ScenarioInfo struct {
	ScenarioDigest
	Apps []*AppInfo
}
type AppInfo struct {
	Name        string
	Description string
	NodeApps    []string
}

func InitScenarioMap() {
	ScenarioMap = make(map[string]*ScenarioModel)
	// go func() {
	// 	for {
	// 		time.Sleep(time.Second)
	// 		ScenLock.Lock()
	// 		for _, scenario := range ScenarioMap {
	// 			// if node.State == NodeStateDisconn && node.OfflineTs+int64(config.GlobalConfig.TentacleFace.RecordTimeout) < time.Now().Unix() {
	// 			// 	// logger.Tentacle.Print("MarkDeadNode", nodename)
	// 			// 	node.State = NodeStateDead
	// 			// }
	// 		}
	// 		ScenLock.Unlock()
	// 	}
	// }()
}

func AddScenario(name, description string) bool {
	ScenLock.Lock()
	defer ScenLock.Unlock()

	var scen *ScenarioModel
	if _, found := ScenarioMap[name]; found {
		return false
	} else {
		scen = &ScenarioModel{
			Id:            uuid.New().ID(),
			Name:          name,
			Description:   description,
			Versions:      []*ScenarioVersionModel{},
			newversionbuf: []*AppModel{},
			modified:      false,
		}
		ScenarioMap[name] = scen
	}
	return true
}

func UpdateScenario(name, message string) bool {
	ScenLock.Lock()
	defer ScenLock.Unlock()

	var scen *ScenarioModel
	var found bool
	if scen, found = ScenarioMap[name]; !found {
		return false
	} else {
		if scen.modified {
			versionhash := sha1.Sum([]byte(message + time.Now().String()))
			scen.Versions = append(scen.Versions, &ScenarioVersionModel{
				AppVersionModel: AppVersionModel{
					Version:   hex.EncodeToString(versionhash[:]),
					Message:   message,
					Timestamp: time.Now().Unix(),
				},
				Apps: scen.newversionbuf,
			})
			// triger save?
		}
		scen.newversionbuf = nil
		scen.modified = false
		logger.Brain.Println("Len of v = ", len(scen.Versions))
		return true
	}
}

func DelScenario(name string) {
	ScenLock.Lock()
	defer ScenLock.Unlock()

	delete(ScenarioMap, name)
}

type NodeAppItem struct {
	AppName  string
	ScenName string
	NodeName string
}

func GetNodeApps(name string) []NodeAppItem {
	ScenLock.Lock()
	defer ScenLock.Unlock()

	list := []NodeAppItem{}
	var scen *ScenarioModel
	scen, found := ScenarioMap[name]
	if !found {
		return list
	}
	for _, app := range scen.Versions[len(scen.Versions)-1].Apps {
		for _, node := range app.NodeApp {
			list = append(list, NodeAppItem{
				AppName:  app.Name,
				ScenName: scen.Name,
				NodeName: node.Name,
			})
		}

	}
	return list
}

func GetScenariosDigestAll() ([]*ScenarioDigest, bool) {
	ScenLock.RLock()
	defer ScenLock.RUnlock()

	if len(ScenarioMap) == 0 {
		return nil, false
	}
	res := make([]*ScenarioDigest, len(ScenarioMap))
	idx := 0
	for _, val := range ScenarioMap {
		res[idx].Name = val.Name
		res[idx].Description = val.Description
		res[idx].Timestamp = val.Versions[len(val.Versions)-1].Timestamp
		res[idx].Version = val.Versions[len(val.Versions)-1].Version
		res[idx].Message = val.Versions[len(val.Versions)-1].Message
		idx++
	}
	return res, true
}

func GetScenarioInfoByName(name string) (*ScenarioInfo, bool) {
	ScenLock.RLock()
	defer ScenLock.RUnlock()

	if scen, found := ScenarioMap[name]; found {
		// init but empty
		if len(scen.Versions) == 0 {
			return nil, true
		}
		res := &ScenarioInfo{
			ScenarioDigest: ScenarioDigest{
				Name:        scen.Name,
				Description: scen.Description,
				Version:     scen.Versions[len(scen.Versions)-1].Version,
				Message:     scen.Versions[len(scen.Versions)-1].Message,
				Timestamp:   scen.Versions[len(scen.Versions)-1].Timestamp,
			},
			Apps: []*AppInfo{},
		}
		for _, v := range scen.Versions[len(scen.Versions)-1].Apps {
			app := &AppInfo{
				Name:        v.Name,
				Description: v.Description,
				NodeApps:    []string{},
			}
			for _, n := range v.NodeApp {
				app.NodeApps = append(app.NodeApps, fmt.Sprintf("(%s - %s)", n.Name, n.Version))
			}
			res.Apps = append(res.Apps, app)
		}
		return res, true
	}
	return nil, false
}

// func startAddScenApp(scenario string) {

// }

// func endAddScenApp(scenario string) {

// }

// each time: add one (node, version) pair to the app
func AddScenNodeApp(scenario, app, description, node, version string, modified bool) bool {
	ScenLock.RLock()
	defer ScenLock.RUnlock()

	var scen *ScenarioModel
	var ok bool

	if scen, ok = ScenarioMap[scenario]; !ok {
		return false
	}
	for _, application := range scen.newversionbuf {
		if application.Name == app {
			for _, worknode := range application.NodeApp {
				if worknode.Name == node {
					logger.Brain.Println("NodeApp Cover");
					worknode.Version = version
					return true
				}
			}
			// build a new node app
			application.NodeApp = append(application.NodeApp, &NodeAppModel{
				Name:    node,
				Version: version,
			})
			if modified {
				scen.modified = true
			}
			return true
		}
	}
	// build a new application
	scen.newversionbuf = append(scen.newversionbuf, &AppModel{
		Name: app,
		Description: description,
	})
	scen.newversionbuf[len(scen.newversionbuf)-1].NodeApp = append(scen.newversionbuf[len(scen.newversionbuf)-1].NodeApp, &NodeAppModel{
		Name: node,
		Version: version,
	})
	if modified {
		scen.modified = true
	}

	return true
}

// func DelScenApp() {

// }

func ResetScenario() {

}

func ResetNodeApp() {

}

func SaveScenario() {

}

func LoadScenario() {

}
