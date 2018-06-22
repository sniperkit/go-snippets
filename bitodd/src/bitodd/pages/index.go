package pages

import (
	"bitodd/config"
	"bitodd/util"
	"net/http"
)

const indexURL = "/"

var indexTmpl = getTemplate("templates/index.html")

// Index Handler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmplModel := &util.TemplateModel{Params: make(map[interface{}]interface{}, 0)}

	tmplModel.Params["Keywords"] = config.GetConfig().Keywords

	util.RenderTemplate(indexTmpl, w, tmplModel)
}
