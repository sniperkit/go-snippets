package database

import (
	"context"
	"errors"
	"time"

	"github.com/Xe/ln"
	"github.com/Xe/uuid"
	"github.com/asdine/storm"
	"golang.org/x/crypto/bcrypt"
)

// User is an individual user of pawd.
type User struct {
	ID             string `storm:"id"`
	Email          string `storm:"unique"`
	PasswordHash   []byte // bcrypt
	TOTPSecret     []byte
	CreationDate   time.Time
	LastAccessDate time.Time

	// Suspended if set true disables every API action by this user.
	Suspended       bool
	SuspendedDate   time.Time
	SuspendedReason string
	SuspendedAdmin  string

	Admin bool
}

// F ields for logging.
func (u User) F() ln.F {
	f := ln.F{
		"user_id":               u.ID,
		"user_email":            u.Email,
		"user_creation_date":    u.CreationDate,
		"user_last_access_date": u.LastAccessDate,
		"user_admin":            u.Admin,
	}

	if u.Suspended {
		f["user_suspended"] = true
		f["user_suspended_date"] = u.SuspendedDate
		f["user_suspended_reason"] = u.SuspendedReason
		f["user_suspended_admin"] = u.SuspendedAdmin
	}

	return f
}

// Users are the datastore calls for user management.
type Users interface {
	// Create s a new user with the given email, password and totp settings.
	// If wantTOTP is set to true, the user MUST be given the secret only once.
	Create(ctx context.Context, email, password string, wantTOTP bool) (*User, error)

	// UpdateLastSeen bumps the last seen time of a user by ID, kind of like the unix
	// "touch" command.
	UpdateLastSeen(ctx context.Context, userID string) error

	// CheckPassword checks a given email address, password and optional TOTP challenge
	// against the values in the database and returns true if they match. This calls UpdateLastSeen.
	CheckPassword(ctx context.Context, email, password, totpChallenge string) (*User, error)

	// Get fetches an individual user by ID.
	Get(ctx context.Context, userID string) (*User, error)
}

// interface compliance checks
var (
	_ Users = &usersStorm{}
)

// NewUsersStorm returns a new instance of the Users dao hooked up to a storm database.
func NewUsersStorm(db *storm.DB, admins []string) Users {
	return &usersStorm{
		db:          db,
		adminEmails: admins,
	}
}

type usersStorm struct {
	db          *storm.DB
	adminEmails []string
}

func (u *usersStorm) Create(ctx context.Context, email, password string, wantTOTP bool) (*User, error) {
	if wantTOTP {
		return nil, errors.New("not implemented")
	}

	hsh, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	id := uuid.New()

	nu := &User{
		ID:             id,
		Email:          email,
		PasswordHash:   hsh,
		CreationDate:   now,
		LastAccessDate: now,
	}

	for _, eml := range u.adminEmails {
		if eml == email {
			nu.Admin = true
		}
	}

	err = u.db.Save(nu)
	if err != nil {
		return nil, err
	}

	return nu, nil
}

func (u *usersStorm) UpdateLastSeen(ctx context.Context, userID string) error {
	var uu User
	err := u.db.One("ID", userID, &uu)
	if err != nil {
		return err
	}

	uu.LastAccessDate = time.Now()

	err = u.db.Save(&uu)
	if err != nil {
		return err
	}

	return nil
}

func (u *usersStorm) CheckPassword(ctx context.Context, email, password, totpChallenge string) (*User, error) {
	var uu User
	err := u.db.One("Email", email, &uu)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(uu.PasswordHash, []byte(password))
	if err != nil {
		return nil, err
	}

	return &uu, u.UpdateLastSeen(ctx, uu.ID)
}

func (u *usersStorm) Get(ctx context.Context, userID string) (*User, error) {
	var uu User
	err := u.db.One("ID", userID, &uu)
	if err != nil {
		return nil, err
	}

	return &uu, nil
}
