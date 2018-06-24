package main

import (
	"database/sql"
	"fmt"
	"github.com/datatogether/core"
	"github.com/datatogether/sql_datastore"
	"github.com/gchaincl/dotsql"
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
	appDB *sql.DB

	//
	store = sql_datastore.DefaultStore

	sqlCmds *dotsql.DotSql

	room *Room
)

func init() {
	log.Out = os.Stdout
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

	connectToAppDb()
	sql_datastore.SetDB(appDB)
	sql_datastore.Register(
		&core.Url{},
		&core.Link{},
		&core.Primer{},
		&core.Source{},
		&core.Collection{},
		&core.CollectionItem{},
	)

	go func() {
		if err := SubscribeTaskProgress(); err != nil {
			log.Infoln("task progress error:", err.Error())
		}
	}()

	room = newRoom()
	go room.run()

	s := &http.Server{}
	// connect mux to server
	s.Handler = NewServerRoutes()

	// print notable config settings
	printConfigInfo()

	// fire it up!
	log.Infof("🌎 starting server on port %s in %s mode", cfg.Port, cfg.Mode)

	// start server wrapped in a log.Fatal b/c http.ListenAndServe will not
	// return unless there's an error
	log.Fatal(StartServer(cfg, s))
}

// NewServerRoutes returns a Muxer that has all API routes.
// This makes for easy testing using httptest
func NewServerRoutes() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/.well-known/acme-challenge/", CertbotHandler)
	m.Handle("/profile", middleware(UserProfileHandler))
	m.Handle("/healthcheck", middleware(HealthCheckHandler))

	m.Handle("/", middleware(WebappHandler))
	m.Handle("/url", middleware(WebappHandler))
	m.Handle("/content/", middleware(WebappHandler))
	m.Handle("/metadata/", middleware(WebappHandler))
	m.Handle("/settings", middleware(WebappHandler))
	m.Handle("/settings/keys", middleware(WebappHandler))
	m.Handle("/users/", middleware(WebappHandler))
	m.Handle("/signup", middleware(WebappHandler))
	m.Handle("/login", middleware(WebappHandler))
	m.Handle("/login/forgot", middleware(WebappHandler))
	m.Handle("/primers", middleware(WebappHandler))
	m.Handle("/primers/", middleware(WebappHandler))
	m.Handle("/subprimers", middleware(WebappHandler))
	m.Handle("/subprimers/", middleware(WebappHandler))
	m.Handle("/tasks", middleware(WebappHandler))
	m.Handle("/tasks/", middleware(WebappHandler))

	m.Handle("/ws", middleware(HandleWebsocketUpgrade))

	return m
}
