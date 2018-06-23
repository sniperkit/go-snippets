package database

import (
	"testing"
)

func TestUsersStorm(t *testing.T) {
	ctx, tdb := newTestingDB(t)
	defer tdb.Close()

	us := NewUsersStorm(tdb.DB, []string{"azurediamond@itsonlystarsto.me"})

	var (
		normalUserID string
	)

	t.Run("create_user_good", func(t *testing.T) {
		u, err := us.Create(ctx, "me@me.me", "hunter2", false)
		if err != nil {
			t.Fatal(err)
		}
		normalUserID = u.ID
	})

	t.Run("create_user_bad_totp", func(t *testing.T) {
		_, err := us.Create(ctx, "", "", true)
		if err == nil {
			t.Fatal("expected err to be non-nil, got nil")
		}
	})

	t.Run("create_admin_user", func(t *testing.T) {
		u, err := us.Create(ctx, "azurediamond@itsonlystarsto.me", "hunter2", false)
		if err != nil {
			t.Fatal(err)
		}

		if u.Admin == false {
			t.Fatal("expected user to be admin but they are not")
		}
	})

	t.Run("update_last_seen", func(t *testing.T) {
		err := us.UpdateLastSeen(ctx, normalUserID)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("check_password_valid", func(t *testing.T) {
		u, err := us.CheckPassword(ctx, "me@me.me", "hunter2", "")
		if err != nil {
			t.Fatal(err)
		}
		if u == nil {
			t.Fatal("password mismatch")
		}
	})

	t.Run("check_password_invalid", func(t *testing.T) {
		u, err := us.CheckPassword(ctx, "me@me.me", "transformers are cool", "")
		if err == nil {
			t.Fatal(err)
		}
		if u != nil {
			t.Fatal("password should not have matched but did")
		}
	})

	t.Run("get_user_valid", func(t *testing.T) {
		_, err := us.Get(ctx, normalUserID)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get_user_invalid", func(t *testing.T) {
		_, err := us.Get(ctx, "not a user id lol")
		if err == nil {
			t.Fatal(err)
		}
	})
}
