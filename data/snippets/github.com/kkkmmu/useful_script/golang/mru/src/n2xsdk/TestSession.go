package n2xsdk

type TestSession struct {
 Handler string
}

func(np *TestSession) GetInterfaceByName ()(string, error) {
 //parameters: InterfaceName
 //AgtTestSession GetInterfaceByName
 return "", nil
}

func(np *TestSession) ListInterfaces ()(string, error) {
 //parameters: 
 //AgtTestSession ListInterfaces
 return "", nil
}

func(np *TestSession) GetInterfaceDetails ()(string, error) {
 //parameters: InterfaceName
 //AgtTestSession GetInterfaceDetails
 return "", nil
}

func(np *TestSession) OpenSession () error {
 //parameters: SessionType Version SessionMode
 //AgtTestSession OpenSession
 return nil
}

func(np *TestSession) CloseSession () error {
 //parameters: 
 //AgtTestSession CloseSession
 return nil
}

func(np *TestSession) CloseSessionForce () error {
 //parameters: 
 //AgtTestSession CloseSessionForce
 return nil
}

func(np *TestSession) GetHandle ()(string, error) {
 //parameters: 
 //AgtTestSession GetHandle
 return "", nil
}

func(np *TestSession) GetType ()(string, error) {
 //parameters: 
 //AgtTestSession GetType
 return "", nil
}

func(np *TestSession) GetMode ()(string, error) {
 //parameters: 
 //AgtTestSession GetMode
 return "", nil
}

func(np *TestSession) GetContext ()(string, error) {
 //parameters: 
 //AgtTestSession GetContext
 return "", nil
}

func(np *TestSession) GetErrorMessageText ()(string, error) {
 //parameters: ErrorCode
 //AgtTestSession GetErrorMessageText
 return "", nil
}

func(np *TestSession) ListSaveableInterfaces ()(string, error) {
 //parameters: 
 //AgtTestSession ListSaveableInterfaces
 return "", nil
}

func(np *TestSession) ListDependencies ()(string, error) {
 //parameters: InterfaceName
 //AgtTestSession ListDependencies
 return "", nil
}

func(np *TestSession) ListAllDependencies ()(string, error) {
 //parameters: InterfaceName
 //AgtTestSession ListAllDependencies
 return "", nil
}

func(np *TestSession) GetSaveableInterfaceDescription ()(string, error) {
 //parameters: InterfaceName
 //AgtTestSession GetSaveableInterfaceDescription
 return "", nil
}

func(np *TestSession) SaveInterfaces () error {
 //parameters: FileName Count InterfaceNames
 //AgtTestSession SaveInterfaces
 return nil
}

func(np *TestSession) SaveSession () error {
 //parameters: FileName
 //AgtTestSession SaveSession
 return nil
}

func(np *TestSession) RestoreSession () error {
 //parameters: FileName
 //AgtTestSession RestoreSession
 return nil
}

func(np *TestSession) ListSavedInterfaces ()(string, error) {
 //parameters: FileName
 //AgtTestSession ListSavedInterfaces
 return "", nil
}

func(np *TestSession) RestoreInterfaces () error {
 //parameters: FileName Count InterfaceNames
 //AgtTestSession RestoreInterfaces
 return nil
}

func(np *TestSession) ResetSession () error {
 //parameters: 
 //AgtTestSession ResetSession
 return nil
}

func(np *TestSession) ResetInterfaces () error {
 //parameters: Count InterfaceNames
 //AgtTestSession ResetInterfaces
 return nil
}

func(np *TestSession) GetNumPorts ()(string, error) {
 //parameters: 
 //AgtTestSession GetNumPorts
 return "", nil
}

func(np *TestSession) RestoreSessionOnPorts () error {
 //parameters: FileName Count PortNames
 //AgtTestSession RestoreSessionOnPorts
 return nil
}

func(np *TestSession) RestoreInterfacesOnPorts () error {
 //parameters: FileName InterfaceCount InterfaceNames PortCount PortNames
 //AgtTestSession RestoreInterfacesOnPorts
 return nil
}

func(np *TestSession) RestoreInterfacesOnSomePorts () error {
 //parameters: FileName InterfaceCount InterfaceNames PortCount PortMapping
 //AgtTestSession RestoreInterfacesOnSomePorts
 return nil
}

func(np *TestSession) ListSavedPorts ()(string, error) {
 //parameters: FileName
 //AgtTestSession ListSavedPorts
 return "", nil
}

func(np *TestSession) GetSavedPortType ()(string, error) {
 //parameters: FileName Label
 //AgtTestSession GetSavedPortType
 return "", nil
}

func(np *TestSession) GetSavedPortComment ()(string, error) {
 //parameters: FileName Label
 //AgtTestSession GetSavedPortComment
 return "", nil
}

func(np *TestSession) GetVersion ()(string, error) {
 //parameters: 
 //AgtTestSession GetVersion
 return "", nil
}

func(np *TestSession) GetBuildNumber ()(string, error) {
 //parameters: 
 //AgtTestSession GetBuildNumber
 return "", nil
}

func(np *TestSession) GetSavedSessionDetails ()(string, error) {
 //parameters: FileName
 //AgtTestSession GetSavedSessionDetails
 return "", nil
}

func(np *TestSession) GetLastRestoreError ()(string, error) {
 //parameters: 
 //AgtTestSession GetLastRestoreError
 return "", nil
}

