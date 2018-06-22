package models

type Nutrient struct {
	FoodID           int32   `json:"-"`
	NutrientID       int32   `json:"nutrient_id"`
	NutritionValue   float32 `json:"nutrient_value"`
	Min              float32 `json:"min"`
	Max              float32 `json:"max"`
	DegreesOfFreedom int32   `json:"-"`
	LowerErrorBound  float32 `json:"-"`
	UpperErrorBound  float32 `json:"-"`
	NutrientDefinition
}

// ScanExtendedNutrient scans nutrient, then nutrient_definition
func ScanNutrient(s StructScanner) (Nutrient, error) {
	var n Nutrient
	var err = s.Scan(&n.FoodID, &n.NutrientID, &n.NutritionValue, &n.Min, &n.Max, &n.DegreesOfFreedom, &n.LowerErrorBound, &n.UpperErrorBound, &n.NutrientDefinition.NutrientID, &n.Units, &n.Tagname, &n.Description, &n.DecimalPlaces)
	if err != nil {
		return n, err
	}
	return n, nil
}
