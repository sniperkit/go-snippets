package rankings

import (
	"github.com/gophergala/GopherKombat/common/dba"
	"github.com/gophergala/GopherKombat/common/user"
	"gopkg.in/mgo.v2"
)

type Ranking struct {
	*user.User
	Value string
}

func GetDaily() []*Ranking {
	var results []*Ranking
	dba.Execute("users", func(col *mgo.Collection) {
		// TODO
	})
	return results
}

func GetMonthly() []*Ranking {
	var results []*Ranking
	dba.Execute("users", func(col *mgo.Collection) {
		// TODO
	})
	return results
}

func GetAllTime() []*Ranking {
	var results []*Ranking
	dba.Execute("users", func(col *mgo.Collection) {
		// TODO
	})
	return results
}
