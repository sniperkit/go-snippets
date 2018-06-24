package user

// import (
// 	"github.com/pborman/uuid"
// 	"testing"
// )

// func TestValidUsername(t *testing.T) {
// 	cases := []struct {
// 		handle string
// 		expect bool
// 	}{
// 		{"bive", true},
// 		{" b5", false},
// 		{"489hwjvwv044378$%%@%$", false},
// 		{"ghj1234123412341234123412341234123asdfa", false},
// 	}

// 	for i, c := range cases {
// 		if got := validUsername(c.handle); got != c.expect {
// 			t.Errorf("case %d failed. %s should be %t.", i, c.handle, c.expect)
// 		}
// 	}
// }

// func TestValidEmail(t *testing.T) {
// 	cases := []struct {
// 		email  string
// 		expect bool
// 	}{
// 		{"b5", false},
// 		{"b5@gmail.com", true},
// 		{"test_email@gmail.co.uk", true},
// 		// {"test_email@gmail.co.uk.", false},
// 	}

// 	for i, c := range cases {
// 		if got := validEmail(c.email); got != c.expect {
// 			t.Errorf("case %d failed. %s should be %t", i, c.email, c.expect)
// 		}
// 	}
// }

// func TestValidSlug(t *testing.T) {
// 	cases := []struct {
// 		slug   string
// 		expect bool
// 	}{
// 		{"b5", true},
// 		{"reallylongslugsshouldbejustfine", true},
// 		{" spacesarenttho", false},
// 		{"spaces arenttho", false},
// 		{"undserscores_are_cool", true},
// 		{"noCapitals", false},
// 		{"test_email@gmail.co.uk", false},
// 	}

// 	for i, c := range cases {
// 		if got := validSlug(c.slug); got != c.expect {
// 			t.Errorf("case %d failed. %s should be %t", i, c.slug, c.expect)
// 		}
// 	}
// }

// func TestValidPath(t *testing.T) {
// 	cases := []struct {
// 		path   string
// 		expect bool
// 	}{
// 		{"b5/", true},
// 		{"mustendinaslash", false},
// 		{"one/two/three/four/five/", true},
// 		// {"one/two/three/four/five//", false},
// 		// {"one//two", false},
// 	}

// 	for i, c := range cases {
// 		if got := validPath(c.path); got != c.expect {
// 			t.Errorf("case %d failed. %s should be %t", i, c.path, c.expect)
// 		}
// 	}
// }

// func TestUsernameTaken(t *testing.T) {
// 	cases := []struct {
// 		handle string
// 		exists bool
// 		err    error
// 	}{
// 		// {TestData.Users.Brendan.Username, true, nil},
// 		{"taken", false, nil},
// 	}

// 	// TODO - restore
// 	// for i, c := range cases {
// 	// exists, err := UsernameTaken(appDB, c.handle)
// 	// if exists != c.exists {
// 	// 	t.Errorf("case %d exist failed. expected to be %t", i, c.exists)
// 	// }
// 	// if exists != c.exists {
// 	// 	t.Errorf("case %d exist failed. expected error: %s, got: %s", i, c.err, err)
// 	// }
// 	// }
// }

// func TestValidUser(t *testing.T) {
// 	cases := []struct {
// 		user *User
// 		err  error
// 	}{
// 		// {TestData.Users.Brendan, nil},
// 		{&User{Id: ""}, ErrInvalidUser},
// 		{&User{Id: uuid.New()}, ErrUserNotFound},
// 	}

// 	// TODO - restore
// 	// for i, c := range cases {
// 	// if got := ValidUser(appDB, c.user); got != c.err {
// 	// 	t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, got)
// 	// }
// 	// }
// }

// func TestValidUrlString(t *testing.T) {
// 	cases := []struct {
// 		rawurl string
// 		result string
// 		err    error
// 	}{
// 		{"apple.com", "http://apple.com", nil},
// 		{"http://localhost:3000", "http://localhost:3000", nil},
// 	}

// 	// TODO - restore
// 	// for i, c := range cases {
// 	// result, err := ValidUrlString(c.rawurl)
// 	// if err != c.err {
// 	// 	t.Errorf("case %d error mismatch. %s != %s", i, c.err, err)
// 	// }
// 	// if result != c.result {
// 	// 	t.Errorf("case %d result mismatch. %s != %s", i, c.result, result)
// 	// }
// 	// }
// }
