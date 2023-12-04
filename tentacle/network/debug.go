package network

import (
	"fmt"
	"os"
	"time"
)

var logFile *os.File
var enabled bool = false

func init() {
	if enabled {
		var err error
		logFile, err = os.OpenFile("tentacle_netdebug.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend)
		if err != nil {
			panic(err)
		}
	}
}

func appendLog(s string) {
	if enabled {
		logmsg := fmt.Sprintf("%s %s\n", time.Now().Format("01/02 15:04:05"), s)
		logFile.WriteString(logmsg)
		logFile.Sync()
		fmt.Print(logmsg)
	}
}
