package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
)

func MainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/index.html", "template/footer.html", "template/header.html")
	//t, err := template.ParseGlob("template/*")
	if err != nil {
		log.Println(err)
		io.WriteString(w, err.Error())
		return
	}

	data := struct {
		Title string
		Items []string
	}{
		Title: "First Page",
		Items: []string{
			"My photos",
			"My movie",
			"My video",
			"My audio",
		},
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
	}

}

func main() {
	http.HandleFunc("/", MainPage)
	http.ListenAndServe(":8080", nil)
}
