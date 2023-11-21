package ostp

import "time"

func ExtimateExecTs(maxDelay int64) int64 {
	return time.Now().UnixMilli() + 2*maxDelay
}
