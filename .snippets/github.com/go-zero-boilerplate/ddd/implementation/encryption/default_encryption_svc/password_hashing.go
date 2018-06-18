package default_encryption_svc

//Code obtained from https://trihackeat.wordpress.com/2014/10/11/758/

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

const (
	cDEFAULT_SALT_LENGTH  = 64
	CDEFAULT_ENCRYPT_COST = 10
)

// This is returned when a new hash + salt combo is generated
type password struct {
	hash []byte
	salt []byte
}

func (p *password) GetHashedPasswordHex() string {
	return hex.EncodeToString(p.hash)
}

func (p *password) GetSaltHex() string {
	return hex.EncodeToString(p.salt)
}

// Handles create a new hash/salt combo from a raw password as inputted
// by the user
func CreatePassword(raw_pass string, saltLength, encryptCost int) (*password, error) {
	password := new(password)
	salt, err := generateSalt(saltLength)
	if err != nil {
		return nil, err
	}
	password.salt = salt

	salted_pass := combine(password.salt, []byte(raw_pass))
	hash, err := hashPassword(salted_pass, encryptCost)
	if err != nil {
		return nil, err
	}
	password.hash = hash

	return password, nil
}

func CreatePassword_DefaultComplexity(raw_pass string) (*password, error) {
	return CreatePassword(raw_pass, cDEFAULT_SALT_LENGTH, CDEFAULT_ENCRYPT_COST)
}

// this handles taking a raw user password and making in into something safe for
// storing in our DB
func hashPassword(salted_pass []byte, encryptCost int) ([]byte, error) {
	hashed_pass, err := bcrypt.GenerateFromPassword(salted_pass, encryptCost)
	if err != nil {
		return nil, err
	}
	return hashed_pass, nil
}

// Handles merging together the salt and the password
func combine(salt []byte, raw_pass []byte) []byte {
	// concat salt + password
	combined := []byte{}
	combined = append(combined, raw_pass...)
	combined = append(combined, salt...)
	return combined
}

// Generates a random salt using DevNull
func generateSalt(saltLength int) ([]byte, error) {
	// Read in data
	data := make([]byte, saltLength)
	_, err := rand.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Checks whether or not the correct password has been provided
func passwordMatch(guess string, salt, hashedPassword []byte) bool {
	salted_guess := combine(salt, []byte(guess))

	// compare to the real deal
	if bcrypt.CompareHashAndPassword(hashedPassword, salted_guess) != nil {
		return false
	}

	return true
}

func PasswordMatchHex(guess string, hexSalt, hexHashedPassword string) (bool, error) {
	salt, err := hex.DecodeString(hexSalt)
	if err != nil {
		return false, err
	}
	hashedPassword, err := hex.DecodeString(hexHashedPassword)
	if err != nil {
		return false, err
	}
	return passwordMatch(guess, salt, hashedPassword), nil
}
