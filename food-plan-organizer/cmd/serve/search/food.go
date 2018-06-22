package search

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gophergala/food-plan-organizer/models"
)

type searchServer struct {
	DB *sql.DB
}

func (s *searchServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	params, _ := url.ParseQuery(req.URL.RawQuery)

	query, ok := params["q"]
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rows, err := s.DB.Query(`SELECT foods.*, food_groups.* FROM foods INNER JOIN food_groups ON food_groups.id = foods.food_group_id WHERE foods.name LIKE $1 OR short_name LIKE $1 OR common_name LIKE $1 OR scientific_name LIKE $1`, fmt.Sprintf("%%%v%%", query[0]))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var foods []models.Food
	for rows.Next() {
		var food models.Food
		var err error
		if food, err = models.ScanFood(rows); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		foods = append(foods, food)
	}

	var enc = json.NewEncoder(rw)
	enc.Encode(foods)
}

func NewFoodSearchServer(db *sql.DB) http.Handler {
	return &searchServer{
		DB: db,
	}
}
