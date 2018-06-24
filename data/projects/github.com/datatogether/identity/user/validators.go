package user

import (
	"database/sql"
	"github.com/datatogether/errors"
	"github.com/datatogether/sqlutil"
	"net/url"
	"regexp"
	"strings"

	"github.com/pborman/uuid"
)

var (
	// alphanumeric must start with a letter and contian only letters & numbers
	alphaNumericRegex = regexp.MustCompile(`^[A-Za-z0-9_-]{2,35}$`)
	titleRegex        = regexp.MustCompile(`^[\sa-z0-9_-]{1,200}$`)
	// yes, this is just ripped from the internet somewhere. Yes it should be improved. TODO - validate emails the right way
	emailRegex = regexp.MustCompile(`(?i)[A-Z0-9!#$%&'*+/=?^_{|}~-]+(?:\.[A-Z0-9!#$%&'*+/=?^_{|}~-]+)*@(?:[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?\.)+[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?`)
	slugRegex  = regexp.MustCompile(`^[a-z0-9-_]+$`)
	pathRegex  = regexp.MustCompile(`^[a-z0-9-_/]+/$`)
)

// make sure a username contains only alphanumeric chars,_,-, and starts with a letter
// and it can't be a uuid b/c that'll confuse the json unmarshaller, also, a username that's a uuid sounds, like, phishy
func validUsername(username string) bool {
	return alphaNumericRegex.MatchString(username) && !validUuid(username)
}

// check email against regex
func validEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// check slug against regex
func validSlug(slug string) bool {
	return slugRegex.MatchString(slug)
}

// check path against regex
func validPath(path string) bool {
	return pathRegex.MatchString(path)
}

// see if a string is in fact a UUID
func validUuid(id string) bool {
	return uuid.Parse(id) != nil
}

// check if a username is taken, also checking against
// organization namespace to avoid collisions
// TODO - refactor to only return an error if taken
func UsernameTaken(db sqlutil.Queryable, username string) (taken bool, err error) {
	e := db.QueryRow("SELECT exists(SELECT 1 FROM(SELECT lower(username) FROM users WHERE username = $1 AND deleted=false) AS existing)", strings.ToLower(username)).Scan(&taken)

	if e == sql.ErrNoRows {
		taken = false
	} else if e != nil {
		err = errors.New500Error(e.Error())
	}

	return
}

// check if an email is taken
func EmailTaken(db sqlutil.Queryable, email string) (taken bool, err error) {
	e := db.QueryRow(`SELECT exists(SELECT 1 FROM users WHERE email = $1 AND deleted=false)`, email).Scan(&taken)

	if e == sql.ErrNoRows {
		taken = false
	} else if e != nil {
		err = errors.New500Error(e.Error())
	}
	return
}

// check if a dataset path is taken
func PathTaken(db sqlutil.Queryable, path string) (taken bool, err error) {
	e := db.QueryRow("SELECT exists(SELECT 1 FROM datasets WHERE path = $1 AND deleted=false)", path).Scan(&taken)
	if e == sql.ErrNoRows {
		taken = false
	} else if e != nil {
		err = errors.New500Error(e.Error())
	}
	return
}

// check if dataset exists in a given dataset
func DatasetExists(db sqlutil.Queryable, datasetId string) (exists bool, err error) {
	e := db.QueryRow("SELECT exists(SELECT 1 FROM datasets WHERE id = $1 and deleted = false)", datasetId).Scan(&exists)
	if e == sql.ErrNoRows {
		exists = false
	} else if e != nil {
		err = errors.New500Error(e.Error())
	}

	return
}

// check if a user exists on a given database
func ValidUser(db sqlutil.Queryable, u *User) (err error) {
	if u == nil {
		return errors.ErrInvalidUser
	}

	if !validUuid(u.Id) {
		return errors.ErrInvalidUser
	}

	exists := false
	err = db.QueryRow("SELECT exists(SELECT 1 FROM users WHERE id = $1 and deleted = false)", u.Id).Scan(&exists)
	if err == sql.ErrNoRows || !exists {
		err = errors.ErrUserNotFound
	}

	return
}

// ValidUrlString adds
func ValidUrlString(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	if u.Scheme == "" {
		// assume http if no scheme exists
		u.Scheme = "http"
	}

	return u.String(), nil
}
