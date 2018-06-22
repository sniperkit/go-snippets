package models

type NutrientDefinition struct {
	NutrientID    int32  `json:"id"`
	Units         string `json:"unit"`
	Tagname       string `json:"tagname"`
	Description   string `json:"description"`
	DecimalPlaces int32  `json:"decimal_places"`
}
