package router

import (
	"encoding/json"
	"github.com/gophergala/not_golang_experts/model"
	"io"
	"net/http"
)

type Registration struct {
	User UserRegistration `json:"user"`
}

type UserRegistration struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func RegisterSession(res http.ResponseWriter, req *http.Request) {
	email, password, password_confirmation, err := parseRegistrationsRequest(req.Body)

	PanicIf(err, res)

	model.RegisterUser(email, password, password_confirmation, func(token string) {
		respondWith(map[string]string{"token": token}, 201, res)
	}, func(message string) {
		respondWith(map[string]string{"error": message}, 422, res)
	})
}

func parseRegistrationsRequest(body io.Reader) (string, string, string, error) {
	registration := Registration{}

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&registration)

	return registration.User.Email, registration.User.Password, registration.User.PasswordConfirmation, err
}
