package controllers

import (
	"fmt"
	"github.com/rsmohamad/comp4321/database"
	"github.com/rsmohamad/comp4321/models"
	"github.com/rsmohamad/comp4321/retrieval"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"
)

var resultTemplate *template.Template
var keywordsTemplate *template.Template
var prefixes []string
var keywords map[string][]string

type KeywordsView struct {
	Prefixes []string
	Keywords map[string][]string
}

func loadKeywords() (map[string][]string, []string) {
	v, _ := database.LoadViewer("index.db")
	defer v.Close()

	k := v.GetKeywords()
	prefixes := make([]string, 0)
	keywords := make(map[string][]string)

	for _, word := range k {
		firstLetter := string(word[0])
		if keywords[firstLetter] == nil {
			keywords[firstLetter] = make([]string, 0)
			prefixes = append(prefixes, firstLetter)
		}
		keywords[firstLetter] = append(keywords[firstLetter], word)
	}

	sort.Strings(prefixes)
	return keywords, prefixes
}

func keywordsHandler(w http.ResponseWriter, r *http.Request) {
	if keywords == nil {
		keywords, prefixes = loadKeywords()
	}

	keywordsTemplate.Execute(w, KeywordsView{prefixes, keywords})
}

func nestedHandler(w http.ResponseWriter, r *http.Request) {
	viewModel := models.ResultView{}
	haystack := r.URL.Query().Get("haystack")
	needle := r.URL.Query().Get("needle")

	startSearch := time.Now()
	se := retrieval.NewSearchEngine("index.db")
	viewModel.Query = fmt.Sprintf("<%s> INSIDE <%s>", needle, haystack)
	viewModel.Results = se.RetrieveNested(haystack, needle)
	viewModel.TotalResults = len(viewModel.Results)
	se.Close()
	elapsed := time.Since(startSearch)

	log.Println(fmt.Sprintf("[%s] [%s] [%s]", r.RemoteAddr, viewModel.Query, elapsed))
	resultTemplate.ExecuteTemplate(w, "resultView", viewModel)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	viewModel := models.ResultView{}
	queries := r.URL.Query().Get("keywords")
	pagerank := r.URL.Query().Get("pagerank")
	fmt.Println(pagerank)

	userId := database.GetCookieInstance().GetCookieId(r)
	database.GetCookieInstance().SetCookieResponse(userId, w)
	database.GetCookieInstance().AddQuery(userId, queries)

	startSearch := time.Now()
	se := retrieval.NewSearchEngine("index.db")
	viewModel.Query = queries

	if pagerank == "on" {
		viewModel.Results = se.RetrievePageRank(queries)
	} else {
		viewModel.Results = se.RetrievePhrase(queries)
	}

	viewModel.TotalResults = len(viewModel.Results)
	se.Close()
	elapsed := time.Since(startSearch)

	log.Println(fmt.Sprintf("[%s] [%s] [%s]", r.RemoteAddr, queries, elapsed))
	resultTemplate.ExecuteTemplate(w, "resultView", viewModel)
}

func LoadSearch() {
	keywordsTemplate, _ = template.ParseFiles("views/keywordsView.html")
	resultTemplate, _ = template.ParseFiles("views/resultView.html", "views/documentView.html")
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/search/nested/", nestedHandler)
	http.HandleFunc("/search/keywords/", keywordsHandler)
}
