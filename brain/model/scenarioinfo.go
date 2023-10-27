package model

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/piaodazhu/Octopoda/brain/logger"

	"github.com/google/uuid"
)

var ScenarioMap map[string]*ScenarioModel
var ScenLock sync.RWMutex

const diskFileName = "scenarios.json"

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
	BasicVersionModel
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
	// Versions []*BasicVersionModel
}

// appversion: a snapshot of a NodeApp
type BasicVersionModel struct {
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

func AddScenario(name, description string) bool {
	ScenLock.Lock()
	defer ScenLock.Unlock()
	defer saveNoLock()

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

func UpdateScenario(name, message string) (bool, bool) {
	ScenLock.Lock()
	defer ScenLock.Unlock()
	defer saveNoLock()

	var scen *ScenarioModel
	var found bool
	if scen, found = ScenarioMap[name]; !found {
		return false, false
	} else {
		hasModified := scen.modified
		if hasModified {
			versionhash := sha1.Sum([]byte(message + time.Now().String()))
			scen.Versions = append(scen.Versions, &ScenarioVersionModel{
				BasicVersionModel: BasicVersionModel{
					Version:   hex.EncodeToString(versionhash[:]),
					Message:   message,
					Timestamp: time.Now().UnixMilli(),
				},
				Apps: scen.newversionbuf,
			})
			// triger save?
			scen.newversionbuf = cloneLayer(scen.newversionbuf) // must deep copy!
			scen.modified = false
		}
		// scen.newversionbuf = cloneLayer(scen.newversionbuf) // must deep copy!
		// scen.modified = false
		// logger.Brain.Println("Len of v = ", len(scen.Versions))
		return hasModified, true
	}
}

func cloneLayer(prototype []*AppModel) []*AppModel {
	res := []*AppModel{}
	for _, app := range prototype {
		nodeapps := []*NodeAppModel{}
		for _, nodeapp := range app.NodeApp {
			nodeapps = append(nodeapps, &NodeAppModel{
				Name:    nodeapp.Name,
				Version: nodeapp.Version,
			})
		}
		res = append(res, &AppModel{
			Id:          app.Id,
			Name:        app.Name,
			Description: app.Description,
			NodeApp:     nodeapps,
		})
	}
	return res
}

func DelScenario(name string) {
	ScenLock.Lock()
	defer ScenLock.Unlock()
	defer saveNoLock()

	delete(ScenarioMap, name)
}

type NodeAppItem struct {
	AppName  string
	ScenName string
	NodeName string
	Version  string
}

func GetNodeApps(name string, version string) []NodeAppItem {
	ScenLock.RLock()
	defer ScenLock.RUnlock()

	list := []NodeAppItem{}
	var scen *ScenarioModel
	scen, found := ScenarioMap[name]
	if !found {
		return list
	}
	if len(scen.Versions) == 0 {
		logger.Exceptions.Println("GetNodeApps Empty Scenario")
		return list
	}

	// find the version index of the scenario
	idx := -1
	if version == "" {
		// default: get latest version
		idx = len(scen.Versions) - 1
	} else {
		// for given version
		for i, v := range scen.Versions {
			if version == v.Version {
				idx = i
			}
		}

	}
	// if not found
	if idx < 0 {
		logger.Exceptions.Println("GetNodeApps Invalid Version")
		return list
	}

	for _, app := range scen.Versions[idx].Apps {
		for _, node := range app.NodeApp {
			list = append(list, NodeAppItem{
				AppName:  app.Name,
				ScenName: scen.Name,
				NodeName: node.Name,
				Version:  node.Version,
			})
		}

	}
	return list
}

func GetScenariosDigestAll() ([]ScenarioDigest, bool) {
	ScenLock.RLock()
	defer ScenLock.RUnlock()

	if len(ScenarioMap) == 0 {
		return []ScenarioDigest{}, true
	}
	res := make([]ScenarioDigest, len(ScenarioMap))
	idx := 0
	for _, val := range ScenarioMap {
		res[idx].Name = val.Name
		res[idx].Description = val.Description
		if len(val.Versions) != 0 {
			res[idx].Timestamp = val.Versions[len(val.Versions)-1].Timestamp
			res[idx].Version = val.Versions[len(val.Versions)-1].Version
			res[idx].Message = val.Versions[len(val.Versions)-1].Message
		}
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
			logger.Exceptions.Println("GetScenarioInfoByName Empty Scenario")
			return &ScenarioInfo{
				ScenarioDigest: ScenarioDigest{
					Name:        scen.Name,
					Description: scen.Description,
				}}, true
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

func GetScenarioVersionByName(name string) []BasicVersionModel {
	ScenLock.RLock()
	defer ScenLock.RUnlock()

	versions := []BasicVersionModel{}
	if scen, found := ScenarioMap[name]; found {
		for _, v := range scen.Versions {
			versions = append(versions, v.BasicVersionModel)
		}
	}
	return versions
}

// each time: add one (node, version) pair to the app
func AddScenNodeApp(scenario, app, description, node, version string, modified bool) bool {
	ScenLock.Lock()
	defer ScenLock.Unlock()

	var scen *ScenarioModel
	var ok bool

	if scen, ok = ScenarioMap[scenario]; !ok {
		return false
	}
	for _, application := range scen.newversionbuf {
		if application.Name == app {
			for _, worknode := range application.NodeApp {
				if worknode.Name == node {
					worknode.Version = version
					if modified {
						scen.modified = true
					}
					return true
				}
			}
			// build a new node app
			logger.SysInfo.Println("NodeApp Not Cover")
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
	logger.SysInfo.Printf("Application <%s> Not Cover. New...", app)
	scen.newversionbuf = append(scen.newversionbuf, &AppModel{
		Name:        app,
		Description: description,
	})
	scen.newversionbuf[len(scen.newversionbuf)-1].NodeApp = append(scen.newversionbuf[len(scen.newversionbuf)-1].NodeApp, &NodeAppModel{
		Name:    node,
		Version: version,
	})
	if modified {
		scen.modified = true
	}

	return true
}

func ResetScenario(scenario, version, message string) bool {
	ScenLock.Lock()
	defer ScenLock.Unlock()
	defer saveNoLock()

	var scen *ScenarioModel
	var ok bool

	if scen, ok = ScenarioMap[scenario]; !ok {
		return false
	}

	// find the history version index
	idx := -1
	for i := range scen.Versions {
		if version == scen.Versions[i].Version {
			idx = i
			break
		}
	}
	if idx < 0 {
		return false
	}

	// found. Is ok to only append the reference of this version's apps?
	scen.Versions = append(scen.Versions, &ScenarioVersionModel{
		BasicVersionModel: BasicVersionModel{
			Version:   version,
			Message:   message,
			Timestamp: time.Now().UnixMilli(),
		},
		Apps: scen.Versions[idx].Apps,
	})

	// update newversionbuf
	scen.newversionbuf = cloneLayer(scen.Versions[idx].Apps)
	scen.modified = false
	return true
}
