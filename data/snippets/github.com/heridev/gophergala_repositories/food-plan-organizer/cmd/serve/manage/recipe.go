package manage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gophergala/food-plan-organizer/models"
)

type recipeServer struct {
	userDB *sql.DB
	sr27DB *sql.DB
}

func loadNutrients(fID int32, DB *sql.DB) ([]models.Nutrient, error) {
	var rows, err = DB.Query(`SELECT * FROM nutrients INNER JOIN nutrient_definitions ON nutrient_definitions.nutrient_id = nutrients.nutrient_id WHERE food_id = $1 AND nutrient_value > 0`, fID)
	if err != nil {
		return nil, err
	}

	var nutrients []models.Nutrient
	var n models.Nutrient
	for rows.Next() {
		if n, err = models.ScanNutrient(rows); err != nil {
			return nil, err
		}
		nutrients = append(nutrients, n)
	}
	return nutrients, nil
}

func (rs *recipeServer) CreateRecipe(rw http.ResponseWriter, req *http.Request) {
	var dec = json.NewDecoder(req.Body)
	var recipe models.Recipe

	if err := dec.Decode(&recipe); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var enc = json.NewEncoder(rw)
	rs.userDB.Exec(`INSERT INTO recipes (name, description) VALUES ($1,$2);`, recipe.Name, recipe.Description)
	if err := rs.userDB.QueryRow(`SELECT last_insert_rowid() FROM recipes`).Scan(&recipe.ID); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]string{"sql_error": fmt.Sprintf("%v", err)})
		return
	}

	for i, ingredient := range recipe.Ingredients {
		rs.userDB.Exec(`INSERT INTO ingredients (recipe_id, food_id, unit, volume) VALUES ($1,$2,$3,$4)`, recipe.ID, ingredient.FoodID, ingredient.Unit, ingredient.Volume)
		rs.userDB.QueryRow(`SELECT last_insert_rowid() FROM ingredients`).Scan(&recipe.Ingredients[i].ID)
		rs.sr27DB.QueryRow(`SELECT name, nitrogen_factor, protein_factor, fat_factor, carbonhydrate_factor FROM foods WHERE id = $1`, ingredient.FoodID).Scan(
			&recipe.Ingredients[i].Name,
			&recipe.Ingredients[i].NitrogenFactor,
			&recipe.Ingredients[i].ProteinFactor,
			&recipe.Ingredients[i].FatFactor,
			&recipe.Ingredients[i].CarbonhydrateFactor)

		recipe.Ingredients[i].Nutrients, _ = loadNutrients(ingredient.FoodID, rs.sr27DB)
	}

	enc.Encode(recipe)
}

func (rs *recipeServer) ListRecipes(rw http.ResponseWriter, req *http.Request) {
	var recipes []models.Recipe
	var enc = json.NewEncoder(rw)
	var rows, err = rs.userDB.Query(`SELECT * FROM recipes`)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]string{"sql_error": fmt.Sprintf("%v", err)})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Recipe
		if err = rows.Scan(&r.ID, &r.Name, &r.Description); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]string{"sql_error": fmt.Sprintf("%v", err)})
			return
		}
		recipes = append(recipes, r)
	}

	enc.Encode(recipes)
}

func (rs *recipeServer) DeleteRecipe(rw http.ResponseWriter, req *http.Request) {
	params, _ := url.ParseQuery(req.URL.RawQuery)
	id, ok := params["id"]
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rs.userDB.Exec(`DELETE FROM ingredients WHERE recipe_id = $1`, id[0])
	rs.userDB.Exec(`DELETE FROM recipes WHERE id = $1`, id[0])
}

func (rs *recipeServer) ShowRecipe(rw http.ResponseWriter, req *http.Request) {
	params, _ := url.ParseQuery(req.URL.RawQuery)
	id, ok := params["id"]
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var r models.Recipe
	var enc = json.NewEncoder(rw)
	if err := rs.userDB.QueryRow(`SELECT * FROM recipes WHERE id = $1`, id[0]).Scan(&r.ID, &r.Name, &r.Description); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]string{"sql_error": fmt.Sprintf("%v", err)})
		return
	}

	var rows, err = rs.userDB.Query(`SELECT ingredients.unit, ingredients.volume, ingredients.food_id FROM ingredients WHERE ingredients.recipe_id = $1`, id[0])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]string{"sql_error": fmt.Sprintf("%v", err)})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var ing models.Ingredient
		if err = rows.Scan(&ing.Unit, &ing.Volume, &ing.FoodID); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]string{"sql_error": fmt.Sprintf("%v", err)})
			return
		}
		if err = rs.sr27DB.QueryRow(`SELECT name, nitrogen_factor, protein_factor, fat_factor, carbonhydrate_factor FROM foods WHERE id = $1`, ing.FoodID).Scan(&ing.Name, &ing.NitrogenFactor, &ing.ProteinFactor, &ing.FatFactor, &ing.CarbonhydrateFactor); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]string{"sql_error": fmt.Sprintf("%v", err)})
			return
		}
		if ing.Nutrients, err = loadNutrients(ing.FoodID, rs.sr27DB); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]string{"sql_error": fmt.Sprintf("%v", err)})
			return
		}
		r.Ingredients = append(r.Ingredients, ing)
	}

	enc.Encode(r)
}

func (rs *recipeServer) UpdateRecipe(rw http.ResponseWriter, req *http.Request) {
	params, _ := url.ParseQuery(req.URL.RawQuery)
	id, ok := params["id"]
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var dec = json.NewDecoder(req.Body)
	var recipe models.Recipe

	if err := dec.Decode(&recipe); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var enc = json.NewEncoder(rw)
	var tx, _ = rs.userDB.Begin()
	if _, err := tx.Exec(`UPDATE recipes SET name = $1, description = $2 WHERE id = $3;`, recipe.Name, recipe.Description, id[0]); err != nil {
		tx.Rollback()
		rw.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]string{"sql_error": fmt.Sprintf("%v", err)})
		return
	}

	if _, err := tx.Exec(`DELETE FROM ingredients WHERE recipe_id = $1`, id[0]); err != nil {
		tx.Rollback()
		rw.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]string{"sql_error": fmt.Sprintf("%v", err)})
		return
	}
	for i, ingredient := range recipe.Ingredients {
		if _, err := tx.Exec(`INSERT INTO ingredients (recipe_id, food_id, unit, volume) VALUES ($1,$2,$3,$4)`, id[0], ingredient.FoodID, ingredient.Unit, ingredient.Volume); err != nil {
			tx.Rollback()
			rw.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]string{"sql_error": fmt.Sprintf("%v", err)})
			return
		}
		if err := rs.sr27DB.QueryRow(`SELECT name, nitrogen_factor, protein_factor, fat_factor, carbonhydrate_factor FROM foods WHERE id = $1`, ingredient.FoodID).Scan(
			&recipe.Ingredients[i].Name,
			&recipe.Ingredients[i].NitrogenFactor,
			&recipe.Ingredients[i].ProteinFactor,
			&recipe.Ingredients[i].FatFactor,
			&recipe.Ingredients[i].CarbonhydrateFactor); err != nil {
			tx.Rollback()
			rw.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]string{"sql_error": fmt.Sprintf("%v", err)})
			return
		}
	}
	tx.Commit()

	enc.Encode(recipe)
}

func (rs *recipeServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	params, _ := url.ParseQuery(req.URL.RawQuery)

	if req.Method == "GET" {
		if _, ok := params["id"]; ok {
			rs.ShowRecipe(rw, req)
		} else {
			rs.ListRecipes(rw, req)
		}
	} else if req.Method == "POST" {
		rs.CreateRecipe(rw, req)
	} else if req.Method == "PUT" {
		rs.UpdateRecipe(rw, req)
	} else if req.Method == "DELETE" {
		rs.DeleteRecipe(rw, req)
	}
}

func NewRecipeServer(userDB, sr27DB *sql.DB) http.Handler {
	return &recipeServer{userDB, sr27DB}
}
