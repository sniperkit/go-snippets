package etl

import (
	"fmt"
	"io"
)

// LanguaL Factors Description File (File name = LANGDESC).
//
// This file (Table 7) is a support file to the LanguaL Factor file and contains the descriptions for only those factors used in coding the selected food items codes in this release of SR.
//
// - Links to the LanguaL Factor File by the Factor_Code field
//
// Factor_Code  A 5*  The Langual factor from the Thesaurus.
// Description  A 140  The description of the LanguaL Factor Code from the thesaurus
type LanguageDescription struct {
	FactorCode  string
	Description string
}

func parseLanguageDescription(r *SR27Reader) (LanguageDescription, error) {
	var s, err = r.Read()
	if err != nil {
		return LanguageDescription{}, err
	}

	var ld = LanguageDescription{
		FactorCode:  s[0],
		Description: s[1],
	}
	return ld, nil
}

type LangDescriptionExtractor struct{}

func (lde LangDescriptionExtractor) Parse(r io.Reader, parsed chan<- interface{}) error {
	var sr27 = newSR27Reader(r)
	sr27.FieldsPerRecord = 2

	var rows = 0
	for {
		var ld, err = parseLanguageDescription(sr27)
		if ld.FactorCode != "" {
			parsed <- ld
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		rows = rows + 1
	}
	fmt.Printf("parsed %d language definition rows\n", rows)
	return nil
}
