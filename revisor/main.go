package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var (
	PathList map[string]bool
	PathRoot string
)

type Page struct {
	Title     string
	Body      []byte
	Directory *map[string]bool
}

func (p *Page) save() error {
	filename := p.Title // + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title //+ ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Received: %s", r.URL.Path[1:])
	//I need to prevent navigation up a dir here with some regex
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		//fmt.Println("not found")
		p = &Page{Title: title}
	}
	//p.Directory = make([]string, len(PathList))
	//p.Directory = PathList[0:]
	p.Directory = &PathList
	t, _ := template.ParseFiles("index.html")
	//fmt.Println(p.Title)
	//fmt.Println(err)
	t.Execute(w, p)
}

func dirHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("dir.html")

	p := Page{}
	p.Title = r.URL.Path[len("/dir/"):]
	p.Directory = &PathList

	t.Execute(w, p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("directory rebuild failed")
		}
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}()

	filepath.Walk(PathRoot, visit)
}

func visit(path string, f os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("File or directory traversal error: %v\n", err)
	}

	//fmt.Printf("Visited: %s | %v\n", path, f.IsDir())
	//Files will link to the edit handler, directories to a different handler
	PathList[path] = f.IsDir() //[PathIndex] = path
	return nil
}

func main() {
	fmt.Println("So it goes")

	//should add configuration options to append to this path
	PathList = make(map[string]bool)
	PathRoot := "./website"
	err := filepath.Walk(PathRoot, visit)
	fmt.Printf("filepath.Walk() returned %v\n", err)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/dir/", dirHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
