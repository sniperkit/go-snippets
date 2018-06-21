package main

import (
	"net/smtp"
	"strings"
)

func sendMail(to, subject, body, mailType string) (err error) {
	user := "wuhialin@vip.qq.com"
	host := "smtp.qq.com:25"
	_, err = smtp.Dial(host)
	if nil != err {
		return
	}
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, "pwoqyrrnvoflbfbf", hp[0])
	contentType := "Content-Type: text/plain" + "; charset=UTF-8"
	if "html" == mailType {
		contentType = "Content-Type: text/html" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ";")
	err = smtp.SendMail(host, auth, user, sendTo, msg)
	return err
}

func main() {
	err := sendMail("wuhailin@globalegrow.com", "test", "test", "")
	if nil != err {
		panic(err)
	}
}
