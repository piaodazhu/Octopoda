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
	_, resbuf, err := RecvMessage(*conn)  // mtype ??
	if err != nil {
		return nil, err
	}
	return resbuf, nil
}
