package logger

import (
	"brain/config"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var Tentacle *log.Logger
var Brain *log.Logger

func InitLogger() {
	var sb strings.Builder
	sb.WriteString(config.GlobalConfig.Logger.Path)
	sb.WriteByte('/')
	sb.WriteString(config.GlobalConfig.Logger.NamePrefix)
	sb.WriteString(time.Now().Format("_2006_01_02.log"))
	filename := sb.String()

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

	Tentacle = log.New(writer, "[Tentacle]", log.LstdFlags|log.Lshortfile)
	Tentacle.Print("Tentacle Logger Started.")

	Brain = log.New(writer, "[Brain]", log.LstdFlags|log.Lshortfile)
	Brain.Print("Brain Logger Started.")

	sb.Reset()
	sb.WriteString(config.GlobalConfig.Logger.Path)
	sb.WriteByte('/')
	sb.WriteString(config.GlobalConfig.Logger.NamePrefix)
	sb.WriteString(time.Now().AddDate(0, 0, -config.GlobalConfig.Logger.RollDays).Format("_2006_01_02.log"))
	if _, err := os.Stat(sb.String()); err == nil {
		os.Remove(sb.String())
	}
}
