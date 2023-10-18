package network

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"tentacle/config"
	"tentacle/heartbeat"
	"tentacle/logger"
	"tentacle/message"
	"tentacle/nameclient"
	"tentacle/snp"
	"time"
)

func KeepAlive() {
	time.Sleep(time.Second)
	go func() {
		retry := 0
	reconnect:
		nameclient.ResolveBrain()
		for retry < config.GlobalConfig.Heartbeat.RetryTime {
			conn, err := Dial(nameclient.BrainHeartAddr)
			if err != nil {
				logger.Network.Print("Cannot connect to brain. retry = ", retry, err)
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				retry++
			} else {
				retry = 0
				ipstr := conn.LocalAddr().(*net.TCPAddr).IP.String()
				fmt.Println("Make NodeJoin, LocalAddr is ", ipstr)
				err := message.SendMessageUnique(conn, message.TypeNodeJoin, snp.GenSerial(), heartbeat.MakeNodeJoin(ipstr))
				if err != nil {
					conn.Close()
					time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
					goto reconnect
				}

				_, raw, err := message.RecvMessageUnique(conn)
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

				err = SynchronizeTime(joinResponse.Ts)
				if err != nil {
					logger.Network.Print(err)
					// conn.Close()
					// time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
					// goto reconnect
				}

				err = LoopHeartbeat(conn, joinResponse.NewNum)
				if err != nil {
					logger.Network.Print(err)
					conn.Close()
					time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
					goto reconnect
				}
			}
		}

		logger.Network.Print("Cannot connect to brain.")
		if config.GlobalConfig.Heartbeat.AutoRestart {
			exec.Command("reboot").Run()
			wg.Done()
		} else {
			logger.Exceptions.Print("Dead but wont restart.")
			retry = 0
			goto reconnect
		}

	}()
}

func LoopHeartbeat(conn net.Conn, randNum uint32) error {
	joinwg.Done()
	defer joinwg.Add(1)
	for {
		err := message.SendMessageUnique(conn, message.TypeHeartbeat, snp.GenSerial(), heartbeat.MakeHeartbeat(randNum))
		if err != nil {
			return err
		}

		mtype, raw, err := message.RecvMessageUnique(conn)
		if err != nil || mtype != message.TypeHeartbeatResponse {
			return err
		}
		response, err := heartbeat.ParseHeartbeatResponse(raw)
		if err != nil || response.Msg != "pong" {
			return err
		}
		randNum = response.NewNum

		time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.SendInterval))
	}
}

func SynchronizeTime(ts int64) error {
	// seems that it may not work. To be fixed.
	var err error = nil
	if runtime.GOOS == "linux" {
		cmd := exec.Command("date", "-s", time.UnixMicro(ts).Format("01/02/2006 15:04:05.999"))
		err = cmd.Run()
	}
	return err
}
