package n2xsdk

type MplsLabelManager struct {
 Handler string
}

func(np *MplsLabelManager) SetLabelRange () error {
 //parameters: hPort Minimum Maximum Increment
 //AgtMplsLabelManager SetLabelRange
 return nil
}

func(np *MplsLabelManager) GetLabelRange ()(string, error) {
 //parameters: hPort
 //AgtMplsLabelManager GetLabelRange
 return "", nil
}

func(np *MplsLabelManager) SetSutImplicitNullLabel () error {
 //parameters: hPort Label
 //AgtMplsLabelManager SetSutImplicitNullLabel
 return nil
}

func(np *MplsLabelManager) GetSutImplicitNullLabel ()(string, error) {
 //parameters: hPort
 //AgtMplsLabelManager GetSutImplicitNullLabel
 return "", nil
}

func(np *MplsLabelManager) EnableSutImplicitNullLabel () error {
 //parameters: hPort
 //AgtMplsLabelManager EnableSutImplicitNullLabel
 return nil
}

func(np *MplsLabelManager) DisableSutImplicitNullLabel () error {
 //parameters: hPort
 //AgtMplsLabelManager DisableSutImplicitNullLabel
 return nil
}

func(np *MplsLabelManager) IsSutImplicitNullLabelEnabled () error {
 //parameters: hPort
 //AgtMplsLabelManager IsSutImplicitNullLabelEnabled
 return nil
}

