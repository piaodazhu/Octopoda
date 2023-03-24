package logger

import (
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"tentacle/config"
	"time"
)

var Client *log.Logger
var Server *log.Logger
var wg sync.WaitGroup
func InitLogger() {
	var sb strings.Builder
	sb.WriteString(config.GlobalConfig.Logger.Path)
	sb.WriteByte('/')
	sb.WriteString(config.GlobalConfig.Logger.NamePrefix)

	wg.Add(1)
	go logController(sb.String())
	wg.Wait()
}

func logController(prefix string) {
	lastday := time.Now().AddDate(0,0,-1).Day()
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
			if config.Stdout {
				writer = io.MultiWriter(f, os.Stdout)
			} else {
				writer = io.MultiWriter(f)
			}

			Client = log.New(writer, "[Client]", log.LstdFlags|log.Lshortfile)
			Client.Print("Client Logger Started.")

			Server = log.New(writer, "[Server]", log.LstdFlags|log.Lshortfile)
			Server.Print("Client Logger Started.")

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
