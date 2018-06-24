package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// generic error response handler, BYO response code
func ErrRes(w http.ResponseWriter, e error) error {
	if cfg.Environment != TEST_MODE && e != nil {
		log.Info(e.Error())
		if cfg.Environment == PRODUCTION_MODE {
			log.Info(fmt.Sprintf("error: %s\n", e.Error()))
		}
	}

	if err, ok := e.(*Error); ok {
		env := map[string]interface{}{
			"meta": err,
		}
		w.WriteHeader(err.HttpCode)
		return json.NewEncoder(w).Encode(env)
	} else {
		env := map[string]interface{}{
			"meta": map[string]interface{}{
				"status":  500,
				"message": e.Error(),
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(w).Encode(env)
	}
}

func Res(w http.ResponseWriter, envelope bool, data interface{}) error {
	if envelope {
		data = map[string]interface{}{
			"meta": map[string]interface{}{
				"code": http.StatusOK,
			},
			"data": data,
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Info(err)
	}
	return err
}

func MessageResponse(w http.ResponseWriter, message string, data interface{}) error {
	res := map[string]interface{}{
		"meta": map[string]interface{}{
			"code":    http.StatusOK,
			"message": message,
		},
		"data": data,
	}

	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(res)
}

func PageRes(w http.ResponseWriter, data interface{}, p Page) error {
	res := map[string]interface{}{
		"meta": map[string]interface{}{
			"code": http.StatusOK,
		},
		"data":       data,
		"pagination": p,
	}

	return json.NewEncoder(w).Encode(res)
}
