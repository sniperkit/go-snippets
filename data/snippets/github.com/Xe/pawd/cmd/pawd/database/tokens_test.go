package database

import "testing"

func TestTokensStorm(t *testing.T) {
	ctx, tdb := newTestingDB(t)
	defer tdb.Close()

	us := NewUsersStorm(tdb.DB, nil)
	tk := NewTokensStorm(tdb.DB, us)

	var (
		normalUserID string
		upTkn        *Token
	)

	u, err := us.Create(ctx, "me@me.me", "hunter2", false)
	if err != nil {
		t.Fatal(err)
	}
	normalUserID = u.ID

	t.Run("create_token", func(t *testing.T) {
		tkn, err := tk.Create(ctx, normalUserID)
		if err != nil {
			t.Fatal(err)
		}
		upTkn = tkn
	})

	t.Run("create_token_invalid", func(t *testing.T) {
		_, err := tk.Create(ctx, "not a user id")
		if err == nil {
			t.Fatal("err should be non-nil")
		}
	})

	t.Run("update_last_seen", func(t *testing.T) {
		err := tk.UpdateLastSeen(ctx, upTkn.ID)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("check", func(t *testing.T) {
		_, err := tk.Check(ctx, upTkn.Body)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("check_invalid", func(t *testing.T) {
		_, err := tk.Check(ctx, "not a token body")
		if err == nil {
			t.Fatal("invalid token body was accepted")
		}
	})
}
