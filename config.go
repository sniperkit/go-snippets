package main

import (
	"fmt"
	conf "github.com/datatogether/config"
	"html/template"
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
// 		2. config.[server_mode].json <- eg: config.test.json
// 		3. config.json
//
// env variables win, but can only set config who's json is ALL_CAPS
// it's totally fine to not have, say, config.develop.json defined, and just
// rely on a base config.json. But if you're in production mode & config.production.json
// exists, that will be read *instead* of config.json.
//
// configuration is read at startup and cannot be alterd without restarting the server.
type config struct {
	// server mode. one of ["develop","production","test"]
	Mode string
	// path to go source code
	Gopath string
	// port to listen on, will be read from PORT env variable if present.
	Port string
	// Title for html templates
	Title string
	// root url for service
	UrlRoot string

	// url of postgres app db
	PostgresDbUrl string
	// url of redis app db
	RedisUrl string

	// TasksServiceUrl is the url for performing tasks
	TasksServiceUrl string
	// url for identity server
	IdentityServiceUrl string

	// Public Key to use for signing metablocks. required.
	PublicKey string

	// TLS (HTTPS) enable support via LetsEncrypt, default false
	// should be true in production
	TLS bool

	// Content Types to Store
	StoreContentTypes []string

	// read from env variable: AWS_REGION
	// the region your bucket is in, eg "us-east-1"
	AwsRegion string
	// read from env variable: AWS_S3_BUCKET_NAME
	// should be just the name of your bucket, no protocol prefixes or paths
	AwsS3BucketName string
	// read from env variable: AWS_ACCESS_KEY_ID
	AwsAccessKeyId string
	// read from env variable: AWS_SECRET_ACCESS_KEY
	AwsSecretAccessKey string
	// path to store & retrieve data from
	AwsS3BucketPath string

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
	// Segment Analytics API token for server-side analytics
	SegmentApiToken string
	// list of urls to webapp entry point(s)
	WebappScripts []string

	// CertbotResponse is only for doing manual SSL certificate generation
	// via LetsEncrypt.
	CertbotResponse string
}

// initConfig pulls configuration from config.json
func initConfig(mode string) (cfg *config, err error) {
	cfg = &config{Mode: mode}

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

	err = requireConfigStrings(map[string]string{
		"GOPATH":          cfg.Gopath,
		"PORT":            cfg.Port,
		"POSTGRES_DB_URL": cfg.PostgresDbUrl,
		"PUBLIC_KEY":      cfg.PublicKey,
	})

	templates = template.Must(template.ParseFiles(
		packagePath("views/profile.html"),
		packagePath("views/webapp.html"),
		packagePath("views/accessDenied.html"),
		packagePath("views/notFound.html"),
	))

	// output to stdout in dev mode
	if mode == DEVELOP_MODE {
		log.Out = os.Stdout
	}

	return
}

func packagePath(path string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src/github.com/datatogether/patchbay", path)
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
	log.Println("Mode:", cfg.Mode)
	log.Println("Gopath:", cfg.Gopath)
	log.Println("Port:", cfg.Port)
	log.Println("Title:", cfg.Title)
	log.Println("UrlRoot:", cfg.UrlRoot)
	// log.Println("PostgresDbUrl:", cfg.PostgresDbUrl)
	log.Println("RedisUrl:", cfg.RedisUrl)
	log.Println("TasksServiceUrl:", cfg.TasksServiceUrl)
}
