package model

import (
	"crypto/md5"
	"time"
	"fmt"
)

type User struct {
	Model

	Username string `gorm:"unique_index;size:20;not null;default:''"`
	Phone    string `gorm:"unique_index;size:15;not null;default:''"`
	Email    string `gorm:"unique_index;size:30;not null;default:''"`
	Password string `gorm:"size:32;not null;default:''"`
	Salt     string `gorm:"size:32;not null"`
}

func (t *User) GeneratePassword(pass string) string {
	data := []byte(t.Salt + pass)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func (t *User) GenerateRandom() string {
	b := md5.Sum([]byte(fmt.Sprintf("%d", time.Now().Unix())))
	s := fmt.Sprintf("%x", b)
	return s
}

func (t *User) CheckPassword(pass string) bool {
	return t.Password == t.GeneratePassword(pass)
}
