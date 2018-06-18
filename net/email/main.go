package main

import (
	"fmt"
	"log"
	"net/smtp"
)

type UserInfo struct {
	Username, Password string
}

func exampleSimple() {

	var (
		email    = "user@gmail.com"
		password = "123456789"
	)

	hostname := "smtp.gmail.com"
	auth := smtp.PlainAuth(
		"",
		email,
		password,
		hostname,
	)

	err := smtp.SendMail(
		hostname+":465",
		auth,
		email,
		[]string{"user_to@gmail.com"},
		[]byte("This is the email body"),
	)
	checkErr(err)
}

func exampleDial() {

	var (
		email    = "user@gmail.com"
		password = "123456789"
	)

	hostname := "smtp.gmail.com"

	c, err := smtp.Dial(hostname + ":465")
	checkErr(err)

	err = c.Mail(email)
	checkErr(err)

	err = c.Rcpt("user_to@gmail.com")
	checkErr(err)

	auth := smtp.PlainAuth(
		"",
		email,
		password,
		hostname,
	)

	err = c.Auth(auth)
	checkErr(err)

	wc, err := c.Data()
	checkErr(err)
	_, err = fmt.Fprintf(wc, "This is the email body")
	checkErr(err)

	err = wc.Close()
	checkErr(err)

	err = c.Quit()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	exampleSimple()
	exampleDial()
}
