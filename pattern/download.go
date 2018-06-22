package pattern

import (
	"log"
	"net/http"
)

// in this file, we will discuess a pattern for downloading same sources multiple times
// requirement: there should be only one real download request for same source
// ref, http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/

func downloadHandler(sourceURL string) {

	req, err := http.NewRequest("GET", sourceURL, nil)
	if err != nil {
		log.Printf("Create Request Error, err: %v", err)
		return
	}

	download(req)
}

func download(r *http.Request) {
	entry := r.URL.Path

	log.Println(entry)
}
