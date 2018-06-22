package mailer

import "github.com/mailgun/mailgun-go"
import "fmt"

type Mailer struct {
	Domain       string
	ApiKey       string
	PublicApiKey string
}

type Message struct {
	Heading  string
	Body     string
	Receiver string
}

func New(domain, apiKey, publicApiKey string) *Mailer {
	return &Mailer{Domain: domain, ApiKey: apiKey, PublicApiKey: publicApiKey}
}

func (m *Mailer) SendMessage(msg Message) (string, error) {
	fmt.Println("Hackman-Postmaster <hackman-postmaster@" + m.Domain + ">")
	mg := mailgun.NewMailgun(m.Domain, m.ApiKey, m.PublicApiKey)
	mail := mg.NewMessage(
		"Hackman-Postmaster <hackman-postmaster@"+m.Domain+">",
		msg.Heading,
		msg.Body,
		msg.Receiver,
	)
	_, id, err := mg.Send(mail)
	return id, err
}
