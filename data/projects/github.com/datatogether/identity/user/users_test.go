package user

import "testing"

func TestReadUsers(t *testing.T) {
	users, err := ReadUsers(testDB, UserTypeNone, 20, 0)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(users) <= 0 {
		t.Error("expected more than one user")
	}

	users, err = ReadUsers(testDB, UserTypeNone, 20, 100000)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(users) != 0 {
		t.Error("expected 0 users")
	}
}
