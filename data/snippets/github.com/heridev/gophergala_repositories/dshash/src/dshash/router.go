package dshash

import (
	"appengine"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type WebContextHandler func(appengine.Context, http.ResponseWriter, *http.Request, httprouter.Params)

type HandlerAndContext interface {
	Handle(WebContextHandler) httprouter.Handle
}

type HandlerWithWebContext struct {
}

func (hwc HandlerWithWebContext) Handle(h WebContextHandler) httprouter.Handle {
	f := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h(appengine.NewContext(r), w, r, ps)
	}

	return httprouter.Handle(f)
}

func Router(hc HandlerAndContext) *httprouter.Router {
	router := httprouter.New()
	router.GET("/locations/:handler", hc.Handle(getHandler))
	router.POST("/locations", hc.Handle(postHandler))

	return router
}

func init() {
	http.Handle("/", Router(HandlerWithWebContext{}))
}
