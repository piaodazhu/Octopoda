package network

import (
	"fmt"
	"os"
	"time"
)

var logFile *os.File

func init() {
	var err error
	logFile, err = os.OpenFile("tentacle_netdebug.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend)
	if err != nil {
		panic(err)
	}
}

func appendLog(s string) {
	logmsg := fmt.Sprintf("%s %s\n", time.Now().Format("01/02 15:04:05"), s)
	logFile.WriteString(logmsg)
	logFile.Sync()
	fmt.Print(logmsg)
}
