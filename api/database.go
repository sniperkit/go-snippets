package api

import (
	r "github.com/dancannon/gorethink"
	"github.com/vsouza/pug/config"
)

type dbResource struct {
	DBConn *r.Session
}

// NewDBResource is an exported method to estabilish connection with database.
// And returns databse connection session.
func NewDBResource(c *config.Config) (*dbResource, error) {
	var err error
	var session *r.Session
	db := &dbResource{}
	session, err = r.Connect(r.ConnectOpts{
		Address:  c.GetString("db.conn_string"),
		Database: c.GetString("db.name"),
		MaxIdle:  10,
		MaxOpen:  10,
	})
	session.SetMaxOpenConns(c.GetInt("db.max_open_conn"))
	db.DBConn = session
	return db, err
}
