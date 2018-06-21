package n2xsdk

type PimSmSource struct {
 Handler string
}

func(np *PimSmSource) GetState ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSmSource GetState
 return "", nil
}

func(np *PimSmSource) StartRegisterSource () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSmSource StartRegisterSource
 return nil
}

func(np *PimSmSource) StopRegisterSource () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSmSource StopRegisterSource
 return nil
}

func(np *PimSmSource) StartAllRegisterSourcesInSession () error {
 //parameters: SessionHandle
 //AgtPimSmSource StartAllRegisterSourcesInSession
 return nil
}

func(np *PimSmSource) StopAllRegisterSourcesInSession () error {
 //parameters: SessionHandle
 //AgtPimSmSource StopAllRegisterSourcesInSession
 return nil
}

func(np *PimSmSource) StartAllRegisterSources () error {
 //parameters: 
 //AgtPimSmSource StartAllRegisterSources
 return nil
}

func(np *PimSmSource) StopAllRegisterSources () error {
 //parameters: 
 //AgtPimSmSource StopAllRegisterSources
 return nil
}

func(np *PimSmSource) Cancel () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSmSource Cancel
 return nil
}

func(np *PimSmSource) CancelAllInSession () error {
 //parameters: SessionHandle
 //AgtPimSmSource CancelAllInSession
 return nil
}

func(np *PimSmSource) CancelAll () error {
 //parameters: 
 //AgtPimSmSource CancelAll
 return nil
}

func(np *PimSmSource) IsPoolAvailableForAction () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSmSource IsPoolAvailableForAction
 return nil
}

func(np *PimSmSource) IsPoolAvailableForCancel () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSmSource IsPoolAvailableForCancel
 return nil
}

func(np *PimSmSource) SetRegisterTransmissionRate () error {
 //parameters: SessionHandle MembershipPoolHandle NumPdusPerInterval Interval
 //AgtPimSmSource SetRegisterTransmissionRate
 return nil
}

func(np *PimSmSource) GetRegisterTransmissionRate ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSmSource GetRegisterTransmissionRate
 return "", nil
}

func(np *PimSmSource) EnableSourceFlag () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSmSource EnableSourceFlag
 return nil
}

func(np *PimSmSource) DisableSourceFlag () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSmSource DisableSourceFlag
 return nil
}

func(np *PimSmSource) GetSourceEnableFlag ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSmSource GetSourceEnableFlag
 return "", nil
}

func(np *PimSmSource) SetBorderBit () error {
 //parameters: SessionHandle MembershipPoolHandle BorderBit
 //AgtPimSmSource SetBorderBit
 return nil
}

func(np *PimSmSource) GetBorderBit ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSmSource GetBorderBit
 return "", nil
}

