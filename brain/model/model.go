package model

import (
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
	Scenario    *ScenarioModel
}

// NodeApp: a application instence on the node
type NodeAppModel struct {
	Name     string
	Versions []*AppVersionModel
}

// appversion: a snapshot of a NodeApp
type AppVersionModel struct {
	Version   string
	Message   string
	Timestamp int64
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
			Id: uuid.New().ID(),
			Name: name,
			Description: description,
			Versions: []*ScenarioVersionModel{},
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
	AppName string
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
				AppName: app.Name,
				ScenName: scen.Name,
				NodeName: node.Name,
			})
		}
		
	}
	return list
}

func AddScenApp() {

}

func DelScenApp() {

}

func ResetScenario() {

}

func ResetNodeApp() {

}

func SaveScenario() {

}

func LoadScenario() {

}
