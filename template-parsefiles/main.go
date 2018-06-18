package main

import (
	"html/template"
	"log"
	"os"
	"time"
)

type Page struct {
	Title       string
	Description string
	Date        string
	Content     string
}

func main() {

	tpl, err := template.ParseFiles("tpl/index.html")
	checkErr(err)

	home := Page{
		Title:       "Home",
		Description: "SEO Description for Home",
		Date:        time.Now().Local().Format("02.01.2006"),
		Content:     "Welcome on my Homepage"}

	err = tpl.Execute(os.Stdout, home)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
}
