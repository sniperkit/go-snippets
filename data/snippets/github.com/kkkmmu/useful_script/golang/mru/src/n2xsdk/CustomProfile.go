package n2xsdk

type CustomProfile struct {
 Handler string
}

func(np *CustomProfile) GetType ()(string, error) {
 //parameters: Handle
 //AgtCustomProfile GetType
 return "", nil
}

func(np *CustomProfile) GetName ()(string, error) {
 //parameters: Handle
 //AgtCustomProfile GetName
 return "", nil
}

func(np *CustomProfile) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtCustomProfile GetLockCount
 return "", nil
}

func(np *CustomProfile) GetSourcePort ()(string, error) {
 //parameters: ProfileHandle
 //AgtCustomProfile GetSourcePort
 return "", nil
}

func(np *CustomProfile) SetProfileType () error {
 //parameters: ProfileHandle ProfileType
 //AgtCustomProfile SetProfileType
 return nil
}

func(np *CustomProfile) GetProfileType ()(string, error) {
 //parameters: ProfileHandle
 //AgtCustomProfile GetProfileType
 return "", nil
}

func(np *CustomProfile) SetAverageLoad () error {
 //parameters: ProfileHandle AverageLoad Units
 //AgtCustomProfile SetAverageLoad
 return nil
}

func(np *CustomProfile) GetAverageLoad ()(string, error) {
 //parameters: ProfileHandle Units
 //AgtCustomProfile GetAverageLoad
 return "", nil
}

func(np *CustomProfile) SetMode () error {
 //parameters: ProfileHandle Mode
 //AgtCustomProfile SetMode
 return nil
}

func(np *CustomProfile) GetMode ()(string, error) {
 //parameters: ProfileHandle
 //AgtCustomProfile GetMode
 return "", nil
}

func(np *CustomProfile) SetNumberOfPacketsToInject () error {
 //parameters: ProfileHandle NumberOfPackets
 //AgtCustomProfile SetNumberOfPacketsToInject
 return nil
}

func(np *CustomProfile) GetNumberOfPacketsToInject ()(string, error) {
 //parameters: ProfileHandle
 //AgtCustomProfile GetNumberOfPacketsToInject
 return "", nil
}

func(np *CustomProfile) SendSingleInjection () error {
 //parameters: ProfileHandle
 //AgtCustomProfile SendSingleInjection
 return nil
}

func(np *CustomProfile) IsSendSingleInjectionInProgress () error {
 //parameters: ProfileHandle
 //AgtCustomProfile IsSendSingleInjectionInProgress
 return nil
}

func(np *CustomProfile) ListStreamGroups ()(string, error) {
 //parameters: ProfileHandle
 //AgtCustomProfile ListStreamGroups
 return "", nil
}

func(np *CustomProfile) AddStreamGroups () error {
 //parameters: ProfileHandle Count psaStreamGroupHandles
 //AgtCustomProfile AddStreamGroups
 return nil
}

func(np *CustomProfile) GetInterdepartureValues ()(string, error) {
 //parameters: ProfileHandle
 //AgtCustomProfile GetInterdepartureValues
 return "", nil
}

func(np *CustomProfile) GetInterdepartureRepeatCounts ()(string, error) {
 //parameters: ProfileHandle
 //AgtCustomProfile GetInterdepartureRepeatCounts
 return "", nil
}

func(np *CustomProfile) GetInterdepartureValueResolution ()(string, error) {
 //parameters: ProfileHandle
 //AgtCustomProfile GetInterdepartureValueResolution
 return "", nil
}

func(np *CustomProfile) Enable () error {
 //parameters: ProfileHandle
 //AgtCustomProfile Enable
 return nil
}

func(np *CustomProfile) Disable () error {
 //parameters: ProfileHandle
 //AgtCustomProfile Disable
 return nil
}

func(np *CustomProfile) IsEnabled () error {
 //parameters: ProfileHandle
 //AgtCustomProfile IsEnabled
 return nil
}

func(np *CustomProfile) SetInterdepartureValues () error {
 //parameters: ProfileHandle Count psaIDVs
 //AgtCustomProfile SetInterdepartureValues
 return nil
}

func(np *CustomProfile) SetInterdepartureValuesWithRepeatCounts () error {
 //parameters: ProfileHandle Count psaIDVs RcCount ppsaRepeatCounts
 //AgtCustomProfile SetInterdepartureValuesWithRepeatCounts
 return nil
}

