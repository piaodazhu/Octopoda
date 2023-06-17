package network

import (
	"net"
	"tentacle/config"
	"tentacle/heartbeat"
	"tentacle/logger"
	"tentacle/message"
	"tentacle/service"
	"time"
)

func KeepAlive() {
	time.Sleep(time.Second)
	go func() {
		retry := 0
	reconnect:
		for retry < config.GlobalConfig.Heartbeat.RetryTime {
			conn, err := net.Dial("tcp", brainHeartAddr)
			if err != nil {
				logger.Network.Print("Cannot connect to master. retry = ", retry, err)

				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				retry++
			} else {
				retry = 0

				err := message.SendMessage(conn, message.TypeNodeJoin, heartbeat.MakeNodeJoin())
				if err != nil {
					conn.Close()
					time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
					goto reconnect
				}

				_, raw, err := message.RecvMessage(conn)
				if err != nil {
					conn.Close()
					time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
					goto reconnect
				}

				joinResponse, err := heartbeat.ParseNodeJoinResponse(raw)
				if err != nil {
					conn.Close()
					time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
					goto reconnect
				}

				err = LoopHeartbeat(conn, joinResponse.Cnt)
				if err != nil {
					logger.Network.Print(err)
					conn.Close()
					time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
					goto reconnect
				}
			}
		}

		logger.Network.Print("Cannot connect to master.")
		if config.GlobalConfig.Heartbeat.AutoRestart {
			service.Reboot()
			wg.Done()
		} else {
			logger.Exceptions.Print("Dead but wont restart.")
			retry = 0
			goto reconnect
		}

	}()
}

func LoopHeartbeat(conn net.Conn, start int64) error {
	counter := start
	for {
		err := message.SendMessage(conn, message.TypeHeartbeat, heartbeat.MakeHeartbeat(counter))
		if err != nil {
			return err
		}

		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeHeartbeatResponse {
			return err
		}
		response, err := heartbeat.ParseHeartbeatResponse(raw)
		if err != nil || response.Cnt <= counter {
			return err
		}
		counter = response.Cnt

		time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.SendInterval))
	}
}
