package handles

import (
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IndexHandler(c *gin.Context) {
	ctx := pongo2.Context{
		"title": "pairparty.io",
	}
	c.HTML(http.StatusOK, "templates/pages/index.html", ctx)
	/*
		ajax := r.Header.Get("X-PUSH")
		log.Println(context.Get(r, UserKey))
		if ajax != "" {
			// serve the partial
			err := index_partial_tmpl.ExecuteWriter(pongo2.Context{"query": r.FormValue("query")}, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			// serve a page
			err := index_page_tmpl.ExecuteWriter(pongo2.Context{"query": r.FormValue("query")}, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	*/
}
