package router

import (
	"encoding/json"
	"net/http"
)

func PanicIf(err error, res http.ResponseWriter) {
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func respondWith(json_map interface{}, status int, res http.ResponseWriter) {
	json_response, err := json.Marshal(json_map)
	PanicIf(err, res)
	res.WriteHeader(status)
	res.Write(json_response)
}
