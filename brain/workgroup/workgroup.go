package workgroup

import (
	"errors"
	"fmt"
	"strings"

	"github.com/piaodazhu/Octopoda/brain/rdb"
	"github.com/piaodazhu/Octopoda/protocols"
)

func Info(path string) (*protocols.WorkgroupInfo, error) {
	// if not found : return nil, nil
	// if access error: return nil, err
	// else return : info, nil
	passwd, found, err := rdb.GetString(infoKey(path))
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &protocols.WorkgroupInfo{
		Path:     path,
		Password: passwd,
	}, nil
}

func createInfo(path string) error {
	ok, err := rdb.SetStringNX(infoKey(path), "")
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("created path is already exist")
	}
	return nil
}

func Grant(path string, password string) error {
	ok, err := rdb.SetStringXX(infoKey(path), password)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("granted path is not exist")
	}
	return nil
}

func Children(path string) (protocols.WorkgroupChildren, error) {
	// return array
	return rdb.GetSMembers(childrenKey(path))
}

func addChild(parent string, path string) error {
	return rdb.AddSMembers(childrenKey(parent), path)
}

func removeChild(parent string, path string) error {
	cnt, err := rdb.RemoveSMembers(childrenKey(parent), path)
	if err != nil {
		return err
	}
	if cnt == 0 { // not any child has been delete
		return nil
	}

	// BFS delete all children
	deleteQueue := []string{path}
	for len(deleteQueue) > 0 {
		tmp := []string{}
		for _, deletePath := range deleteQueue {
			children, err := Children(deletePath)
			if err != nil {
				fmt.Println("error in RemoveChild: " + err.Error())
			}
			rdb.DelKey(infoKey(deletePath))
			rdb.DelKey(membersKey(deletePath))
			rdb.DelKey(childrenKey(deletePath))
			tmp = append(tmp, children...)
		}
		deleteQueue = tmp
	}

	return nil
}

func Members(path string) (protocols.WorkgroupMembers, error) {
	// return array
	return rdb.GetSMembers(membersKey(path))
}

func AddMembers(parent, path string, names protocols.WorkgroupMembers) error {
	// if not found: set
	if info, err := Info(path); err != nil {
		return err
	} else if info == nil {
		// should create new path
		err = createInfo(path)
		if err != nil {
			return err
		}
		err = addChild(parent, path)
		if err != nil {
			return err
		}
	} 
	return rdb.AddSMembers(membersKey(path), names...)
}

func RemoveMembers(parent, path string, names protocols.WorkgroupMembers) error {
	// if names is empry: delete all
	// delete along the sub trees
	if len(names) == 0 {
		return removeChild(parent, path)
	}

	// 1. delete names from path
	deleteCnt, err := rdb.RemoveSMembers(membersKey(path), names...)
	if err != nil {
		return err 
	}
	if deleteCnt == 0 { // no need delete. finish
		return nil
	}

	// 2. check is path is empty
	cnt, err := rdb.CountSMembers(membersKey(path))
	if err != nil {
		return err 
	}
	if cnt == 0 {
		// path has no members. remove all of it
		return removeChild(parent, path)
	}
	
	// 3. remove names from all children
	children, err := Children(path)
	if err != nil {
		return err 
	}
	for _, child := range children { // recursively delete
		RemoveMembers(path, child, names)
	}

	return nil
}

func IsSubSet(subMembers, members []string) bool {
	return IsInScope(MakeScope(members), subMembers...) 
}

func MakeScope(members []string) map[string]string {
	scope := map[string]string{}
	for i := range members {
		scope[members[i]] = ""
	}
	return scope
}

func IsInScope(scope map[string]string, members ...string) bool {
	for i := range members {
		if _, found := scope[members[i]]; !found {
			return false
		}
	}
	return true
}

func IsSubPath(subPath, path string) bool {
	// prevent end with //
	if !strings.HasPrefix(subPath, path) {
		return false
	}
	if len(subPath) == len(path) {
		return false
	}
	return subPath[len(path)] == '/'
}

func IsSameOrSubPath(subPath, path string) bool {
	// prevent end with //
	if !strings.HasPrefix(subPath, path) {
		return false
	}
	if len(subPath) == len(path) {
		return true
	}
	return subPath[len(path)] == '/'
}

func IsDirectSubPath(subPath, path string) bool {
	// prevent end with //
	return false
}

func infoKey(key string) string {
	return "info:" + key
}

func childrenKey(key string) string {
	return "children:" + key
}

func membersKey(key string) string {
	return "members:" + key
}
