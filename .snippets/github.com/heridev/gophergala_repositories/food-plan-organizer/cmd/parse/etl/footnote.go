package etl

// Footnote File (file name = FOOTNOTE).
//
// This file (Table 13) contains additional information about the food item, household weight, and nutrient value.
//
// - Links to the Food Description file by NDB_No
// - Links to the Nutrient Data file by NDB_No and when applicable, Nutr_No
// - Links to the Nutrient Definition file by Nutr_No, when applicable
//
// NDB_No      A 5    5-digit Nutrient Databank number.
// Footnt_No   A 4    Sequence number. If a given footnote applies to more than one nutrient number, the same footnote number is used. As a result, this file cannot be indexed.
// Footnt_Typ  A 1    Type of footnote:
//                      D = footnote adding information to the food description;
//                      M = footnote adding information to measure description;
//                      N = footnote providing additional information on a nutrient value. If the Footnt_typ = N, the Nutr_No will also be filled in.
// Nutr_No     A 3    Unique 3-digit identifier code for a nutrient to which footnote applies.
// Footnt_Txt  A 200  Footnote text.
type footnote struct {
}
