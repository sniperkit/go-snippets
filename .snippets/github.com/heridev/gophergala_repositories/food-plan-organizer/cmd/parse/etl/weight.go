package etl

import (
	"fmt"
	"io"
)

// Weight File (file name = WEIGHT).
//
// This file (Table 12) contains the weight in grams of a number of common measures for each food item.
//
// - Links to Food Description file by NDB_No
// - Links to Nutrient Data file by NDB_No
//
// NDB_No         A 5*     5-digit Nutrient Databank number.
// Seq            A 2*     Sequence number.
// Amount         N 5.3    Unit modifier (for example, 1 in “1 cup”).
// Msre_Desc      A 84     Description (for example, cup, diced, and 1-inch pieces).
// Gm_Wgt         N 7.1    Gram weight.
// Num_Data_Pts   N 3      Number of data points.
// Std_Dev        N 7.3    Standard deviation.

type Weight struct {
	NutrientID  int32
	Seq         string
	Amount      float32
	Description string
	GramWeight  float32
}

func parseWeight(r *SR27Reader) (Weight, error) {
	var s, err = r.Read()
	if err != nil {
		return Weight{}, err
	}

	var w = Weight{
		NutrientID:  parseInt32(s[0]),
		Seq:         s[1],
		Amount:      parseFloat32(s[2]),
		Description: s[3],
		GramWeight:  parseFloat32(s[4]),
	}
	return w, nil
}

type WeightExtractor struct{}

func (we WeightExtractor) Parse(r io.Reader, parsed chan<- interface{}) error {
	var sr27 = newSR27Reader(r)
	sr27.FieldsPerRecord = 7

	var rows = 0
	for {
		var w, err = parseWeight(sr27)
		if err == io.EOF {
			break
		}
		if w.NutrientID != 0 {
			parsed <- w
		}
		if err != nil {
			return err
		}

		rows = rows + 1
	}
	fmt.Printf("parsed %d weight rows\n", rows)
	return nil
}
