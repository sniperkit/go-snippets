package dba

import (
	"gopkg.in/mgo.v2"
	"os"
)

func Execute(col string, f func(*mgo.Collection)) {
	session, err := mgo.Dial(os.Getenv("MONGOLAB_URL"))
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("gopherkombat").C(col)
	f(c)
}
