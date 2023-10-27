package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/piaodazhu/Octopoda/pakma/config"
)

type UpdatorHistoryItem struct {
	Timestamp int64
	Message   string
}
type UpdatorHistory []UpdatorHistoryItem

var history UpdatorHistory

const timefmt string = "2006-01-02@15:04:05"

func WriteHistory(format string, a ...interface{}) {
	history = append(history, UpdatorHistoryItem{
		Timestamp: time.Now().Unix(),
		Message:   fmt.Sprintf(format, a...),
	})
	dumpHistory()
}

func SearchHistory(ts int64, limit int) []UpdatorHistoryItem {
	idx := sort.Search(len(history), func(i int) bool {
		return history[i].Timestamp >= ts
	})
	if idx == len(history) {
		return history[max(0, len(history)-limit):] // last nlimit records
	}
	return history[idx:min(idx+limit, len(history))] // nlimit records after idx
}

func GetHistoryHandler(ctx *gin.Context) {
	timestr := ctx.Query("time")
	limit := ctx.Query("limit")
	var err error
	var timestamp int64
	if len(timestr) == 0 {
		timestamp = time.Now().Unix()
	} else {
		t, err := time.Parse(timefmt, timestr)
		if err != nil {
			timestamp = time.Now().Unix()
		} else {
			timestamp = t.Unix()
		}
	}

	var nlimit int
	if len(limit) == 0 {
		nlimit = 10
	} else {
		nlimit, err = strconv.Atoi(limit)
		if err != nil {
			nlimit = 10
		}
	}

	history := SearchHistory(timestamp, nlimit)

	res := Response{
		Msg:         "OK",
		HistoryList: []string{},
	}
	for _, h := range history {
		res.HistoryList = append(res.HistoryList, fmt.Sprintf("[%s]: %s", time.Unix(h.Timestamp, 0).Format(timefmt), h.Message))
	}
	ctx.JSON(200, res)
}

func min[T int | int64 | uint | uint64 | float32 | float64](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func max[T int | int64 | uint | uint64 | float32 | float64](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func dumpHistory() {
	serialized, _ := json.Marshal(history)
	if fileExists(config.GlobalConfig.Packma.Root + "pakma_history.json") {
		os.Rename(config.GlobalConfig.Packma.Root+"pakma_history.json", config.GlobalConfig.Packma.Root+"pakma_history.json.bk")
	}
	os.WriteFile(config.GlobalConfig.Packma.Root+"pakma_history.json", serialized, os.ModePerm)
}

func loadHistory() error {
	if !fileExists(config.GlobalConfig.Packma.Root + "pakma_history.json") {
		return fmt.Errorf("pakma_history.json not found")
	}
	serialized, err := os.ReadFile(config.GlobalConfig.Packma.Root + "pakma_history.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(serialized, &history)
	if err != nil {
		return err
	}
	return nil
}
