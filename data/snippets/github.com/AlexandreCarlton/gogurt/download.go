package gogurt

import (
	"io"
	"net/http"
	"log"
	"os"
	"path/filepath"
)

func Download(url string, dest string) error {
	log.Printf("Downloading url '%s' to file '%s'.\n", url, dest)

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	// Follow redirects. TODO - clean up.
	finalURL := response.Request.URL.String()
	redirectedResponse, err := http.Get(finalURL)
	if err != nil {
		return err
	}
	defer redirectedResponse.Body.Close()

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		log.Fatalf("Error creating directory '%s' for download:  %s", filepath.Dir(dest), err.Error())
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, redirectedResponse.Body)
	if err != nil {
		return err
	}
	return nil
}
