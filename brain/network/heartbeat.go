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
	"strings"
	"time"
)

var tentacleListener net.Listener

func InitTentacleFace() {
	var sb strings.Builder
	sb.WriteString(config.GlobalConfig.TentacleFace.Ip)
	sb.WriteByte(':')
	sb.WriteString(fmt.Sprint(config.GlobalConfig.TentacleFace.Port))
	listener, err := net.Listen("tcp", sb.String())
	if err != nil {
		logger.Exceptions.Panic(err)
	}
	tentacleListener = listener
}

func ListenNodeJoin() {
	go func() {
		for {
			conn, err := tentacleListener.Accept()
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
	if !checkTimestamp(joinRequest.Ts) {
		logger.Network.Print("Invalid TimeStamp")
		conn.Close()
		return
	}

	id := model.StoreNode(joinRequest.Name, joinRequest.IP, joinRequest.Port)
	logger.Network.Printf("New node join, name=%s, id=%d\n", joinRequest.Name, id)

	err = message.SendMessage(conn, message.TypeNodeJoinResponse, heartbeat.MakeNodeJoinResponse())
	if err != nil {
		logger.Network.Print(err)
		conn.Close()
		return
	}

	// heartbeat connection established

	timeout := time.Second * time.Duration(config.GlobalConfig.TentacleFace.ActiveTimeout)

	hbchan := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	go ProcessHeartbeat(ctx, hbchan, conn)
	for {
		select {
		case hbstate := <-hbchan:
			if !hbstate {
				// quited
				if !model.DisconnNode(joinRequest.Name) {
					goto errout
				}
				goto errout
			} else {
				if !model.UpdateNode(joinRequest.Name) {
					goto errout
				}
			}
		case <-time.After(timeout):
			// timeout
			if !model.DisconnNode(joinRequest.Name) {
				goto errout
			}
		}
	}
errout:
	cancel()
}

func ProcessHeartbeat(ctx context.Context, c chan bool, conn net.Conn) {
	var mtype int
	var msg []byte
	var health bool
	var hb heartbeat.HeartBeatInfo
	var err error

	for {
		health = true
		mtype, msg, err = message.RecvMessage(conn)
		if err != nil || mtype != message.TypeHeartbeat {
			// logger.Tentacle.Print(err)
			health = false
			goto reportstate
		}

		hb, err = heartbeat.ParseHeartbeat(msg)
		if err != nil || !checkTimestamp(hb.Ts) {
			logger.Network.Print(err)
			health = false
			goto reportstate
		}

		err = message.SendMessage(conn, message.TypeHeartbeatResponse, heartbeat.MakeHeartbeatResponse(ticker.GetTick()))
		if err != nil || !checkTimestamp(hb.Ts) {
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

func checkTimestamp(ts int64) bool {
	return abs(time.Now().Unix()-ts) <= 2
}
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
