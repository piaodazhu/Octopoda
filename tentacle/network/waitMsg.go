package network

import (
	"fmt"
	"net"
	"tentacle/config"
	"tentacle/heartbeat"
	"tentacle/message"
	"tentacle/nameclient"
	"tentacle/service"
	"tentacle/snp"
	"time"
)

func ReadAndServe() {
	go func() {
		// always loop
		for {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("recover from ", err.(error).Error())
					ReadAndServe()
				}
			}()
			joinwg.Wait()
			conn, err := net.Dial("tcp", nameclient.BrainMsgAddr)
			// conn, err := dialWithDevice(nameclient.BrainMsgAddr)
			if err != nil {
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}
			ipstr := conn.LocalAddr().(*net.TCPAddr).IP.String()
			err = message.SendMessageUnique(conn, message.TypeNodeJoin, snp.GenSerial(), heartbeat.MakeNodeJoin(ipstr))
			if err != nil {
				conn.Close()
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}

			_, raw, err := message.RecvMessageUnique(conn)
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
