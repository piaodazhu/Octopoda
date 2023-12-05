package workgroup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

type workgroupClient struct {
	rootPath    string
	password    string
	currentPath string

	client *http.Client
}

func pathClean(path string) string {
	path = filepath.Clean(path)
	return strings.ReplaceAll(filepath.Clean(path), "\\", "/")
}

func newWorkgroupClient(rootGroup, passwd, currentGroup string, client *http.Client) workgroupClient {
	rootGroup = pathClean(rootGroup)
	if rootGroup == "." {
		rootGroup = ""
	}
	if !strings.HasPrefix(rootGroup, "/") {
		rootGroup = "/" + rootGroup
	}
	return workgroupClient{
		rootPath:    strings.TrimSuffix(rootGroup, "/"),
		password:    passwd,
		currentPath: strings.TrimSuffix(currentGroup, "/"),
		client:      client,
	}
}

func (wg *workgroupClient) auth() error {
	if _, err := wg.remoteGetInfo(wg.rootPath); err != nil {
		return err
	}
	return nil
}

func (wg *workgroupClient) grant(groupPath string, password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password must longer than 6")
	}
	groupPath = wg.fixPath(groupPath)
	if !wg.isSubPath(groupPath, wg.rootPath) {
		return fmt.Errorf("root workgroup is %s, cannot access %s", wg.rootPath, groupPath)
	}
	if err := wg.remoteSetGrant(groupPath, password); err != nil {
		return err
	}
	return nil
}

func (wg *workgroupClient) cd(groupPath string) error {
	if groupPath == "~" {
		groupPath = wg.rootPath
	} else {
		groupPath = wg.fixPath(groupPath)
		if !wg.isSubPath(groupPath, wg.rootPath) {
			return fmt.Errorf("root workgroup is %s, cannot access %s", wg.rootPath, groupPath)
		}
	}
	if _, err := wg.remoteGetInfo(groupPath); err != nil {
		return err
	}
	wg.currentPath = groupPath
	return nil
}

func (wg *workgroupClient) ls(groupPath string) ([]string, error) {
	groupPath = wg.fixPath(groupPath)
	if !wg.isSubPath(groupPath, wg.rootPath) {
		return nil, fmt.Errorf("root workgroup is %s, cannot access %s", wg.rootPath, groupPath)
	}
	if children, err := wg.remoteGetChildren(groupPath); err != nil {
		return nil, err
	} else {
		res := make([]string, len(children))
		for i := range children {
			res[i] = strings.TrimPrefix(children[i], groupPath+"/")
		}
		return res, nil
	}
}

func (wg *workgroupClient) pwd() string {
	if wg.currentPath == "" {
		return "/"
	}
	return wg.currentPath
}

func (wg *workgroupClient) root() string {
	if wg.rootPath == "" {
		return "/"
	}
	return wg.rootPath
}

func (wg *workgroupClient) valid() bool {
	return wg.isSubPath(wg.currentPath, wg.rootPath)
}

func (wg *workgroupClient) toRoot() {
	wg.currentPath = wg.rootPath
}

func (wg *workgroupClient) get(groupPath string) ([]string, error) {
	groupPath = wg.fixPath(groupPath)
	if !wg.isSubPath(groupPath, wg.rootPath) {
		return nil, fmt.Errorf("root workgroup is %s, cannot access %s", wg.rootPath, groupPath)
	}
	if members, err := wg.remoteGetMembers(groupPath); err != nil {
		return nil, err
	} else {
		return members, nil
	}
}

func (wg *workgroupClient) addMembers(groupPath string, names ...string) error {
	if names == nil {
		return nil
	}
	groupPath = wg.fixPath(groupPath)
	if !wg.isSubPath(groupPath, wg.rootPath) {
		return fmt.Errorf("root workgroup is %s, cannot access %s", wg.rootPath, groupPath)
	}

	if err := wg.remoteAddMember(groupPath, names); err != nil {
		return err
	}
	return nil
}

func (wg *workgroupClient) removeMembers(groupPath string, names ...string) error {
	// if names == nil, remove this group and its children
	groupPath = wg.fixPath(groupPath)
	if !wg.isSubPath(groupPath, wg.rootPath) {
		return fmt.Errorf("root workgroup is %s, cannot access %s", wg.rootPath, groupPath)
	}
	if err := wg.remoteRemoveMember(groupPath, names); err != nil {
		return err
	}
	return nil
}

func (wg *workgroupClient) fixPath(path string) string {
	var fixed string
	if strings.HasPrefix(path, "/") {
		fixed = pathClean(path)
	} else {
		fixed = pathClean(wg.currentPath + "/" + path)
	}

	if fixed == "/" {
		return ""
	}
	return fixed
}

func (wg *workgroupClient) isSubPath(subPath, path string) bool {
	if !strings.HasPrefix(subPath, path) {
		return false
	}
	if len(subPath) == len(path) {
		return true
	}
	return subPath[len(path)] == '/'
}

func (wg *workgroupClient) setHeader(req *http.Request) {
	req.Header.Set("rootpath", wg.rootPath)
	req.Header.Set("password", wg.password)
	req.Header.Set("currentpath", wg.currentPath)
}

func (wg *workgroupClient) remoteGetInfo(path string) (*protocols.WorkgroupInfo, error) {
	info := protocols.WorkgroupInfo{}
	url := fmt.Sprintf("https://%s/%s%s?path=%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_WorkgroupInfo,
		path,
	)
	req, _ := http.NewRequest("GET", url, nil)
	wg.setHeader(req)
	res, err := wg.client.Do(req)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode == http.StatusOK {
		err := json.Unmarshal(raw, &info)
		if err != nil {
			output.PrintFatalln(err)
			return nil, err
		}
		return &info, nil
	}
	return nil, fmt.Errorf("status code: %d", res.StatusCode)
}

func (wg *workgroupClient) remoteSetGrant(path, password string) error {
	info := protocols.WorkgroupInfo{
		Path:     path,
		Password: password,
	}
	body, _ := json.Marshal(info)
	url := fmt.Sprintf("https://%s/%s%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_WorkgroupInfo,
	)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	wg.setHeader(req)
	res, err := wg.client.Do(req)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return nil
	}
	raw, _ := io.ReadAll(res.Body)
	return fmt.Errorf("status code: %d. msg: %s", res.StatusCode, string(raw))
}

func (wg *workgroupClient) remoteGetChildren(path string) (protocols.WorkgroupChildren, error) {
	children := protocols.WorkgroupChildren{}
	url := fmt.Sprintf("https://%s/%s%s?path=%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_WorkgroupChildren,
		path,
	)
	req, _ := http.NewRequest("GET", url, nil)
	wg.setHeader(req)
	res, err := wg.client.Do(req)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode == http.StatusOK {
		err := json.Unmarshal(raw, &children)
		if err != nil {
			output.PrintFatalln(err)
			return nil, err
		}
		return children, nil
	}
	return nil, fmt.Errorf("status code: %d. msg: %s", res.StatusCode, string(raw))
}

func (wg *workgroupClient) remoteGetMembers(path string) (protocols.WorkgroupMembers, error) {
	members := protocols.WorkgroupMembers{}
	url := fmt.Sprintf("https://%s/%s%s?path=%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_WorkgroupMembers,
		path,
	)
	req, _ := http.NewRequest("GET", url, nil)
	wg.setHeader(req)
	res, err := wg.client.Do(req)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode == http.StatusOK {
		err := json.Unmarshal(raw, &members)
		if err != nil {
			output.PrintFatalln(err)
			return nil, err
		}
		return members, nil
	}
	return nil, fmt.Errorf("status code: %d. msg: %s", res.StatusCode, string(raw))
}

func (wg *workgroupClient) remoteAddMember(path string, names []string) error {
	var err error
	if names, err = NodesParse(names); err != nil {
		return err
	}
	return wg.remoteOperateMember(path, true, names)
}

func (wg *workgroupClient) remoteRemoveMember(path string, names []string) error {
	var err error
	if names, err = NodesParseNoCheck(names); err != nil {
		return err
	}
	return wg.remoteOperateMember(path, false, names)
}

func (wg *workgroupClient) remoteOperateMember(path string, isAdd bool, names []string) error {
	params := protocols.WorkgroupMembersPostParams{
		Path:    path,
		IsAdd:   isAdd,
		Members: names,
	}
	body, _ := json.Marshal(params)
	url := fmt.Sprintf("https://%s/%s%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_WorkgroupMembers,
	)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	wg.setHeader(req)
	res, err := wg.client.Do(req)
	if err != nil {
		emsg := "http get error: " + err.Error()
		output.PrintFatalln(emsg)
		return errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return nil
	}
	raw, _ := io.ReadAll(res.Body)
	return fmt.Errorf("status code: %d. msg: %s", res.StatusCode, string(raw))
}
