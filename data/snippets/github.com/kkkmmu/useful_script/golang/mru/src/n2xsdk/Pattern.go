package n2xsdk

type Pattern struct {
 Handler string
}

func(np *Pattern) GetType ()(string, error) {
 //parameters: Handle
 //AgtPattern GetType
 return "", nil
}

func(np *Pattern) GetName ()(string, error) {
 //parameters: Handle
 //AgtPattern GetName
 return "", nil
}

func(np *Pattern) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtPattern GetLockCount
 return "", nil
}

func(np *Pattern) ClearPattern () error {
 //parameters: PatternHandle
 //AgtPattern ClearPattern
 return nil
}

func(np *Pattern) SetPatternOctets () error {
 //parameters: PatternHandle OctetOffset ValueLength psaValueOctets MaskLength psaMaskOctets
 //AgtPattern SetPatternOctets
 return nil
}

func(np *Pattern) GetPatternOctets ()(string, error) {
 //parameters: PatternHandle
 //AgtPattern GetPatternOctets
 return "", nil
}

func(np *Pattern) GetMaximumPatternLength ()(string, error) {
 //parameters: 
 //AgtPattern GetMaximumPatternLength
 return "", nil
}

func(np *Pattern) GetPatternLengthByPatternType ()(string, error) {
 //parameters: Type
 //AgtPattern GetPatternLengthByPatternType
 return "", nil
}

func(np *Pattern) ListProtocols ()(string, error) {
 //parameters: 
 //AgtPattern ListProtocols
 return "", nil
}

func(np *Pattern) ListProtocolsByPatternType ()(string, error) {
 //parameters: Type
 //AgtPattern ListProtocolsByPatternType
 return "", nil
}

func(np *Pattern) ListFields ()(string, error) {
 //parameters: Protocol
 //AgtPattern ListFields
 return "", nil
}

func(np *Pattern) GetFieldDetails ()(string, error) {
 //parameters: Protocol Field
 //AgtPattern GetFieldDetails
 return "", nil
}

func(np *Pattern) GetMaximumFieldValue ()(string, error) {
 //parameters: Protocol Field
 //AgtPattern GetMaximumFieldValue
 return "", nil
}

func(np *Pattern) GetProtocolOffset ()(string, error) {
 //parameters: PatternHandle Protocol
 //AgtPattern GetProtocolOffset
 return "", nil
}

func(np *Pattern) SetProtocolOffset () error {
 //parameters: PatternHandle Protocol NewOffset
 //AgtPattern SetProtocolOffset
 return nil
}

func(np *Pattern) ListMatchingFields ()(string, error) {
 //parameters: PatternHandle Protocol
 //AgtPattern ListMatchingFields
 return "", nil
}

func(np *Pattern) SetPatternField () error {
 //parameters: PatternHandle Protocol Field FieldBits ValueOctets psaValueOctets MaskOctets psaMaskOctets
 //AgtPattern SetPatternField
 return nil
}

func(np *Pattern) GetPatternField ()(string, error) {
 //parameters: PatternHandle Protocol Field
 //AgtPattern GetPatternField
 return "", nil
}

func(np *Pattern) SetPatternField32 () error {
 //parameters: PatternHandle Protocol Field FieldValue FieldMask
 //AgtPattern SetPatternField32
 return nil
}

func(np *Pattern) GetPatternField32 ()(string, error) {
 //parameters: PatternHandle Protocol Field
 //AgtPattern GetPatternField32
 return "", nil
}

