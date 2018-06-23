package controller

import (
	"../common"
	"net/http"
)

type Home struct {
	http.Handler
}

func (t *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	common.Render().HTML(w, http.StatusOK, "index", nil)
}
