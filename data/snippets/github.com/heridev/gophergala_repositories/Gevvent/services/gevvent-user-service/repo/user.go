package repo

import (
	r "github.com/dancannon/gorethink"

	"github.com/gophergala/Gevvent/services/gevvent-lib/rethinkdb"
	"github.com/gophergala/Gevvent/services/gevvent-user-service/model"
)

func CreateUser(user *model.User) (*model.User, error) {
	s, err := rethinkdb.Session()
	if err != nil {
		return nil, err
	}
	response, err := r.Table("users").Insert(user).RunWrite(s)
	if err != nil {
		return nil, err
	}

	// Find new ID of product if needed
	if user.ID == "" && len(response.GeneratedKeys) == 1 {
		user.ID = response.GeneratedKeys[0]
	}

	return user, nil
}

func GetUser(id string) (*model.User, error) {
	var user *model.User

	s, err := rethinkdb.Session()
	if err != nil {
		return nil, err
	}
	row, err := r.Table("users").Get(id).Run(s)
	if err != nil {
		return nil, err
	}

	if row.IsNil() {
		return nil, nil
	}

	if err = row.One(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func GetUsername(username string) (*model.User, error) {
	var user *model.User

	s, err := rethinkdb.Session()
	if err != nil {
		return nil, err
	}
	row, err := r.Table("users").GetAllByIndex("username", username).Run(s)
	if err != nil {
		return nil, err
	}

	if row.IsNil() {
		return nil, nil
	}

	if err = row.One(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func FindSession(userID string, secret string) (*model.Session, error) {
	s, err := rethinkdb.Session()
	if err != nil {
		return nil, err
	}

	var session *model.Session
	row, err := r.Table("sessions").GetAllByIndex("session", []interface{}{userID, secret}).Run(s)
	if err != nil {
		return session, err
	}

	if err = row.One(&session); err != nil {
		return nil, err
	}

	return session, nil
}

func CreateSession(session *model.Session) (*model.Session, error) {
	s, err := rethinkdb.Session()
	if err != nil {
		return nil, err
	}

	response, err := r.Table("sessions").Insert(session).RunWrite(s)
	if err != nil {
		return nil, err
	}

	// Find new ID of product if needed
	if session.ID == "" && len(response.GeneratedKeys) == 1 {
		session.ID = response.GeneratedKeys[0]
	}

	return session, nil
}

func DeleteSession(id string) error {
	s, err := rethinkdb.Session()
	if err != nil {
		return err
	}

	return r.Table("sessions").Get(id).Delete().Exec(s)
}
