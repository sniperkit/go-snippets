package main

import (
	"database/sql"
	"fmt"
	"github.com/datatogether/identity/jwt"
	"github.com/datatogether/identity/oauth"
	"github.com/datatogether/sqlutil"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var (
	// cfg is the global configuration for the server. It's read in at startup from
	// the config.json file and enviornment variables, see config.go for more info.
	cfg *config

	// log output
	log = logrus.New()

	// application database connection
	appDB = &sql.DB{}

	// cookie session storage
	sessionStore *sessions.CookieStore
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

	oauth.InitOauth(cfg.GithubAppId, cfg.GithubAppSecret)
	jwt.InitKeys(cfg.PublicKey, cfg.PrivateKey)

	sessionStore = sessions.NewCookieStore([]byte(cfg.SessionSecret))
	if cfg.UserCookieDomain != "" {
		sessionStore.Options.Domain = cfg.UserCookieDomain
	}

	go func() {
		if err := sqlutil.ConnectToDb("postgres", cfg.PostgresDbUrl, appDB); err != nil {
			log.Panic(err.Error())
		}
		created, err := sqlutil.EnsureSeedData(appDB, packagePath("sql/schema.sql"), packagePath("sql/test_data.sql"),
			"users",
			"reset_tokens",
			"keys",
			"oauth_tokens",
			"community_users")
		if err != nil {
			log.Infof("error allocating seed data: %s", err.Error())
		}
		if len(created) > 0 {
			log.Infoln("created tables & seed data:", created)
		}
	}()
	go listenRpc()

	s := &http.Server{}
	// connect mux to server
	s.Handler = NewServerRoutes()

	// print notable config settings
	// printConfigInfo()

	// fire it up!
	fmt.Println("starting server on port", cfg.Port)

	// start server wrapped in a log.Fatal b/c http.ListenAndServe will not
	// return unless there's an error
	log.Fatal(StartServer(cfg, s))
}

// NewServerRoutes returns a Muxer that has all API routes.
// This makes for easy testing using httptest
func NewServerRoutes() *http.ServeMux {
	m := http.NewServeMux()

	m.Handle("/", middleware(NotFoundHandler))
	m.Handle("/healthcheck", middleware(HealthCheckHandler))

	m.Handle("/publickey", middleware(PublicKeyHandler))

	m.Handle("/session", middleware(SessionHandler))
	m.Handle("/session/keys", middleware(KeysHandler))
	m.Handle("/session/access_token", middleware(AccessTokenHandler))
	m.Handle("/session/oauth", middleware(SessionUserTokensHandler))
	m.Handle("/session/oauth/github/repoaccess", middleware(GithubRepoAccessHandler))
	m.Handle("/session/communities", middleware(SessionCommunitiesHandler))

	m.Handle("/jwt", middleware(JwtHandler))
	m.Handle("/logout", middleware(LogoutHandler))

	m.Handle("/users", middleware(UsersHandler))
	m.Handle("/users/", middleware(UserHandler))
	m.Handle("/users/search", middleware(UsersSearchHandler))
	m.Handle("/users/communities", middleware(UserCommunitiesHandler))
	m.Handle("/communities/users", middleware(CommunityMembersHandler))

	m.Handle("/oauth/github", middleware(GithubOauthHandler))
	m.Handle("/oauth/github/callback", middleware(GithubOAuthCallbackHandler))

	// m.Handle("/reset", middleware(ResetPasswordHandler))
	// m.Handle("/reset/", middleware(ResetPasswordHandler))

	m.HandleFunc("/.well-known/acme-challenge/", CertbotHandler)

	return m
}
