package model

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var ScenarioMap map[string]*ScenarioModel
var ScenLock sync.RWMutex

// Scenario: a set of applications
type ScenarioModel struct {
	Id          uint32
	Name        string
	Description string
	Versions    []*ScenarioVersionModel
	newlevelbuf []*AppModel
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
	Name     string
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
			Id:          uuid.New().ID(),
			Name:        name,
			Description: description,
			Versions:    []*ScenarioVersionModel{},
			newlevelbuf: []*AppModel{},
		}
		ScenarioMap[name] = scen
	}
	return true
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

// func AddScenApp(scenario, app, node string, version string) bool {
// 	// ScenLock.RLock()
// 	// defer ScenLock.RUnlock()

// 	// var scen *ScenarioModel
// 	// var ok bool

// 	// if scen, ok = ScenarioMap[scenario]; !ok {
// 	// 	return false
// 	// }
// 	// scen.newlevelbuf = append(scen.newlevelbuf, &AppModel{
// 	// 	Name: app,
// 	// 	NodeApp: []*NodeAppModel{},
// 	// })

// }

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
