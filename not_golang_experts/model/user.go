package model

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gophergala/not_golang_experts/conf"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id                int64
	Email             string
	EncryptedPassword []byte
	Token             string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func RegisterUser(email string, password string, password_confirmation string, success func(token string), not_success func(message string)) {
	db := conf.SetupDB()

	if userExists(email) {
		not_success("Email has already been taken")
		return
	}

	if passwordsMatch(password, password_confirmation) {
		encrypted_password, _ := bcrypt.GenerateFromPassword([]byte(password), 5)
		user := User{
			Email:             email,
			EncryptedPassword: encrypted_password,
			Token:             generateAuthToken(),
		}

		db.Create(&user)
		success(user.Token)
	} else {
		not_success("Passwords don't match")
	}
}

func RegisterUserSession(email string, password string, success func(token string), not_success func(message string)) {
	db := conf.SetupDB()
	user := User{}
	db.Where("email = ?", email).First(&user)

	if userExists(email) && passwordValid(user.EncryptedPassword, password) {
		user.Token = generateAuthToken()
		db.Save(&user)
		success(user.Token)
	} else {
		not_success("Invalid email or password")
	}
}

func DestroyUserSession(token string, success func(token string), not_success func(message string)) {
	db := conf.SetupDB()
	user := FindUserByAuthToken(token)

	if userExists(user.Email) {
		user.Token = generateAuthToken()
		db.Save(&user)
		success("Successfully logged out user")
	} else {
		not_success("Not found")
	}
}

func FindUserByAuthToken(value string) User {
	db := conf.SetupDB()
	user := User{}
	db.Where("token = ?", value).First(&user)
	return user
}

func passwordsMatch(password string, password_confirmation string) bool {
	return password == password_confirmation && password != ""
}

func userExists(email string) bool {
	db := conf.SetupDB()
	user := User{}
	db.Where("email = ?", email).First(&user)
	return user.Id != 0
}

func generateAuthToken() string {
	var random_string string
	for {
		random_string, _ = generateRandomString(10)
		u := FindUserByAuthToken(random_string)
		if u.Id == 0 {
			break
		}
	}
	return random_string
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func passwordValid(encrypted_password []byte, password string) bool {
	result := bcrypt.CompareHashAndPassword(encrypted_password, []byte(password))
	return result == nil
}
