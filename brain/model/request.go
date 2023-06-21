package model

import (
	"brain/message"
	"fmt"
	"net"
	"time"
)

func Request(name string, mtype int, payload []byte) ([]byte, error) {
	var conn *net.Conn
	var rcode int
	var rtype int
	var resbuf []byte
	var err error
	retryCnt := 3
retry:
	conn, rcode = GetNodeMsgConn(name)
	if rcode == GetConnNoNode {
		return nil, fmt.Errorf("node %s not exists", name)
	} else if rcode == GetConnNoConn {
		if retryCnt == 0 {
			return nil, fmt.Errorf("node %s msgConn off", name)
		}
		time.Sleep(time.Microsecond * 600)
		retryCnt--
		goto retry
	}

	err = message.SendMessage(*conn, mtype, payload)
	if err != nil {
		if retryCnt == 0 {
			return nil, fmt.Errorf("cannot send to node %s", name)
		}
		time.Sleep(time.Microsecond * 600)
		retryCnt--
		goto retry
	}
	rtype, resbuf, err = message.RecvMessage(*conn)
	if err != nil {
		if retryCnt == 0 {
			return nil, fmt.Errorf("cannot read from node %s", name)
		}
		time.Sleep(time.Microsecond * 600)
		retryCnt--
		goto retry
	} else if rtype != mtype+1 {
		return nil, fmt.Errorf("node %s send malformed response", name)
	}
	return resbuf, nil
}
