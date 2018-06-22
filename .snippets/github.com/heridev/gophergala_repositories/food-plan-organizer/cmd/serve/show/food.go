package show

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gophergala/food-plan-organizer/models"
)

type showServer struct {
	DB *sql.DB
}

func loadNutrients(f *models.Food, DB *sql.DB) ([]models.Nutrient, error) {
	rows, err := DB.Query(`SELECT * FROM nutrients INNER JOIN nutrient_definitions ON nutrient_definitions.nutrient_id = nutrients.nutrient_id WHERE food_id = $1 AND nutrient_value > 0`, f.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ns []models.Nutrient
	for rows.Next() {
		var n models.Nutrient
		if n, err = models.ScanNutrient(rows); err != nil {
			fmt.Printf("err: %v\n", err)
			return nil, err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func (s showServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	params, _ := url.ParseQuery(req.URL.RawQuery)

	id, ok := params["id"]
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var row = s.DB.QueryRow(`SELECT foods.*, food_groups.* FROM foods INNER JOIN food_groups ON food_groups.id = foods.food_group_id WHERE foods.id = $1`, id[0])

	var food models.Food
	var err error
	if food, err = models.ScanFood(row); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if food.Nutrients, err = loadNutrients(&food, s.DB); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var enc = json.NewEncoder(rw)
	enc.Encode(food)
}

func NewFoodShowServer(db *sql.DB) http.Handler {
	return showServer{
		DB: db,
	}
}
