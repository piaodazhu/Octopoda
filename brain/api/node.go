package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/brain/network"
	"github.com/piaodazhu/Octopoda/brain/rdb"
	"github.com/piaodazhu/Octopoda/brain/sys"
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
	ctx.JSON(200, node)
}

func NodesInfo(ctx *gin.Context) {
	targetNodes := ctx.Query("targetNodes")
	names := []string{}
	err := config.Jsoner.Unmarshal([]byte(targetNodes), &names)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	octlFaceIp, _ := network.GetOctlFaceIp()
	nodes := protocols.NodesInfo{
		BrainName:    config.GlobalConfig.Name,
		BrainVersion: buildinfo.String(),
		BrainAddr:    octlFaceIp,
	}
	var ok bool

	if len(names) == 0 {
		if nodes.InfoList, ok = model.GetNodesInfoAll(); !ok {
			ctx.JSON(http.StatusNotFound, struct{}{})
			return
		}
	} else {
		if nodes.InfoList, ok = model.GetNodesInfo(names); !ok {
			ctx.JSON(http.StatusNotFound, struct{}{})
			return
		}
	}

	ctx.JSON(200, nodes)
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
		ctx.JSON(200, sys.LocalStatus())
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
	ctx.JSON(200, status)
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

	if len(names) == 0 {
		if nodes, ok = model.GetNodesInfoAll(); !ok {
			ctx.JSON(http.StatusNotFound, struct{}{})
		}
	} else {
		if nodes, ok = model.GetNodesInfo(names); !ok {
			ctx.JSON(http.StatusNotFound, struct{}{})
			return
		}
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

	ctx.JSON(200, nodesStatus)
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
	model.PruneDeadNode()
	ctx.Status(200)
}

// NodesParse will: 1 parse group to nodes 2 check node health 3 remove duplicate nodes
func NodesParse(ctx *gin.Context) {
	names := []string{}
	if err := ctx.Bind(&names); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	nodeSet := map[string]struct{}{}
	for _, name := range names {
		if len(name) == 0 {
			continue
		}
		if _, ok := nodeSet[name]; ok {
			continue
		}

		if name[0] == '@' {
			res, ok := rdb.GroupGet(name[1:])
			if !ok || len(res) == 0 {
				ctx.String(http.StatusBadRequest, "invalid group %s", name[1:])
				return
			}
			for _, node := range res {
				nodeSet[node] = struct{}{}
			}
		} else {
			res, ok := model.GetNodeState(name)
			if !ok || res != protocols.NodeStateReady {
				ctx.String(http.StatusBadRequest, "invalid node %s", name)
				return
			}
			nodeSet[name] = struct{}{}
		}
	}

	nodes := []string{}
	for node := range nodeSet {
		nodes = append(nodes, node)
	}
	ctx.JSON(200, nodes)
}
