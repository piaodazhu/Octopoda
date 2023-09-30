package api

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"brain/model"
	"brain/network"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type proxyMsg struct {
	Code int
	Msg  string
	Data string
}

func SshLoginInfo(ctx *gin.Context) {
	name := ctx.Query("name")
	if info, found := network.GetSshInfo(name); found {
		ctx.JSON(200, info)
		return 
	}
	ctx.JSON(404, struct{}{})
}

func SshRegister(ctx *gin.Context) {
	name := ctx.PostForm("name")
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if len(name) == 0 || len(username) == 0 {
		ctx.JSON(400, proxyMsg{
			Code: -1,
			Msg: "ERR",
			Data: "invalid arguments",
		})
		return
	}
	services, err := network.ProxyServices()
	if err != nil {
		ctx.JSON(404, proxyMsg{
			Code: -1,
			Msg: "ERR",
			Data: fmt.Sprintf("cannot get proxy services: %s", err.Error()),
		})
		return
	}

	// 存在一致性问题，但是危害不大
	for _, s := range services {
		if name == s.Name {
			ctx.JSON(400, proxyMsg{
				Code: -1,
				Msg: "ERR",
				Data: fmt.Sprintf("service %s already exists", name),
			})
			return
		}
	}
	network.CreateSshInfo(name, username, password)
	if !proxyCmd(ctx, name, message.TypeSshRegister) { // 不成功就删除
		network.DelSshInfo(name)
	}
}

func SshUnregister(ctx *gin.Context) {
	name := ctx.Query("name")
	if proxyCmd(ctx, name, message.TypeSshUnregister) { // 成功就删除
		network.DelSshInfo(name)
	}
}

func proxyCmd(ctx *gin.Context, name string, cmdType int) bool {
	if name == "master" {
		ip, _ := network.GetOctlFaceIp()
		if cmdType == message.TypeSshRegister {
			network.CompleteSshInfo(name, ip, uint32(config.GlobalConfig.OctlFace.SshPort))
		}
		ctx.JSON(200, proxyMsg{
			Code: 0,
			Msg: "OK",
			Data: fmt.Sprintf("%s:%d", ip, config.GlobalConfig.OctlFace.SshPort),
		})
		return true
	}
	if state, ok := model.GetNodeState(name); !ok || state != model.NodeStateReady {
		ctx.JSON(404, proxyMsg{
			Code: -1,
			Msg: "ERR",
			Data: fmt.Sprintf("node %s not found", name),
		})
		return false
	}
	raw, err := model.Request(name, cmdType, []byte{})
	if err != nil {
		logger.Comm.Println(message.MsgTypeString[cmdType], err)
		ctx.JSON(500, proxyMsg{
			Code: -1,
			Msg: "ERR",
			Data: fmt.Sprintf("master request %s error", name),
		})
		return false
	}

	pmsg := proxyMsg{}
	err = json.Unmarshal(raw, &pmsg)
	if err != nil {
		logger.Comm.Println("proxyMsg Unmarshal", err)
		ctx.JSON(500, proxyMsg{
			Code: -1,
			Msg: "ERR",
			Data: fmt.Sprintf("proxyMsg Unmarshal error: %s", err),
		})
		return false
	}
	if pmsg.Code == 0 {
		ctx.JSON(200, pmsg)
		return true
	} else {
		ctx.JSON(500, pmsg)
		return false
	}
}
