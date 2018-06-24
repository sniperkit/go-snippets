package reset_token

import (
	"database/sql"
	"fmt"
	"github.com/datatogether/errors"
	"github.com/datatogether/identity/user"
	"github.com/datatogether/sqlutil"
	"github.com/pborman/uuid"
	"strings"
	"time"
)

// a user reset token
type ResetToken struct {
	Id      string `json:"-" sql:"id"`
	Created int64  `json:"-" sql:"created"`
	Updated int64  `json:"-" sql:"updated"`
	Email   string `json:"-" sql:"email"`
	Used    bool   `json:"-" sql:"used"`
}

func resetTokenColumns() string {
	return "id, created, updated, email, used"
}

// create a reset token
func CreateResetToken(db sqlutil.Execable, email string) (*ResetToken, error) {
	now := time.Now().Unix()
	t := &ResetToken{
		Id:      uuid.New(),
		Created: now,
		Updated: now,
		Email:   email,
		Used:    false,
	}

	if err := t.validate(db); err != nil {
		return t, err
	}

	if _, err := db.Exec("INSERT INTO reset_tokens VALUES ($1, $2, $3, $4, $5)", t.Id, t.Created, t.Updated, t.Email, t.Used); err != nil {
		return t, errors.NewFmtError(500, err.Error())
	}

	return t, nil
}

// validate a reset token
func (r *ResetToken) validate(db sqlutil.Queryable) error {
	r.Email = strings.TrimSpace(r.Email)
	if r.Email == "" {
		return errors.ErrEmailRequired
	}
	// if !emailRegex.MatchString(r.Email) {
	// 	return ErrInvalidEmail
	// }

	var exists bool
	if err := db.QueryRow("SELECT exists(SELECT 1 FROM users WHERE email = $1)", r.Email).Scan(&exists); err != nil {
		return errors.New500Error(err.Error())
	} else if !exists {
		return errors.ErrEmailDoesntExist
	}

	return nil
}

// read a token
func (t *ResetToken) Read(db *sql.DB) error {
	if t.Id == "" {
		return errors.ErrNotFound
	}

	token, err := serializeResetToken(db.QueryRow(fmt.Sprintf("SELECT %s FROM reset_tokens WHERE id=$1", resetTokenColumns()), t.Id))
	if err != nil {
		return err
	}

	*t = *token
	return nil
}

// returns nil if the token is usable, otherwise
// returns an error
func (r *ResetToken) Usable() error {
	if r.Used {
		return errors.ErrTokenAlreadyUsed
	}
	// tokens expire after two days, or if a created date isn't found
	// TODO - should no created value return not found?
	if r.Created == 0 || time.Now().Sub(time.Unix(r.Created, 0)) > time.Hour*48 {
		return errors.ErrTokenExpired
	}

	return nil
}

// use the token to reset the user's password, returning the updated user
func (r *ResetToken) Consume(db sqlutil.Execable, password string) (*user.User, error) {
	if err := r.Usable(); err != nil {
		return nil, err
	}

	u := user.NewUser("")
	u.Email = r.Email
	if err := u.Read(db); err != nil {
		return u, err
	}

	// u.password = password
	if err := u.SavePassword(db, password); err != nil {
		return u, err
	}

	r.Updated = time.Now().Unix()
	r.Used = true
	if _, err := db.Exec("UPDATE reset_tokens SET updated=$2, used=true WHERE id=$1", r.Id, r.Updated); err != nil {
		return u, errors.New500Error(err.Error())
	}

	return u, nil
}

// turn a sql.Row result into a reset token pointer
func serializeResetToken(row sqlutil.Scannable) (*ResetToken, error) {
	var (
		id, email        string
		created, updated int64
		used             bool
	)

	if err := row.Scan(&id, &created, &updated, &email, &used); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	t := &ResetToken{
		Id:      id,
		Created: created,
		Updated: updated,
		Email:   email,
		Used:    used,
	}

	return t, nil
}
