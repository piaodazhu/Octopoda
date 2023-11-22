package network

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/piaodazhu/Octopoda/brain/alert"
	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/heartbeat"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/snp"
)

type hbState struct {
	isHealthy bool
	delay     int64
}

func ProcessHeartbeat(ctx context.Context, c chan hbState, conn net.Conn, randNum uint32) {
	var mtype int
	var msg []byte
	var hs hbState
	var hbinfo protocols.HeartBeatRequest
	var err error

	for {
		hs.isHealthy = true
		mtype, _, msg, err = protocols.RecvMessageUnique(conn)
		if err != nil || mtype != protocols.TypeHeartbeat {
			logger.Network.Print(err) // TODO who?
			hs.isHealthy = false
			goto reportstate
		}

		hbinfo, err = heartbeat.ParseHeartbeat(msg)
		if err != nil || hbinfo.Num != randNum {
			logger.Network.Println(err)
			hs.isHealthy = false
			goto reportstate
		}
		hs.delay = hbinfo.Delay

		randNum = snp.GenSerial()
		err = protocols.SendMessageUnique(conn, protocols.TypeHeartbeatResponse, snp.GenSerial(), heartbeat.MakeHeartbeatResponse(randNum))
		if err != nil {
			logger.Network.Print(err)
			hs.isHealthy = false
			goto reportstate
		}
	reportstate:
		select {
		case c <- hs:
			if !hs.isHealthy {
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

func startHeartbeat(conn net.Conn, name string, randNum uint32) {
	timeout := time.Second * time.Duration(config.GlobalConfig.TentacleFace.ActiveTimeout)
	hbStartTime := time.Now()

	hbchan := make(chan hbState, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go ProcessHeartbeat(ctx, hbchan, conn, randNum)
	for {
		select {
		case hbstate := <-hbchan:
			if !hbstate.isHealthy {
				// quited
				if !model.DisconnNode(name) {
					goto errout
				}
				goto errout
			} else {
				if !model.UpdateNode(name, hbstate.delay) {
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
	brainLive := time.Since(startTime)
	tentacleLive := time.Since(hbStartTime)
	if brainLive > 5*time.Minute && tentacleLive > time.Minute {
		msg := fmt.Sprintf("[TRACE NODESTATE]: node <%s> is offline. Brain has been live for %s, this node has been live for %s.\n", name, brainLive, tentacleLive)
		alert.Alert(msg)
	}
}
