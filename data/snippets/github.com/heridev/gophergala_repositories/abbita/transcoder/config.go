package transcoder

import (
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
)

type DatabaseSession struct {
	*mgo.Session
	databaseName string
}

func NewSession(name string) *DatabaseSession {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	return &DatabaseSession{session, name}
}

func (session *DatabaseSession) Database() martini.Handler {
	return func(context martini.Context) {
		s := session.Clone()
		context.Map(s.DB(session.databaseName))
		defer s.Close()
		context.Next()
	}
}
