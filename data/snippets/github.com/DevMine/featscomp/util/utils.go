// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package util provides various utilities.
package util

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/DevMine/featscomp/config"
)

// OpenDBSession creates a session to the database.
func OpenDBSession(cfg config.DatabaseConfig) (*sql.DB, error) {
	dbURL := fmt.Sprintf(
		"user='%s' password='%s' host='%s' port=%d dbname='%s' sslmode='%s'",
		cfg.UserName, cfg.Password, cfg.HostName, cfg.Port, cfg.DBName, cfg.SSLMode)

	return sql.Open("postgres", dbURL)
}

// FPrintlnErr wraps fmt.Fprintln and use os.Stderr as output stream.
func FPrintlnErr(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

// FindFeatureID returns a feature ID, given its name.
func FindFeatureID(db *sql.DB, featureName string) (int64, error) {
	var id int64
	if err := db.QueryRow(
		`SELECT f.id
         FROM features f
         WHERE f.name = $1`, featureName).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
