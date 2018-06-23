package core

import (
	"net/http"
)

func init() {
	// I expect someone else on the team to know what routing scheme we want to use.
	http.HandleFunc("/submit/", submitHandler)
}
