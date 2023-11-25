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
	go func() {
		// always loop
		for {
			defer func() {
				if err := recover(); err != nil {
					logger.Network.Println("[MSGCONN DBG] PANIC")
					logger.Network.Println("[MSGCONN DBG] panic recover from ", err.(error).Error())
					logger.Network.Println("[MSGCONN DBG] RECOVER")
					ReadAndServe()
				}
			}()
			joinwg.Wait()
			logger.Network.Printf("[NETDBG] msg conn dial %s", nameclient.BrainMsgAddr)
			conn, err := Dial(nameclient.BrainMsgAddr)
			if err != nil {
				logger.Network.Println("[MSGCONN DBG] cannot dial ", nameclient.BrainMsgAddr)
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}
			ipstr := conn.LocalAddr().(*net.TCPAddr).IP.String()
			err = protocols.SendMessageUnique(conn, protocols.TypeNodeJoin, snp.GenSerial(), heartbeat.MakeNodeJoin(ipstr))
			if err != nil {
				conn.Close()
				logger.Network.Println("[MSGCONN DBG] cannot send msgconn join from ", ipstr, ", err=", err.Error())
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}

			_, _, raw, err := protocols.RecvMessageUnique(conn)
			if err != nil {
				conn.Close()
				logger.Network.Println("[MSGCONN DBG] cannot receive msgconn join response err=", err.Error())
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}

			_, err = heartbeat.ParseNodeJoinResponse(raw)
			if err != nil {
				conn.Close()
				logger.Network.Println("[MSGCONN DBG] cannot parse msgconn join response err=", err.Error())
				time.Sleep(time.Second * time.Duration(config.GlobalConfig.Heartbeat.ReconnectInterval))
				continue
			}

			// now we just process commmand one by one...
			for {
				err = service.HandleMessage(conn)
				if err != nil {
					conn.Close()
					logger.Network.Println("[MSGCONN DBG] close msgconn because err=", err.Error())
					break
				}
			}
			logger.Network.Print("[NETDBG] end msgconn")
		}
	}()
}
