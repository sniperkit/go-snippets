package etl

import (
	"fmt"
	"io"
)

// Nutrient Data File (file name = NUT_DATA).
//
// This file (Table 8) contains the nutrient values and information about the values, including expanded statistical information.
//
// - Links to the Food Description file by NDB_No
// - Links to the Food Description file by Ref_NDB_No
// - Links to the Weight file by NDB_No
// - Links to the Footnote file by NDB_No and when applicable, Nutr_No
// - Links to the Nutrient Definition file by Nutr_No
// - Links to the Source Code file by Src_Cd
// - Links to the Derivation Code file by Deriv_Cd
//
// NDB_No         A 5*    5-digit Nutrient Databank number.
// Nutr_No        A 3*    Unique 3-digit identifier code for a nutrient.
// Nutr_Val       N 10.3  Amount in 100 grams, edible portion †.
// Num_Data_Pts   N 5.0   Number of data points (previously called Sample_Ct) is the number of analyses used to calculate the nutrient value. If the number of data points is 0, the value was calculated or imputed.
// Std_Error      N 8.3   Standard error of the mean. Null if cannot be calculated. The standard error is also not given if the number of data points is less than three.
// Src_Cd         A 2     Code indicating type of data.
// Deriv_Cd       A 4     Data Derivation Code giving specific information on how the value is determined. This field is populated only for items added or updated starting with SR14.
// Ref_NDB_No     A 5     NDB number of the item used to calculate a missing value. Populated only for items added or updated starting with SR14.
// Add_Nutr_Mark  A 1     Indicates a vitamin or mineral added for fortification or enrichment. This field is populated for ready-to- eat breakfast cereals and many brand-name hot cereals in food group 8.
// Num_Studies    N 2     Number of studies.
// Min            N 10.3  Minimum value.
// Max            N 10.3  Maximum value.
// DF             N 4     Degrees of freedom.
// Low_EB         N 10.3  Lower 95% error bound.
// Up_EB          N 10.3  Upper 95% error bound.
// Stat_cmt       A 10    Statistical comments
// AddMod_Date    A 10    Indicates when a value was either added to the database or last modified.
// CC             A 1     Confidence Code indicating data quality, based on evaluation of sample plan, sample handling, analytical method, analytical quality control, and number of samples analyzed

type Nutrient struct {
	FoodID           int32
	NutrientID       int32
	NutritionValue   float32
	Min              float32
	Max              float32
	DegreesOfFreedom int32
	LowerErrorBound  float32
	UpperErrorBound  float32
}

func parseNutrient(r *SR27Reader) (Nutrient, error) {
	var s, err = r.Read()
	if err != nil {
		return Nutrient{}, err
	}

	var n = Nutrient{
		FoodID:           parseInt32(s[0]),
		NutrientID:       parseInt32(s[1]),
		NutritionValue:   parseFloat32(s[2]),
		Min:              parseFloat32(s[10]),
		Max:              parseFloat32(s[11]),
		DegreesOfFreedom: parseInt32(s[12]),
		LowerErrorBound:  parseFloat32(s[13]),
		UpperErrorBound:  parseFloat32(s[14]),
	}
	return n, nil
}

type NutrientExtractor struct{}

func (ne NutrientExtractor) Parse(r io.Reader, parsed chan<- interface{}) error {
	var sr27 = newSR27Reader(r)
	sr27.FieldsPerRecord = 18

	var rows = 0
	for {
		var n, err = parseNutrient(sr27)
		if err == io.EOF {
			break
		}
		if n.NutrientID != 0 {
			parsed <- n
		}
		if err != nil {
			return err
		}

		rows = rows + 1
	}
	fmt.Printf("parsed %d nutrient rows\n", rows)
	return nil
}
