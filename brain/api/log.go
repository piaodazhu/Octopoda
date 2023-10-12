package api

import (
	"brain/config"
	"brain/message"
	"brain/model"
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LogParams struct {
	MaxLines      int
	MaxDaysBefore int
	Logs          []string
}

func NodeLog(ctx *gin.Context) {
	name, ok := ctx.GetQuery("name")
	if !ok {
		ctx.Status(404)
		return
	}
	maxlines, maxdaysbefore := 0, 0
	l, ok := ctx.GetQuery("maxlines")
	if !ok {
		maxlines = 30
	} else {
		maxlines, _ = strconv.Atoi(l)
	}
	d, ok := ctx.GetQuery("maxdaysbefore")
	if !ok {
		maxdaysbefore = 0
	} else {
		maxdaysbefore, _ = strconv.Atoi(d)
	}
	lparams := LogParams{
		MaxLines:      maxlines,
		MaxDaysBefore: maxdaysbefore,
		Logs:          []string{},
	}

	if name == "brain" {
		readLogs(&lparams)
	} else {
		ok := readLogsRemote(name, &lparams)
		if !ok {
			ctx.Status(404)
			return
		}
	}
	ctx.JSON(200, lparams)
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

func readLogsRemote(name string, params *LogParams) bool {
	query, _ := config.Jsoner.Marshal(params)
	answer, err := model.Request(name, message.TypeNodeLog, query)
	if err != nil {
		return false
	}
	err = config.Jsoner.Unmarshal(answer, params)
	return err == nil
}
