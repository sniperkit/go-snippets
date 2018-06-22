package user

import (
	"encoding/json"
	"github.com/gophergala/GopherKombat/common/dba"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type User struct {
	Name    string
	Repo    string
	Image   string
	Wins    int
	Matches int
}

func ParseFromJson(content []byte) *User {
	var data map[string]interface{}

	err := json.Unmarshal(content, &data)
	if err != nil {
		log.Printf("Error parsing response: $s", err)
	}
	user, found := Find(data["login"].(string))
	if found {
		return user
	} else {
		return &User{
			Name:  data["login"].(string),
			Repo:  data["html_url"].(string),
			Image: data["avatar_url"].(string),
		}
	}
}

func (u *User) Save() {
	dba.Execute("users", func(col *mgo.Collection) {
		err := col.Insert(u)
		if err != nil {
			panic(err)
		}
	})
}

func Find(name string) (*User, bool) {
	user := &User{}
	found := true
	dba.Execute("users", func(col *mgo.Collection) {
		err := col.Find(bson.M{"name": name}).One(user)
		if err != nil {
			found = false
		}
	})
	return user, found
}

func GetAll() []*User {
	var results []*User
	dba.Execute("users", func(col *mgo.Collection) {
		col.Find(nil).All(&results)
	})
	return results
}
