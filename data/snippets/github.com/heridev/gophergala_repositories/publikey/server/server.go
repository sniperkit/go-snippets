package server

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"golang.org/x/crypto/bcrypt"
)

type server struct {
	db            *bolt.DB
	dataBucket    *bolt.Bucket
	sessionBucket *bolt.Bucket
}

type User struct {
	Email        string      `json:"email"`
	PasswordHash []byte      `json:"-"`
	CreatedAt    time.Time   `json:"createdAt"`
	Keys         []PublicKey `json:"publicKeys"`
}

type NewUser struct {
	*User
	ApiKey string `json:"apiKey"`
}

type PublicKey struct {
	Public     bool      `json:"public"`
	Value      string    `json:"value"`
	Expiration time.Time `json:"expiresAt"`
}

var Server = &server{}

func Serve(port string, dbfile string) {
	setupDatabase(dbfile)
	m := martini.Classic()

	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		r.JSON(200, struct{}{})
	})

	m.Post("/keys", addKey)
	m.Group("/users", func(r martini.Router) {
		r.Group("/:email", func(k martini.Router) {
			r.Get("/keys", getKeys)
		})
		r.Post("/new", register)
	})

	m.RunOnAddr(":" + port)
}

func setupDatabase(dbfile string) {
	db, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	Server.db = db

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Data"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		Server.dataBucket = b
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Sessions"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		Server.sessionBucket = b
		return nil
	})
}

func getKeys(r render.Render, hr *http.Request, a martini.Params) {
	var u User
	err := Server.db.View(func(tx *bolt.Tx) error {
		email := a["email"]
		b := tx.Bucket([]byte("Data"))
		v := b.Get([]byte(email))

		err := json.Unmarshal(v, &u)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, u)
}

func addKey(r render.Render, hr *http.Request) {
	apiKey := hr.Header.Get("X-PUBLIKEY-API-KEY")
	var user User

	err := Server.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Sessions"))
		db := tx.Bucket([]byte("Data"))
		v := b.Get([]byte(apiKey))

		if v != nil {
			dv := db.Get(v)
			if dv != nil {
				err := json.Unmarshal(dv, &user)
				return err
			}
		}
		return errors.New("not found")
	})

	if err != nil {
		r.Status(http.StatusUnauthorized)
		return
	}
	user.Keys = append(user.Keys, PublicKey{Public: true, Value: hr.FormValue("key")})

	err = Server.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Data"))
		json, err := json.Marshal(user)
		err = b.Put([]byte(user.Email), []byte(json))

		return err
	})

	r.JSON(http.StatusOK, user)
}

func register(r render.Render, hr *http.Request) {
	email := hr.FormValue("email")
	password := hr.FormValue("password")

	err := Server.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Data"))
		v := b.Get([]byte(email))
		if v != nil {
			return errors.New("already exists")
		}
		return nil
	})

	if err != nil || email == "" || password == "" {
		r.Status(http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		panic(err)
	}

	user := &User{
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		Keys:         []PublicKey{},
	}

	var apiKey string

	err = Server.db.Update(func(tx *bolt.Tx) error {
		data, _ := json.Marshal(user)
		dataBucket := tx.Bucket([]byte("Data"))
		dataBucket.Put([]byte(email), []byte(data))

		apiKey = generateApiKey(email)

		sessionBucket := tx.Bucket([]byte("Sessions"))
		sessionBucket.Put([]byte(apiKey), []byte(email))

		return nil
	})

	if err != nil {

		r.Error(http.StatusBadRequest)
	}
	r.JSON(http.StatusOK, NewUser{User: user, ApiKey: apiKey})

}

func generateApiKey(email string) string {
	length := 128
	b := make([]byte, length)
	rand.Read(b)
	hasher := sha1.New()
	hasher.Write([]byte(email))
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
}
