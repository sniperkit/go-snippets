package etl

import (
	"fmt"
	"io"
)

// Nutrient Definition File (file name = NUTR_DEF).
//
// This file (Table 9) is a support file to the Nutrient Data file. It provides the 3-digit nutrient code, unit of measure, INFOODS tagname, and description.
//
// - Links to the Nutrient Data file by Nutr_No
//
// Nutr_No   A 3*  Unique 3-digit identifier code for a nutrient.
// Units     A 7   Units of measure (mg, g, μg, and so on).
// Tagname   A 20  International Network of Food Data Systems (INFOODS) Tagnames.
// NutrDesc  A 60  Name of nutrient/food component.
// Num_Dec   A 1   Number of decimal places to which a nutrient value is rounded.
// SR_Order  N 6   Used to sort nutrient records in the same order as various reports produced from SR.

type NutrientDefinition struct {
	NutrientID    int32
	Units         string
	Tagname       string
	Description   string
	DecimalPlaces int32
}

func parseNutrientDefinition(r *SR27Reader) (NutrientDefinition, error) {
	var s, err = r.Read()
	if err != nil {
		return NutrientDefinition{}, err
	}

	var nd = NutrientDefinition{
		NutrientID:    parseInt32(s[0]),
		Units:         s[1],
		Tagname:       s[2],
		Description:   s[3],
		DecimalPlaces: parseInt32(s[4]),
	}
	return nd, nil
}

type NutrientDefinitionExtractor struct{}

func (nde NutrientDefinitionExtractor) Parse(r io.Reader, parsed chan<- interface{}) error {
	var sr27 = newSR27Reader(r)
	sr27.FieldsPerRecord = 6

	var rows = 0
	for {
		var nd, err = parseNutrientDefinition(sr27)
		if err == io.EOF {
			break
		}
		if nd.NutrientID != 0 {
			parsed <- nd
		}
		if err != nil {
			return err
		}

		rows = rows + 1
	}
	fmt.Printf("parsed %d nutrient definition rows\n", rows)
	return nil
}
