package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type cntTmpltCouple struct {
	cnt func(http.ResponseWriter, *http.Request) Model
	tmp string
}

// templateDependencies lists the necessary templates that need to be rendered
// for each template
var templateDependencies = map[string][]string{
	"about.html":            {"base.html"},
	"datasets.html":         {"base.html"},
	"datasets_view.html":    {"base.html"},
	"tasks.html":            {"base.html"},
	"error.html":            {"base.html"},
	"sm_heatmap.html":       {"base.html"},
	"coords_visual.html":    {"base.html"},
	"model_visual.html":     {"base.html"},
	"model_comparison.html": {"base.html"},

	// The rest are popups
	"forms/new_sm_form.html":      {},
	"forms/new_op_form.html":      {},
	"forms/new_mds_form.html":     {},
	"forms/new_dataset_form.html": {},
	"coords_view.html":            {},
	"forms/new_model_form.html":   {},
}

// routingControllerTemplates hold the controller and the respective template
// that need to be rendered for each possible path
var routingControllerTemplates = map[string]cntTmpltCouple{
	"datasets/":           {controllerDatasetList, "datasets.html"},
	"datasets/view":       {controllerDatasetView, "datasets_view.html"},
	"tasks/":              {controllerTasksList, "tasks.html"},
	"sm/visual":           {controllerSMVisual, "sm_heatmap.html"},
	"coords/visual":       {controllerCoordsVisual, "coords_visual.html"},
	"modeling/visual":     {controllerModelVisual, "model_visual.html"},
	"modeling/comparison": {controllerModelComparison, "model_comparison.html"},

	// forms
	"datasets/newsm": {controllerDatasetNewSM, "forms/new_sm_form.html"},
	"datasets/new":   {controllerDatasetNew, "forms/new_dataset_form.html"},
	"datasets/newop": {controllerDatasetNewOP, "forms/new_op_form.html"},
	"mds/run":        {controllerMDSRun, "forms/new_mds_form.html"},
	"coords/view":    {controllerCoordsView, "coords_view.html"},
	"modeling/new":   {controllerModelNew, "forms/new_model_form.html"},

	// No GUI urls
	"download/":       {controllerDownload, ""},
	"sm/csv":          {controllerSMtoCSV, ""},
	"sm/delete":       {controllerSMDelete, ""},
	"operator/run":    {controllerOperatorRun, ""},
	"operator/delete": {controllerOperatorDelete, ""},
	"modeling/delete": {controllerModelDelete, ""},
	"scores/text":     {controllerScoresText, ""},
	"datasets/delete": {controllerDatasetDelete, ""},

	// TODO: implement these URLs
	"about/":  {nil, "about.html"},
	"search/": {nil, ""}, // does nothing for now
}

func uiHandler(w http.ResponseWriter, r *http.Request) {
	cnt, t := selectControllerAndTemplate(r.URL.Path)
	var m Model
	if cnt != nil {
		m = cnt(w, r)
	}
	if t != nil {
		if t.Lookup("error.html") != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		t.Execute(w, m)
	}
}

func restHandler(w http.ResponseWriter, r *http.Request) {
	cnt, _ := selectControllerAndTemplate(r.URL.Path)
	r.URL.Path = strings.Replace(r.URL.Path, "/api", "", 1)
	var m Model
	if cnt != nil {
		m = cnt(w, r)
	} else {
		w.WriteHeader(http.StatusNotImplemented)
	}
	if m != nil {
		res, err := json.Marshal(m)
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func selectControllerAndTemplate(url string) (func(http.ResponseWriter, *http.Request) Model, *template.Template) {
	model, id, cmd := parseURL(url)
	if id != "" && cmd == "" { // default action is view
		cmd = "view"
	}
	route := model + "/" + cmd

	tmplt := "error.html"
	var cnt func(http.ResponseWriter, *http.Request) Model

	if coup, ok := routingControllerTemplates[route]; ok {
		cnt = coup.cnt
		tmplt = coup.tmp
	}
	t := loadTemplate(tmplt)
	return cnt, t
}

// loadTemplate attempts to load the specified template, else returns the
// error page
func loadTemplate(templateName string) *template.Template {
	if templateName == "" { // no template is needed
		return nil
	}
	deps, ok := templateDependencies[templateName]
	if !ok { // error
		deps = templateDependencies["error.html"]
		templateName = "error.html"
	}
	newDeps := make([]string, len(deps))
	copy(newDeps, deps)
	deps = newDeps
	for i := range deps {
		deps[i] = Conf.Server.Dirs.Templates + "/" + deps[i]
	}
	deps = append(deps, Conf.Server.Dirs.Templates+"/"+templateName)
	t, err := template.ParseFiles(deps...)
	if err != nil {
		log.Println(err)
	}
	return t
}

func parseURL(url string) (string, string, string) {
	url = strings.TrimPrefix(url, "/api")
	if url == "/" {
		url = "/datasets/"
	}
	arr := strings.Split(url, "/")[1:]
	if arr[len(arr)-1] == "" {
		arr = arr[0 : len(arr)-1]
	}
	model, id, cmd := arr[0], "", ""
	if len(arr) > 1 {
		id = arr[1]
	}
	if len(arr) > 2 {
		cmd = arr[2]
	}
	return model, id, cmd
}
