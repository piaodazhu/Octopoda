package network

import (
	"brain/config"
	"brain/heartbeat"
	"brain/logger"
	"brain/message"
	"brain/model"
	"context"
	"fmt"
	"net"
	"time"
)

func ProcessHeartbeat(ctx context.Context, c chan bool, conn net.Conn) {
	var mtype int
	var msg []byte
	var health bool
	var hbinfo heartbeat.HeartBeatInfo
	var err error

	for {
		health = true
		mtype, msg, err = message.RecvMessage(conn)
		if err != nil || mtype != message.TypeHeartbeat {
			// logger.Tentacle.Print(err)
			fmt.Println("tag1", mtype, msg, err)
			health = false
			goto reportstate
		}

		hbinfo, err = heartbeat.ParseHeartbeat(msg)
		if err != nil || hbinfo.Msg != "ping" {
			fmt.Println("tag2", err)
			logger.Network.Print(err)
			health = false
			goto reportstate
		}

		err = message.SendMessage(conn, message.TypeHeartbeatResponse, heartbeat.MakeHeartbeatResponse("pong"))
		if err != nil {
			fmt.Println("tag3", err)
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
