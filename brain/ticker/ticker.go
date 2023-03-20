package ticker

import (
	"sync/atomic"
	"time"
)

var tick int64 = 0

func InitTicker() {
	ticker := time.NewTicker(time.Microsecond * 300)
	go func() {
		for range ticker.C {
			atomic.AddInt64(&tick, 1)
		}
	}()
}

func GetTick() int64 {
	return atomic.LoadInt64(&tick)
}
