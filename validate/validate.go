package validate

import (
	"errors"
	"net/mail"

	"github.com/ttacon/libphonenumber"
)

// Error messages
var (
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
	ErrInvalidEmail       = errors.New("invalid email")
)

// Phone validates a phonenumber using libphonenumber
func Phone(num string) error {
	p, err := libphonenumber.Parse(num, "")
	if err != nil {
		return err
	}
	if !libphonenumber.IsValidNumber(p) {
		return ErrInvalidPhoneNumber
	}
	return nil
}

// Email validates an email address using RFC 5322
func Email(addr string) error {
	_, err := mail.ParseAddress(addr)
	return err
}
