package api

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"brain/model"
	"net"
	"sync"

	"github.com/gin-gonic/gin"
)

func NodeInfo(ctx *gin.Context) {
	var name string
	var ok bool
	var node *model.NodeInfo
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	if node, ok = model.GetNodeInfoByName(name); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	ctx.JSON(200, node)
}

func NodeState(ctx *gin.Context) {
	var name string
	var conn *net.Conn
	var ok bool
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	if conn, ok = model.GetNodeMsgConn(name); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	err := message.SendMessage(*conn, message.TypeNodeStatus, []byte{})
	if err != nil {
		logger.Comm.Println("NodeState", err)
		ctx.JSON(404, struct{}{})
		return
	}

	mtype, raw, err := message.RecvMessage(*conn)
	if err != nil || mtype != message.TypeNodeStatusResponse {
		logger.Comm.Println("NodeState", err)
		ctx.JSON(404, struct{}{})
		return
	}
	ctx.Data(200, "application/json", raw)
}

func NodesInfo(ctx *gin.Context) {
	var nodes []*model.NodeInfo
	var ok bool

	if nodes, ok = model.GetNodesInfoAll(); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	ctx.JSON(200, nodes)
}

func NodesState(ctx *gin.Context) {
	var nodes []*model.NodeInfo
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
	var conn *net.Conn
	var ok bool
	var state model.State = model.State{Name: name}
	var err error
	var raw []byte
	var mtype int

	if conn, ok = model.GetNodeMsgConn(name); !ok {
		goto sendres
	}
	err = message.SendMessage(*conn, message.TypeNodeStatus, []byte{})
	if err != nil {
		goto sendres
	}
	mtype, raw, err = message.RecvMessage(*conn)
	if err != nil || mtype != message.TypeNodeStatusResponse {
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
