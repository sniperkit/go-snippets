// Copyright (c) 2018 Jef Oliver. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbabstract

import (
	"database/sql"
	"fmt"
)

const (
	// TCP constant for connection type of tcp4 and tcp6
	TCP = "tcp"
	// TCP4 constant for connection type of tcp4
	TCP4 = "tcp4"
	// TCP6 constant for connection type of tcp6
	TCP6 = "tcp6"
	// UNIX constant for connection type of unix socket
	UNIX = "unix"
)

var holderFactories = make(map[string]dbHolderFactory)

// DBOpts holds database configuration options
type DBOpts struct {
	DataDir     string // path to the database data directory (only applicable with sqlite3)
	Driver      string // database driver to use
	ConnectType string // tcp or unix (tcp must set host and port, unix must set path)
	Host        string // address that database server is listening on
	Port        int    // port that the database server is listening on
	SocketPath  string // path to unix socket to use
	Username    string // username to use when connecting to the database server
	Password    string // password to use when connecting to the database server
	DBName      string // database name to connect to (for sqlite3, this is the filename without the .db extension)
}

// DBHolder holds the open database connection and methods for interacting with the database.
type DBHolder interface {
	// Close is a wrapper function to sql.DB.Close()
	Close() error
	// DB returns the pointer to the open SQL database connection
	DB() *sql.DB
	// Driver returns the configured driver
	Driver() string
	// Format formats an SQL queries argument identifier for the underlying driver. This needs to be run
	// before Prepare().
	//
	// Replaces ? with $ or vice versa
	Format(query string) string
	// Path returns the path used to connect to the database. This is useful for debug purposes.
	// the username and password are not stored.
	Path() string
	// TableExists checks for the existence of a table in the database for that specified driver.
	// logger should be the function you wish to have used for logging (ie log.Debug)
	// The query statement for checking if the table exists will be logged
	// If logging isn't desired, logger should be nil
	TableExists(table string, logger func(args ...interface{})) (bool, error)
}

// dbHolderFactory is called to register a database driver
type dbHolderFactory func(opts DBOpts) (DBHolder, error)

// NewDBHolder returns an initialized DBHolder for operating on a configured database
func NewDBHolder(opts DBOpts) (DBHolder, error) {
	holder, ok := holderFactories[opts.Driver]
	if !ok {
		return nil, fmt.Errorf("unsupported database type: %s", opts.Driver)
	}

	return holder(opts)
}
