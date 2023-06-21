package api

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"brain/model"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type NodeInfoText struct {
	Name         string `json:"name"`
	Health       string `json:"health"`
	MsgConnState string `json:"msg_conn"`
	OnlineTime   string `json:"online_time,omitempty"`
	OfflineTime  string `json:"offline_time,omitempty"`
	LastOnline   string `json:"last_active,omitempty"`
}

type NodesInfoText struct {
	NodeInfoList []*NodeInfoText `json:"nodes"`
	Total        int             `json:"total"`
	Active       int             `json:"active"`
	Offline      int             `json:"offline"`
}

func nodeInfoToText(node *model.NodeModel) *NodeInfoText {
	res := &NodeInfoText{
		Name: node.Name,
	}
	switch node.State {
	case 0:
		res.Health = "Healthy"
		res.OnlineTime = time.Since(time.UnixMilli(node.OnlineTs)).String()
	case 1:
		res.Health = "Disconnect"
		res.LastOnline = time.UnixMilli(node.ActiveTs).Format("2006-01-02 15:04:05")
	case 2:
		res.Health = "Offline"
		res.OfflineTime = time.Since(time.UnixMilli(node.OfflineTs)).String()
	}
	if node.MsgConn == nil {
		res.MsgConnState = "Off"
	} else {
		res.MsgConnState = "On"
	}
	return res
}

func nodesInfoToText(nodes []*model.NodeModel) *NodesInfoText {
	res := &NodesInfoText{
		NodeInfoList: make([]*NodeInfoText, len(nodes)),
	}
	for i, node := range nodes {
		res.Total++
		if node.State == 0 {
			res.Active++
		} else if node.State == 2 {
			res.Offline++
		}
		res.NodeInfoList[i] = nodeInfoToText(node)
	}
	return res
}

func NodeInfo(ctx *gin.Context) {
	var name string
	var ok bool
	var node *model.NodeModel
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	if node, ok = model.GetNodeInfoByName(name); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	fmt.Println(nodeInfoToText(node))
	ctx.JSON(200, nodeInfoToText(node))
}

func NodeState(ctx *gin.Context) {
	var name string
	var ok bool
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	raw, err := model.Request(name, message.TypeNodeStatus, []byte{})
	if err != nil {
		logger.Comm.Println("NodeState", err)
		ctx.JSON(404, struct{}{})
		return
	}
	ctx.Data(200, "application/json", raw)
}

func NodesInfo(ctx *gin.Context) {
	var nodes []*model.NodeModel
	var ok bool

	if nodes, ok = model.GetNodesInfoAll(); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	ctx.JSON(200, nodesInfoToText(nodes))
}

func NodesState(ctx *gin.Context) {
	var nodes []*model.NodeModel
	var states []model.State
	var ok bool

	if nodes, ok = model.GetNodesInfoAll(); !ok {
		ctx.JSON(404, struct{}{})
	}

	channel := make(chan model.State, len(nodes))
	var wg sync.WaitGroup
	wg.Add(len(nodes))
	for _, node := range nodes {
		go getNodeState(node.Name, channel, &wg)
	}
	wg.Wait()
	close(channel)
	for v := range channel {
		states = append(states, v)
	}
	ctx.JSON(200, states)
}

func getNodeState(name string, channel chan<- model.State, wg *sync.WaitGroup) {
	defer wg.Done()
	var state model.State = model.State{Name: name}
	var err error
	var raw []byte

	raw, err = model.Request(name, message.TypeNodeStatus, []byte{})
	if err != nil {
		logger.Comm.Println("getNodeState", err)
		goto sendres
	}

	err = config.Jsoner.Unmarshal(raw, &state)
	if err != nil {
		logger.Exceptions.Println("UnmarshalNodeState", err)
		goto sendres
	}
sendres:
	channel <- state
}

func NodePrune(ctx *gin.Context) {
	model.PruneDeadNode()
	ctx.Status(200)
}
