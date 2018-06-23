package etl

import (
	"fmt"
	"io"
)

// Food Description File (file name = FOOD_DES).
//
// This file (Table 4) contains long and short descriptions and food group designators for all food items,
// along with common names, manufacturer name, scientific name, percentage and description of refuse, and factors used for
// calculating protein and kilocalories, if applicable.
// Items used in the FNDDS are also identified by value of “Y” in the Survey field.
//
// - Links to the Food Group Description file by the FdGrp_Cd field
// - Links to the Nutrient Data file by the NDB_No field
// - Links to the Weight file by the NDB_No field
// - Links to the Footnote file by the NDB_No field
// - Links to the LanguaL Factor file by the NDB_No field
//
// NDB_No  	   A 5*   5-digit Nutrient Databank number that uniquely identifies a food item. If this field is defined as numeric, the leading zero will be lost
// FdGrp_Cd    A 4    4-digit code indicating food group to which a food item belongs.
// Long_Desc   A 200  200-character description of food item.
// Shrt_Desc   A 60   60-character abbreviated description of food item. Generated from the 200-character description
// ComName     A 100  Other names commonly used to describe a food,
// ManufacName A 65   Indicates the company that manufactured the product, when appropriate.
// Survey      A 1    Indicates if the food item is used in the USDA Food and Nutrient Database for Dietary Studies
// Ref_desc    A 135  Description of inedible parts of a food item (refuse), such as seeds or bone.
// Refuse      N 2    Percentage of refuse.
// SciName     A 65   Scientific name of the food item. Given for the least processed form of the food (usually raw), if applicable.
// N_Factor    N 4.2  Factor for converting nitrogen to protein
// Pro_Factor  N 4.2  Factor for calculating calories from protein
// Fat_Factor  N 4.2  Factor for calculating calories from fat
// CHO_Factor  N 4.2  Factor for calculating calories from carbohydrate

type Food struct {
	ID                 int32
	FoodGroupID        int32
	Name               string
	ShortName          string
	CommonName         string
	ScientificName     string
	NitrogenFactor     float32
	ProteinFactor      float32
	FatFactor          float32
	CarbohydrateFactor float32
}

func parseFood(r *SR27Reader) (Food, error) {
	var s, err = r.Read()
	if err != nil {
		return Food{}, err
	}

	var f = Food{
		ID:                 parseInt32(s[0]),
		FoodGroupID:        parseInt32(s[1]),
		Name:               s[2],
		ShortName:          s[3],
		CommonName:         s[4],
		ScientificName:     s[9],
		NitrogenFactor:     parseFloat32(s[10]),
		ProteinFactor:      parseFloat32(s[11]),
		FatFactor:          parseFloat32(s[12]),
		CarbohydrateFactor: parseFloat32(s[13]),
	}
	return f, nil
}

type FoodExtractor struct{}

// TODO(rr): store in database
func (fe FoodExtractor) Parse(r io.Reader, parsed chan<- interface{}) error {
	var sr27 = newSR27Reader(r)
	sr27.FieldsPerRecord = 14

	var rows = 0
	for {
		var food, err = parseFood(sr27)
		if food.ID != 0 {
			parsed <- food
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		rows = rows + 1
	}
	fmt.Printf("parsed %d food rows\n", rows)
	return nil
}
