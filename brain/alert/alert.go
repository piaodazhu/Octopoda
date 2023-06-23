package alert

import (
	"brain/config"
	"brain/logger"
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

var mailBuf chan string

func InitAlert() {
	if !config.GlobalConfig.MailAlert.Enabled {
		return
	}
	mailDialer := gomail.NewDialer(config.GlobalConfig.MailAlert.SmtpServer, int(config.GlobalConfig.MailAlert.SmtpPort), config.GlobalConfig.MailAlert.Sender, config.GlobalConfig.MailAlert.Password)
	mailDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	mailTemplate := gomail.NewMessage()
	mailTemplate.SetHeader("From", config.GlobalConfig.MailAlert.Sender)
	mailTemplate.SetHeader("Subject", "Octopoda Alert Message")

	mailBuf = make(chan string, 10)
	go func() {
		for msg := range mailBuf {
			mailTemplate.SetBody("text/plain", msg)
			for _, r := range config.GlobalConfig.MailAlert.Receivers {
				mailTemplate.SetHeader("To", r)
				if err := mailDialer.DialAndSend(mailTemplate); err != nil {
					logger.Exceptions.Printf("cannot mail alert to %s: %s. \nmessage:\n%s", r, err.Error(), msg)
				}
			}
		}
	}()
}

func Alert(msg string) {
	if config.GlobalConfig.MailAlert.Enabled {
		mailBuf <- msg
	} else {
		logger.Exceptions.Print(msg)
	}
}
