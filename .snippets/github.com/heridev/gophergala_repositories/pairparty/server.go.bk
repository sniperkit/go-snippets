package main

import (
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/logie17/Project-V/config"
	"github.com/logie17/Project-V/model"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

type key int

const UserKey key = 0

func main() {
	r := mux.NewRouter()
	configuration := config.LoadConfig()

	db, err := setupDB(configuration)

	if err != nil {
		log.Println("Uh, can't open the database: %s", err.Error())
	}

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/flex", flexHandler)
	r.HandleFunc("/webrtc", webrtcHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	http.Handle("/", Middleware(r, db))

	fmt.Printf("Starting on Port: [::]:%v\n", configuration.Port)
	http.ListenAndServe(fmt.Sprintf("[::]:%s", configuration.Port), nil)
}

// Apparently gorrilla doesn't support some sort
// of route chaining or middleware :-(
func Middleware(h http.Handler, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := model.User{Name: "Jinzhu", Username: "foo", CreatedAt: time.Now()}
		db.NewRecord(user) // => true
		db.Save(&user)
		context.Set(r, UserKey, user.Name)
		log.Println("middleware begin")
		h.ServeHTTP(w, r)
		log.Println("middleware end")
	})
}

func setupDB(configuration *config.Configuration) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", configuration.Database)

	if err != nil {
		return &db, err
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.SingularTable(true)

	db.CreateTable(&model.User{})

	return &db, err

}
