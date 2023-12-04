package ostp

import (
	"time"
)

func ExtimateExecTs(maxDelay int64) int64 {
	now := time.Now().UnixMilli()
	res := now + maxDelay + 200
	// logger.Network.Printf("[DBG] max delay is %d and exec ts is now(%d) + %d = %d", maxDelay, now, maxDelay + 200, res)
	return res
}
