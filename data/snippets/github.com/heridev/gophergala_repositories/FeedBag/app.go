package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tommy351/gin-cors"
)

var (
	configPort   = flag.String("port", "3000", "Port to run the server on")
	templatesDir = flag.String("templates", "./templates", "Path to your templates directory")
	indexFile    = flag.String("index-file", "./public/index.html", "Path to the index template")
	dbmap        = setupDb()
	templates    = setupTemplates()
	activityChan = make(chan []Activity)
)

func main() {
	//Parse flags
	flag.Parse()

	//Setup gin
	r := gin.Default()

	r.Use(cors.Middleware(cors.Options{AllowCredentials: true}))

	// Close the database connection if we fail
	defer dbmap.Db.Close()

	// Setup Goth Authentication
	goth.UseProviders(
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:3000/auth/github/callback"),
	)

	// Setup Socket.io server and related activity fetching
	socketServer, err := SetupSocketIO()
	checkErr(err, "Problem starting socket.io server:")

	SetupRoutes(r, socketServer)

	err = StartSocketPusher(socketServer, activityChan)
	checkErr(err, "Problem starting socket goroutine:")

	err = StartExistingUsers(activityChan)
	checkErr(err, "Problem starting user goroutines:")

	//Configure port for server to run on
	port := *configPort
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}

	// Listen and Serve on port from ENV or flag
	r.Run(fmt.Sprintf(":%s", port))
}

//Raw Json type for Type Converter
type RawJson map[string]interface{}

type TypeConverter struct{}

func (me TypeConverter) ToDb(val interface{}) (interface{}, error) {
	switch t := val.(type) {
	case RawJson:
		b, err := json.Marshal(t)
		if err != nil {
			return "", err
		}
		return string(b), nil

	case ActivityPayload:
		b, err := json.Marshal(t)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	return val, nil
}

func (me TypeConverter) FromDb(target interface{}) (gorp.CustomScanner, bool) {
	switch target.(type) {
	case *RawJson:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New("FromDb: Unable to convert Json to *string")
			}
			b := []byte(*s)
			return json.Unmarshal(b, target)
		}
		return gorp.CustomScanner{new(string), target, binder}, true

	case *ActivityPayload:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New("FromDb: Unable to convert Json to *string")
			}
			b := []byte(*s)
			return json.Unmarshal(b, target)
		}
		return gorp.CustomScanner{new(string), target, binder}, true
	}

	return gorp.CustomScanner{}, false
}

func setupDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("sqlite3", "./tmp/feedbag.bin")
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.TypeConverter = TypeConverter{}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(Activity{}, "activities").SetKeys(true, "Id")
	dbmap.AddTableWithName(ActivityPayload{}, "events").SetKeys(true, "Id").ColMap("GithubId").SetUnique(true)

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func setupTemplates() []*Template {
	// Parse templates
	templates, err := ParseTemplatesDir(*templatesDir)
	if err != nil {
		checkErr(err, "Problem parsing templates")
	}
	log.Println(fmt.Sprintf("Found %d valid templates", len(templates)))
	return templates
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
