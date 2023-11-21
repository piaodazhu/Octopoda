package ostp

import (
	"time"
)

var Delay int64
var TsDiff int64

func EstimateDelay(t1, t2, tB int64) {
	Delay = (t2 - t1) >> 1
	TsDiff = time.Now().UnixMilli() + Delay - tB
	// fmt.Printf("[DBG] Delay=%dms, tsDiff=%dms (now=%d, tb=%d)\n", Delay, TsDiff, time.Now().UnixMilli(), tB)
}

func SleepForExec(execTs int64) bool {
	localExecDefer := TsDiff + execTs - time.Now().UnixMilli()
	// fmt.Println("[DBG] SleepForExec ", localExecDefer)
	if localExecDefer < 0 {
		return false
	}
	if localExecDefer > 1000 {
		// fmt.Println("[DBG] localExecDefer too long: ", localExecDefer)
		return false
	}
	time.Sleep(time.Millisecond * time.Duration(localExecDefer))
	return true
}
