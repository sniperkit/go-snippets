package store

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID         int
	Email      string
	Management bool

	// mixins

	CryptPassword []byte
	cached_conn   *Conn
}

type Deck struct {
	ID          int
	Name        string
	Description string

	Private   bool
	FullGame  bool
	GameType  string
	MinPlayer int

	AccountID   int
	cached_conn *Conn
}

type Card struct {
	ID int

	Name string
	Type string
	Data string `type:"text"`

	DeckID      int
	cached_conn *Conn
}

func (sp *Account) SetPassword(password string) {
	var err error
	sp.CryptPassword, err = bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		log.Println("SetPassword", err)
	}
}

func (sp Account) ComparePassword(password string) bool {
	if bcrypt.CompareHashAndPassword(sp.CryptPassword, []byte(password)) == nil {
		return true
	}
	return false
}
