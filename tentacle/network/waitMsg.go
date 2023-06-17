package network

import (
	"fmt"
	"net"
	"tentacle/config"
	"tentacle/heartbeat"
	"tentacle/message"
	"tentacle/service"
	"time"
)

func ReadAndServe() {
	go func() {
		// always loop
		for {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("recover from ", err.(error).Error())
				}
			}()

			conn, err := net.Dial("tcp", brainMsgAddr)
			if err != nil {
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}
			err = message.SendMessage(conn, message.TypeNodeJoin, heartbeat.MakeNodeJoin())
			if err != nil {
				conn.Close()
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}

			_, raw, err := message.RecvMessage(conn)
			if err != nil {
				conn.Close()
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}

			_, err = heartbeat.ParseNodeJoinResponse(raw)
			if err != nil {
				conn.Close()
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}

			// now we just process commmand one by one...
			for {
				err = service.HandleMessage(conn)
				if err != nil {
					conn.Close()
					break
				}
			}
		}
	}()
}
