package n2xsdk

type IsisNetwork struct {
 Handler string
}

func(np *IsisNetwork) Clear () error {
 //parameters: NetworkHandle
 //AgtIsisNetwork Clear
 return nil
}

func(np *IsisNetwork) AddReachablePrefix () error {
 //parameters: NetworkHandle NetworkPrefix PrefixLength
 //AgtIsisNetwork AddReachablePrefix
 return nil
}

func(np *IsisNetwork) GetReachablePrefixCount ()(string, error) {
 //parameters: NetworkHandle
 //AgtIsisNetwork GetReachablePrefixCount
 return "", nil
}

func(np *IsisNetwork) GetReachablePrefix ()(string, error) {
 //parameters: NetworkHandle PrefixIndex
 //AgtIsisNetwork GetReachablePrefix
 return "", nil
}

func(np *IsisNetwork) SetReachablePrefix () error {
 //parameters: NetworkHandle PrefixIndex NetworkPrefix
 //AgtIsisNetwork SetReachablePrefix
 return nil
}

func(np *IsisNetwork) GetReachablePrefixLength ()(string, error) {
 //parameters: NetworkHandle PrefixIndex
 //AgtIsisNetwork GetReachablePrefixLength
 return "", nil
}

func(np *IsisNetwork) SetReachablePrefixLength () error {
 //parameters: NetworkHandle PrefixIndex PrefixLength
 //AgtIsisNetwork SetReachablePrefixLength
 return nil
}

func(np *IsisNetwork) GenerateReachablePrefixes () error {
 //parameters: NetworkHandle StartPrefix PrefixLength Increment Count
 //AgtIsisNetwork GenerateReachablePrefixes
 return nil
}

func(np *IsisNetwork) ListReachablePrefixes ()(string, error) {
 //parameters: NetworkHandle
 //AgtIsisNetwork ListReachablePrefixes
 return "", nil
}

func(np *IsisNetwork) ListReachablePrefixLengths ()(string, error) {
 //parameters: NetworkHandle
 //AgtIsisNetwork ListReachablePrefixLengths
 return "", nil
}

func(np *IsisNetwork) GetMaximumNumberOfPrefixes ()(string, error) {
 //parameters: 
 //AgtIsisNetwork GetMaximumNumberOfPrefixes
 return "", nil
}

