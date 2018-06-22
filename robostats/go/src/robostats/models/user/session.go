package user

import (
	"crypto/sha1"
	"fmt"
	"robostats/errmsg"
	"robostats/storage"
	"time"

	"gopkg.in/mgo.v2/bson"
	"upper.io/db"
	"upper.io/i/v1/session/tokener"
)

func hash(s string) string {
	return fmt.Sprintf(`%x`, sha1.Sum([]byte(s)))
}

// Session represents user sessions.
type Session struct {
	UserID    bson.ObjectId `bson:"user_id" json:"user_id"`
	TokenHash string        `bson:"token_hash" json:"-"`
	Token     string        `bson:"-" json:"token"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

const (
	userSessionCollectionName = "user_session"
	sessionTokenLength        = 40
)

var (
	// SessionCollection is the actual storage reference for user sessions.
	SessionCollection db.Collection
)

func init() {
	SessionCollection = storage.C(userSessionCollectionName)
}

// Constraint defines a compound key for session.
func (s *Session) Constraint() db.Cond {
	return db.Cond{
		"token_hash": s.TokenHash,
		"user_id":    s.UserID,
	}
}

// NewSession creates a new unique token that represents a user session.
func NewSession(userID bson.ObjectId) (*Session, error) {
	// TODO: Limit the number of possible sessions for the same user ID.
	var token string

	for {
		token = tokener.String(sessionTokenLength)
		if _, err := RetrieveSession(token); err == errmsg.ErrNoSuchSession {
			s := &Session{
				UserID:    userID,
				Token:     token,
				TokenHash: hash(token),
				CreatedAt: time.Now(),
			}
			if _, err := SessionCollection.Append(s); err != nil {
				return nil, err
			}
			return s, nil
		}
	}

}

// RetrieveSession retrieves a session using a token.
func RetrieveSession(token string) (*Session, error) {
	var err error

	res := SessionCollection.Find(db.Cond{
		"token_hash": hash(token),
	})

	if c, _ := res.Count(); c < 1 {
		return nil, errmsg.ErrNoSuchSession
	}

	var s Session
	if err = res.One(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Remove deletes a session from the list of active sessions.
func (s *Session) Remove() error {
	res := SessionCollection.Find(s.Constraint())

	if c, _ := res.Count(); c < 1 {
		return errmsg.ErrNoSuchSession
	}

	return res.Remove()
}
