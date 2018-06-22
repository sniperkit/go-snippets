package model

import (
	"fmt"
	"time"

	"github.com/asim/go-micro/store"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `gorethink:"id,omitempty", json:"id"`
	Username string `gorethink:"username", json:"username"`
	PassHash string `gorethink:"passhash", json:"-"`
}

func NewUser(username, password string) (*User, error) {
	hp, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return nil, err
	}

	return &User{
		Username: username,
		PassHash: string(hp),
	}, nil
}

type Session struct {
	ID     string `gorethink:"id,omitempty"`
	UserID string `gorethink:"user_id", json:"user_id"`
	Secret string `gorethink:"secret", json:"-"`
}

func (s Session) Token() (tokenString string, err error) {
	// Create the token
	signingKeyID, err := store.Get("auth/signing_key_id")
	if err != nil {
		return "", err
	}

	signingKey, err := store.Get(fmt.Sprintf("auth/signing_key/%s", signingKeyID.Value()))
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["kid"] = string(signingKeyID.Value())
	token.Claims["uid"] = s.UserID
	token.Claims["sec"] = s.Secret
	token.Header["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	tokenString, err = token.SignedString(signingKey.Value())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
