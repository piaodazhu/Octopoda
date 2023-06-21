package network

import (
	"brain/config"
	"brain/heartbeat"
	"brain/logger"
	"brain/message"
	"brain/model"
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
)

var engine *gin.Engine
var listenaddr string

func InitBrainFace() {
	listenaddr = fmt.Sprintf("%s:%d", config.GlobalConfig.OctlFace.Ip, config.GlobalConfig.OctlFace.Port)

	// gin.SetMode(gin.DebugMode)
	// engine = gin.Default()

	gin.SetMode(gin.ReleaseMode)
	engine = gin.New()
	engine.Use(gin.Recovery())

	initRouter(engine)
}

var heartbeatListener net.Listener
var messagerListener net.Listener

func InitTentacleFace() {
	var err error
	addr1 := fmt.Sprintf("%s:%d", config.GlobalConfig.TentacleFace.Ip, config.GlobalConfig.TentacleFace.HeartbeatPort)
	heartbeatListener, err = net.Listen("tcp", addr1)
	if err != nil {
		logger.Exceptions.Panic(err)
	}

	addr2 := fmt.Sprintf("%s:%d", config.GlobalConfig.TentacleFace.Ip, config.GlobalConfig.TentacleFace.MessagePort)
	messagerListener, err = net.Listen("tcp", addr2)
	if err != nil {
		logger.Exceptions.Panic(err)
	}
}

func WaitNodeJoin() {
	acceptNodeJoin(heartbeatListener)
	acceptNodeJoin(messagerListener)
}

func acceptNodeJoin(listener net.Listener) {
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Network.Print(err)
				continue
			}
			go ProcessNodeJoin(conn)
		}
	}()
}

func ProcessNodeJoin(conn net.Conn) {
	mtype, msg, err := message.RecvMessage(conn)
	if err != nil || mtype != message.TypeNodeJoin {
		logger.Comm.Print(err, message.TypeNodeJoin)
		conn.Close()
		return
	}
	joinRequest, err := heartbeat.ParseNodeJoin(msg)
	if err != nil {
		logger.Network.Print(err)
		conn.Close()
		return
	}

	err = message.SendMessage(conn, message.TypeNodeJoinResponse, heartbeat.MakeNodeJoinResponse())
	if err != nil {
		logger.Network.Print(err)
		conn.Close()
		return
	}

	_, port, _ := net.SplitHostPort(conn.LocalAddr().String())
	if port == fmt.Sprint(config.GlobalConfig.TentacleFace.HeartbeatPort) {
		// heartbeat connection established
		model.StoreNode(joinRequest.Name, nil)
		logger.Network.Printf("New node join, name=%s\n", joinRequest.Name)
		startHeartbeat(conn, joinRequest.Name)
	} else {
		model.StoreNode(joinRequest.Name, &conn)
		logger.Network.Printf("establish msg conn, name=%s\n", joinRequest.Name)
	}
}

func ListenCommand() {
	engine.Run(listenaddr)
}
