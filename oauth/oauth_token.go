package oauth

import (
	"database/sql"
	"github.com/datatogether/errors"
	"github.com/datatogether/identity/user"
	"github.com/datatogether/sqlutil"
	"golang.org/x/oauth2"
	"time"
)

type UserOauthToken struct {
	User    *user.User
	Service string
	Token   *oauth2.Token
}

func UserOauthTokens(db sqlutil.Queryable, u *user.User) ([]*UserOauthToken, error) {
	res, err := db.Query(qUserOauthTokensForUser, u.Id)
	if err != nil {
		return nil, err
	}

	return UnmarshalTokens(res)
}

func (t *UserOauthToken) ReadUser(db sqlutil.Queryable) (*user.User, error) {
	row, err := db.Query(qUserOauthTokenUser, t.Service, t.Token)
	if err != nil {
		return nil, err
	}

	u := &user.User{}
	if err := u.UnmarshalSQL(row); err != nil {
		return nil, err
	}

	return u, nil
}

// func (t *UserOauthToken) UserService() (OauthUserService, error) {
// 	switch t.Service {
// 	case OauthServiceGithub:
// 		return NewGithub(t.token), nil
// 	default:
// 		return nil, fmt.Errorf("invalid service name")
// 	}
// }

func (t *UserOauthToken) Read(db sqlutil.Queryable) error {
	// first try to read by token id
	if t.Token != nil {
		if err := t.UnmarshalSQL(db.QueryRow(qUserOauthTokenByAccessToken, t.Token.AccessToken)); err == nil {
			return nil
		}
	}
	if t.User == nil || t.Service == "" {
		return errors.ErrNotFound
	}
	return t.UnmarshalSQL(db.QueryRow(qUserOauthTokenByUserAndService, t.User.Id, t.Service))
}

func (t *UserOauthToken) Save(db sqlutil.Execable) error {
	prev := &UserOauthToken{User: t.User, Service: t.Service}
	if err := prev.Read(db); err != nil {
		if err == errors.ErrNotFound {
			if _, err := db.Exec(qUserOauthTokenInsert, t.SQLArgs()...); err != nil {
				// return NewFmtError(500, err.Error())
				return err
			}
		} else {
			return err
		}
	} else {
		if _, err := db.Exec(qUserOauthTokenUpdate, t.SQLArgs()...); err != nil {
			return err
		}
	}
	return nil
}

func (g *UserOauthToken) Delete(db *sql.DB) error {
	if g.User == nil {
		return errors.ErrNotFound
	}
	_, err := db.Exec(qUserOauthTokenDelete, g.User.Id, g.Service)
	return err
}

// turn a sql.Row result into a reset token pointer
func (g *UserOauthToken) UnmarshalSQL(row sqlutil.Scannable) error {
	var (
		userId, service, access, tType, refresh string
		expiry                                  time.Time
	)

	if err := row.Scan(&userId, &service, &access, &tType, &refresh, &expiry); err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrNotFound
		}
		return err
	}

	tok := &UserOauthToken{
		User:    &user.User{Id: userId},
		Service: service,
		Token: &oauth2.Token{
			AccessToken:  access,
			TokenType:    tType,
			RefreshToken: refresh,
			Expiry:       expiry,
		},
	}

	*g = *tok
	return nil
}

func (t *UserOauthToken) SQLArgs() []interface{} {
	if t.Token == nil {
		t.Token = &oauth2.Token{}
	}
	return []interface{}{
		t.User.Id,
		t.Service,
		t.Token.AccessToken,
		t.Token.Type(),
		t.Token.RefreshToken,
		t.Token.Expiry,
	}
}
