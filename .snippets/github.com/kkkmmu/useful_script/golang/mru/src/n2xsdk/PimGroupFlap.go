package n2xsdk

type PimGroupFlap struct {
 Handler string
}

func(np *PimGroupFlap) SetFlapTimers () error {
 //parameters: SessionHandle MembershipPoolHandle JoinToLeaveDelay LeaveToJoinDelay
 //AgtPimGroupFlap SetFlapTimers
 return nil
}

func(np *PimGroupFlap) GetFlapTimers ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimGroupFlap GetFlapTimers
 return "", nil
}

func(np *PimGroupFlap) StartFlap () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimGroupFlap StartFlap
 return nil
}

func(np *PimGroupFlap) StopFlap () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimGroupFlap StopFlap
 return nil
}

func(np *PimGroupFlap) StartAllFlapsInSession () error {
 //parameters: SessionHandle
 //AgtPimGroupFlap StartAllFlapsInSession
 return nil
}

func(np *PimGroupFlap) StopAllFlapsInSession () error {
 //parameters: SessionHandle
 //AgtPimGroupFlap StopAllFlapsInSession
 return nil
}

func(np *PimGroupFlap) StartAllFlaps () error {
 //parameters: 
 //AgtPimGroupFlap StartAllFlaps
 return nil
}

func(np *PimGroupFlap) StopAllFlaps () error {
 //parameters: 
 //AgtPimGroupFlap StopAllFlaps
 return nil
}

func(np *PimGroupFlap) Cancel () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimGroupFlap Cancel
 return nil
}

func(np *PimGroupFlap) CancelAllInSession () error {
 //parameters: SessionHandle
 //AgtPimGroupFlap CancelAllInSession
 return nil
}

func(np *PimGroupFlap) CancelAll () error {
 //parameters: 
 //AgtPimGroupFlap CancelAll
 return nil
}

func(np *PimGroupFlap) IsPoolAvailableForAction () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimGroupFlap IsPoolAvailableForAction
 return nil
}

func(np *PimGroupFlap) IsPoolAvailableForCancel () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimGroupFlap IsPoolAvailableForCancel
 return nil
}

func(np *PimGroupFlap) EnableFlapFlag () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimGroupFlap EnableFlapFlag
 return nil
}

func(np *PimGroupFlap) DisableFlapFlag () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimGroupFlap DisableFlapFlag
 return nil
}

func(np *PimGroupFlap) GetFlapEnableFlag ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimGroupFlap GetFlapEnableFlag
 return "", nil
}

