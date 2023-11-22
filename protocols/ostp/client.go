package ostp

import (
	"time"

	"github.com/piaodazhu/Octopoda/tentacle/logger"
)

var Delay int64
var TsDiff int64

var delayAverager, tsDiffAverager *averager

func init() {
	delayAverager = newAverager(4)
	tsDiffAverager = newAverager(60)
}

func EstimateDelay(t1, t2, tB int64) {
	thisDelay := (t2 - t1) >> 1
	now := time.Now().UnixMilli()
	thisTsDiff := now + thisDelay - tB

	Delay = delayAverager.Update(thisDelay)
	TsDiff = tsDiffAverager.Update(thisTsDiff)

	// logger.Exceptions.Printf("[DBG] avr tsDiff=%dms (now=%d, avr delay=%d, tb=%d)\n", TsDiff, now, Delay, tB)
}

func SleepForExec(execTs int64) bool {
	now := time.Now().UnixMilli()
	localExecDefer := TsDiff + execTs - now
	// fmt.Println("[DBG] SleepForExec ", localExecDefer)
	if localExecDefer < 0 {
		logger.Exceptions.Printf("localExecDefer %d < 0. (TsDiff=%d, execTs=%d, now=%d)", localExecDefer, TsDiff, execTs, now)
		return false
	}
	if localExecDefer > 3000 {
		logger.Exceptions.Printf("localExecDefer %d too long. (TsDiff=%d, execTs=%d, now=%d)", localExecDefer, TsDiff, execTs, now)
		return false
	}
	time.Sleep(time.Millisecond * time.Duration(localExecDefer))
	return true
}
