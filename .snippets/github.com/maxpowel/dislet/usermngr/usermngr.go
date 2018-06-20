package usermngr

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha1"
	"crypto/rand"
	"encoding/base64"
	"bytes"
	"fmt"
	"github.com/maxpowel/dislet"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);not null;unique"`
	Password string `gorm:"not null"`
	Salt string `gorm:"not null"`
	Email string `gorm:"type:varchar(50);not null;unique"`
}


func NewUser() *User {
	// Just initialize a random salt
	saltSize := 16
	salt := make([]byte, saltSize)
	rand.Read(salt)
	return &User{
		Salt:base64.URLEncoding.EncodeToString(salt),
	}
}

func PlainPassword(user *User, plainPassword string) (error){
	// Salt must exists. If you created the user with NewUser there will not be any problem
	salt, err := base64.URLEncoding.DecodeString(user.Salt)
	if err != nil {
		return err
	}

	dk := pbkdf2.Key([]byte(plainPassword), salt, 4096, 32, sha1.New)
	user.Password = base64.URLEncoding.EncodeToString(dk)
	return nil
}

func CheckPassword(user *User, plainPassword string) (error) {
	salt, err := base64.URLEncoding.DecodeString(user.Salt)
	if err != nil {
		return err
	}

	password, err := base64.URLEncoding.DecodeString(user.Password)
	if err != nil {
		return err
	}
	dk := pbkdf2.Key([]byte(plainPassword), salt, 4096, 32, sha1.New)
	if bytes.Compare(dk, password) == 0 {
		return nil
	} else {
		return fmt.Errorf("Password does not match")
	}
}


type Manager struct {
	k *dislet.Kernel
}

func (m* Manager) LoadUser(userId uint) (*User, error) {
	db := m.k.Container.MustGet("database").(*gorm.DB)
	user := User{}
	if db.First(&user, userId).RecordNotFound() {
		return nil, fmt.Errorf("User not found")
	} else {
		return &user, nil
	}
}

func (m* Manager) FindUser(username string) (*User, error) {
	db := m.k.Container.MustGet("database").(*gorm.DB)
	user := User{}
	if db.Where(&User{Username: username}).First(&user).RecordNotFound() {
		return nil, fmt.Errorf("User not found")
	} else {
		return &user, nil
	}
}

func Bootstrap(k *dislet.Kernel) {
	//mapping := k.Config.Mapping
	iny := func() *Manager{
		return &Manager{k: k}
	}
	k.Container.RegisterType("user_manager", iny)

	/*var baz dislet.OnKernelReady = func(k *dislet.Kernel){

	}
	k.Subscribe(baz)*/
}
