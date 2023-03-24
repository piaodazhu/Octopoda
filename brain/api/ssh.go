package api

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"brain/model"
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
)

type sshInfo struct {
	Addr     string
	Username string
	Password string
}

func SSHInfo(ctx *gin.Context) {
	var name, addr string
	var ok bool
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	if name == "master" {
		sinfo := sshInfo{
			Addr: fmt.Sprintf("%s:%d", config.GlobalConfig.Sshinfo.Ip, config.GlobalConfig.Sshinfo.Port),
			Username: config.GlobalConfig.Sshinfo.Username,
			Password: config.GlobalConfig.Sshinfo.Password,
		}
		ctx.JSON(200, sinfo)
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
		message.SendMessage(conn, message.TypeCommandSSH, []byte{})
		mtype, payload, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeCommandResponse {
			logger.Tentacle.Println("NodeReboot", err)
			ctx.JSON(404, struct{}{})
			return
		}
		ctx.Data(200, "application/json", payload)
	}
}
