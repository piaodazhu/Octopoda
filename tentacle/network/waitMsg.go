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

// var localListener net.Listener

// func InitListener() {
// 	var sb strings.Builder
// 	sb.WriteString(config.GlobalConfig.Ip)
// 	sb.WriteByte(':')
// 	sb.WriteString(fmt.Sprint(config.GlobalConfig.Port))
// 	listener, err := net.Listen("tcp", sb.String())
// 	if err != nil {
// 		logger.Exceptions.Panic(err)
// 	}
// 	localListener = listener
// }

// func ListenAndServe() {
// 	defer localListener.Close()
// 	logger.SysInfo.Println("Listening on", localListener.Addr())
// 	for {
// 		conn, err := localListener.Accept()
// 		if err != nil {
// 			logger.Exceptions.Println(err)
// 		}
// 		go service.HandleConn(conn)
// 	}
// }

func ReadAndServe() {
	addr := fmt.Sprintf("%s:%d", config.GlobalConfig.Brain.Ip, config.GlobalConfig.Brain.MessagePort)

	go func() {
		// always loop
		for {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("recover from ", err.(error).Error())
				}
			}()

			conn, err := net.Dial("tcp", addr)
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
