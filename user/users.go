package user

import (
	"database/sql"
	"fmt"
	"github.com/datatogether/errors"
	"github.com/datatogether/sqlutil"
)

// ReadUsers reads a page of users
func ReadUsers(db sqlutil.Queryable, userType UserType, limit, offset int) (users []*User, err error) {
	users = make([]*User, 0)
	// TODO - make this not bad
	query := fmt.Sprintf("SELECT %s FROM users WHERE deleted=false ORDER BY created DESC LIMIT $1 OFFSET $2", userColumns())
	if userType == UserTypeCommunity {
		query = fmt.Sprintf("SELECT %s FROM users WHERE deleted=false AND type = 2 ORDER BY created DESC LIMIT $1 OFFSET $2", userColumns())
	} else if userType == UserTypeUser {
		query = fmt.Sprintf("SELECT %s FROM users WHERE deleted=false AND type = 1 ORDER BY created DESC LIMIT $1 OFFSET $2", userColumns())
	}

	rows, e := db.Query(query, limit, offset)
	if e != nil {
		if e == sql.ErrNoRows {
			return []*User{}, nil
		}
		return nil, errors.New500Error(e.Error())
	}
	defer rows.Close()
	if us, e := scanUsers(rows); e != nil {
		return nil, errors.New500Error(e.Error())
	} else {
		return us, nil
	}
}

// scan a table of users results into a slice of user pointers.
func scanUsers(rows *sql.Rows) ([]*User, error) {
	us := make([]*User, 0)
	for rows.Next() {
		u := &User{}
		err := u.UnmarshalSQL(rows)
		if err != nil {
			return nil, err
		}
		us = append(us, u)
	}

	return us, nil
}

func UsersSearch(db sqlutil.Queryable, query string, limit, offset int) ([]*User, error) {
	q := fmt.Sprintf("%%%s%%", query)
	rows, err := db.Query(qUsersSearch, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanUsers(rows)
}

func CommunityUsers(db sqlutil.Queryable, community *User, order string, limit, offset int) ([]*User, error) {
	rows, e := db.Query(qCommunityMembers, community.Id, order, limit, offset)
	if e != nil {
		if e == sql.ErrNoRows {
			return []*User{}, nil
		}
		return nil, errors.New500Error(e.Error())
	}
	defer rows.Close()
	if us, e := scanUsers(rows); e != nil {
		return nil, errors.New500Error(e.Error())
	} else {
		return us, nil
	}
}

func UserCommunities(db sqlutil.Queryable, user *User, order string, limit, offset int) ([]*User, error) {
	rows, e := db.Query(qUserCommunities, user.Id, order, limit, offset)
	if e != nil {
		if e == sql.ErrNoRows {
			return []*User{}, nil
		}
		return nil, errors.New500Error(e.Error())
	}
	defer rows.Close()
	if us, e := scanUsers(rows); e != nil {
		return nil, errors.New500Error(e.Error())
	} else {
		return us, nil
	}
}
