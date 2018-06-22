package gorgonzola

import (
	"net/http"

	"github.com/gorilla/mux"
)

const firstPageJobsLimit = 100

func (g *Gorgonzola) indexHandler(w http.ResponseWriter, r *http.Request) error {
	tm := NewTemplate(w)
	jobs, err := g.storage.GetJobs(r, firstPageJobsLimit)
	if err != nil {
		return err
	}
	tm.set("jobs", jobs)
	return tm.render("templates/layout.html", "templates/index.html")
}

func (g *Gorgonzola) jobHandler(w http.ResponseWriter, r *http.Request) error {
	tm := NewTemplate(w)
	vars := mux.Vars(r)
	job, err := g.storage.GetJob(r, vars["key"])
	if err != nil {
		return err
	}
	tm.set("job", job)
	return tm.render("templates/layout.html", "templates/job.html")
}

func (g *Gorgonzola) addHandler(w http.ResponseWriter, r *http.Request) error {
	tm := NewTemplate(w)
	if r.Method == "POST" {
		url := r.FormValue("url")
		if url != "" {
			if err := validateURL(url); err != nil {
				return err
			}
			err := g.storage.AddURL(r, url)
			if err != nil {
				return err
			}
			http.Redirect(w, r, "/thankyou.html", http.StatusFound)
		}
	}
	return tm.render("templates/layout.html", "templates/add.html")
}

func (g *Gorgonzola) thankyouHandler(w http.ResponseWriter, r *http.Request) error {
	tm := NewTemplate(w)
	return tm.render("templates/layout.html", "templates/thankyou.html")
}

func (g *Gorgonzola) taskUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return g.storage.Update(r)
}
