// Task Management manages tasks, including tracking the state of tasks as they move through a queue
// As tasks are completed task_mgmt updates records of when tasks started, stopped, etc.
package main

import (
	"database/sql"
	"fmt"
	"github.com/datatogether/sql_datastore"
	"github.com/datatogether/sqlutil"
	"github.com/datatogether/task_mgmt/source"
	"github.com/datatogether/task_mgmt/tasks"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var (
	// cfg is the global configuration for the server. It's read in at startup from
	// the config.json file and enviornment variables, see config.go for more info.
	cfg *config
	// log output
	log = logrus.New()
	// application database connection
	appDB = &sql.DB{}
	// hoist default store
	store = sql_datastore.DefaultStore
)

func init() {
	log.Out = os.Stderr
	log.Level = logrus.InfoLevel
	log.Formatter = &logrus.TextFormatter{
		ForceColors: true,
	}
}

func main() {
	var err error
	cfg, err = initConfig(os.Getenv("GOLANG_ENV"))
	if err != nil {
		// panic if the server is missing a vital configuration detail
		panic(fmt.Errorf("server configuration error: %s", err.Error()))
	}
	configureTasks()

	go initPostgres()
	go listenRpc()
	go connectRedis()

	stop, err := acceptTasks()
	if err != nil {
		panic(err.Error())
	}

	s := &http.Server{}
	// connect mux to server
	s.Handler = NewServerRoutes()

	// print notable config settings
	// printConfigInfo()

	// fire it up!
	log.Infoln("starting server on port", cfg.Port)

	// start server wrapped in a log.Fatal b/c http.ListenAndServe will not
	// return unless there's an error
	log.Fatal(StartServer(cfg, s))

	// lol will never happen, left here as a reminder
	// that we should be able to stop accepting new tasks
	// at any point without issue
	stop <- true
}

// NewServerRoutes returns a Muxer that has all API routes.
// This makes for easy testing using httptest
func NewServerRoutes() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/.well-known/acme-challenge/", CertbotHandler)
	m.Handle("/", middleware(NotFoundHandler))
	m.Handle("/healthcheck", middleware(HealthCheckHandler))

	m.Handle("/tasks", middleware(TasksHandler))
	m.Handle("/tasks/", middleware(TaskHandler))
	// TODO - restore this:
	// m.Handle("/tasks/cancel/", middleware(CancelTaskHandler))

	// Example of individual task routing:
	m.HandleFunc("/ipfs/add", middleware(EnqueueIpfsAddHandler))

	m.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("public/js"))))
	m.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("public/css"))))

	return m
}

func initPostgres() {
	log.Infoln("connecting to postgres db")
	if err := sqlutil.ConnectToDb("postgres", cfg.PostgresDbUrl, appDB); err != nil {
		panic(err)
	}
	log.Infoln("connected to postgres db")
	created, err := sqlutil.EnsureTables(appDB, packagePath("sql/schema.sql"),
		"tasks")
	if err != nil {
		log.Infoln(err)
	}
	if len(created) > 0 {
		log.Infoln("created tables:", created)
	}

	sql_datastore.SetDB(appDB)
	store.Register(
		&tasks.Task{},
		&source.Source{},
	)
}
