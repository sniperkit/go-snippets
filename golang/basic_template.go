package main

import (
	"html/template"
	"log"
	"net/http"
)

func MainPage(w http.ResponseWriter, r *http.Request) {
	tpl := `<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<title>{{.Title}}</title>
		</head>
		<body>
			{{range .Items}}<div>{{.}}</div>{{else}}<div>no row</div>{{end}}
		</body>
	</html>`

	t, err := template.New("welcome").Parse(tpl)
	if err != nil {
		log.Println(err)
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
	log.Println("Starting Server")
	http.HandleFunc("/", MainPage)
	http.ListenAndServe(":8080", nil)
}
