package network

import (
	"fmt"
	"net"
	"protocols"
	"protocols/snp"
	"tentacle/config"
	"tentacle/heartbeat"
	"tentacle/nameclient"
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
					ReadAndServe()
				}
			}()
			joinwg.Wait()
			conn, err := Dial(nameclient.BrainMsgAddr)
			if err != nil {
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}
			ipstr := conn.LocalAddr().(*net.TCPAddr).IP.String()
			err = protocols.SendMessageUnique(conn, protocols.TypeNodeJoin, snp.GenSerial(), heartbeat.MakeNodeJoin(ipstr))
			if err != nil {
				conn.Close()
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}

			_, _, raw, err := protocols.RecvMessageUnique(conn)
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
