package model

import (
	"brain/alert"
	"brain/message"
	"fmt"
	"net"
	"strings"
	"time"
)

func Request(name string, mtype int, payload []byte) ([]byte, error) {
	var conn *net.Conn
	var rcode int
	var rtype int
	var resbuf []byte
	var err error
	retryCnt := 3
	trace := []string{}
	defer func() {
		if resbuf == nil || retryCnt != 3 {
			msg := fmt.Sprintf("[TRACE REQUEST]: resbuf is nil?%t, retryCnt=%d\n%s", resbuf==nil, retryCnt, strings.Join(trace, "\n"))
			alert.Alert(msg)
		}
	}()
retry:
	conn, rcode = GetNodeMsgConn(name)
	trace = append(trace, fmt.Sprint(time.Now().Format("01-02 15:04:05 "), name, ", ",
		message.MsgTypeString[mtype], ", ", string(payload),
		"retry=", 3-retryCnt, ", GetNodeMsgConn(", name, "): conn is nil?",
		conn == nil, ", rcode=", rcode))
	if rcode == GetConnNoNode {
		return nil, fmt.Errorf("node %s not exists", name)
	} else if rcode == GetConnNoConn {
		if retryCnt == 0 {
			return nil, fmt.Errorf("node %s msgConn off", name)
		}
		time.Sleep(time.Millisecond * 600)
		retryCnt--
		goto retry
	}

	err = message.SendMessage(*conn, mtype, payload)
	trace = append(trace, fmt.Sprint(time.Now().Format("01-02 15:04:05 "), name, ", ",
		message.MsgTypeString[mtype], ", ", string(payload),
		"retry=", 3-retryCnt, ", SendMessage(conn is nil?", conn == nil, ", ",
		message.MsgTypeString[mtype], ", ", string(payload), "): err=", err))
	if err != nil {
		ResetNodeMsgConn(name)
		if retryCnt == 0 {
			return nil, fmt.Errorf("cannot send to node %s", name)
		}
		time.Sleep(time.Millisecond * 600)
		retryCnt--
		goto retry
	}
	rtype, resbuf, err = message.RecvMessage(*conn)
	trace = append(trace, fmt.Sprint(time.Now().Format("01-02 15:04:05 "),
		name, ", ", message.MsgTypeString[mtype], ", ", string(payload),
		"retry=", 3-retryCnt, ", RecvMessage(conn is nil?", conn == nil, "): ",
		message.MsgTypeString[rtype], ", ", string(resbuf), ", ", err))
	if err != nil {
		ResetNodeMsgConn(name)
		if retryCnt == 0 {
			return nil, fmt.Errorf("cannot read from node %s", name)
		}
		time.Sleep(time.Millisecond * 600)
		retryCnt--
		goto retry
	} else if rtype != mtype+1 {
		return nil, fmt.Errorf("node %s send malformed response", name)
	}

	return resbuf, nil
}
