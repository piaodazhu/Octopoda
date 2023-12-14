package model

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/protocols/san"

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
	Versions      []*san.ScenarioVersionModel
	newversionbuf []*san.AppModel
	modified      bool
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
			Versions:      []*san.ScenarioVersionModel{},
			newversionbuf: []*san.AppModel{},
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
			scen.Versions = append(scen.Versions, &san.ScenarioVersionModel{
				Version: san.Version{
					Hash: hex.EncodeToString(versionhash[:]),
					Msg:  message,
					Time: time.Now().UnixMilli(),
				},
				Apps: scen.newversionbuf,
			})
			// triger save?
			scen.newversionbuf = cloneLayer(scen.newversionbuf) // must deep copy!
			scen.modified = false
		}

		return hasModified, true
	}
}

func cloneLayer(prototype []*san.AppModel) []*san.AppModel {
	res := []*san.AppModel{}
	for _, app := range prototype {
		nodeapps := []*san.NodeAppModel{}
		for _, nodeapp := range app.NodeApp {
			nodeapps = append(nodeapps, &san.NodeAppModel{
				Name:    nodeapp.Name,
				Version: nodeapp.Version,
			})
		}
		res = append(res, &san.AppModel{
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

func GetNodeApps(name string, version string) []san.NodeAppItem {
	ScenLock.RLock()
	defer ScenLock.RUnlock()

	list := []san.NodeAppItem{}
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
			if version == v.Hash {
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
			list = append(list, san.NodeAppItem{
				AppName:  app.Name,
				ScenName: scen.Name,
				NodeName: node.Name,
				Version:  node.Version,
			})
		}

	}
	return list
}

func GetScenariosDigestAll() ([]san.ScenarioDigest, bool) {
	ScenLock.RLock()
	defer ScenLock.RUnlock()

	if len(ScenarioMap) == 0 {
		return []san.ScenarioDigest{}, true
	}
	res := make([]san.ScenarioDigest, len(ScenarioMap))
	idx := 0
	for _, val := range ScenarioMap {
		res[idx].Name = val.Name
		res[idx].Description = val.Description
		if len(val.Versions) != 0 {
			res[idx].Timestamp = val.Versions[len(val.Versions)-1].Time
			res[idx].Version = val.Versions[len(val.Versions)-1].Hash
			res[idx].Message = val.Versions[len(val.Versions)-1].Msg
		}
		idx++
	}
	return res, true
}

func GetScenarioInfoByName(name string) (*san.ScenarioInfo, bool) {
	ScenLock.RLock()
	defer ScenLock.RUnlock()

	if scen, found := ScenarioMap[name]; found {
		// init but empty
		if len(scen.Versions) == 0 {
			logger.Exceptions.Println("GetScenarioInfoByName Empty Scenario")
			return &san.ScenarioInfo{
				ScenarioDigest: san.ScenarioDigest{
					Name:        scen.Name,
					Description: scen.Description,
				}}, true
		}
		res := &san.ScenarioInfo{
			ScenarioDigest: san.ScenarioDigest{
				Name:        scen.Name,
				Description: scen.Description,
				Version:     scen.Versions[len(scen.Versions)-1].Hash,
				Message:     scen.Versions[len(scen.Versions)-1].Msg,
				Timestamp:   scen.Versions[len(scen.Versions)-1].Time,
			},
			Apps: []*san.AppInfo{},
		}
		for _, v := range scen.Versions[len(scen.Versions)-1].Apps {
			app := &san.AppInfo{
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

func GetScenarioVersionByName(name string) []san.Version {
	ScenLock.RLock()
	defer ScenLock.RUnlock()

	versions := []san.Version{}
	if scen, found := ScenarioMap[name]; found {
		for _, v := range scen.Versions {
			versions = append(versions, v.Version)
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
			// logger.SysInfo.Println("NodeApp Not Cover")
			application.NodeApp = append(application.NodeApp, &san.NodeAppModel{
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
	// logger.SysInfo.Printf("Application <%s> Not Cover. New...", app)
	scen.newversionbuf = append(scen.newversionbuf, &san.AppModel{
		Name:        app,
		Description: description,
	})
	scen.newversionbuf[len(scen.newversionbuf)-1].NodeApp = append(scen.newversionbuf[len(scen.newversionbuf)-1].NodeApp, &san.NodeAppModel{
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
		if version == scen.Versions[i].Hash {
			idx = i
			break
		}
	}
	if idx < 0 {
		return false
	}

	// found. Is ok to only append the reference of this version's apps?
	scen.Versions = append(scen.Versions, &san.ScenarioVersionModel{
		Version: san.Version{
			Hash: version,
			Msg:  message,
			Time: time.Now().UnixMilli(),
		},
		Apps: scen.Versions[idx].Apps,
	})

	// update newversionbuf
	scen.newversionbuf = cloneLayer(scen.Versions[idx].Apps)
	scen.modified = false
	return true
}
