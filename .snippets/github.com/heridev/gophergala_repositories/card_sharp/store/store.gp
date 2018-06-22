package store

import (
  "golang.org/x/crypto/bcrypt"
)

type Account table {
  ID int
  Email string
  Management bool

  SecurePassword

  relation {
    []Deck
  }

  index {
    Email
  }
}

type Deck table {
  ID int
  Name string
  Description string

  Private bool
  FullGame bool
  GameType string
  MinPlayer int

  AccountID int
  relation {
    Account
    []Card
  }

  index {
    AccountID
  }
}

type Card table {
  ID int

  Name string
  Type string
  Data string `type:"text"`

  DeckID int
  relation {
    Deck
  }

  index {
    DeckID
  }
}


// mixins
type SecurePassword mixin {
  CryptPassword []byte
}

func (sp *SecurePassword) SetPassword(password string) {
  var err error
  sp.CryptPassword, err = bcrypt.GenerateFromPassword([]byte(password), 0)
  if err != nil {
    log.Println("SetPassword", err)
  }
}

func (sp SecurePassword) ComparePassword(password string) bool {
  if bcrypt.CompareHashAndPassword(sp.CryptPassword, []byte(password)) == nil {
    return true
  }
  return false
}
