package passwords

import (
	"fmt"

	passlib "gopkg.in/hlandau/passlib.v1"
)

//Verify will verify the password
func Verify(password, hashedPassword string) error {
	err := passlib.VerifyNoUpgrade(password, hashedPassword)
	if err != nil {
		return fmt.Errorf("Password verification failed, error: %s", err.Error())
	}
	return nil
}
