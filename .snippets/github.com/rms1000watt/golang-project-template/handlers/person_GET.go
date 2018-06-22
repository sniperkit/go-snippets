package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// PersonHandlerGET retrieves person data based on ID
func PersonHandlerGET(w http.ResponseWriter, r *http.Request) {
	log.Debug("Starting PersonHandlerGET")
	defer log.Debug("Finished PersonHandlerGET")

	in := Person{}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		log.Debug("Error decoding r.Body:", err)
		http.Error(w, "Malformed input", http.StatusInternalServerError)
		return
	}

	out := PersonMap[in.ID]

	jsonBytes, err := json.Marshal(out)
	if err != nil {
		log.Error("Failed marshalling to JSON:", err)
		http.Error(w, "JSON Marshal Error", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(jsonBytes); err != nil {
		log.Error("Failed writing to response writer:", err)
		http.Error(w, "Failed writing to output", http.StatusInternalServerError)
		return
	}
}
