package passwords

import (
	"fmt"

	passlib "gopkg.in/hlandau/passlib.v1"
)

//Hash will hash the password
func Hash(password string) (string, error) {
	hash, err := passlib.Hash(password)
	if err != nil {
		return "", fmt.Errorf("Unable hash password, error: %s", err.Error())
	}

	return hash, nil
}
