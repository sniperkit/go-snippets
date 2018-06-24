package user

// import (
// 	"encoding/json"
// 	"fmt"
// 	"github.com/datatogether/errors"
// 	"testing"
// )

// var egUserOne = &User{Name: "Brendan O'Brien (test user)", Id: "3fe7d2cc-a8dc-4da0-ac37-c3061d067ae7", Username: "brendan", accessToken: "1234567890ABCDE", Type: UserTypeUser, Email: "test_user_brendan@qri.io"}

// func CompareUsers(a, b *User, strict bool) error {
// 	if a == nil {
// 		return fmt.Errorf("a is nil")
// 	} else if b == nil {
// 		return fmt.Errorf("b is nil")
// 	}

// 	if strict {
// 		if a.Id != b.Id {
// 			return fmt.Errorf("ids don't match")
// 		}
// 		if a.Created != b.Created {
// 			return fmt.Errorf("created doesn't match")
// 		}
// 		if a.Updated != b.Updated {
// 			return fmt.Errorf("updated doesn't match")
// 		}
// 	}
// 	if a.Username != b.Username {
// 		return fmt.Errorf("Username mismatch")
// 	}
// 	if a.Email != b.Email {
// 		return fmt.Errorf("Email mismatch")
// 	}
// 	if a.Name != b.Name {
// 		return fmt.Errorf("Name mismatch: %s != %s", a.Name, b.Name)
// 	}
// 	return nil
// }

// func TestNewUser(t *testing.T) {
// 	if NewUser("one").Id != "one" {
// 		t.Error("NewUser didn't fill in Id field")
// 	}
// }

// func TestNewUserFromString(t *testing.T) {
// 	cases := []struct {
// 		s string
// 		u *User
// 	}{
// 		{egUserOne.Id, &User{Id: egUserOne.Id}},
// 		{egUserOne.Username, &User{Username: egUserOne.Username}},
// 		{"", &User{}},
// 	}

// 	for i, c := range cases {
// 		u := NewUserFromString(c.s)
// 		if err := CompareUsers(u, c.u, true); err != nil {
// 			t.Errorf("case %d user mismatch error: %s", i, err.Error())
// 		}
// 	}
// }

// // func TestUserMarshalJSON(t *testing.T) {
// // 	userBytes, err := json.Marshal(_user(*egUserOne))
// // 	if err != nil {
// // 		t.Errorf("encoding example user bytes failed: %s", err.Error())
// // 	}

// // 	cases := []struct {
// // 		user   *User
// // 		expect string
// // 		err    error
// // 	}{
// // 		{NewUser("asdfghjkl"), `"asdfghjkl"`, nil},
// // 		{egUserOne, string(userBytes), nil},
// // 	}

// // 	for i, c := range cases {
// // 		data, err := c.user.MarshalJSON()
// // 		if string(data) != c.expect {
// // 			t.Errorf("case %d results mismatch. expected: %s, got: %s", i, c.expect, string(data))
// // 		}
// // 		if err != c.err {
// // 			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, err)
// // 		}
// // 	}
// // }

// func TestUserUnmarshalJSON(t *testing.T) {
// 	// var userBytes []byte
// 	// json.Marshal(egUserOne, userBytes)

// 	cases := []struct {
// 		data   []byte
// 		expect *User
// 		err    error
// 	}{
// 		{[]byte(`"` + egUserOne.Id + `"`), NewUser(egUserOne.Id), nil},
// 		{[]byte(`"` + egUserOne.Username + `"`), &User{Username: egUserOne.Username}, nil},
// 		{[]byte(`{ "id": "` + egUserOne.Id + `" }`), NewUser(egUserOne.Id), nil},
// 		// {fipsBytes, TestData.Datasets.Fips, nil},
// 	}

// 	for i, c := range cases {
// 		got := NewUser("")
// 		if err := json.Unmarshal(c.data, &got); err != c.err {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s ", i, c.err, err)
// 		}

// 		if err := CompareUsers(got, c.expect, true); err != nil {
// 			t.Errorf("case %d dataset mismatch. %s ", i, err)
// 		}
// 	}
// }

// func TestUserRead(t *testing.T) {
// 	cases := []struct {
// 		in, out *User
// 		expect  error
// 	}{
// 		{&User{}, nil, errors.ErrNotFound},
// 		{&User{Id: egUserOne.Id}, egUserOne, nil},
// 		{&User{Username: egUserOne.Username}, egUserOne, nil},
// 		{&User{Email: egUserOne.Email}, egUserOne, nil},
// 		{&User{accessToken: egUserOne.accessToken}, egUserOne, nil},
// 	}

// 	for i, c := range cases {
// 		if got := c.in.Read(testDB); c.expect != got {
// 			t.Errorf("case %i error mismatch. expected: %s, got: %s", i, c.expect, got)
// 		}

// 		if c.out != nil {
// 			if err := CompareUsers(c.in, c.out, false); err != nil {
// 				t.Errorf("case %d user mismatch. %s", i, err)
// 			}
// 		}
// 	}
// }

// func TestUserSave(t *testing.T) {
// 	cases := []struct {
// 		u      *User
// 		expect error
// 	}{
// 		{&User{Username: ""}, errors.ErrUsernameRequired},
// 		{&User{Username: "#@%!@%!#"}, errors.ErrInvalidUsername},
// 		// TODO - email no longer required b/c oauth
// 		// {&User{Username: "test_user", Email: " "}, errors.ErrEmailRequired},

// 		// TODO
// 		// {&User{Username: "test_user", Email: "Cap@@tian.email@stuf.stuff.com"}, errors.ErrInvalidEmail},
// 		{&User{Username: egUserOne.Username, Email: "test@qri.io"}, errors.ErrUsernameTaken},
// 		// {&User{Username: CaCensus.Username, Email: "test@qri.io"}, errors.ErrUsernameTaken},

// 		// TODO - password no longer required b/c oauth
// 		// {&User{Username: "test_user", Email: "test@qri.io", password: " "}, errors.ErrPasswordRequired},
// 		// {&User{Username: "test_user", Email: "test@qri.io", password: "13"}, errors.ErrPasswordTooShort},

// 		// TODO
// 		// {&User{Username: "test_user", Email: "test@qri.io", password: "13 asdfj90e9 9a0"}, ErrInvalidPassword},
// 		{&User{Username: "test_user", Email: "test@qri.io", password: "password"}, nil},
// 		{&User{Id: egUserOne.Id, Name: "test!"}, nil},
// 		{&User{Id: egUserOne.Id, Name: egUserOne.Name}, nil},

// 		{&User{Id: egUserOne.Id, Username: "b5"}, nil},
// 		{&User{Id: egUserOne.Id, Email: "bob@qri.io"}, nil},
// 	}

// 	for i, c := range cases {
// 		if got := c.u.Save(testDB); got != c.expect {
// 			t.Errorf("case %d unexpected error. expected: %s, got: %s", i, c.expect, got)
// 		}
// 	}
// }

// func TestUserSavePassword(t *testing.T) {
// 	cases := []struct {
// 		u      *User
// 		pass   string
// 		expect error
// 	}{
// 		{egUserOne, "", errors.ErrPasswordRequired},
// 		{egUserOne, "asd", errors.ErrPasswordTooShort},
// 		{egUserOne, "98hbuqw9cbq9wc0wae0dE", nil},
// 	}

// 	for i, c := range cases {
// 		c.u.password = c.pass
// 		if got := c.u.SavePassword(testDB, c.pass); got != c.expect {
// 			t.Errorf("case %d unexpected error. expected: %s, got: %s", i, c.expect, got)
// 		}
// 	}

// 	// resetTestData(testDB, TypeUser)
// }

// func TestUserDelete(t *testing.T) {
// 	// TODO - make tests independent
// 	// read test_user from save test
// 	testUser := &User{Username: "test_user", Email: "test@qri.io", password: "password"}
// 	if err := testUser.Read(testDB); err != nil {
// 		t.Error("couldn't save test_user")
// 		return
// 	}

// 	cases := []struct {
// 		u      *User
// 		expect error
// 	}{
// 		{&User{}, errors.ErrNotFound},
// 		{testUser, nil},
// 	}

// 	for i, c := range cases {
// 		if got := c.u.Delete(testDB); got != c.expect {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s", c.expect, got)
// 		}

// 		// try to read, should return not found
// 		if err := c.u.Read(testDB); err != errors.ErrNotFound {
// 			t.Errorf("case %d didn't return not found on read after delete", i)
// 		}
// 	}
// }

// func TestCreateUser(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping CreateUser In Short Mode")
// 	}
// 	cases := []struct {
// 		username, email, name, password string
// 		t                               UserType
// 		u                               *User
// 		expect                          error
// 	}{
// 		{"test_user_create", "email@gmail.com", "word", "testUserPassw0rd", UserTypeUser, &User{Username: "test_user_create", Email: "email@gmail.com", Name: "word"}, nil},
// 	}

// 	for i, c := range cases {
// 		// username, email, name, password string, t UserType
// 		if user, got := CreateUser(testDB, c.username, c.email, c.name, c.password, c.t); got != c.expect {
// 			t.Errorf("cases %i error mismatch. expected: %s, got: %s", i, c.expect, got)
// 		} else if c.u != nil {
// 			if err := CompareUsers(c.u, user, false); err != nil {
// 				t.Errorf("case %d user mismatch. %s", i, err)
// 			}
// 		}
// 	}
// }

// func TestAuthenticateUser(t *testing.T) {
// 	cases := []struct {
// 		username, password string
// 		u                  *User
// 		expect             error
// 	}{
// 		{"not_a_user", "testUserPassw0rd", &User{Username: "test_user_create", Email: "email@gmail.com", Name: "word"}, errors.ErrAccessDenied},
// 		{"test_user_create", "wrongPassword", &User{Username: "test_user_create", Email: "email@gmail.com", Name: "word"}, errors.ErrAccessDenied},
// 		{"test_user_create", "testUserPassw0rd", &User{Username: "test_user_create", Email: "email@gmail.com", Name: "word"}, nil},
// 	}

// 	for i, c := range cases {
// 		if user, got := AuthenticateUser(testDB, c.username, c.password); got != c.expect {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.expect, got)
// 		} else if c.u != nil && c.expect == nil {
// 			if err := CompareUsers(c.u, user, false); err != nil {
// 				t.Errorf("case %d user mismatch. %s", i, err)
// 			}
// 		}
// 	}
// }
