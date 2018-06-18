package default_encryption_svc

import (
	"github.com/go-zero-boilerplate/ddd/interfaces/encryption"
)

type defaultEncryption struct{}

func (d *defaultEncryption) PasswordMatchHex(guessedPassword, hexSalt, hexHashedPassword string) (bool, error) {
	return PasswordMatchHex(guessedPassword, hexSalt, hexHashedPassword)
}

func (d *defaultEncryption) CreatePassword(raw_pass string) (encryption.Password, error) {
	return CreatePassword_DefaultComplexity(raw_pass)
}

func New() encryption.Service {
	return &defaultEncryption{}
}
