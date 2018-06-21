package main

import (
	"errors"
	"fmt"
)

type MailSender interface {
	Send(to, subject, body string) error
	SendFrom(from, to, subject, body string) error
}

type MockSender struct {
	SendFunc     func(to, subject, body string) error
	SendFromFunc func(from, to, subject, body string) error
}

func (m MockSender) Send(to, subject, body string) error {
	return m.SendFunc(to, subject, body)
}
func (m MockSender) SendFrom(from, to, subject, body string) error {
	return m.SendFromFunc(from, to, subject, body)
}

func main() {
	errTest := errors.New("nope")
	var msg string

	sender := MockSender{
		SendFunc: func(to, subject, body string) error {
			fmt.Printf("(%s) %s: %s", to, subject, body)
			return nil
		},
		SendFromFunc: func(from, to, subject, body string) error {
			return errTest
		},
	}

	SendWelcomeEmail(sender, "to", "subject", "body")

	if msg != "(to) subject: body" {
		fmt.Errorf("SendWelcomeEmail: %s", msg)
	}
}

func SendWelcomeEmail(sender MailSender, to, subject, body string) {
	sender.Send(to, subject, body)
}
