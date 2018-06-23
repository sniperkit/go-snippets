package main

import (
	"davebalmain.com/noshgobble/src/nutdb"
	"fmt"
	"github.com/blevesearch/bleve"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	IndexDirectory = "foods.bleve"
)

var index bleve.Index

type Nutrient struct {
	Name     string
	Quantity float32
	Units    string
}

func renderTemplate(w http.ResponseWriter, tmpl string, d interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl + ".html")
	if err != nil {
		log.Fatalf("Error loading template: %v\n", err)
	}
	t.Execute(w, d)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprintf(w, "Type a query in the search bar")
	}
	foodId, err := strconv.ParseInt(r.URL.Path[1:], 10, 32)
	if err != nil {
		qstring := r.URL.Path[1:]
		fmt.Fprintf(w, "Trying a search for %s instead\n", qstring)
		query := bleve.NewQueryStringQuery(qstring)
		searchRequest := bleve.NewSearchRequest(query)
		searchResult, _ := index.Search(searchRequest)
		fmt.Fprintf(w, "The result is %+v\n", searchResult)
	} else {
		if foodId < 0 || int(foodId) >= len(nutdb.FoodDb) {
			fmt.Fprintf(w, "Error retrieving food. FoodId %d does not exist!", foodId)
			fmt.Fprintf(w, "FoodId must be between 0 and %d!", len(nutdb.FoodDb))
		} else {
			food := nutdb.FoodDb[int32(foodId)]
			nutrients := make([]Nutrient, len(food.NutrientQtys))
			for i, nq := range food.NutrientQtys {
				nutrients[i].Name = nutdb.NutrientDb[nq.NutrientId].Description
				nutrients[i].Quantity = nq.Quantity
				nutrients[i].Units = nutdb.NutrientDb[nq.NutrientId].Units
			}
			d := struct {
				Food      nutdb.Food
				Nutrients []Nutrient
			}{
				food,
				nutrients,
			}
			renderTemplate(w, "view", d)
		}
	}
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "analyze", nil)
}

func main() {
	nutdb.InitializeFoodDb()
	if _, err := os.Stat(IndexDirectory); err != nil && os.IsNotExist(err) {
		indexFoodDb()
	}
	index, _ = bleve.Open(IndexDirectory)
	fmt.Println("Now listening at http://localhost:8080/")
	http.HandleFunc("/", handler)
	http.HandleFunc("/analyze/", analyzeHandler)
	http.ListenAndServe(":8080", nil)
}
