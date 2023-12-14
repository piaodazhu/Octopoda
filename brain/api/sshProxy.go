package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/brain/network"
	"github.com/piaodazhu/Octopoda/brain/workgroup"
	"github.com/piaodazhu/Octopoda/protocols"
)

func SshLoginInfo(ctx *gin.Context) {
	name := ctx.Query("name")
	if !workgroup.IsInScope(ctx.GetStringMapString("octopoda_scope"), name) {
		ctx.JSON(http.StatusBadRequest, struct{}{})
		return
	}

	if info, found := network.GetSshInfo(name); found {
		ctx.JSON(http.StatusOK, info)
		return
	}
	ctx.JSON(http.StatusNotFound, struct{}{})
}

func SshRegister(ctx *gin.Context) {
	name := ctx.PostForm("name")
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if len(name) == 0 || len(username) == 0 {
		ctx.JSON(http.StatusBadRequest, network.ProxyMsg{
			Code: -1,
			Msg:  "ERR",
			Data: "invalid arguments",
		})
		return
	}

	if !workgroup.IsInScope(ctx.GetStringMapString("octopoda_scope"), name) {
		ctx.JSON(http.StatusBadRequest, network.ProxyMsg{
			Code: -1,
			Msg:  "ERR",
			Data: "some nodes are invalid or out of scope",
		})
		return
	}

	services, err := network.ProxyServices()
	if err != nil {
		ctx.JSON(http.StatusNotFound, network.ProxyMsg{
			Code: -1,
			Msg:  "ERR",
			Data: fmt.Sprintf("cannot get proxy services: %s", err.Error()),
		})
		return
	}

	// 存在一致性问题，但是危害不大
	for _, s := range services {
		if name == s.Name {
			ctx.JSON(http.StatusBadRequest, network.ProxyMsg{
				Code: -1,
				Msg:  "ERR",
				Data: fmt.Sprintf("service %s already exists", name),
			})
			return
		}
	}
	network.CreateSshInfo(name, username, password)
	if !proxyCmd(ctx, name, protocols.TypeSshRegister) { // 不成功就删除
		network.DelSshInfo(name)
	}
}

func SshUnregister(ctx *gin.Context) {
	name := ctx.Query("name")
	if !workgroup.IsInScope(ctx.GetStringMapString("octopoda_scope"), name) {
		ctx.JSON(http.StatusBadRequest, network.ProxyMsg{
			Code: -1,
			Msg:  "ERR",
			Data: "some nodes are invalid or out of scope",
		})
		return
	}

	network.DelSshInfo(name)
	if !proxyCmd(ctx, name, protocols.TypeSshUnregister) {
		logger.SysInfo.Println("cannot unregister ssh proxy of node " + name)
	}
}

func proxyCmd(ctx *gin.Context, name string, cmdType int) bool {
	if name == "brain" {
		ip, _ := network.GetOctlFaceIp()
		if cmdType == protocols.TypeSshRegister {
			network.CompleteSshInfo(name, ip, uint32(config.GlobalConfig.OctlFace.SshPort))
		}
		ctx.JSON(http.StatusOK, network.ProxyMsg{
			Code: 0,
			Msg:  "OK",
			Data: fmt.Sprintf("%s:%d", ip, config.GlobalConfig.OctlFace.SshPort),
		})
		return true
	}
	if state, ok := model.GetNodeState(name); !ok || state != protocols.NodeStateReady {
		ctx.JSON(http.StatusNotFound, network.ProxyMsg{
			Code: -1,
			Msg:  "ERR",
			Data: fmt.Sprintf("node %s not found", name),
		})
		return false
	}
	raw, err := model.Request(name, cmdType, []byte{})
	if err != nil {
		logger.Comm.Println(protocols.MsgTypeString[cmdType], err)
		ctx.JSON(http.StatusInternalServerError, network.ProxyMsg{
			Code: -1,
			Msg:  "ERR",
			Data: fmt.Sprintf("brain request %s error", name),
		})
		return false
	}

	pmsg := network.ProxyMsg{}
	err = json.Unmarshal(raw, &pmsg)
	if err != nil {
		logger.Comm.Println("network.ProxyMsg Unmarshal", err)
		ctx.JSON(http.StatusInternalServerError, network.ProxyMsg{
			Code: -1,
			Msg:  "ERR",
			Data: fmt.Sprintf("network.ProxyMsg Unmarshal error: %s", err),
		})
		return false
	}
	if pmsg.Code == 0 {
		ctx.JSON(http.StatusOK, pmsg)
		return true
	} else {
		ctx.JSON(http.StatusInternalServerError, pmsg)
		return false
	}
}
