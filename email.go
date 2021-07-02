package main

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func sendEmail(from, to, subject, body, attachmentFilePath string) bool {
	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	m.Attach(attachmentFilePath)

	d := gomail.NewDialer("smtp.gmail.com", 587, "sender-mail-id", "password")

	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
