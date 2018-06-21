package n2xsdk

type LacpSession struct {
 Handler string
}

func(np *LacpSession) SelectPort () error {
 //parameters: SessionHandle
 //AgtLacpSession SelectPort
 return nil
}

func(np *LacpSession) UnselectPort () error {
 //parameters: SessionHandle
 //AgtLacpSession UnselectPort
 return nil
}

func(np *LacpSession) StandbyPort () error {
 //parameters: SessionHandle
 //AgtLacpSession StandbyPort
 return nil
}

func(np *LacpSession) ReadyPort () error {
 //parameters: SessionHandle
 //AgtLacpSession ReadyPort
 return nil
}

func(np *LacpSession) UnreadyPort () error {
 //parameters: SessionHandle
 //AgtLacpSession UnreadyPort
 return nil
}

func(np *LacpSession) ListUpdateSelectedActions ()(string, error) {
 //parameters: SessionHandle
 //AgtLacpSession ListUpdateSelectedActions
 return "", nil
}

func(np *LacpSession) SetUpdateSelectedActions () error {
 //parameters: SessionHandle UpdateSelectedActions
 //AgtLacpSession SetUpdateSelectedActions
 return nil
}

func(np *LacpSession) SetDefaultSelectedValue () error {
 //parameters: SessionHandle Selected
 //AgtLacpSession SetDefaultSelectedValue
 return nil
}

func(np *LacpSession) GetDefaultSelectedValue ()(string, error) {
 //parameters: SessionHandle
 //AgtLacpSession GetDefaultSelectedValue
 return "", nil
}

func(np *LacpSession) GetSelectedValue ()(string, error) {
 //parameters: SessionHandle
 //AgtLacpSession GetSelectedValue
 return "", nil
}

func(np *LacpSession) EnableDefaultReadyFlag () error {
 //parameters: SessionHandle
 //AgtLacpSession EnableDefaultReadyFlag
 return nil
}

func(np *LacpSession) DisableDefaultReadyFlag () error {
 //parameters: SessionHandle
 //AgtLacpSession DisableDefaultReadyFlag
 return nil
}

func(np *LacpSession) IsDefaultReadyFlagEnabled () error {
 //parameters: SessionHandle
 //AgtLacpSession IsDefaultReadyFlagEnabled
 return nil
}

func(np *LacpSession) IsReadyFlagEnabled () error {
 //parameters: SessionHandle
 //AgtLacpSession IsReadyFlagEnabled
 return nil
}

func(np *LacpSession) GetState ()(string, error) {
 //parameters: SessionHandle
 //AgtLacpSession GetState
 return "", nil
}

func(np *LacpSession) GetReceiveState ()(string, error) {
 //parameters: SessionHandle
 //AgtLacpSession GetReceiveState
 return "", nil
}

func(np *LacpSession) SetTimerValue () error {
 //parameters: SessionHandle Timer Seconds
 //AgtLacpSession SetTimerValue
 return nil
}

func(np *LacpSession) GetTimerValue ()(string, error) {
 //parameters: SessionHandle Timer
 //AgtLacpSession GetTimerValue
 return "", nil
}

func(np *LacpSession) GetCollectorMaxDelay ()(string, error) {
 //parameters: SessionHandle
 //AgtLacpSession GetCollectorMaxDelay
 return "", nil
}

func(np *LacpSession) SetCollectorMaxDelay () error {
 //parameters: SessionHandle TensOfMicroSeconds
 //AgtLacpSession SetCollectorMaxDelay
 return nil
}

func(np *LacpSession) SetActorPortAggregatorId () error {
 //parameters: SessionHandle Id
 //AgtLacpSession SetActorPortAggregatorId
 return nil
}

func(np *LacpSession) GetActorPortAggregatorId ()(string, error) {
 //parameters: SessionHandle
 //AgtLacpSession GetActorPortAggregatorId
 return "", nil
}

func(np *LacpSession) GetActorAdminParameter ()(string, error) {
 //parameters: SessionHandle ParameterType
 //AgtLacpSession GetActorAdminParameter
 return "", nil
}

func(np *LacpSession) SetActorAdminParameter () error {
 //parameters: SessionHandle ParameterType ParameterValue
 //AgtLacpSession SetActorAdminParameter
 return nil
}

func(np *LacpSession) GetPartnerParameter ()(string, error) {
 //parameters: SessionHandle ParameterType
 //AgtLacpSession GetPartnerParameter
 return "", nil
}

func(np *LacpSession) EnablePartnerAdminParametersLearning () error {
 //parameters: SessionHandle
 //AgtLacpSession EnablePartnerAdminParametersLearning
 return nil
}

func(np *LacpSession) DisablePartnerAdminParametersLearning () error {
 //parameters: SessionHandle
 //AgtLacpSession DisablePartnerAdminParametersLearning
 return nil
}

func(np *LacpSession) IsPartnerAdminParametersLearningEnabled () error {
 //parameters: SessionHandle
 //AgtLacpSession IsPartnerAdminParametersLearningEnabled
 return nil
}

func(np *LacpSession) EnablePartnerAdminParametersRelearning () error {
 //parameters: SessionHandle
 //AgtLacpSession EnablePartnerAdminParametersRelearning
 return nil
}

func(np *LacpSession) DisablePartnerAdminParametersRelearning () error {
 //parameters: SessionHandle
 //AgtLacpSession DisablePartnerAdminParametersRelearning
 return nil
}

func(np *LacpSession) IsPartnerAdminParametersRelearningEnabled () error {
 //parameters: SessionHandle
 //AgtLacpSession IsPartnerAdminParametersRelearningEnabled
 return nil
}

func(np *LacpSession) GetPartnerAdminParameter ()(string, error) {
 //parameters: SessionHandle ParameterType
 //AgtLacpSession GetPartnerAdminParameter
 return "", nil
}

func(np *LacpSession) SetPartnerAdminParameter () error {
 //parameters: SessionHandle ParameterType ParameterValue
 //AgtLacpSession SetPartnerAdminParameter
 return nil
}

func(np *LacpSession) IsActorPortStateFlagEnabled () error {
 //parameters: SessionHandle StateFlag
 //AgtLacpSession IsActorPortStateFlagEnabled
 return nil
}

func(np *LacpSession) EnableActorAdminPortStateFlag () error {
 //parameters: SessionHandle StateFlag
 //AgtLacpSession EnableActorAdminPortStateFlag
 return nil
}

func(np *LacpSession) DisableActorAdminPortStateFlag () error {
 //parameters: SessionHandle StateFlag
 //AgtLacpSession DisableActorAdminPortStateFlag
 return nil
}

func(np *LacpSession) IsActorAdminPortStateFlagEnabled () error {
 //parameters: SessionHandle StateFlag
 //AgtLacpSession IsActorAdminPortStateFlagEnabled
 return nil
}

func(np *LacpSession) IsPartnerPortStateFlagEnabled () error {
 //parameters: SessionHandle StateFlag
 //AgtLacpSession IsPartnerPortStateFlagEnabled
 return nil
}

func(np *LacpSession) EnablePartnerAdminPortStateFlag () error {
 //parameters: SessionHandle StateFlag
 //AgtLacpSession EnablePartnerAdminPortStateFlag
 return nil
}

func(np *LacpSession) DisablePartnerAdminPortStateFlag () error {
 //parameters: SessionHandle StateFlag
 //AgtLacpSession DisablePartnerAdminPortStateFlag
 return nil
}

func(np *LacpSession) IsPartnerAdminPortStateFlagEnabled () error {
 //parameters: SessionHandle StateFlag
 //AgtLacpSession IsPartnerAdminPortStateFlagEnabled
 return nil
}

func(np *LacpSession) GetActorPortStateValue ()(string, error) {
 //parameters: SessionHandle
 //AgtLacpSession GetActorPortStateValue
 return "", nil
}

func(np *LacpSession) GetPartnerPortStateValue ()(string, error) {
 //parameters: SessionHandle
 //AgtLacpSession GetPartnerPortStateValue
 return "", nil
}

func(np *LacpSession) GetActorAdminPortStateValue ()(string, error) {
 //parameters: SessionHandle
 //AgtLacpSession GetActorAdminPortStateValue
 return "", nil
}

func(np *LacpSession) GetPartnerAdminPortStateValue ()(string, error) {
 //parameters: SessionHandle
 //AgtLacpSession GetPartnerAdminPortStateValue
 return "", nil
}

