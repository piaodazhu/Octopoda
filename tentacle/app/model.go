package app

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"sync"
	"tentacle/config"
	"tentacle/logger"
	"time"
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
		// logger.Server.Panic(err)
		nLock.Lock()
		nodeApps.NodeVersion = 0
		saveNoLock()
		nLock.Unlock()
	} else {
		defer f.Close()
		content, _ := io.ReadAll(f)
		nLock.Lock()
		if err := json.Unmarshal(content, &nodeApps); err != nil {
			logger.Server.Fatal("Invalid nodeapps file!")
		}
		nLock.Unlock()
	}
	go autoFixAll()
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

func Delete(appname, scenario string) bool {
	name := appname + "@" + scenario
	nLock.Lock()
	defer nLock.Unlock()
	// find the target
	idx := -1
	for i := range nodeApps.Apps {
		if nodeApps.Apps[i].Name == name {
			idx = i
			break
		}
	}
	if idx >= 0 {
		nodeApps.Apps[idx], nodeApps.Apps[len(nodeApps.Apps) - 1] = nodeApps.Apps[len(nodeApps.Apps) - 1], nodeApps.Apps[idx]
		nodeApps.Apps = nodeApps.Apps[:len(nodeApps.Apps) - 1]
	}
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
			Name:        nodeApps.Apps[i].Name,
			Discription: nodeApps.Apps[i].Discription,
			CurVersion:  nodeApps.Apps[i].Versions[nodeApps.Apps[i].VersionPtr],
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

func CurVersion(name string) Version {
	nLock.Lock()
	defer nLock.Unlock()
	for i := range nodeApps.Apps {
		if nodeApps.Apps[i].Name == name {
			return nodeApps.Apps[i].Versions[nodeApps.Apps[i].VersionPtr]
		}
	}
	return Version{}
}

// --------------------------------------
// only for experimental
func Fix(appname string) bool {
	// versions is in time desc order
	versions, err := gitLogs(appname, 5)
	if err != nil || len(versions) == 0 {
		return false
	}

	nLock.Lock()
	defer nLock.Unlock()
	idx := -1
	for i := range nodeApps.Apps {
		if nodeApps.Apps[i].Name == appname {
			idx = i
			break
		}
	}

	if idx < 0 {
		return false
	}

	oldlist := nodeApps.Apps[idx].Versions
	// Case 1: no missing commit
	if oldlist[nodeApps.Apps[idx].VersionPtr].Hash == versions[len(versions)-1].Hash {
		// need not fix
		return true
	}

	// Case 2: has a missing reset
	for i := range oldlist {
		if oldlist[i].Hash == versions[len(versions)-1].Hash {
			nodeApps.Apps[idx].VersionPtr = i
			saveNoLock()
			return true
		}
	}

	// Case 3: has many missing commits. We can compensate at most N missing commits.
	// versions: Actual latest N commits
	// oldlist: Octopoda currently record
	// if a version not exists in oldlist, we add this version to Octopoda app versions list
	for i := 0; i < len(versions); i++ {
		exists := false
		for j := range oldlist {
			if oldlist[j].Hash == versions[i].Hash {
				exists = true
				break
			}
		}
		if !exists {
			// all version newer than this verison is considered NEW, so add all of them
			nodeApps.Apps[idx].Versions = append(nodeApps.Apps[idx].Versions, versions[i:]...)
			nodeApps.Apps[idx].VersionPtr = len(nodeApps.Apps[idx].Versions) - 1
			saveNoLock()
			break
		}
	}
	return true
}

// not optimized. bad loop
func FixAll() {
	for i := range nodeApps.Apps {
		Fix(nodeApps.Apps[i].Name)
	}
}

func autoFixAll() {
	for {
		time.Sleep(time.Second * 3)
		FixAll()
	}
}

func saveNoLock() {
	var file strings.Builder
	file.WriteString(config.GlobalConfig.Workspace.Root)
	file.WriteString(consistFileName)
	_, err := os.Stat(file.String())
	if os.IsExist(err) {
		os.Rename(file.String(), file.String()+".bk")
	}
	nodeApps.NodeVersion++
	serialized, _ := json.Marshal(&nodeApps)
	err = os.WriteFile(file.String(), serialized, os.ModePerm)
	if err != nil {
		logger.Server.Print("cannot WriteFile")
	}
}
