package n2xsdk

type MsdpGroup struct {
 Handler string
}

func(np *MsdpGroup) SourceActive () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtMsdpGroup SourceActive
 return nil
}

func(np *MsdpGroup) SourceActiveOnAllPoolsInSession () error {
 //parameters: SessionHandle
 //AgtMsdpGroup SourceActiveOnAllPoolsInSession
 return nil
}

func(np *MsdpGroup) SourceActiveOnAllPools () error {
 //parameters: 
 //AgtMsdpGroup SourceActiveOnAllPools
 return nil
}

func(np *MsdpGroup) SourceActiveRequest () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtMsdpGroup SourceActiveRequest
 return nil
}

func(np *MsdpGroup) SourceActiveRequestOnAllPoolsInSession () error {
 //parameters: SessionHandle
 //AgtMsdpGroup SourceActiveRequestOnAllPoolsInSession
 return nil
}

func(np *MsdpGroup) SourceActiveRequestOnAllPools () error {
 //parameters: 
 //AgtMsdpGroup SourceActiveRequestOnAllPools
 return nil
}

func(np *MsdpGroup) Cancel () error {
 //parameters: SessionHandle MembershipPoolHandle Action
 //AgtMsdpGroup Cancel
 return nil
}

func(np *MsdpGroup) CancelAllInSession () error {
 //parameters: SessionHandle Action
 //AgtMsdpGroup CancelAllInSession
 return nil
}

func(np *MsdpGroup) CancelAll () error {
 //parameters: Action
 //AgtMsdpGroup CancelAll
 return nil
}

func(np *MsdpGroup) IsPoolAvailableForAction () error {
 //parameters: SessionHandle MembershipPoolHandle Action
 //AgtMsdpGroup IsPoolAvailableForAction
 return nil
}

func(np *MsdpGroup) IsPoolAvailableForCancel () error {
 //parameters: SessionHandle MembershipPoolHandle Action
 //AgtMsdpGroup IsPoolAvailableForCancel
 return nil
}

func(np *MsdpGroup) EnableSourceActiveResponse () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtMsdpGroup EnableSourceActiveResponse
 return nil
}

func(np *MsdpGroup) EnableSourceActiveResponseOnAllPoolsInSession () error {
 //parameters: SessionHandle
 //AgtMsdpGroup EnableSourceActiveResponseOnAllPoolsInSession
 return nil
}

func(np *MsdpGroup) EnableSourceActiveResponseOnAllPools () error {
 //parameters: 
 //AgtMsdpGroup EnableSourceActiveResponseOnAllPools
 return nil
}

func(np *MsdpGroup) DisableSourceActiveResponse () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtMsdpGroup DisableSourceActiveResponse
 return nil
}

func(np *MsdpGroup) DisableSourceActiveResponseOnAllPoolsInSession () error {
 //parameters: SessionHandle
 //AgtMsdpGroup DisableSourceActiveResponseOnAllPoolsInSession
 return nil
}

func(np *MsdpGroup) DisableSourceActiveResponseOnAllPools () error {
 //parameters: 
 //AgtMsdpGroup DisableSourceActiveResponseOnAllPools
 return nil
}

func(np *MsdpGroup) IsSourceActiveResponseEnabled () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtMsdpGroup IsSourceActiveResponseEnabled
 return nil
}

func(np *MsdpGroup) SetMcastDataEncapsulation () error {
 //parameters: SessionHandle MembershipPoolHandle EncapsulationFlag
 //AgtMsdpGroup SetMcastDataEncapsulation
 return nil
}

func(np *MsdpGroup) GetMcastDataEncapsulation ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtMsdpGroup GetMcastDataEncapsulation
 return "", nil
}

func(np *MsdpGroup) SetMcastDataEncapsulationInAllPoolsInSession () error {
 //parameters: SessionHandle EncapsulationFlag
 //AgtMsdpGroup SetMcastDataEncapsulationInAllPoolsInSession
 return nil
}

func(np *MsdpGroup) SetMcastDataEncapsulationInAllPools () error {
 //parameters: EncapsulationFlag
 //AgtMsdpGroup SetMcastDataEncapsulationInAllPools
 return nil
}

func(np *MsdpGroup) SetTlvTransmissionRate () error {
 //parameters: SessionHandle MembershipPoolHandle Action NumTlvsPerInterval Interval
 //AgtMsdpGroup SetTlvTransmissionRate
 return nil
}

func(np *MsdpGroup) GetTlvTransmissionRate ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle Action
 //AgtMsdpGroup GetTlvTransmissionRate
 return "", nil
}

func(np *MsdpGroup) SetPoolActionAtSessionEnable () error {
 //parameters: SessionHandle MembershipPoolHandle Action
 //AgtMsdpGroup SetPoolActionAtSessionEnable
 return nil
}

func(np *MsdpGroup) GetPoolActionAtSessionEnable ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtMsdpGroup GetPoolActionAtSessionEnable
 return "", nil
}

func(np *MsdpGroup) GetPoolStatistics ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtMsdpGroup GetPoolStatistics
 return "", nil
}

