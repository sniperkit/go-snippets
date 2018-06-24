package main

import (
	"fmt"
	conf "github.com/datatogether/config"
	"os"
	"path/filepath"
)

// server modes
const (
	DEVELOP_MODE    = "develop"
	PRODUCTION_MODE = "production"
	TEST_MODE       = "test"
)

// config holds all configuration for the server. It pulls from three places (in order):
// 		1. environment variables
// 		2. .[MODE].env OR .env
//
// globally-set env variables win.
// it's totally fine to not have, say, .env.develop defined, and just
// rely on a base ".env" file. But if you're in production mode & ".env.production"
// exists, that will be read *instead* of .env
//
// configuration is read at startup and cannot be alterd without restarting the server.
type config struct {
	// configurable environment variable
	Environment string

	// port to listen on, will be read from PORT env variable if present.
	Port string

	// port to listen on for RPC calls
	RpcPort string

	// root url for service
	UrlRoot string

	// url of postgres app db. required.
	PostgresDbUrl string

	// Private Key to use for signing thinks like JWTs
	// to create: openssl genrsa -out app.rsa keysize
	PrivateKey string
	// Public Key to use for signing things like JWTs
	// openssl rsa -in app.rsa -pubout > app.rsa.pub
	PublicKey string

	// TLS (HTTPS) enable support via LetsEncrypt, default false
	// should be true in production
	TLS bool

	// because identity is an api-only server, FrontendUrl specifies
	// a user-friendly URL to redirect to
	FrontendUrl string

	// secret for encrypted session cookies. required.
	SessionSecret string
	// Key to store user's cookie under
	UserCookieKey string
	// domain to attach cookie to.
	UserCookieDomain string

	// support CORS signing from a list of origins
	AllowedOrigins []string

	// setting HTTP_AUTH_USERNAME & HTTP_AUTH_PASSWORD
	// will enable basic http auth for the server. This is a single
	// username & password that must be passed in with every request.
	// leaving these values blank will disable http auth
	// read from env variable: HTTP_AUTH_USERNAME
	HttpAuthUsername string
	// read from env variable: HTTP_AUTH_PASSWORD
	HttpAuthPassword string

	// if true, requests that have X-Forwarded-Proto: http will be redirected
	// to their https variant
	ProxyForceHttps bool
	// CertbotResponse is only for doing manual SSL certificate generation
	// via LetsEncrypt.
	CertbotResponse string

	GithubAppId     string
	GithubAppSecret string
}

// initConfig pulls configuration from config.json
func initConfig(mode string) (cfg *config, err error) {
	cfg = &config{}

	if path := configFilePath(mode, cfg); path != "" {
		log.Infof("loading config file: %s", filepath.Base(path))
		if err := conf.Load(cfg, path); err != nil {
			log.Info("error loading config:", err)
		}
	} else {
		if err := conf.Load(cfg); err != nil {
			log.Info("error loading config:", err)
		}
	}

	// make sure port is set
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.UserCookieKey == "" {
		// TODO - this should be set based on on url root value
		cfg.UserCookieKey = "identity.user"
	}

	err = requireConfigStrings(map[string]string{
		"PORT":            cfg.Port,
		"POSTGRES_DB_URL": cfg.PostgresDbUrl,
		"PUBLIC_KEY":      cfg.PublicKey,
		"PRIVATE_KEY":     cfg.PrivateKey,
	})

	// output to stdout in dev mode
	if mode == DEVELOP_MODE {
		log.Out = os.Stdout
	}

	return
}

func packagePath(path string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src/github.com/datatogether/identity", path)
}

// requireConfigStrings panics if any of the passed in values aren't set
func requireConfigStrings(values map[string]string) error {
	for key, value := range values {
		if value == "" {
			return fmt.Errorf("%s env variable or config key must be set", key)
		}
	}

	return nil
}

// checks for .[mode].env file to read configuration from if the file exists
// defaults to .env, returns "" if no file is present
func configFilePath(mode string, cfg *config) string {
	fileName := packagePath(fmt.Sprintf(".%s.env", mode))
	if !fileExists(fileName) {
		fileName = packagePath(".env")
		if !fileExists(fileName) {
			return ""
		}
	}
	return fileName
}

// Does this file exist?
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// outputs any notable settings to stdout
func printConfigInfo() {
	// TODO
}
