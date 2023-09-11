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

func SshRegister(ctx *gin.Context) {
	proxyCmd(ctx, ctx.PostForm("name"), message.TypeSshRegister)
}

func SshUnregister(ctx *gin.Context) {
	proxyCmd(ctx, ctx.Query("name"), message.TypeSshUnregister)
}

func proxyCmd(ctx *gin.Context, name string, cmdType int) {
	if name == "master" {
		ip, _ := network.GetOctlFaceIp()
		ctx.JSON(200, proxyMsg{
			Code: 0,
			Msg: "OK",
			Data: fmt.Sprintf("%s:%d", ip, config.GlobalConfig.OctlFace.SshPort),
		})
		return
	}
	if state, ok := model.GetNodeState(name); !ok || state != model.NodeStateReady {
		ctx.JSON(404, struct{}{})
		return
	}
	raw, err := model.Request(name, cmdType, []byte{})
	if err != nil {
		logger.Comm.Println(message.MsgTypeString[cmdType], err)
		ctx.JSON(500, struct{}{})
		return
	}

	pmsg := proxyMsg{}
	err = json.Unmarshal(raw, &pmsg)
	if err != nil {
		logger.Comm.Println("proxyMsg Unmarshal", err)
		ctx.JSON(500, struct{}{})
		return
	}
	if pmsg.Code == 0 {
		ctx.JSON(200, pmsg)
	} else {
		ctx.JSON(500, pmsg)
	}
}
