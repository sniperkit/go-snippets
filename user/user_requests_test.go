package user

// import "testing"

// func TestUsersRequest(t *testing.T) {
// 	cases := []struct {
// 		req *UsersRequest
// 		err error
// 	}{
// 		{&UsersRequest{User: TestData.Users.Brendan, Page: Page{1, 50}}, nil},
// 	}

// 	for i, c := range cases {
// 		if _, err := c.req.Exec(); err != c.err {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, err)
// 		}
// 	}
// }

// func TestUserRequest(t *testing.T) {
// 	cases := []struct {
// 		req *UserRequest
// 		err error
// 	}{
// 		{&UserRequest{User: TestData.Users.Brendan, Subject: &User{Username: TestData.Users.Janelle.Username}}, nil},
// 		{&UserRequest{User: TestData.Users.Brendan, Subject: &User{Id: TestData.Users.Janelle.Id}}, nil},
// 	}

// 	for i, c := range cases {
// 		if _, err := c.req.Exec(); err != c.err {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, err)
// 		}
// 	}
// }

// func TestCreateUserRequest(t *testing.T) {
// 	cases := []struct {
// 		req *CreateUserRequest
// 		err error
// 	}{
// 		{&CreateUserRequest{User: TestData.Users.Carl}, nil},
// 	}

// 	for i, c := range cases {
// 		if _, err := c.req.Exec(); err != c.err {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, err)
// 		}
// 	}
// }

// func TestSaveUserRequest(t *testing.T) {
// 	cases := []struct {
// 		req *SaveUserRequest
// 		err error
// 	}{
// 		{&SaveUserRequest{User: TestData.Users.Brendan, Subject: TestData.Users.Brendan}, nil},
// 	}

// 	for i, c := range cases {
// 		if _, err := c.req.Exec(); err != c.err {
// 			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, err)
// 		}
// 	}
// }
