package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func NewExtractorAPI(extractor *Extractor) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.Header().Set("Allow", "GET")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/extract", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Set("Allow", "POST")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		postBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errlog.LogFromClientRequest(map[string]interface{}{
				"error":  fmt.Sprintf("Error reading post body: %v", err),
				"status": http.StatusInternalServerError,
			}, r)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		matchedTermIds := extractor.Extract(string(postBody))
		if matchedTermIds == nil {
			errlog.LogFromClientRequest(map[string]interface{}{
				"error":  "matchedTermIds was nil",
				"status": http.StatusInternalServerError,
			}, r)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var marshalled []byte
		marshalled, err = json.Marshal(matchedTermIds)
		if err != nil {
			errlog.LogFromClientRequest(map[string]interface{}{
				"error":  fmt.Sprintf("Failed to marshal matched terms to Json: %v", err),
				"status": http.StatusInternalServerError,
			}, r)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(marshalled)
	})

	return mux
}
