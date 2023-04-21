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

	Versions []Version
}

type NodeApps struct {
	NodeVersion int64
	Apps        []App
}

const diskFileName = "nodeapps.json"

var nodeApps NodeApps
var nLock sync.Mutex

func InitAppModel() {
	var file strings.Builder
	file.WriteString(config.GlobalConfig.Workspace.Root)
	file.WriteString(diskFileName)
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
			logger.Exceptions.Fatal("Invalid nodeapps file!")
		}
		nLock.Unlock()
	}
	go autoFixAll()
}

func Save() {
	var file strings.Builder
	file.WriteString(config.GlobalConfig.Workspace.Root)
	file.WriteString(diskFileName)
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
		logger.Exceptions.Print("cannot WriteFile")
	}
}

func ConvertHash(appname, scenario, shorthash string) (string, bool) {
	name := appname + "@" + scenario
	nLock.Lock()
	defer nLock.Unlock()
	for i := range nodeApps.Apps {
		if nodeApps.Apps[i].Name == name {
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
		nodeApps.Apps[idx], nodeApps.Apps[len(nodeApps.Apps)-1] = nodeApps.Apps[len(nodeApps.Apps)-1], nodeApps.Apps[idx]
		nodeApps.Apps = nodeApps.Apps[:len(nodeApps.Apps)-1]
	}
	return true
}

func Update(appname, scenario string, version Version) bool {
	name := appname + "@" + scenario
	nLock.Lock()
	defer nLock.Unlock()
	// update must be existed
	for i := range nodeApps.Apps {
		if nodeApps.Apps[i].Name == name {
			nodeApps.Apps[i].Versions = append(nodeApps.Apps[i].Versions, version)
			// logger.Server.Println("update version", version.Hash[:4], ", then:\n", nodeApps.Apps[i].Versions)
			return true
		}
	}
	return false
}

func Reset(appname, scenario, versionhash, message string) bool {
	name := appname + "@" + scenario
	nLock.Lock()
	defer nLock.Unlock()
	// reset must be existed
	for i := range nodeApps.Apps {
		if nodeApps.Apps[i].Name == name {
			for j := range nodeApps.Apps[i].Versions {
				if nodeApps.Apps[i].Versions[j].Hash[:len(versionhash)] == versionhash {
					// nodeApps.Apps[i].VersionPtr = j
					nodeApps.Apps[i].Versions = append(nodeApps.Apps[i].Versions, Version{
						Time: time.Now().Unix(),
						Hash: nodeApps.Apps[i].Versions[j].Hash,
						Msg:  message,
					})
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
			CurVersion:  nodeApps.Apps[i].Versions[len(nodeApps.Apps[i].Versions)-1],
		})
	}
	serialized, _ := json.Marshal(digest)
	return serialized
}

func Versions(appname, scenario string) []byte {
	name := appname + "@" + scenario
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

	if idx >= 0 {
		serialized, _ = json.Marshal(&nodeApps.Apps[idx])
	}
	return serialized
}

func CurVersion(appname, scenario string) Version {
	name := appname + "@" + scenario
	nLock.Lock()
	defer nLock.Unlock()
	for i := range nodeApps.Apps {
		if nodeApps.Apps[i].Name == name {
			if len(nodeApps.Apps[i].Versions) == 0 {
				logger.Exceptions.Println("CurVersion is NULL")
				return Version{}
			}
			return nodeApps.Apps[i].Versions[len(nodeApps.Apps[i].Versions)-1]
		}
	}
	return Version{}
}

// --------------------------------------
// only for experimental
func Fix(fullname string) bool {
	// versions is in time desc order
	versions, err := gitLogs(fullname, 5)
	if err != nil || len(versions) == 0 {
		return false
	}
	// fmt.Println("\nFix Start--")
	// fmt.Println("Disks:")
	// fmt.Println(versions)

	nLock.Lock()
	defer nLock.Unlock()
	idx := -1
	for i := range nodeApps.Apps {
		if nodeApps.Apps[i].Name == fullname {
			idx = i
			break
		}
	}

	if idx < 0 {
		return false
	}

	oldlist := nodeApps.Apps[idx].Versions

	// fmt.Println("Mems:")
	// fmt.Println(oldlist)

	// fmt.Println(">>case1")
	// Case 1: no missing commit
	if oldlist[len(oldlist)-1].Hash == versions[len(versions)-1].Hash {
		// need not fix
		return true
	}

	// fmt.Println(">>case2")
	// Case 2: has a missing reset
	for i := range oldlist {
		if oldlist[i].Hash == versions[len(versions)-1].Hash {
			nodeApps.Apps[idx].Versions = append(nodeApps.Apps[idx].Versions, versions[len(versions)-1])
			saveNoLock()
			return true
		}
	}

	// fmt.Println(">>case3")
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
	file.WriteString(diskFileName)
	_, err := os.Stat(file.String())
	if os.IsExist(err) {
		os.Rename(file.String(), file.String()+".bk")
	}
	nodeApps.NodeVersion++
	serialized, _ := json.Marshal(&nodeApps)
	err = os.WriteFile(file.String(), serialized, os.ModePerm)
	if err != nil {
		logger.Exceptions.Print("cannot WriteFile")
	}
}
