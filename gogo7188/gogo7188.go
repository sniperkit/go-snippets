package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/flosch/pongo2"
	"github.com/hhatto/klip"
	gzip "github.com/lidashuang/goji_gzip"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var CLIPPINGS_FILE string
var ROOT string

type KindleClippings []*klip.KindleClipping

func (c KindleClippings) Len() int {
	return len(c)
}

func (c KindleClippings) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c KindleClippings) Less(i, j int) bool {
	return c[i].AddedOn.Before(c[j].AddedOn)
}

func getBook(c web.C, w http.ResponseWriter, r *http.Request) {
	k, _ := klip.Load(CLIPPINGS_FILE)
	var clips KindleClippings
	bookTitle := c.URLParams["name"]
	for i := range k {
		if k[i].Meta.Type == klip.Highlight && k[i].Title == bookTitle {
			clips = append(clips, &k[i])
		}
	}
	sort.Sort(sort.Reverse(clips))

	tpl, _ := pongo2.DefaultSet.FromFile(ROOT + "/templates/index.tpl")
	tpl.ExecuteWriter(pongo2.Context{"clips": clips, "title": bookTitle}, w)
}

func getIndex(c web.C, w http.ResponseWriter, r *http.Request) {
	k, _ := klip.Load(CLIPPINGS_FILE)
	var clips KindleClippings
	for i := range k {
		if k[i].Meta.Type == klip.Highlight {
			clips = append(clips, &k[i])
		}
	}
	sort.Sort(sort.Reverse(clips))

	tpl, _ := pongo2.DefaultSet.FromFile(ROOT + "/templates/index.tpl")
	tpl.ExecuteWriter(pongo2.Context{"clips": clips}, w)
}

func main() {
	flag.StringVar(&CLIPPINGS_FILE, "filename", "./clippings.txt", "My Clippings.txt file path")
	ROOT, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	goji.Use(gzip.GzipHandler)

	goji.Get("/book/:name", getBook)
	goji.Get("/", getIndex)
	goji.Get("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir(ROOT+"/static"))))
	goji.Serve()
}
