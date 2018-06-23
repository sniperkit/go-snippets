package user

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

var (
	fakeUserID bson.ObjectId
	fakeToken  string
)

func init() {
	fakeUserID = bson.NewObjectId()
}

func TestSessionCreate(t *testing.T) {
	// Creating a session.
	s, err := NewSession(fakeUserID)

	if err != nil {
		t.Fatal(err)
	}

	if s == nil {
		t.Fatal("Expecting a valid session.")
	}

	if s.Token == "" {
		t.Fatal("Token must not be empty.")
	}

	if s.TokenHash != hash(s.Token) {
		t.Fatal("Token hash must actually match hash of token.")
	}

	// Filling for later usage.
	fakeToken = s.Token

	// Creating another session
	var z *Session
	z, err = NewSession(fakeUserID)

	if err != nil {
		t.Fatal(err)
	}

	if z.Token == s.Token {
		t.Fatal("Tokens must be different.")
	}
}

func TestSessionRetrieve(t *testing.T) {
	// Getting the same session we created before.
	s, err := RetrieveSession(fakeToken)

	if err != nil {
		t.Fatal(err)
	}

	if s.UserID != fakeUserID {
		t.Fatal("Expecting the same user ID we saved before.")
	}
}

func TestSessionRemove(t *testing.T) {
	var err error
	var s *Session

	// Getting the same session we created before.
	if s, err = RetrieveSession(fakeToken); err != nil {
		t.Fatal(err)
	}

	// Removing it.
	if err = s.Remove(); err != nil {
		t.Fatal(err)
	}

	// Attempt to get the same session (but it does not exists anymore!)
	if _, err = RetrieveSession(fakeToken); err == nil {
		t.Fatal("Expecting an error!")
	}
}
