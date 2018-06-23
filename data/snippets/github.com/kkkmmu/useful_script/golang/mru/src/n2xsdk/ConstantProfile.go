package n2xsdk

type ConstantProfile struct {
 Handler string
}

func(np *ConstantProfile) GetType ()(string, error) {
 //parameters: Handle
 //AgtConstantProfile GetType
 return "", nil
}

func(np *ConstantProfile) GetName ()(string, error) {
 //parameters: Handle
 //AgtConstantProfile GetName
 return "", nil
}

func(np *ConstantProfile) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtConstantProfile GetLockCount
 return "", nil
}

func(np *ConstantProfile) GetSourcePort ()(string, error) {
 //parameters: ProfileHandle
 //AgtConstantProfile GetSourcePort
 return "", nil
}

func(np *ConstantProfile) SetProfileType () error {
 //parameters: ProfileHandle ProfileType
 //AgtConstantProfile SetProfileType
 return nil
}

func(np *ConstantProfile) GetProfileType ()(string, error) {
 //parameters: ProfileHandle
 //AgtConstantProfile GetProfileType
 return "", nil
}

func(np *ConstantProfile) SetAverageLoad () error {
 //parameters: ProfileHandle AverageLoad Units
 //AgtConstantProfile SetAverageLoad
 return nil
}

func(np *ConstantProfile) GetAverageLoad ()(string, error) {
 //parameters: ProfileHandle Units
 //AgtConstantProfile GetAverageLoad
 return "", nil
}

func(np *ConstantProfile) SetMode () error {
 //parameters: ProfileHandle Mode
 //AgtConstantProfile SetMode
 return nil
}

func(np *ConstantProfile) GetMode ()(string, error) {
 //parameters: ProfileHandle
 //AgtConstantProfile GetMode
 return "", nil
}

func(np *ConstantProfile) SetNumberOfPacketsToInject () error {
 //parameters: ProfileHandle NumberOfPackets
 //AgtConstantProfile SetNumberOfPacketsToInject
 return nil
}

func(np *ConstantProfile) GetNumberOfPacketsToInject ()(string, error) {
 //parameters: ProfileHandle
 //AgtConstantProfile GetNumberOfPacketsToInject
 return "", nil
}

func(np *ConstantProfile) SendSingleInjection () error {
 //parameters: ProfileHandle
 //AgtConstantProfile SendSingleInjection
 return nil
}

func(np *ConstantProfile) IsSendSingleInjectionInProgress () error {
 //parameters: ProfileHandle
 //AgtConstantProfile IsSendSingleInjectionInProgress
 return nil
}

func(np *ConstantProfile) ListStreamGroups ()(string, error) {
 //parameters: ProfileHandle
 //AgtConstantProfile ListStreamGroups
 return "", nil
}

func(np *ConstantProfile) AddStreamGroups () error {
 //parameters: ProfileHandle Count psaStreamGroupHandles
 //AgtConstantProfile AddStreamGroups
 return nil
}

func(np *ConstantProfile) GetInterdepartureValues ()(string, error) {
 //parameters: ProfileHandle
 //AgtConstantProfile GetInterdepartureValues
 return "", nil
}

func(np *ConstantProfile) GetInterdepartureRepeatCounts ()(string, error) {
 //parameters: ProfileHandle
 //AgtConstantProfile GetInterdepartureRepeatCounts
 return "", nil
}

func(np *ConstantProfile) GetInterdepartureValueResolution ()(string, error) {
 //parameters: ProfileHandle
 //AgtConstantProfile GetInterdepartureValueResolution
 return "", nil
}

func(np *ConstantProfile) Enable () error {
 //parameters: ProfileHandle
 //AgtConstantProfile Enable
 return nil
}

func(np *ConstantProfile) Disable () error {
 //parameters: ProfileHandle
 //AgtConstantProfile Disable
 return nil
}

func(np *ConstantProfile) IsEnabled () error {
 //parameters: ProfileHandle
 //AgtConstantProfile IsEnabled
 return nil
}

