package show

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gophergala/food-plan-organizer/models"
)

type showNutrientServer struct {
	DB *sql.DB
}

func loadNutrientDefinitions(DB *sql.DB) ([]models.NutrientDefinition, error) {
	rows, err := DB.Query(`SELECT * FROM nutrient_definitions`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ns []models.NutrientDefinition
	for rows.Next() {
		var n models.NutrientDefinition
		if err = rows.Scan(&n.NutrientID, &n.Units, &n.Tagname, &n.Description, &n.DecimalPlaces); err != nil {
			return nil, err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func (s *showNutrientServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var enc = json.NewEncoder(rw)
	var ns []models.NutrientDefinition
	var err error
	if ns, err = loadNutrientDefinitions(s.DB); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	enc.Encode(ns)
}

func NewNutrientShowServer(db *sql.DB) http.Handler {
	return &showNutrientServer{
		DB: db,
	}
}
