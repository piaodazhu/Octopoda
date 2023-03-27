package app

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"sync"
	"tentacle/config"
	"tentacle/logger"
)

type Version struct {
	Time int64
	Hash string
	Msg  string
}

type App struct {
	Name        string
	Discription string

	Versions   []Version
	VersionPtr int
}

type NodeApps struct {
	NodeVersion int64
	Apps        []App
}

const consistFileName = "nodeapps.json"

var nodeApps NodeApps
var nLock sync.Mutex

func InitAppModel() {
	var file strings.Builder
	file.WriteString(config.GlobalConfig.Workspace.Root)
	file.WriteString(consistFileName)
	f, err := os.Open(file.String())
	if err != nil {
		logger.Server.Print(err)
	} else {
		defer f.Close()
		content, _ := io.ReadAll(f)
		nLock.Lock()
		defer nLock.Unlock()
		if err := json.Unmarshal(content, &nodeApps); err != nil {
			logger.Server.Fatal("Invalid nodeapps file!")
		}
	}
}

func Save() {
	var file strings.Builder
	file.WriteString(config.GlobalConfig.Workspace.Root)
	file.WriteString(consistFileName)
	_, err := os.Stat(file.String())
	if os.IsExist(err) {
		os.Rename(file.String(), file.String()+".bk")
	}
	nodeApps.NodeVersion++
	nLock.Lock()
	defer nLock.Unlock()
	serialized, _ := json.Marshal(&nodeApps)
	err = os.WriteFile(file.String(), serialized, os.ModePerm)
	if err != nil {
		logger.Server.Print("cannot WriteFile")
	}
}

func ConvertHash(appname, shorthash string) (string, bool) {
	nLock.Lock()
	defer nLock.Unlock()
	for i := range nodeApps.Apps {
		if nodeApps.Apps[i].Name == appname {
			for j := range nodeApps.Apps[i].Versions {
				if nodeApps.Apps[i].Versions[j].Hash[:len(shorthash)] == shorthash {
					return nodeApps.Apps[i].Versions[j].Hash, true
				}
			}
			break
		}
	}
	return "", false
}

func Exists(appname, scenario string) bool {
	name := appname + "@" + scenario
	nLock.Lock()
	defer nLock.Unlock()
	for i := range nodeApps.Apps {
		if nodeApps.Apps[i].Name == name {
			return true
		}
	}
	return false
}


func Create(appname, scenario, description string) bool {
	name := appname + "@" + scenario
	nLock.Lock()
	defer nLock.Unlock()
	// create must be new
	for i := range nodeApps.Apps {
		if nodeApps.Apps[i].Name == name {
			return false
		}
	}
	nodeApps.Apps = append(nodeApps.Apps, App{
		Name:        name,
		Discription: description,
		Versions:    nil,
		VersionPtr:  0,
	})
	return true
}

func Update(appname string, version Version) bool {
	nLock.Lock()
	defer nLock.Unlock()
	// update must be existed
	for i := range nodeApps.Apps {
		if nodeApps.Apps[i].Name == appname {
			nodeApps.Apps[i].Versions = append(nodeApps.Apps[i].Versions, version)
			nodeApps.Apps[i].VersionPtr = len(nodeApps.Apps[i].Versions) - 1
			return true
		}
	}
	return false
}

func Reset(appname, versionhash string) bool {
	nLock.Lock()
	defer nLock.Unlock()
	// reset must be existed
	for i := range nodeApps.Apps {
		if nodeApps.Apps[i].Name == appname {
			for j := range nodeApps.Apps[i].Versions {
				if nodeApps.Apps[i].Versions[j].Hash[:len(versionhash)] == versionhash {
					nodeApps.Apps[i].VersionPtr = j
					return true
				}
			}
			break
		}
	}
	return false
}

type NodeAppsDigest struct {
	Apps []AppDigest
}
type AppDigest struct {
	Name        string
	Discription string
	CurVersion  Version
}

func Digest() []byte {
	digest := &NodeAppsDigest{}
	nLock.Lock()
	defer nLock.Unlock()
	for i := range nodeApps.Apps {
		digest.Apps = append(digest.Apps, AppDigest{
			Name: nodeApps.Apps[i].Name, 
			Discription: nodeApps.Apps[i].Discription,
			CurVersion: nodeApps.Apps[i].Versions[nodeApps.Apps[i].VersionPtr],
		})
	}
	serialized, _ := json.Marshal(digest)
	return serialized
}

func Versions(name string) []byte {
	nLock.Lock()
	defer nLock.Unlock()
	idx := -1
	serialized := []byte{}
	for i := range nodeApps.Apps {
		if nodeApps.Apps[i].Name == name {
			idx = i 
			break
		}
	}
	if idx > 0 {
		serialized, _ = json.Marshal(&nodeApps.Apps[idx]) 
	}
	return serialized
}

func Fix(appname string) {

}

func FixAll(appname string) {

}
