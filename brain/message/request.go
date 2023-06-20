package message

import (
	"net"
)

// abstract:  getconn + sendmsg + receivemsg + retry = ?
func Request(conn *net.Conn, mtype int, payload []byte) ([]byte, error) {
	err := SendMessage(*conn, mtype, payload)
	if err != nil {
		return nil, err
	}
	rtype, resbuf, err := RecvMessage(*conn) 
	if err != nil || rtype != mtype+1 {
		return nil, err
	}
	return resbuf, nil
}
