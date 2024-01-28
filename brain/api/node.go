package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/brain/network"
	"github.com/piaodazhu/Octopoda/brain/sys"
	"github.com/piaodazhu/Octopoda/brain/workgroup"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/buildinfo"
)

func NodeInfo(ctx *gin.Context) {
	var name string
	var ok bool
	var node protocols.NodeInfo
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(http.StatusNotFound, struct{}{})
		return
	}
	if node, ok = model.GetNodeInfoByName(name); !ok {
		ctx.JSON(http.StatusNotFound, struct{}{})
		return
	}
	node.BrainTs = time.Now().Unix()
	ctx.JSON(http.StatusOK, node)
}

func NodesInfo(ctx *gin.Context) {
	targetNodes := ctx.Query("targetNodes")
	names := []string{}
	err := config.Jsoner.Unmarshal([]byte(targetNodes), &names)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	fmt.Println(names)
	if !workgroup.IsInScope(ctx.GetStringMapString("octopoda_scope"), names...) {
		emsg := "some nodes are invalid or out of scope"
		ctx.JSON(http.StatusBadRequest, errors.New(emsg))
		return
	}

	octlFaceIp, _ := network.GetOctlFaceIp()
	nodes := protocols.NodesInfo{
		BrainName:    config.GlobalConfig.Name,
		BrainVersion: buildinfo.String(),
		BrainAddr:    octlFaceIp,
	}
	var ok bool

	if nodes.InfoList, ok = model.GetNodesInfo(names); !ok {
		ctx.JSON(http.StatusOK, struct{}{})
		return
	}

	currTs := time.Now().Unix()
	for i := range nodes.InfoList {
		nodes.InfoList[i].BrainTs = currTs
	}

	ctx.JSON(http.StatusOK, nodes)
}

func NodeStatus(ctx *gin.Context) {
	var name string
	var ok bool
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(http.StatusNotFound, struct{}{})
		return
	}
	var status protocols.Status
	if name == "brain" {
		ctx.JSON(http.StatusOK, sys.LocalStatus())
		return
	}
	if state, ok := model.GetNodeState(name); !ok || state != protocols.NodeStateReady {
		ctx.JSON(http.StatusNotFound, struct{}{})
		return
	}
	raw, err := model.Request(name, protocols.TypeNodeStatus, []byte{})
	if err != nil {
		logger.Comm.Println("NodeStatus", err)
		ctx.JSON(http.StatusNotFound, struct{}{})
		return
	}

	err = json.Unmarshal(raw, &status)
	if err != nil {
		logger.Comm.Println("NodeStatus Unmarshal", err)
		ctx.JSON(http.StatusNotFound, struct{}{})
		return
	}
	ctx.JSON(http.StatusOK, status)
}

func NodesState(ctx *gin.Context) {
	targetNodes := ctx.Query("targetNodes")
	var nodes, aliveNodes []protocols.NodeInfo
	var nodesStatus protocols.NodesStatus
	var ok bool

	names := []string{}
	err := config.Jsoner.Unmarshal([]byte(targetNodes), &names)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if !workgroup.IsInScope(ctx.GetStringMapString("octopoda_scope"), names...) {
		emsg := "ERROR: some nodes are invalid or out of scope."
		ctx.String(http.StatusBadRequest, emsg)
		return
	}
	if nodes, ok = model.GetNodesInfo(names); !ok {
		ctx.JSON(http.StatusNotFound, struct{}{})
		return
	}

	for _, n := range nodes {
		if n.State == protocols.NodeStateReady {
			aliveNodes = append(aliveNodes, n)
		}
	}

	channel := make(chan protocols.Status, len(aliveNodes))
	var wg sync.WaitGroup
	wg.Add(len(aliveNodes))
	for _, node := range aliveNodes {
		go getNodeStatus(node.Name, channel, &wg)
	}
	wg.Wait()
	close(channel)

	for v := range channel {
		nodesStatus.StatusList = append(nodesStatus.StatusList, v)
	}

	ctx.JSON(http.StatusOK, nodesStatus)
}

func getNodeStatus(name string, channel chan<- protocols.Status, wg *sync.WaitGroup) {
	defer wg.Done()
	var state protocols.Status
	var err error
	var raw []byte

	raw, err = model.Request(name, protocols.TypeNodeStatus, []byte{})
	if err != nil {
		logger.Comm.Println("getNodeStatus", err)
		goto sendres
	}

	err = config.Jsoner.Unmarshal(raw, &state)
	if err != nil {
		logger.Exceptions.Println("UnmarshalNodeStatus", err)
		goto sendres
	}
sendres:
	channel <- state
}

func NodePrune(ctx *gin.Context) {
	targetNodes := ctx.Query("targetNodes")
	names := []string{}
	err := config.Jsoner.Unmarshal([]byte(targetNodes), &names)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if !workgroup.IsInScope(ctx.GetStringMapString("octopoda_scope"), names...) {
		emsg := "some nodes are invalid or out of scope"
		ctx.JSON(http.StatusBadRequest, errors.New(emsg))
		return
	}

	pruned := model.PruneDeadNode(names)
	if len(pruned) > 0 {
		err = workgroup.RemoveMembers("", "", pruned...)
		if err != nil {
			logger.Exceptions.Printf("remove members when prune error: " + err.Error())
		}
	}
	ctx.Status(200)
}

// NodesParse will: 1 parse group to nodes 2 check node health 3 remove duplicate nodes
func NodesParse(ctx *gin.Context) {
	names := []string{}
	if err := ctx.Bind(&names); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if len(names) == 0 {
		names = append(names, "@")
	}

	nodeSet := map[string]struct{}{}
	invalidNameSet := map[string]struct{}{}
	unhealthyNodeSet := map[string]struct{}{}

	nodeNames := []string{}

	// first loop: only process group names
	currentPath := ctx.GetHeader("currentpath")
	currentPath = strings.TrimSuffix(currentPath, "/")
	for _, name := range names {
		if len(name) == 0 {
			continue
		}
		if name[0] == '@' { // group name
			groupName := name[1:]
			groupPath := fmt.Sprintf("%s/%s", currentPath, groupName)
			if groupName == "" || groupName == "ALL" || groupName == "all" {
				groupPath = currentPath
			}

			members, err := workgroup.Members(groupPath)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
				return
			}
			if len(members) == 0 {
				invalidNameSet[name] = struct{}{}
				continue
			}
			nodeNames = append(nodeNames, members...)
		} else {
			nodeNames = append(nodeNames, name)
		}
	}

	// second loop: get rid of duplicate
	scope := ctx.GetStringMapString("octopoda_scope")
	for _, name := range nodeNames {
		if len(name) == 0 {
			continue
		}
		if _, ok := nodeSet[name]; ok {
			continue
		}

		if !workgroup.IsInScope(scope, name) {
			invalidNameSet[name] = struct{}{}
			continue
		}

		state, ok := model.GetNodeState(name)
		if !ok {
			invalidNameSet[name] = struct{}{}
			continue
		}

		if state != protocols.NodeStateReady {
			unhealthyNodeSet[name] = struct{}{}
			continue
		}
		nodeSet[name] = struct{}{}
	}

	res := protocols.NodeParseResult{}
	for node := range nodeSet {
		res.OutputNames = append(res.OutputNames, node)
	}
	for node := range unhealthyNodeSet {
		res.UnhealthyNodes = append(res.UnhealthyNodes, node)
	}
	for node := range invalidNameSet {
		res.InvalidNames = append(res.InvalidNames, node)
	}
	ctx.JSON(http.StatusOK, res)
}
