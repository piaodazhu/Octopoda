package api

import (
	"brain/logger"
	"brain/message"
	"brain/model"
	"net"

	"github.com/gin-gonic/gin"
)

func SSHInfo(ctx *gin.Context) {
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
