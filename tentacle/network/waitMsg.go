package network

import (
	"net"
	"time"

	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/snp"
	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/heartbeat"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
	"github.com/piaodazhu/Octopoda/tentacle/nameclient"
	"github.com/piaodazhu/Octopoda/tentacle/service"
)

func ReadAndServe() {
	var err error
	go func() {
		// always loop
		for {
			defer func() {
				if err := recover(); err != nil {
					appendLog("waitmsg: recover")
					logger.Network.Println("[MSGCONN DBG] PANIC")
					logger.Network.Println("[MSGCONN DBG] panic recover from ", err.(error).Error())
					logger.Network.Println("[MSGCONN DBG] RECOVER")
					ReadAndServe()
				}
			}()
			appendLog("waitmsg: wait hb loop start")
			joinwg.Wait()
			appendLog("waitmsg: hb loop start, try dial")
			logger.Network.Printf("[NETDBG] msg conn dial %s", nameclient.BrainMsgAddr)
			msgConn, err = Dial(nameclient.BrainMsgAddr)
			if err != nil {
				appendLog("waitmsg: failed to dial because: " + err.Error())
				logger.Network.Println("[MSGCONN DBG] cannot dial ", nameclient.BrainMsgAddr)
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}

			appendLog("waitmsg: try to send node join for msgconn")
			ipstr := msgConn.LocalAddr().(*net.TCPAddr).IP.String()
			err = protocols.SendMessageUnique(msgConn, protocols.TypeNodeJoin, snp.GenSerial(), heartbeat.MakeNodeJoin(ipstr))
			if err != nil {
				appendLog("waitmsg: failed to send node join for msgconn because: " + err.Error())
				msgConn.Close()
				logger.Network.Println("[MSGCONN DBG] cannot send msgconn join from ", ipstr, ", err=", err.Error())
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}

			appendLog("waitmsg: try to recv node join for msgconn")
			_, _, raw, err := protocols.RecvMessageUnique(msgConn)
			if err != nil {
				appendLog("waitmsg: failed to recv node join for msgconn because: " + err.Error())
				msgConn.Close()
				logger.Network.Println("[MSGCONN DBG] cannot receive msgconn join response err=", err.Error())
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}

			_, err = heartbeat.ParseNodeJoinResponse(raw)
			if err != nil {
				appendLog("waitmsg: failed to parse recv node join for msgconn response: " + err.Error())
				msgConn.Close()
				logger.Network.Println("[MSGCONN DBG] cannot parse msgconn join response err=", err.Error())
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}

			// now we just process commmand one by one...
			appendLog("waitmsg: start msgconn loop")
			for {
				err = service.HandleMessage(msgConn)
				if err != nil {
					appendLog("waitmsg: failed to handle a message because: " + err.Error())
					msgConn.Close()
					logger.Network.Println("[MSGCONN DBG] close msgconn because err=", err.Error())
					break
				}
				appendLog("waitmsg: succeed to handle a message. continue")
			}
			logger.Network.Print("[NETDBG] end msgconn")
			appendLog("waitmsg: end msgconn. continue")
		}
	}()
}
