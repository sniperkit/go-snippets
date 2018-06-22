package database

import (
	"context"
	"errors"
	"time"

	"github.com/Xe/ln"
	"github.com/Xe/uuid"
	"github.com/asdine/storm"
)

// Token is an individual API access token.
type Token struct {
	ID           string `storm:"id"`
	UserID       string `storm:"index"`
	Body         string `storm:"index"`
	CreationDate time.Time
	LastUseDate  time.Time

	// Deleted flags if the token shows up to users or not. If this is set to
	// true, this token cannot be used in any API calls.
	Deleted bool
}

// F ields for logging.
func (t Token) F() ln.F {
	f := ln.F{
		"token_id":            t.ID,
		"token_userid":        t.UserID,
		"token_creation_date": t.CreationDate,
	}

	return f
}

type Tokens interface {
	Create(ctx context.Context, userID string) (*Token, error)
	UpdateLastSeen(ctx context.Context, tokenID string) error
	Check(ctx context.Context, body string) (*Token, error)
}

// interface compliance checks
var (
	_ Tokens = &tokensStorm{}
)

func NewTokensStorm(db *storm.DB, us Users) Tokens {
	return &tokensStorm{db: db, us: us}
}

type tokensStorm struct {
	db *storm.DB
	us Users
}

func (t *tokensStorm) Create(ctx context.Context, userID string) (*Token, error) {
	_, err := t.us.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	tkn := &Token{
		ID:           uuid.New(),
		UserID:       userID,
		Body:         uuid.New() + uuid.New(),
		CreationDate: time.Now(),
		LastUseDate:  time.Now(),
	}

	err = t.db.Save(tkn)
	if err != nil {
		return nil, err
	}

	return tkn, nil
}

func (t *tokensStorm) UpdateLastSeen(ctx context.Context, tokenID string) error {
	var tkn Token
	err := t.db.One("ID", tokenID, &tkn)
	if err != nil {
		return err
	}

	tkn.LastUseDate = time.Now()

	err = t.db.Save(&tkn)
	if err != nil {
		return err
	}

	return nil
}

func (t *tokensStorm) Check(ctx context.Context, body string) (*Token, error) {
	var tkn Token
	err := t.db.One("Body", body, &tkn)
	if err != nil {
		return nil, err
	}

	if tkn.Deleted {
		return nil, errors.New("token deleted")
	}

	t.UpdateLastSeen(ctx, tkn.ID)

	return &tkn, nil
}
