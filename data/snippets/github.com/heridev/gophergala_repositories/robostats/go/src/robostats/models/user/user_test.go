package user

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var (
	fakeUserEmail    string
	fakeUserPassword string
)

func init() {
	rand.Seed(time.Now().UnixNano())

	// Using a random e-mail to reduce the chance of hitting an existent e-mail.
	fakeUserEmail = fmt.Sprintf("user-%d@example.com", rand.Int63())
	fakeUserPassword = fmt.Sprintf("password-%s", rand.Int63())
}

func TestUserCreateAccount(t *testing.T) {
	u := User{
		Email:    fakeUserEmail,
		Password: fakeUserPassword,
	}

	// Attempt to create user.
	if err := u.Create(); err != nil {
		t.Fatal(err)
	}

	// New user must have a new ID.
	if u.ID.Valid() == false {
		t.Fatal("Expecting an ID.")
	}

	// Plain password must not be visible.
	if u.Password != "" {
		t.Fatal("Passwor dmust not be visible.")
	}

	// User must have a session.
	if u.Session == nil {
		t.Fatal("A new user must have a session.")
	}

	fakeUserID = u.ID
}

func TestUserRetrievebyID(t *testing.T) {
	var u *User
	var err error

	// Attempt to retrieve an user by ID.
	if u, err = GetByID(fakeUserID); err != nil {
		t.Fatal(err)
	}

	if u.ID != fakeUserID {
		t.Fatal("Expecting the IDs to be equal.")
	}

	if u.Email != fakeUserEmail {
		t.Fatal("Expecting the e-mails to be equal.")
	}

	if u.Password != "" {
		t.Fatal("Password can't contain anything.")
	}

	if u.Session != nil {
		t.Fatal("Session must be nil.")
	}

}

func TestUserLogin(t *testing.T) {
	var u *User
	var err error

	// Attempt to retrieve an user by login and password.
	if u, err = Login(fakeUserEmail, fakeUserPassword); err != nil {
		t.Fatal(err)
	}

	if u.ID != fakeUserID {
		t.Fatal("Expecting the IDs to be equal.")
	}

	if u.Email != fakeUserEmail {
		t.Fatal("Expecting the e-mails to be equal.")
	}

	if u.Password != "" {
		t.Fatal("Password can't contain anything.")
	}

	if u.Session == nil {
		t.Fatal("Session must NOT be nil.")
	}

	if u.Session.Token == "" {
		t.Fatal("Missing session token.")
	}
}

func TestUserTruncateAndFillup(t *testing.T) {

	UserCollection.Truncate()

	u := &User{
		Email:    "user@example.com",
		Password: "pass",
	}

	if err := u.Create(); err != nil {
		t.Fatal(err)
	}
}
