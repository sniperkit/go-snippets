package main

import (
	"github.com/scorredoira/email"
	"log"
	"net/mail"
	"net/smtp"
	"strings"
)

func main() {
	user := "wuhialin@vip.qq.com"
	host := "smtp.qq.com:25"
	_, err := smtp.Dial(host)
	if nil != err {
		log.Fatalln(err)
	}
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, "pwoqyrrnvoflbfbf", hp[0])
	m := email.NewMessage("test", "test")
	m.From = mail.Address{Name: "wuhialin", Address: "wuhialin@vip.qq.com"}
	m.To = append(m.To, "wuhailin@globalegrow.com")
	m.To = append(m.To, "wuhailin@vip.qq.com")
	err = m.Attach("E:/Data/sql/colourlife.7z")
	if err != nil {
		log.Fatalln(err)
	}
	err = email.Send(host, auth, m)
	if nil != err {
		log.Fatalln(err)
	}
}
