package hyperdb

import (
	"database/sql"
	"fmt"

	// This package is only valid with MySQL. It's meant to be convenient and
	// Therefor demanding that you do a blank import to use this package is
	// a thoughtless user (in this case the developer is the user) experience
	_ "github.com/go-sql-driver/mysql"
)

// Server is the unit of configuration for HyperDB.  From this structure we can
// determine which servers in which locations you can, or should, run your query
// against.
type Server struct {
	Dataset       string `json:"ds"`
	Datacenter    string `json:"dc"`
	ReadPriority  int    `json:"read"`
	WritePriority int    `json:"write"`
	WANAddress    string `json:"wan"`
	LANAddress    string `json:"lan"`
	Database      string `json:"db"`
	Username      string `json:"user"`
	Password      string `json:"pw"`
	db            *sql.DB
	hyperdb       *DB
}

func (s *Server) dsn() string {
	// note to self: need to get server DC versus running DC here for comparison...
	var address string
	if s.Datacenter == s.hyperdb.datacenter {
		address = s.LANAddress
	} else {
		address = s.WANAddress
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", s.Username, s.Password, address, s.Database)
}

// getDB attempts to get and validate a *sql.DB for the server
func (s *Server) getDB() *sql.DB {
	if s.db == nil {
		db, err := sql.Open("mysql", s.dsn())
		if err != nil {
			return nil
		}
		if err := db.Ping(); err != nil {
			return nil
		}
		s.db = db
		return db
	}
	if err := s.db.Ping(); err != nil {
		return nil
	}
	return s.db
}
