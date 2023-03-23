package api

import (
	"brain/logger"
	"brain/message"
	"brain/model"
	"encoding/json"
	"net"
	"sync"

	"github.com/gin-gonic/gin"
)

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
	}
	ctx.JSON(200, node)
}

func NodeState(ctx *gin.Context) {
	var name, addr string
	var ok bool
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	if addr, ok = model.GetNodeAddress(name); !ok {
		ctx.JSON(404, struct{}{})
		return
	}

	if conn, err := net.Dial("tcp", addr); err != nil {
		ctx.JSON(404, struct{}{})
	} else {
		defer conn.Close()
		message.SendMessage(conn, message.TypeNodeState, []byte{})
		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeNodeStateResponse {
			logger.Tentacle.Println("NodeState", err)
			ctx.JSON(404, struct{}{})
			return
		}
		ctx.Data(200, "application/json", raw)
	}
}

func NodesInfo(ctx *gin.Context) {
	var nodes []*model.NodeModel
	var ok bool

	if nodes, ok = model.GetNodesInfoAll(); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	ctx.JSON(200, nodes)
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
		go getNodeState(node.Addr, channel, &wg)
	}
	wg.Wait()
	close(channel)
	for v := range channel {
		states = append(states, v)
	}
	ctx.JSON(200, states)
}

func getNodeState(addr string, channel chan<- model.State, wg *sync.WaitGroup) {
	defer wg.Done()

	if conn, err := net.Dial("tcp", addr); err != nil {
		return
	} else {
		defer conn.Close()

		message.SendMessage(conn, message.TypeNodeState, []byte{})
		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeNodeStateResponse {
			logger.Tentacle.Println("getNodeState", err)
			return
		}

		var state model.State
		err = json.Unmarshal(raw, &state)
		if err != nil {
			logger.Tentacle.Println("UnmarshalNodeState", err)
			return
		}
		channel <- state
	}
}


func NodeReboot(ctx *gin.Context) {
	var name, addr string
	var ok bool
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	if addr, ok = model.GetNodeAddress(name); !ok {
		ctx.JSON(404, struct{}{})
		return
	}

	if conn, err := net.Dial("tcp", addr); err != nil {
		ctx.JSON(404, struct{}{})
	} else {
		defer conn.Close()
		message.SendMessage(conn, message.TypeCommandReboot, []byte{})
		mtype, _, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeCommandResponse {
			logger.Tentacle.Println("NodeReboot", err)
			ctx.JSON(404, struct{}{})
			return
		}
		ctx.JSON(200, struct{}{})
	}
}
