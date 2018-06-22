package etl

import (
	"fmt"
	"io"
)

// Food Group Description File (file name = FD_GROUP).
//
// This file (Table 5) is a support file to the Food Description file and contains a list of food groups used in SR27 and their descriptions.
//
// - Links to the Food Description file by FdGrp_Cd
//
// FdGrp_Cd    A 4*  4-digit code identifying a food group
// FdGrp_Desc  A 60  Name of food group

type FoodGroup struct {
	GroupID int32
	Name    string
}

func parseFoodGroup(r *SR27Reader) (FoodGroup, error) {
	var s, err = r.Read()
	if err != nil {
		return FoodGroup{}, err
	}

	var fg = FoodGroup{
		GroupID: parseInt32(s[0]),
		Name:    s[1],
	}
	return fg, nil
}

type FoodGroupExtractor struct{}

func (fge FoodGroupExtractor) Parse(r io.Reader, parsed chan<- interface{}) error {
	var sr27 = newSR27Reader(r)
	sr27.FieldsPerRecord = 2

	var rows = 0
	for {
		var group, err = parseFoodGroup(sr27)

		if group.GroupID != 0 {
			parsed <- group
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		rows = rows + 1
	}
	fmt.Printf("parsed %d food group rows\n", rows)
	return nil
}
