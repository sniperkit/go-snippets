package api

import (
	"github.com/infinitbyte/framework/core/api/router"
	"github.com/infinitbyte/framework/core/stats"
	"net/http"
)

// StatsAction return stats information
func (handler *API) StatsAction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	m := stats.StatsAll()
	handler.WriteJSONHeader(w)
	handler.Write(w, *m)
}

func (handler *API) FaviconAction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Write([]byte("."))
}
