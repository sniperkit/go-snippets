// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package apiutil provides various functions meant to be used by handlers
// from the api package.
package apiutil

import (
	"database/sql"
	"errors"

	"github.com/golang/glog"

	"github.com/DevMine/api-server/model"
	"github.com/DevMine/api-server/util/pgutil"
	"github.com/DevMine/api-server/util/typeutil"
)

// CreateGhOrgsFromPGArray generates a slice of GhOrganization from a
// PostgresSQL array from gh_organizations table. Elements MUST be in
// table order.
// FIXME Find a safer way to deal with this
func CreateGhOrgsFromPGArray(pgArray string) []*model.GhOrganization {
	var ghOrgs []*model.GhOrganization

	for _, tmp := range pgutil.ParsePgArray(pgArray) {
		var gho model.GhOrganization
		org := pgutil.ParsePgRow(tmp)

		if len(org) != 13 {
			glog.Error(errors.New("GhOrganization pgArray must contain 13 elements"))
			return nil
		}

		// Fields are expected in table order
		id, err := typeutil.StrToInt(org[0])
		if err != nil {
			glog.Error(err)
			// id is a non nullable field so there is a problem somewhere
			return nil
		}
		gho.ID = &id

		ghid, err := typeutil.StrToInt(org[1])
		if err != nil {
			glog.Error(err)
			// github_id is a non nullable field so there is a problem somewhere
			return nil
		}
		gho.GithubID = &ghid

		if len(org[2]) > 0 {
			gho.Login = &org[2]
		} else {
			glog.Error(errors.New("encountered NULL login value"))
			// login is a non nullable field
			return nil
		}

		if len(org[3]) > 0 {
			gho.AvatarURL = &org[3]
		}

		if len(org[4]) > 0 {
			gho.HTMLURL = &org[4]
		}

		if len(org[5]) > 0 {
			gho.Name = &org[5]
		}

		if len(org[6]) > 0 {
			gho.Company = &org[6]
		}

		if len(org[7]) > 0 {
			gho.Blog = &org[7]
		}

		if len(org[8]) > 0 {
			gho.Location = &org[8]
		}

		if len(org[9]) > 0 {
			gho.Email = &org[9]
		}

		collCount, err := typeutil.StrToInt(org[10])
		if err == nil {
			gho.CollaboratorsCount = &collCount
		}

		if len(org[11]) > 0 {
			tmp, err := pgutil.TimestampTZToTime(org[11])
			if err == nil {
				gho.CreatedAt = &tmp
			}
		}

		if len(org[12]) > 0 {
			tmp, err := pgutil.TimestampTZToTime(org[12])
			if err == nil {
				gho.UpdatedAt = &tmp
			}
		}

		ghOrgs = append(ghOrgs, &gho)
	}

	return ghOrgs
}

// FetchUser retrieves a user from the database.
func FetchUser(db *sql.DB, id *int64) *model.User {
	if id == nil || db == nil {
		return nil
	}

	var u model.User
	err := db.QueryRow(`
        SELECT id, username, name, email
        FROM users
        WHERE id=$1`,
		*id).Scan(&u.ID, &u.Username, &u.Name, &u.Email)
	if err != nil {
		panic(err)
	}

	return &u
}

// FetchRepository retrieves a repository from the database.
func FetchRepository(db *sql.DB, id *int64) *model.Repository {
	if id == nil || db == nil {
		return nil
	}

	var r model.Repository
	err := db.QueryRow(`
        SELECT id, name, primary_language, clone_url, clone_path, vcs
        FROM repositories
        WHERE id=$1`,
		*id).Scan(&r.ID, &r.Name, &r.PrimaryLanguage, &r.CloneURL, &r.ClonePath, &r.VCS)
	if err != nil {
		panic(err)
	}

	return &r
}
