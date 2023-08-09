package model

import (
	"brain/alert"
	"brain/message"
	"brain/snp"
	"errors"
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
	serialNum := snp.GenSerial()
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
		message.MsgTypeString[mtype], ", ", string(payload[:min(len(payload), 100)]),
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

	err = message.SendMessageUnique(*conn, mtype, serialNum, payload)
	trace = append(trace, fmt.Sprint(time.Now().Format("01-02 15:04:05 "), name, ", ",
		message.MsgTypeString[mtype], ", ", string(payload[:min(len(payload), 100)]),
		"retry=", 3-retryCnt, ", SendMessageUnique(conn is nil?", conn == nil, ", ",
		message.MsgTypeString[mtype], ", ", string(payload[:min(len(payload), 100)]), "): err=", err))
	if err != nil {
		ResetNodeMsgConn(name)
		if retryCnt == 0 {
			return nil, fmt.Errorf("cannot send to node %s", name)
		}
		time.Sleep(time.Millisecond * 600)
		retryCnt--
		goto retry
	}
	rtype, resbuf, err = message.RecvMessageUnique(*conn)
	trace = append(trace, fmt.Sprint(time.Now().Format("01-02 15:04:05 "),
		name, ", ", message.MsgTypeString[mtype], ", ", string(payload[:min(len(payload), 100)]),
		"retry=", 3-retryCnt, ", RecvMessageUnique(conn is nil?", conn == nil, "): ",
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
		return nil, fmt.Errorf("node %s send malformed response. (%s->%s)", name, message.MsgTypeString[mtype],  message.MsgTypeString[rtype])
	}

	return resbuf, nil
}

func RequestWithTimeout(name string, mtype int, payload []byte, timeout time.Duration) (res []byte, err error) {
	ch := make(chan struct{}, 1)
	go func() {
		res, err = Request(name, mtype, payload)
		ch <- struct{}{}
		close(ch)
	} ()

	select {
	case <- time.After(timeout):
		return nil, errors.New("request timeout")
	case <- ch:
		return
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}