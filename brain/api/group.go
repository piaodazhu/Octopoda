package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/brain/rdb"
)

type GroupInfo struct {
	Name  string   `json:"name" binding:"required"`
	Nodes []string `json:"nodes" binding:"required"`
	// NoCheck can be in request
	NoCheck bool `json:"nocheck" binding:"omitempty"`

	// Size and Unhealthy will be in response
	Size      int      `json:"size" binding:"omitempty"`
	Unhealthy []string `json:"unhealthy" binding:"omitempty"`
}

func GroupGetGroup(ctx *gin.Context) {
	var name string
	var ok bool
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(400, struct{}{})
		return
	}

	// get nodes from redis
	nodes, ok := rdb.GroupGet(name)
	if !ok {
		ctx.JSON(404, struct{}{})
		return
	}

	ginfo := GroupInfo{
		Name:      name,
		Size:      len(nodes),
		Nodes:     nodes,
		Unhealthy: []string{},
	}
	// check nodes state
	for _, node := range nodes {
		if state, ok := model.GetNodeState(node); !ok || state != model.NodeStateReady {
			ginfo.Unhealthy = append(ginfo.Unhealthy, node)
		}
	}
	// write response
	ctx.JSON(200, ginfo)
}

func GroupSetGroup(ctx *gin.Context) {
	var ginfo GroupInfo
	err := ctx.ShouldBind(&ginfo)
	if err != nil {
		ctx.JSON(400, struct{}{})
		return
	}
	// unique
	nodes := []string{}
	seen := map[string]struct{}{}
	for _, node := range ginfo.Nodes {
		seen[node] = struct{}{}
	}
	for node := range seen {
		nodes = append(nodes, node)
	}

	// check nodes
	if !ginfo.NoCheck {
		if rdb.GroupExist(ginfo.Name) {
			ctx.String(404, "group %s already exists", ginfo.Name)
			return
		}
		invalid := []string{}
		for _, node := range nodes {
			if state, ok := model.GetNodeState(node); !ok || state != model.NodeStateReady {
				invalid = append(invalid, node)
			}
		}
		if len(invalid) > 0 {
			ctx.String(404, "group %s has unhealthy nodes. reject: %s", ginfo.Name, strings.Join(invalid, ", "))
			return
		}
	}

	// add group nodes
	if ok := rdb.GroupAdd(ginfo.Name, nodes); !ok {
		ctx.String(404, "set group %s failed", ginfo.Name)
		return
	}

	// write response
	ctx.String(200, "OK")
}

func GroupDeleteGroup(ctx *gin.Context) {
	var name string
	var ok bool
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(400, struct{}{})
		return
	}

	// group must exists
	_, ok = rdb.GroupGet(name)
	if !ok {
		ctx.JSON(404, struct{}{})
		return
	} else {
		rdb.GroupDel(name)
	}

	// write response
	ctx.JSON(200, struct{}{})
}

func GroupGetAll(ctx *gin.Context) {
	ctx.JSON(200, rdb.GroupGetAll())
}
