package main

import (
	"io/ioutil"
	"log"

	"github.com/matcornic/hermes"
	"github.com/skratchdot/open-golang/open"
)

const previewFile = "hermespreview.html"

// Sample use case:
// Auto preview using python `entr`.
// `find . -name \*.go | entr sh -c "go run hermespreview.go"`
func main() {
	log.Println("Generating email preview..")

	h := hermesInit()
	emailBody, _ := h.GenerateHTML(email(payload()))
	ioutil.WriteFile("preview.html", []byte(emailBody), 0644)
	open.Start(previewFile)
	log.Println("✔ Done")
}

// 1. Configure and initialize hermes
func hermesInit() hermes.Hermes {
	return hermes.Hermes{
		// Optional Theme
		Theme: new(hermes.Default),
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "Hermes",
			Link: "https://example-hermes.com/",
			// Optional product logo
			Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}
}

// 2. Payload set up
type MailPayload struct {
	user_name string
}

func payload() MailPayload {
	return MailPayload{
		"John Doe",
	}
}

// 3. Craft email content
func email(payload MailPayload) hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: payload.user_name,
			Intros: []string{
				"Welcome to Hermes! We're very excited to have you on board.",
			},
			Dictionary: []hermes.Entry{
				{Key: "Firstname", Value: "Jon"},
				{Key: "Lastname", Value: "Snow"},
				{Key: "Birthday", Value: "01/01/283"},
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with Hermes, please click here:",
					Button: hermes.Button{
						Text: "Confirm your account",
						Link: "https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}
}
