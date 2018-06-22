package etl

import (
	"fmt"
	"io"
)

// LanguaL Factor File (File name = LANGUAL).
//
// This file (Table 6) is a support file to the Food Description file and contains the factors from the LanguaL Thesaurus used to code a particular food.
//
// - Links to the Food Description file by the NDB_No field
// - Links to LanguaL Factors Description file by the Factor_Code field
//
// NDB_No       A 5*  5-digit Nutrient Databank number that uniquely identifies a food item.
// Factor_Code  A 5*  The Langual factor from the Thesaurus

type Lang struct {
	NutrientID string
	FactorCode string
}

func parseLang(r *SR27Reader) (Lang, error) {
	var s, err = r.Read()
	if err != nil {
		return Lang{}, err
	}

	var ld = Lang{
		NutrientID: s[0],
		FactorCode: s[1],
	}
	return ld, nil
}

type LangExtractor struct{}

func (le LangExtractor) Parse(r io.Reader, parsed chan<- interface{}) error {
	var sr27 = newSR27Reader(r)
	sr27.FieldsPerRecord = 2

	var rows = 0
	for {
		var lang, err = parseLang(sr27)
		if err == io.EOF {
			break
		}
		if lang.NutrientID != "" {
			parsed <- lang
		}
		if err != nil {
			return err
		}

		rows = rows + 1
	}
	fmt.Printf("parsed %d language definition mapping rows\n", rows)
	return nil
}
