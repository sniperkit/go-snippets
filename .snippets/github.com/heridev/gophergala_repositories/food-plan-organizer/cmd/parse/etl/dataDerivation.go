package etl

// Data Derivation Code Description File (file name = DERIV_CD).
//
// This file (Table 11) provides information on how the nutrient values were determined. The file contains the derivation codes and their descriptions.
//
// - Links to the Nutrient Data file by Deriv_Cd
//
// Deriv_Cd    A 4*    Derivation Code
// Deriv_Desc  A 120   Description of derivation code giving specific information on how the value was determined.

type dataDerivation struct {
	Code        string
	Description string
}
