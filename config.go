package main

import (
	"fmt"
	conf "github.com/datatogether/config"
	"github.com/datatogether/core"
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
	// port to listen on, will be read from PORT env variable if present.
	Port string

	// root url for service
	UrlRoot string

	// url of postgres app db
	PostgresDbUrl string

	// Public Key to use for signing metablocks. required.
	PublicKey string

	// TLS (HTTPS) enable support via LetsEncrypt, default false
	// should be true in production
	TLS bool

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

	// seed        = flag.String("seed", "", "seed URL")
	// cancelAfter = flag.Duration("cancelafter", 0, "automatically cancel the fetchbot after a given time")
	// cancelAtURL = flag.String("cancelat", "", "automatically cancel the fetchbot at a given URL")
	// stopAfter   = flag.Duration("stopafter", 0, "automatically stop the fetchbot after a given time")
	// stopAtURL   = flag.String("stopat", "", "automatically stop the fetchbot at a given URL")
	// memStats    = flag.Duration("memstats", 0, "display memory statistics at a given interval")

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
	// CertbotResponse is only for doing manual SSL certificate generation via LetsEncrypt.
	CertbotResponse string
}

// initConfig pulls configuration from config.json
// func initConfig(mode string) (cfg *config, err error) {
// 	cfg = &config{}

// 	if err := loadConfigFile(mode, cfg); err != nil {
// 		return cfg, err
// 	}

// 	// override config settings with env settings, passing in the current configuration
// 	// as the default. This has the effect of leaving the config.json value unchanged
// 	// if the env variable is empty
// 	cfg.Port = readEnvString("PORT", cfg.Port)
// 	cfg.UrlRoot = readEnvString("URL_ROOT", cfg.UrlRoot)
// 	cfg.PublicKey = readEnvString("PUBLIC_KEY", cfg.PublicKey)
// 	cfg.TLS = readEnvBool("TLS", cfg.TLS)
// 	cfg.PostgresDbUrl = readEnvString("POSTGRES_DB_URL", cfg.PostgresDbUrl)
// 	cfg.HttpAuthUsername = readEnvString("HTTP_AUTH_USERNAME", cfg.HttpAuthUsername)
// 	cfg.HttpAuthPassword = readEnvString("HTTP_AUTH_PASSWORD", cfg.HttpAuthPassword)
// 	cfg.AwsAccessKeyId = readEnvString("AWS_ACCESS_KEY_ID", cfg.AwsAccessKeyId)
// 	cfg.AwsSecretAccessKey = readEnvString("AWS_SECRET_ACCESS_KEY", cfg.AwsSecretAccessKey)
// 	cfg.AwsRegion = readEnvString("AWS_REGION", cfg.AwsRegion)
// 	cfg.AwsS3BucketName = readEnvString("AWS_S3_BUCKET_NAME", cfg.AwsS3BucketName)
// 	cfg.AwsS3BucketPath = readEnvString("AWS_S3_BUCKET_PATH", cfg.AwsS3BucketPath)
// 	cfg.CertbotResponse = readEnvString("CERTBOT_RESPONSE", cfg.CertbotResponse)
// 	// cfg.StaleDuration = readEnvInt("STALE_DURATION", cfg.StaleDuration)

// 	// make sure port is set
// 	if cfg.Port == "" {
// 		cfg.Port = "8080"
// 	}

// 	err = requireConfigStrings(map[string]string{
// 		"PORT":            cfg.Port,
// 		"POSTGRES_DB_URL": cfg.PostgresDbUrl,
// 		"PUBLIC_KEY":      cfg.PublicKey,
// 	})

// 	// transfer settings to core library
// 	core.AwsRegion = cfg.AwsRegion
// 	core.AwsAccessKeyId = cfg.AwsAccessKeyId
// 	core.AwsS3BucketName = cfg.AwsS3BucketName
// 	core.AwsS3BucketPath = cfg.AwsS3BucketPath
// 	core.AwsSecretAccessKey = cfg.AwsSecretAccessKey

// 	return
// }

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

	err = requireConfigStrings(map[string]string{
		"PORT":            cfg.Port,
		"POSTGRES_DB_URL": cfg.PostgresDbUrl,
	})

	// transfer settings to core library
	core.AwsRegion = cfg.AwsRegion
	core.AwsAccessKeyId = cfg.AwsAccessKeyId
	core.AwsS3BucketName = cfg.AwsS3BucketName
	core.AwsS3BucketPath = cfg.AwsS3BucketPath
	core.AwsSecretAccessKey = cfg.AwsSecretAccessKey

	// output to stdout in dev mode
	if mode == DEVELOP_MODE {
		log.Out = os.Stdout
	}

	return
}

func packagePath(path string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src/github.com/datatogether/content", path)
}

// readEnvString reads key from the environment, returns def if empty
func readEnvString(key, def string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return def
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
