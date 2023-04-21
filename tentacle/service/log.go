package service

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"strings"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/message"
	"time"
)

type LogParams struct {
	MaxLines      int
	MaxDaysBefore int
	Logs          []string
}

func NodeLog(conn net.Conn, raw []byte) {
	lparams := LogParams{}
	json.Unmarshal(raw, &lparams)
	if lparams.MaxLines <= 0 {
		lparams.MaxLines = 30
	}
	if lparams.MaxDaysBefore <= 0 {
		lparams.MaxDaysBefore = 0
	}
	
	readLogs(&lparams)

	response, _ := json.Marshal(&lparams)

	err := message.SendMessage(conn, message.TypeNodeLogResponse, response)
	if err != nil {
		logger.Comm.Print("NodeLog")
	}
}

func readLogs(params *LogParams) {
	offset := 0
	daysBefore := 0
	validdate, validlines := 0, 0
	params.Logs = []string{}

	var sb strings.Builder
	sb.WriteString(config.GlobalConfig.Logger.Path)
	sb.WriteByte('/')
	sb.WriteString(config.GlobalConfig.Logger.NamePrefix)
	prefix := sb.String()

	for offset < params.MaxLines && daysBefore <= params.MaxDaysBefore {
		fname := prefix + time.Now().AddDate(0, 0, -daysBefore).Format("_2006_01_02.log")
		daysBefore++
		f, err := os.Open(fname)
		if err != nil {
			logger.Exceptions.Print(fname, err)
			continue
		}

		validdate = daysBefore - 1

		reader := bufio.NewReader(f)
		logbuf := []string{}
		cnt := 0
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			logbuf = append(logbuf, line[:len(line)-1])
			cnt++
		}
		f.Close()
		if offset + cnt > params.MaxLines {
			cnt = params.MaxLines - offset
		}
		logbuf = logbuf[len(logbuf) - cnt:]

		params.Logs = append(logbuf, params.Logs...)
		offset += cnt

		validlines = offset
	}
	params.MaxLines = validlines
	params.MaxDaysBefore = validdate
}
