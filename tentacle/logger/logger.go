package logger

import (
	"io"
	"log"
	"tentacle/config"
	"os"
	"strings"
	"time"
)

var Client *log.Logger
var Server *log.Logger

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

	Client = log.New(writer, "[Client]", log.LstdFlags|log.Lshortfile)
	Client.Print("Client Logger Started.")

	Server = log.New(writer, "[Server]", log.LstdFlags|log.Lshortfile)
	Server.Print("Client Logger Started.")

	sb.Reset()
	sb.WriteString(config.GlobalConfig.Logger.Path)
	sb.WriteByte('/')
	sb.WriteString(config.GlobalConfig.Logger.NamePrefix)
	sb.WriteString(time.Now().AddDate(0, 0, -config.GlobalConfig.Logger.RollDays).Format("_2006_01_02.log"))
	if _, err := os.Stat(sb.String()); err == nil {
		os.Remove(sb.String())
	}
}
