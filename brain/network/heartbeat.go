package network

import (
	"brain/config"
	"brain/heartbeat"
	"brain/logger"
	"brain/message"
	"brain/model"
	"brain/ticker"
	"context"
	"fmt"
	"net"
	"time"
)

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

func ListenNodeJoin() {
	go func() {
		for {
			conn, err := heartbeatListener.Accept()
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
		id := model.StoreNode(joinRequest.Name, nil)
		logger.Network.Printf("New node join, name=%s, id=%d\n", joinRequest.Name, id)
		startHeartbeat(conn, joinRequest.Name)
	} else {
		id := model.StoreNode(joinRequest.Name, &conn)
		logger.Network.Printf("establish msg conn, name=%s, id=%d\n", joinRequest.Name, id)
		startMessager(conn, joinRequest.Name)
	}
}

func ProcessHeartbeat(ctx context.Context, c chan bool, conn net.Conn) {
	var mtype int
	var msg []byte
	var health bool
	// var hb heartbeat.HeartBeatInfo
	var err error

	for {
		health = true
		mtype, msg, err = message.RecvMessage(conn)
		if err != nil || mtype != message.TypeHeartbeat {
			// logger.Tentacle.Print(err)
			health = false
			goto reportstate
		}

		_, err = heartbeat.ParseHeartbeat(msg)
		if err != nil {
			logger.Network.Print(err)
			health = false
			goto reportstate
		}

		err = message.SendMessage(conn, message.TypeHeartbeatResponse, heartbeat.MakeHeartbeatResponse(ticker.GetTick()))
		if err != nil {
			logger.Network.Print(err)
			health = false
			goto reportstate
		}
	reportstate:
		select {
		case c <- health:
			if !health {
				goto closeconnection
			}
		case <-ctx.Done():
			goto closeconnection
		}
	}
closeconnection:
	close(c)
	conn.Close()
}

// func checkTimestamp(ts int64) bool {
// 	// return abs(time.Now().Unix()-ts) <= 2
// 	return true
// }
// func abs(x int64) int64 {
// 	if x < 0 {
// 		return -x
// 	}
// 	return x
// }

func startHeartbeat(conn net.Conn, name string) {
	timeout := time.Second * time.Duration(config.GlobalConfig.TentacleFace.ActiveTimeout)

	hbchan := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	go ProcessHeartbeat(ctx, hbchan, conn)
	for {
		select {
		case hbstate := <-hbchan:
			if !hbstate {
				// quited
				if !model.DisconnNode(name) {
					goto errout
				}
				goto errout
			} else {
				if !model.UpdateNode(name) {
					goto errout
				}
			}
		case <-time.After(timeout):
			// timeout
			if !model.DisconnNode(name) {
				goto errout
			}
		}
	}
errout:
	cancel()
}

func startMessager(conn net.Conn, name string) {

}
