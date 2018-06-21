package n2xsdk

type TestTopologyParameters struct {
 Handler string
}

func(np *TestTopologyParameters) GetLinkLayers ()(string, error) {
 //parameters: SessionHandle
 //AgtTestTopologyParameters GetLinkLayers
 return "", nil
}

func(np *TestTopologyParameters) SetLinkLayers () error {
 //parameters: SessionHandle EAgtTopologyParameterCategories
 //AgtTestTopologyParameters SetLinkLayers
 return nil
}

func(np *TestTopologyParameters) SelectParameter () error {
 //parameters: SessionHandle Category Parameter
 //AgtTestTopologyParameters SelectParameter
 return nil
}

func(np *TestTopologyParameters) DeselectParameter () error {
 //parameters: SessionHandle Category Parameter
 //AgtTestTopologyParameters DeselectParameter
 return nil
}

func(np *TestTopologyParameters) IsParameterSelected () error {
 //parameters: SessionHandle Category Parameter
 //AgtTestTopologyParameters IsParameterSelected
 return nil
}

func(np *TestTopologyParameters) SetParameterValue () error {
 //parameters: SessionHandle Category Parameter Value
 //AgtTestTopologyParameters SetParameterValue
 return nil
}

func(np *TestTopologyParameters) GetParameterValue ()(string, error) {
 //parameters: SessionHandle Category Parameter
 //AgtTestTopologyParameters GetParameterValue
 return "", nil
}

func(np *TestTopologyParameters) SetParameterValueFixed () error {
 //parameters: SessionHandle Category Parameter Value
 //AgtTestTopologyParameters SetParameterValueFixed
 return nil
}

func(np *TestTopologyParameters) SetIncrementingParameterValue () error {
 //parameters: SessionHandle Category Parameter StartValue Repeat Increment
 //AgtTestTopologyParameters SetIncrementingParameterValue
 return nil
}

func(np *TestTopologyParameters) GetIncrementingParameterValue ()(string, error) {
 //parameters: SessionHandle Category Parameter
 //AgtTestTopologyParameters GetIncrementingParameterValue
 return "", nil
}

func(np *TestTopologyParameters) SetIncrementingParameterValueAndCount () error {
 //parameters: SessionHandle Category Parameter StartValue Repeat Increment Count
 //AgtTestTopologyParameters SetIncrementingParameterValueAndCount
 return nil
}

func(np *TestTopologyParameters) GetIncrementingParameterValueAndCount ()(string, error) {
 //parameters: SessionHandle Category Parameter
 //AgtTestTopologyParameters GetIncrementingParameterValueAndCount
 return "", nil
}

func(np *TestTopologyParameters) SetIncrementingParameterValueAndOffset () error {
 //parameters: SessionHandle Category Parameter StartValue Repeat Increment Offset
 //AgtTestTopologyParameters SetIncrementingParameterValueAndOffset
 return nil
}

func(np *TestTopologyParameters) GetIncrementingParameterValueAndOffset ()(string, error) {
 //parameters: SessionHandle Category Parameter
 //AgtTestTopologyParameters GetIncrementingParameterValueAndOffset
 return "", nil
}

