package util

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

func HandleError(rw http.ResponseWriter, req *http.Request, err, apperr error, isAppErr bool) bool {
	if err == nil {
		return false
	}

	if status := http.StatusBadRequest; isAppErr {
		http.Error(rw, apperr.Error(), status)
	} else {
		status = http.StatusInternalServerError

		logError(err, req)

		msg := fmt.Sprintf("INTERNAL ERROR: %s", err.Error())

		http.Error(rw, msg, status)
	}

	return true
}

func HandlePanic(err interface{}, rw http.ResponseWriter, req *http.Request) {
	logError(err, req)

	http.Error(rw, "SYSTEM ERROR", http.StatusInternalServerError)
}

func logError(err interface{}, req *http.Request) {
	var stack [4096]byte

	runtime.Stack(stack[:], false)

	log.Printf(
		"Handler for %s[url:%s, data:%s] has failed with %s\nStack trace:\n %s\n",
		req.Method, req.URL, req.Form, err, stack[:])
}
