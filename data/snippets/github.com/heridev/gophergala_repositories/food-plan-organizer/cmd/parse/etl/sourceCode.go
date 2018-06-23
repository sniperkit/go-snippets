package etl

// Source Code File (file name = SRC_CD).
//
// This file (Table 10) contains codes indicating the type of data (analytical, calculated, assumed zero, and so on) in the Nutrient Data file.
// To improve the usability of the database and to provide values for the FNDDS, NDL staff imputed nutrient values for a number of
// proximate components, total dietary fiber, total sugar, and vitamin and mineral values.
//
// - Links to the Nutrient Data file by Src_Cd
//
// Src_Cd        A 2*    2-digit code.
// SrcCd_Desc    A 60    Description of source code that identifies the type of nutrient data.

type sourceCode struct {
	Code        string
	Description string
}
