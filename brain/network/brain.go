package network

import (
	"brain/config"
	"brain/heartbeat"
	"brain/logger"
	"brain/message"
	"brain/model"
	"brain/snp"
	"time"

	"fmt"
	"net"
)

var heartbeatListener net.Listener
var messagerListener net.Listener

var startTime time.Time

func InitTentacleFace() {
	startTime = time.Now()
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
	mtype, msg, err := message.RecvMessageUnique(conn)
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

	err = message.SendMessageUnique(conn, message.TypeNodeJoinResponse, snp.GenSerial(), heartbeat.MakeNodeJoinResponse())
	if err != nil {
		logger.Network.Print(err)
		conn.Close()
		return
	}

	_, port, _ := net.SplitHostPort(conn.LocalAddr().String())
	if port == fmt.Sprint(config.GlobalConfig.TentacleFace.HeartbeatPort) {
		// heartbeat connection established
		model.StoreNode(joinRequest.Name, joinRequest.Version, joinRequest.Addr, nil)
		logger.Network.Printf("New node join, name=%s\n", joinRequest.Name)
		startHeartbeat(conn, joinRequest.Name)
	} else {
		model.StoreNode(joinRequest.Name, joinRequest.Version, joinRequest.Addr, &conn)
		logger.Network.Printf("establish msg conn, name=%s\n", joinRequest.Name)
	}
}
