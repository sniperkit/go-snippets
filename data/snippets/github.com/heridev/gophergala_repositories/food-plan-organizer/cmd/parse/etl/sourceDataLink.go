package etl

// Sources of Data Link File (file name = DATSRCLN)
//
// This file (Table 14) is used to link the Nutrient Data file with the Sources of Data table. It is needed to resolve the many-to- many relationship between the two tables.
//
// - Links to the Nutrient Data file by NDB No. and Nutr_No
// - Links to the Nutrient Definition file by Nutr_No
// - Links to the Sources of Data file by DataSrc_ID
//
// NDB_No      A 5*  5-digit Nutrient Databank number.
// Nutr_No     A 3*  Unique 3-digit identifier code for a nutrient.
// DataSrc_ID  A 6*  Unique ID identifying the reference/source.

type sourceDataLink struct {
}
