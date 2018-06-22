package user

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"edigophers/utils"

	"github.com/gorilla/sessions"
	"gopkg.in/mgo.v2/bson"
)

// User type
type User struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string
	Location  Location
	Interests Interests
}

// Location type
type Location struct {
	Longitude float64
	Latitude  float64
}

// Interests type
type Interests []Interest

// ByRatingDesc sorting type for Interests
type ByRatingDesc []Interest

func (a ByRatingDesc) Len() int           { return len(a) }
func (a ByRatingDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRatingDesc) Less(i, j int) bool { return a[i].Rating > a[j].Rating }

// Interest type
type Interest struct {
	Name     string
	Rating   float64
	Distance float64
}

// Repository type
type Repository interface {
	GetUsers() ([]User, error)
	GetUser(name string) (*User, error)
	SaveUser(user User) error
}

// NewUser creates a new user
func NewUser(name string, interests Interests) *User {
	return &User{Name: name, Interests: interests}
}

// NewInterest creates a new interest
func NewInterest(name string, rating float64) *Interest {
	return &Interest{Name: name, Rating: rating, Distance: -1.0}
}

// AsMap creates a map from the users interests
func (ui Interests) AsMap() map[interface{}]float64 {
	interestMap := make(map[interface{}]float64)
	for _, i := range ui {
		interestMap[i.Name] = i.Rating
	}
	return interestMap
}

// GetSessionUser gets the user stored in the session if there is one
func GetSessionUser(w http.ResponseWriter, r *http.Request, repo Repository, store *sessions.CookieStore) (*User, error) {
	session, err := store.Get(r, "user-session")
	utils.CheckError(err)
	username, ok := session.Values["username"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil, errors.New("No user in the session")
	}
	user, err := repo.GetUser(username.(string))
	if err != nil {
		if err := LogOutSessionUser(w, r, store); err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil, err
	}
	return user, nil
}

// SetSessionUser sets the user in the session
func SetSessionUser(w http.ResponseWriter, r *http.Request, username string, store *sessions.CookieStore) error {
	session, err := store.Get(r, "user-session")
	utils.CheckError(err)
	if strings.Trim(username, " ") == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	session.Values["username"] = username
	return session.Save(r, w)
}

// LogOutSessionUser logs out the user
func LogOutSessionUser(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore) error {
	session, err := store.Get(r, "user-session")
	utils.CheckError(err)
	delete(session.Values, "username")
	return session.Save(r, w)
}
