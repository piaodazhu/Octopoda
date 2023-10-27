package model

import (
	"fmt"
	"net"
	"protocols"
	"strings"
	"sync"
)

type ConnMsg struct {
	SerialNum uint32
	Mtype     int
	Raw       []byte
}

type ConnInfo struct {
	ConnState string
	Conn      net.Conn
	Messages  sync.Map
}

func CreateConnInfo(conn net.Conn) ConnInfo {
	connState := "On"
	if conn == nil {
		connState = "Off"
	}
	return ConnInfo{
		Conn:      conn,
		ConnState: connState,
		Messages:  sync.Map{},
	}
}

func (c *ConnInfo) Close() {
	if c.Conn != nil {
		c.Conn.Close()
		c.Conn = nil
	}
	c.ConnState = "Off"

	pendingList := []uint32{}
	c.Messages.Range(func(key, value any) bool {
		pendingList = append(pendingList, key.(uint32))
		return true
	})
	for _, pendingSerialNum := range pendingList {
		if value, exist := c.Messages.LoadAndDelete(pendingSerialNum); exist {
			mchan := value.(chan *ConnMsg)
			mchan <- nil
			close(mchan)
		}
	}
}

func (c *ConnInfo) Fresh(conn net.Conn) {
	c.Close()
	c.Conn = conn
	if c.Conn == nil {
		c.ConnState = "Off"
	} else {
		c.ConnState = "On"
	}
}

func (c *ConnInfo) StartReceive() {
	go func() {
		// fmt.Println("[DEBUG] start receving...")
		for {
			conn := c.Conn
			if conn == nil {
				return
			}
			mtype, serialNum, raw, err := protocols.RecvMessageUnique(conn)
			if err != nil {
				// TODO error reason
				// fmt.Println("[DEBUG] receive error: ", err)
				c.Close()
				return
			}
			c.NotifyMsg(ConnMsg{
				SerialNum: serialNum,
				Mtype:     mtype,
				Raw:       raw,
			})
		}
	}()
}

func (c *ConnInfo) ListenMsg(serialNum uint32) bool {
	mchan := make(chan *ConnMsg, 1)
	if _, exist := c.Messages.LoadOrStore(serialNum, mchan); exist {
		// duplicated message
		// fmt.Println("[DEBUG] listen failed: ", serialNum)
		close(mchan)
		return false
	}
	// fmt.Println("[DEBUG] listen done: ", serialNum)
	return true
}

func (c *ConnInfo) NotifyMsg(msg ConnMsg) bool {
	if value, exist := c.Messages.LoadAndDelete(msg.SerialNum); exist {
		mchan := value.(chan *ConnMsg)
		mchan <- &msg
		close(mchan)
		// fmt.Println("[DEBUG] notify done: ", msg)
		return true
	}

	// fmt.Println("[DEBUG] notify failed: ", msg)
	return false
}

func (c *ConnInfo) WaitMsg(serialNum uint32) (msg *ConnMsg, ok bool) {
	debug_waitmsg := strings.Builder{}
	defer func() {
		debug_waitmsg.WriteByte('\n')
		if !ok {
			fmt.Println(debug_waitmsg.String())
		}
	} ()

	debug_waitmsg.WriteString(fmt.Sprintf("serialNum %d--> ", serialNum))
	if value, exist := c.Messages.Load(serialNum); exist {
		mchan := value.(chan *ConnMsg)
		// fmt.Println("[DEBUG] waiting...: ", serialNum)
		debug_waitmsg.WriteString("Waiting --> ")
		msg = <-mchan
		if msg == nil {
			// fmt.Println("[DEBUG] wait canceled: ", serialNum)
			debug_waitmsg.WriteString("msg is nil --> ")
			ok = false
			return
		}
		// fmt.Println("[DEBUG] wait done: ", msg)
		debug_waitmsg.WriteString("msg get! --> ")
		ok = true
		return
	}
	// fmt.Println("[DEBUG] wait failed: ", serialNum)
	debug_waitmsg.WriteString("serailNum not found --> ")

	msg = nil 
	ok = false
	return
}

// It may be attacked by flooding...
