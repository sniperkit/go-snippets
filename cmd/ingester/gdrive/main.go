package main

import (
	"flag"
	"io/ioutil"

	"github.com/keelerh/omniscience/cmd/ingester/common"
	"github.com/keelerh/omniscience/cmd/ingester/gdrive/lib"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
)

var (
	fGoogleServiceAccountFilePath = flag.String(
		"google_service_account_file_path",
		"",
		"The path to the Google Drive service account JSON file.")
)

func readGoogleServiceAccountCfg() (*jwt.Config, error) {
	secret, err := ioutil.ReadFile(*fGoogleServiceAccountFilePath)
	if err != nil {
		return nil, err
	}

	config, err := google.JWTConfigFromJSON(secret, drive.DriveReadonlyScope)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func createGDriveFetcher() (common.DocumentFetcher, error) {
	cfg, err := readGoogleServiceAccountCfg()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read Google service account configuration: %v", *fGoogleServiceAccountFilePath)
	}

	return lib.NewGoogleDrive(cfg), nil
}

func main() {
	common.CreateIngesterCLI(createGDriveFetcher)
}
