package network

import (
	"fmt"
	"net"
	"os/exec"
	"time"

	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/snp"
	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/heartbeat"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
	"github.com/piaodazhu/Octopoda/tentacle/nameclient"
	"github.com/piaodazhu/Octopoda/protocols/ostp"
)

func KeepAlive() {
	time.Sleep(time.Second)
	go func() {
		retry := 0
	reconnect:
		nameclient.ResolveBrain()

		// TODO: ugly code
		resolvRetry := 0
		for len(nameclient.BrainHeartAddr) == 0 {
			time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
			nameclient.ResolveBrain()
			resolvRetry++
			if resolvRetry > 10 {
				break
			}
		}

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
				err := protocols.SendMessageUnique(conn, protocols.TypeNodeJoin, snp.GenSerial(), heartbeat.MakeNodeJoin(ipstr))
				if err != nil {
					conn.Close()
					time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
					goto reconnect
				}
				t1 := time.Now().UnixMilli()

				_, _, raw, err := protocols.RecvMessageUnique(conn)
				if err != nil {
					conn.Close()
					logger.Network.Println("[HBCONN DBG] close heartbeat because err=", err.Error())
					time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
					goto reconnect
				}
				t2 := time.Now().UnixMilli()

				joinResponse, err := heartbeat.ParseNodeJoinResponse(raw)
				if err != nil {
					conn.Close()
					time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
					goto reconnect
				}

				ostp.EstimateDelay(t1, t2, joinResponse.Ts)

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
			exec.Command("bash", "-c", config.GlobalConfig.Heartbeat.AutoRestartCommand).Run()
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
		err := protocols.SendMessageUnique(conn, protocols.TypeHeartbeat, snp.GenSerial(), heartbeat.MakeHeartbeat(randNum))
		if err != nil {
			return err
		}
		t1 := time.Now().UnixMilli()

		mtype, _, raw, err := protocols.RecvMessageUnique(conn)
		if err != nil || mtype != protocols.TypeHeartbeatResponse {
			return err
		}
		t2 := time.Now().UnixMilli()

		response, err := heartbeat.ParseHeartbeatResponse(raw)
		if err != nil || response.Msg != "pong" {
			return err
		}
		randNum = response.NewNum
		ostp.EstimateDelay(t1, t2, response.Ts)

		time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.SendInterval))
	}
}
