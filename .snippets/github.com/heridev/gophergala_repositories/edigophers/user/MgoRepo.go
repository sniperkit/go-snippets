package user

import (
	"edigophers/utils"
	"sort"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MgoRepo is a repository based on a MongoDb
type MgoRepo struct {
	url        string
	database   string
	collection string
}

const collection = "users"

func (r MgoRepo) getSession() *mgo.Session {
	s, err := mgo.Dial(r.url)
	utils.CheckError(err)
	return s
}

// GetUser gets a user per his username
func (r MgoRepo) GetUser(name string) (*User, error) {
	s := r.getSession()
	defer s.Close()
	c := s.DB(r.database).C(r.collection)
	user := User{}
	err := c.Find(bson.M{"name": name}).One(&user)
	if err != nil {
		return nil, err
	}
	sort.Sort(ByRatingDesc(user.Interests))
	return &user, nil
}

//GetUsers is a function returning a list of users
func (r MgoRepo) GetUsers() ([]User, error) {
	s := r.getSession()
	defer s.Close()
	c := s.DB(r.database).C(r.collection)
	users := []User{}
	err := c.Find(nil).All(&users)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		sort.Sort(ByRatingDesc(user.Interests))
	}
	return users, nil
}

// SaveUser saves a user to the database
func (r MgoRepo) SaveUser(usr User) error {
	s := r.getSession()
	defer s.Close()
	c := s.DB(r.database).C(r.collection)
	return c.UpdateId(usr.ID, usr)
}

// NewMgoRepo creates a new Mongo database repository
func NewMgoRepo(url, database string) Repository {
	repo := MgoRepo{url: url, database: database, collection: collection}
	return Repository(repo)
}
