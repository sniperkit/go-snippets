package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"

	"github.com/GopherGala/i_love_indexes/api"
	"github.com/GopherGala/i_love_indexes/crawler"
	"github.com/Scalingo/go-workers"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/shaoshing/train"
)

func main() {
	isWeb := flag.Bool("web", false, "run web server")
	isCrawler := flag.Bool("crawler", false, "run async crawler")
	flag.Parse()
	if *isWeb {
		mainWebServer()
	} else if *isCrawler {
		mainCrawler()
	} else {
		log.Fatalln("Invalid type of process, precis -web, -crawler")
	}
}

func setupTrain(router *mux.Router) {
	train.SetFileServer()
	router.Handle(train.Config.AssetsUrl+"{any:.*}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		train.ServeRequest(w, r)
	}))
}

func mainWebServer() {
	configureWorkers()
	router := mux.NewRouter()
	router.HandleFunc("/", handleIndex)
	router.HandleFunc("/about", handleAbout)
	router.Handle("/api/{any:.*}", api.NewAPI())
	setupTrain(router)

	staticHandler := negroni.Classic()
	router.Handle("/{any:.*}", staticHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	log.Println("Listening on", port)
	log.Fatalln(http.ListenAndServe(":"+port, router))
}

func mainCrawler() {
	configureWorkers()
	workers.Process("index-crawler", crawler.CrawlWorker, 10)
	workers.Run()
}

func configureWorkers() {
	redisURL := "redis://localhost:6379"
	if os.Getenv("REDIS_URL") != "" {
		redisURL = os.Getenv("REDIS_URL")
	}
	u, err := url.Parse(redisURL)
	if err != nil {
		log.Fatalln(err)
	}
	password := ""
	if u.User != nil {
		p, ok := u.User.Password()
		if ok {
			password = p
		}
	}

	workers.Configure(map[string]string{
		"server":   u.Host,
		"password": password,
		"database": "0",
		"pool":     "30",
		"process":  "1",
	})
}

func handleIndex(res http.ResponseWriter, req *http.Request) {
	tpl := template.New("index.html")

	// Adding train helpers
	tpl.Funcs(template.FuncMap{
		"javascript_tag":            train.JavascriptTag,
		"stylesheet_tag":            train.StylesheetTag,
		"stylesheet_tag_with_param": train.StylesheetTagWithParam,
	})
	_, err := tpl.ParseFiles("views/index.html")
	if err != nil {
		res.WriteHeader(500)
		fmt.Fprintln(res, err)
		return
	}

	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, nil)
	if err != nil {
		res.WriteHeader(500)
		fmt.Fprintln(res, err)
		return
	}

	res.WriteHeader(200)
	buffer.WriteTo(res)
}

func handleAbout(res http.ResponseWriter, req *http.Request) {
	tpl := template.New("about.html")

	// Adding train helpers
	tpl.Funcs(template.FuncMap{
		"javascript_tag":            train.JavascriptTag,
		"stylesheet_tag":            train.StylesheetTag,
		"stylesheet_tag_with_param": train.StylesheetTagWithParam,
	})
	_, err := tpl.ParseFiles("views/about.html")
	if err != nil {
		res.WriteHeader(500)
		fmt.Fprintln(res, err)
		return
	}

	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, nil)
	if err != nil {
		res.WriteHeader(500)
		fmt.Fprintln(res, err)
		return
	}

	res.WriteHeader(200)
	buffer.WriteTo(res)
}
