package n2xsdk

type RipFlapProfile struct {
 Handler string
}

func(np *RipFlapProfile) GetAdvertiseToWithdrawDelay ()(string, error) {
 //parameters: SessionHandle
 //AgtRipFlapProfile GetAdvertiseToWithdrawDelay
 return "", nil
}

func(np *RipFlapProfile) SetAdvertiseToWithdrawDelay () error {
 //parameters: SessionHandle WithdrawDelay
 //AgtRipFlapProfile SetAdvertiseToWithdrawDelay
 return nil
}

func(np *RipFlapProfile) GetWithdrawToAdvertiseDelay ()(string, error) {
 //parameters: SessionHandle
 //AgtRipFlapProfile GetWithdrawToAdvertiseDelay
 return "", nil
}

func(np *RipFlapProfile) SetWithdrawToAdvertiseDelay () error {
 //parameters: SessionHandle AdvertiseDelay
 //AgtRipFlapProfile SetWithdrawToAdvertiseDelay
 return nil
}

func(np *RipFlapProfile) StartFlap () error {
 //parameters: SessionHandle
 //AgtRipFlapProfile StartFlap
 return nil
}

func(np *RipFlapProfile) StopFlap () error {
 //parameters: SessionHandle
 //AgtRipFlapProfile StopFlap
 return nil
}

func(np *RipFlapProfile) GetFlapState ()(string, error) {
 //parameters: SessionHandle
 //AgtRipFlapProfile GetFlapState
 return "", nil
}

func(np *RipFlapProfile) GetRoutePools ()(string, error) {
 //parameters: SessionHandle
 //AgtRipFlapProfile GetRoutePools
 return "", nil
}

func(np *RipFlapProfile) SetRoutePools () error {
 //parameters: SessionHandle Count pRoutePoolHandles
 //AgtRipFlapProfile SetRoutePools
 return nil
}

