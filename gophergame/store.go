package main

import (
	"errors"
	"github.com/gorilla/sessions"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
)

var Database *mgo.Database
var Sessions = sessions.NewCookieStore([]byte("7c2c23a802b04a65c28b1ca1e55dced3"))

func init() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	Database = session.DB("gopherquiz")
	//create an index for the username field on the users collection
	if err := Database.C("users").EnsureIndex(mgo.Index{
		Key:    []string{"username", "id"},
		Unique: true,
	}); err != nil {
		log.Fatal(err)
	}
}

type User struct {
	Id        string
	Username  string
	Highscore int
}

func AddUser(id, username string) error {
	user := &User{
		Id:       id,
		Username: username,
	}
	return Database.C("users").Insert(user)
}

func GetUser(id string) *User {
	var u *User
	err := Database.C("users").Find(bson.M{"id": id}).One(&u)
	if err != nil {
		log.Print(err)
		return nil
	}
	return u
}

func LoginUser(w http.ResponseWriter, r *http.Request, id, username, token string) {
	s, _ := Sessions.Get(r, "gopherquiz")
	s.Values["id"] = id
	s.Values["username"] = username
	s.Values["token"] = token
	s.Save(r, w)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	s, _ := Sessions.Get(r, "gopherquiz")
	user := CurrentUser(r)
	if user != nil {
		delete(s.Values, "username")
		delete(s.Values, "id")
		delete(s.Values, "token")
		s.Save(r, w)
	}
}

func CurrentUser(r *http.Request) *User {
	s, _ := Sessions.Get(r, "gopherquiz")
	id, ok := s.Values["id"]
	if ok {
		return GetUser(id.(string))
	}
	return nil
}

func GetHighScores() []*User {
	var users []*User
	err := Database.C("users").Find(nil).Select(bson.M{"_id": 0, "id": 0}).Sort("-highscore").Limit(20).All(&users)
	if err != nil {
		return nil
	}
	return users
}

func SaveHighScore(score string, r *http.Request) (bool, error) {
	s, _ := strconv.Atoi(score)
	if u := CurrentUser(r); u != nil {
		if s <= u.Highscore {
			return false, nil
		}
		u.Highscore = s
		return true, Database.C("users").Update(bson.M{"id": u.Id}, u)
	}
	return false, errors.New("unauthorized")
}
