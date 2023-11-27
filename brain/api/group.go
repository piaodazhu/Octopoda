package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/brain/rdb"
	"github.com/piaodazhu/Octopoda/protocols"
)

func GroupGetGroup(ctx *gin.Context) {
	var name string
	var ok bool
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(http.StatusBadRequest, struct{}{})
		return
	}

	// get nodes from redis
	nodes, ok := rdb.GroupGet(name)
	if !ok {
		ctx.JSON(http.StatusNotFound, struct{}{})
		return
	}

	ginfo := protocols.GroupInfo{
		Name:      name,
		Size:      len(nodes),
		Nodes:     nodes,
		Unhealthy: []string{},
	}
	// check nodes state
	for _, node := range nodes {
		if state, ok := model.GetNodeState(node); !ok || state != protocols.NodeStateReady {
			ginfo.Unhealthy = append(ginfo.Unhealthy, node)
		}
	}
	// write response
	ctx.JSON(http.StatusOK, ginfo)
}

func GroupSetGroup(ctx *gin.Context) {
	var ginfo protocols.GroupInfo
	err := ctx.ShouldBind(&ginfo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, struct{}{})
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
			ctx.String(http.StatusNotFound, "group %s already exists", ginfo.Name)
			return
		}
		invalid := []string{}
		for _, node := range nodes {
			if state, ok := model.GetNodeState(node); !ok || state != protocols.NodeStateReady {
				invalid = append(invalid, node)
			}
		}
		if len(invalid) > 0 {
			ctx.String(http.StatusNotFound, "group %s has unhealthy nodes. reject: %s", ginfo.Name, strings.Join(invalid, ", "))
			return
		}
	}

	// add group nodes
	if ok := rdb.GroupAdd(ginfo.Name, nodes); !ok {
		ctx.String(http.StatusNotFound, "set group %s failed", ginfo.Name)
		return
	}

	// write response
	ctx.String(http.StatusOK, "OK")
}

func GroupDeleteGroup(ctx *gin.Context) {
	var name string
	var ok bool
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(http.StatusBadRequest, struct{}{})
		return
	}

	// group must exists
	_, ok = rdb.GroupGet(name)
	if !ok {
		ctx.JSON(http.StatusNotFound, struct{}{})
		return
	} else {
		rdb.GroupDel(name)
	}

	// write response
	ctx.JSON(http.StatusOK, struct{}{})
}

func GroupGetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, rdb.GroupGetAll())
}
