package main

import (
	"flag"
	"fmt"
	"github.com/gophergala/NextMatch/social/instagram"
	"github.com/gophergala/NextMatch/updater/xmlstats"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	shortf = "20060102"
	lonfg  = "2006-01-2T15:04:05-07:00"
)

var port = flag.String("p", "80", "the port on wich we're serving")

func init() {
	addTfunc("parse", time.Parse)
	addTfunc("now", time.Now)
	addTfunc("lower", strings.ToLower)
}

func handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sport, ok := vars[`sport`]
	if !ok {
		// default to nba
		sport = `nba`
	}

	e, err := xmlstats.BySport(sport)
	if err != nil {
		log.Printf("Didn't get data :( err = %v", err)
	}

	renderArgs := args{
		"events": e,
		"title":  "Home",
		`sport`:  sport,
	}

	execT(w, "home", renderArgs)
}

func reload(rw http.ResponseWriter, req *http.Request) {
	loadTmpl()
	http.Redirect(rw, req, "/", 302)
}

func main() {

	xmlstats.Token = os.Getenv("XMLSTATS_TOKEN")
	if len(xmlstats.Token) == 0 {
		log.Fatal("Specify XMLSTATS_TOKEN environment variable")
	}

	flag.Parse()
	loadTmpl()
	r := mux.NewRouter()
	r.HandleFunc(`/sport/{sport}`, handler)
	r.HandleFunc(`/details/{sport}/{id}`, eventScore)
	r.HandleFunc(`/{sport}/{home}-vs-{away}`, showGameDetails)
	r.HandleFunc(`/sport/{name}/{date}`, sportHandle)
	r.HandleFunc(`/sport/{name}`, sportHandle)
	r.HandleFunc("/refresh", reload)
	r.HandleFunc("/", handler)
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	r.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
	    	execT(w, "about", nil)
	 })

	r.HandleFunc("/404", thatsA404)
	http.Handle("/static/", static(http.FileServer(http.Dir("."))))
	http.Handle(`/`, r)
	r.NotFoundHandler = new(fof)
	log.Println("Yes M'Lord. I'm ready.")

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func static(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.Replace(r.URL.Path, "/static", "/resources", 1)
		h.ServeHTTP(w, r)
	}
}

func sportHandle(w http.ResponseWriter, req *http.Request) {
	e, err := getSport(req)
	if err != nil {
		log.Printf("Didn't get data :( err = %v", err)
	}

	renderArgs := args{"events": e}

	execT(w, "events", renderArgs)
}

func showGameDetails(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	sport := vars["sport"]
	home := vars["home"]
	away := vars["away"]
	e, err := xmlstats.BySport(sport)
	//log.Println(e)
	if err != nil {
		log.Printf("Didn't get data :( err = %v", err)
	}
	tag := fmt.Sprintf("%svs%s", strings.Split(home, " ")[0], strings.Split(away, " ")[0])
	o, err := instagram.ByTag(tag)
	if err != nil {
		log.Printf("Awwm :( dint' get anything for %s, error: %v", tag, err)
	}
	//log.Println(o)
	renderArgs := args{
		"events": e,
		"images": o,
		"title":  fmt.Sprintf("%s vs %s", home, away),
	}

	execT(w, "details", renderArgs)
}

func getSport(req *http.Request) (xmlstats.Events, error) {
	vars := mux.Vars(req)

	name := vars[`name`]
	date, ok := vars[`date`]
	if !ok {
		date = time.Now().Format(shortf)
	}

	return xmlstats.BySport(name, date)
}

// 404 shenanigans
func thatsA404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	execT(w, "404", nil)
}

type fof int

func (f *fof) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	thatsA404(w, r)
}

func eventScore(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	id := vars[`id`]
	sport := vars[`sport`]

	r, err := xmlstats.Score(sport, id)
	if err != nil {
		log.Printf("err on details: ", err)
	}

	e, err := xmlstats.BySport(sport)
	if err != nil {
		log.Println(`event list error `, err)
	}

	slug := fmt.Sprintf("%svs%s", r.AwayTeam.LastName, r.HomeTeam.LastName)
	if slug == `vs` {
		o := e.ById(id)
		if o == nil {
			slug = ``
		} else {
			slug = fmt.Sprintf("%svs%s", o.AwayTeam.LastName, o.HomeTeam.LastName)
		}
	}

	images, err := instagram.ByTag(slug)
	if err != nil {
		log.Printf("instagram err on details: %+v", err)
	}

	renderArgs := args{
		`event`:  r,
		`events`: e,
		`title`:  `details`,
		`images`: images,
		`sport`:  sport,
	}

	execT(w, `details`, renderArgs)
}
