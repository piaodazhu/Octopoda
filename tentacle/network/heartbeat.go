package network

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/ostp"
	"github.com/piaodazhu/Octopoda/protocols/snp"
	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/heartbeat"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
	"github.com/piaodazhu/Octopoda/tentacle/nameclient"
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
			logger.Network.Printf("[NETDBG] want dial %s", nameclient.BrainHeartAddr)

			appendLog("heartbeat: dail " + nameclient.BrainHeartAddr)
			conn, err := Dial(nameclient.BrainHeartAddr)
			if err != nil {
				appendLog("heartbeat: dail " + nameclient.BrainHeartAddr + " failed because: " + err.Error())
				logger.Network.Print("Cannot connect to brain. retry = ", retry, err)
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				retry++
			} else {
				appendLog("heartbeat: dail ok :" + nameclient.BrainHeartAddr)
				retry = 0
				ipstr := conn.LocalAddr().(*net.TCPAddr).IP.String()
				appendLog("heartbeat: try send node join request")
				err := protocols.SendMessageUnique(conn, protocols.TypeNodeJoin, snp.GenSerial(), heartbeat.MakeNodeJoin(ipstr))
				if err != nil {
					appendLog("heartbeat: failed to send node join request: " + err.Error())
					conn.Close()
					time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
					goto reconnect
				}
				t1 := time.Now().UnixMilli()

				appendLog("heartbeat: try recv node join response")
				_, _, raw, err := protocols.RecvMessageUnique(conn)
				if err != nil {
					appendLog("heartbeat: failed to recv node join response: " + err.Error())
					conn.Close()
					logger.Network.Println("[HBCONN DBG] close heartbeat because err=", err.Error())
					time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
					goto reconnect
				}
				t2 := time.Now().UnixMilli()

				joinResponse, err := heartbeat.ParseNodeJoinResponse(raw)
				if err != nil {
					appendLog("heartbeat: failed to parse recv node join response: " + err.Error())
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
		if canBeRestart() && config.GlobalConfig.Heartbeat.AutoRestart {
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
	appendLog("heartbeat: start heartbeatloop")
	defer appendLog("heartbeat: stop heatbeatloop")
	joinwg.Done()
	defer joinwg.Add(1)
	logger.Network.Print("[NETDBG] start heartbeatloop")
	defer logger.Network.Print("[NETDBG] end heartbeatloop")

	msgConnOffCnt := 0
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
		if err != nil {
			return err
		}
		randNum = response.NewNum
		ostp.EstimateDelay(t1, t2, response.Ts)

		if !response.IsMsgConnOn {
			// trigger
			appendLog("heartbeat: brain response that msgconn is off")
			msgConnOffCnt++
			if msgConnOffCnt > 3 && msgConn != nil {
				appendLog("heartbeat: detect msgconn is really off. close it.")
				msgConnOffCnt = 0
				msgConn.Close()
			}
		} else {
			msgConnOffCnt = 0
		}

		time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.SendInterval))
	}
}

func canBeRestart() bool {
	var cnt int
	buf, err := os.ReadFile("/etc/octopoda/tentacle/restartcnt")
	if err != nil {
		cnt = 1
	} else {
		cnt, err = strconv.Atoi(string(buf))
		if err != nil {
			cnt = 1
		} else {
			cnt++
		}
	}

	os.WriteFile("/etc/octopoda/tentacle/restartcnt", []byte(fmt.Sprint(cnt)), os.ModePerm)
	if cnt > 3 {
		return false 
	} else {
		return true
	}
}
