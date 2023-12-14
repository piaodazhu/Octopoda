package service

import (
	"bufio"
	"net"
	"os"
	"strings"
	"time"

	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
)

func NodeLog(conn net.Conn, serialNum uint32, raw []byte) {
	lparams := protocols.LogParams{}
	config.Jsoner.Unmarshal(raw, &lparams)
	if lparams.MaxLines <= 0 {
		lparams.MaxLines = 30
	}
	if lparams.MaxDaysBefore <= 0 {
		lparams.MaxDaysBefore = 0
	}

	readLogs(&lparams)

	response, _ := config.Jsoner.Marshal(&lparams)

	err := protocols.SendMessageUnique(conn, protocols.TypeNodeLogResponse, serialNum, response)
	if err != nil {
		logger.Comm.Print("NodeLog")
	}
}

func readLogs(params *protocols.LogParams) {
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
		if offset+cnt > params.MaxLines {
			cnt = params.MaxLines - offset
		}
		logbuf = logbuf[len(logbuf)-cnt:]

		params.Logs = append(logbuf, params.Logs...)
		offset += cnt

		validlines = offset
	}
	params.MaxLines = validlines
	params.MaxDaysBefore = validdate
}
