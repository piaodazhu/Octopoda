package network

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/heartbeat"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/brain/workgroup"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/snp"

	"fmt"
	"net"
)

var heartbeatListener net.Listener
var messagerListener net.Listener

var startTime time.Time

func InitTentacleFace() {
	startTime = time.Now()
	var err error

	heartbeatListener, err = newTlsListener(fmt.Sprintf("%s:%d", config.GlobalConfig.TentacleFace.Ip, config.GlobalConfig.TentacleFace.HeartbeatPort))
	if err != nil {
		logger.Exceptions.Panic(err)
	}

	messagerListener, err = newTlsListener(fmt.Sprintf("%s:%d", config.GlobalConfig.TentacleFace.Ip, config.GlobalConfig.TentacleFace.MessagePort))
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

func newTlsListener(address string) (net.Listener, error) {
	lcfg := net.ListenConfig{}
	listener, err := lcfg.Listen(context.Background(), "tcp", address)
	if err != nil {
		return nil, err
	}
	return tls.NewListener(listener, config.TLSConfig), nil
}

func ProcessNodeJoin(conn net.Conn) {
	mtype, _, msg, err := protocols.RecvMessageUnique(conn)
	if err != nil || mtype != protocols.TypeNodeJoin {
		logger.Comm.Print(err, protocols.TypeNodeJoin)
		conn.Close()
		return
	}
	joinRequest, err := heartbeat.ParseNodeJoin(msg)
	if err != nil {
		logger.Network.Print(err)
		conn.Close()
		return
	}

	randNum := snp.GenSerial()
	err = protocols.SendMessageUnique(conn, protocols.TypeNodeJoinResponse, snp.GenSerial(), heartbeat.MakeNodeJoinResponse(randNum))
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
		workgroup.AddMembers("", "", joinRequest.Name)
		startHeartbeat(conn, joinRequest.Name, randNum)
	} else {
		model.StoreNode(joinRequest.Name, joinRequest.Version, joinRequest.Addr, conn)
		logger.Network.Printf("establish msg conn, name=%s\n", joinRequest.Name)
	}
}
