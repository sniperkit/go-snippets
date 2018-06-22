// Package user provides logic and authentication routines for robostats users.
package user

import (
	"code.google.com/p/go.crypto/bcrypt"

	"robostats/errmsg"
	"robostats/storage"
	"time"

	"gopkg.in/mgo.v2/bson"
	"upper.io/db"
)

// User represents system users.
type User struct {
	// Basic properties.
	ID           bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Email        string        `bson:"email" json:"email"`
	Password     string        `bson:"-" json:"password"`
	PasswordHash string        `bson:"password_hash" json:"-"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
	// Session.
	Session *Session `bson:"-" json:"session,omitempty"`
}

const (
	userCollectionName = "user"
)

var (
	// UserCollection is the actual storage reference for users.
	UserCollection db.Collection
)

func init() {
	UserCollection = storage.C(userCollectionName)
}

// Constraint returns the user collection key.
func (u *User) Constraint() db.Cond {
	return db.Cond{
		"email": u.Email,
	}
}

// GetByToken returns the user associated with a given token.
func GetByToken(token string) (u *User, err error) {
	var s *Session

	if s, err = RetrieveSession(token); err != nil {
		return nil, err
	}

	if u, err = GetByID(s.UserID); err != nil {
		return nil, err
	}

	u.Session = s

	return u, nil
}

// GetByID returns an user by the given ID.
func GetByID(id bson.ObjectId) (*User, error) {
	var err error
	var u User

	src := UserCollection.Find(db.Cond{
		"_id": id,
	})

	if c, _ := src.Count(); c < 1 {
		return nil, errmsg.ErrNoSuchUser
	}

	if err = src.One(&u); err != nil {
		return nil, err
	}

	return &u, err
}

// Login exchanges email and password for an user pointer and a new session.
func Login(email, password string) (*User, error) {
	var err error

	u := User{
		Email: email,
	}

	// Attempt to find user.
	res := UserCollection.Find(u.Constraint())

	if c, _ := res.Count(); c < 1 {
		return nil, errmsg.ErrNoSuchUser
	}

	if err = res.One(&u); err != nil {
		return nil, err
	}

	// Verify password.
	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errmsg.ErrPasswordsDoNotMatch
	}

	// Create a new session.
	if u.Session, err = NewSession(u.ID); err != nil {
		return nil, err
	}

	return &u, nil
}

// Create adds a new user to the database.
func (u *User) Create() error {
	var err error

	if u.Password == "" {
		return errmsg.ErrMissingPassword
	}
	if u.Email == "" {
		return errmsg.ErrMissingEmail
	}

	mailExists := UserCollection.Find(db.Cond{
		"email": u.Email,
	})

	if c, _ := mailExists.Count(); c > 0 {
		return errmsg.ErrUserAlreadyexists
	}

	u.CreatedAt = time.Now()

	var passwordHash []byte

	if passwordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost); err != nil {
		return err
	}

	u.PasswordHash = string(passwordHash)
	u.Password = ""

	u.ID = bson.ObjectId("")

	if err = u.save(); err != nil {
		return err
	}

	// Create a new session.
	if u.Session, err = NewSession(u.ID); err != nil {
		return err
	}

	return nil
}

// save updates or appends a user.
func (u *User) save() error {

	if u.ID.Valid() {
		res := UserCollection.Find(db.Cond{
			"_id": u.ID,
		})
		return res.Update(u)
	}

	id, err := UserCollection.Append(u)

	if err != nil {
		return err
	}

	u.ID = id.(bson.ObjectId)

	return nil
}
