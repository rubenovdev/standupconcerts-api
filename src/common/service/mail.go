package service

import (
	"log"
	"os"

	gomail "gopkg.in/mail.v2"
)

var client gomail.Dialer


func init() {
	login := os.Getenv("MAIL_LOGIN")
	password := os.Getenv("MAIL_PASSWORD")

	client = *gomail.NewDialer("smtp.mail.ru", 587, login, password)
}

func SendMail(messageData MessageData) error {
	message := gomail.NewMessage()
	login := os.Getenv("MAIL_LOGIN")

	message.SetHeader("From",  login)
	message.SetHeader("To", messageData.To)
	message.SetHeader("Subject", messageData.Subject)
	message.SetBody("text/plain", messageData.Body)

	if err := client.DialAndSend(message); err != nil {
		log.Panic(err)
		return err
	}

	return nil
}

type MessageData struct {
	To      string
	Subject string
	Body    string
}
