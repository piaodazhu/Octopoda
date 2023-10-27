package logger

import (
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/piaodazhu/Octopoda/tentacle/config"
)

var Request *log.Logger
var Comm *log.Logger
var Network *log.Logger
var Exceptions *log.Logger
var SysInfo *log.Logger

var wg sync.WaitGroup

func InitLogger(stdout bool) {
	Request = log.New(nil, "[ Request  ]", log.LstdFlags|log.Lshortfile)
	Comm = log.New(nil, "[   Comm   ]", log.LstdFlags|log.Lshortfile)
	Network = log.New(nil, "[ Network  ]", log.LstdFlags|log.Lshortfile)
	Exceptions = log.New(nil, "[Exceptions]", log.LstdFlags|log.Lshortfile)
	SysInfo = log.New(nil, "[ SysInfo  ]", log.LstdFlags|log.Lshortfile)

	var sb strings.Builder
	sb.WriteString(config.GlobalConfig.Logger.Path)
	sb.WriteByte('/')
	sb.WriteString(config.GlobalConfig.Logger.NamePrefix)

	wg.Add(1)
	go logController(sb.String(), stdout)
	wg.Wait()
}

func logController(prefix string, stdout bool) {
	lastday := time.Now().AddDate(0, 0, -1).Day()
	var lastf *os.File
	once := true

	for {
		if time.Now().Day() == lastday {
			time.Sleep(time.Second)
		} else {
			lastday = time.Now().Day()

			filename := prefix + time.Now().Format("_2006_01_02.log")

			f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
			if err != nil {
				panic("Cannot init logger!" + err.Error())
			}

			var writer io.Writer
			if stdout {
				writer = io.MultiWriter(f, os.Stdout)
			} else {
				writer = io.MultiWriter(f)
			}

			Request.SetOutput(writer)
			Comm.SetOutput(writer)
			Network.SetOutput(writer)
			Exceptions.SetOutput(writer)
			SysInfo.SetOutput(writer)

			SysInfo.Println("Logger Updated")

			lastf.Close()
			lastf = f

			if once {
				wg.Done()
				once = false
			}

			deletefilename := prefix + time.Now().AddDate(0, 0, -config.GlobalConfig.Logger.RollDays).Format("_2006_01_02.log")
			if _, err := os.Stat(deletefilename); err == nil {
				os.Remove(deletefilename)
			}
		}
	}
}
