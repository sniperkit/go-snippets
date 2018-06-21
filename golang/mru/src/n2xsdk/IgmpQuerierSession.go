package n2xsdk

type IgmpQuerierSession struct {
 Handler string
}

func(np *IgmpQuerierSession) GetInterfaceIpAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession GetInterfaceIpAddress
 return "", nil
}

func(np *IgmpQuerierSession) SetSutIpAddress () error {
 //parameters: SessionHandle SutIpAddress
 //AgtIgmpQuerierSession SetSutIpAddress
 return nil
}

func(np *IgmpQuerierSession) GetSutIpAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession GetSutIpAddress
 return "", nil
}

func(np *IgmpQuerierSession) Enable () error {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession Enable
 return nil
}

func(np *IgmpQuerierSession) EnableAllSessions () error {
 //parameters: 
 //AgtIgmpQuerierSession EnableAllSessions
 return nil
}

func(np *IgmpQuerierSession) Disable () error {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession Disable
 return nil
}

func(np *IgmpQuerierSession) DisableAllSessions () error {
 //parameters: 
 //AgtIgmpQuerierSession DisableAllSessions
 return nil
}

func(np *IgmpQuerierSession) GetEnableFlag ()(string, error) {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession GetEnableFlag
 return "", nil
}

func(np *IgmpQuerierSession) GetState ()(string, error) {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession GetState
 return "", nil
}

func(np *IgmpQuerierSession) GetStateSummary ()(string, error) {
 //parameters: 
 //AgtIgmpQuerierSession GetStateSummary
 return "", nil
}

func(np *IgmpQuerierSession) GetPortStateSummary ()(string, error) {
 //parameters: PortHandle
 //AgtIgmpQuerierSession GetPortStateSummary
 return "", nil
}

func(np *IgmpQuerierSession) GetVersion ()(string, error) {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession GetVersion
 return "", nil
}

func(np *IgmpQuerierSession) GetParameter ()(string, error) {
 //parameters: SessionHandle Parameter
 //AgtIgmpQuerierSession GetParameter
 return "", nil
}

func(np *IgmpQuerierSession) GetInfo ()(string, error) {
 //parameters: SessionHandle Info
 //AgtIgmpQuerierSession GetInfo
 return "", nil
}

func(np *IgmpQuerierSession) AddGeneralQueryMessage () error {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession AddGeneralQueryMessage
 return nil
}

func(np *IgmpQuerierSession) RemoveGeneralQueryMessage () error {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession RemoveGeneralQueryMessage
 return nil
}

func(np *IgmpQuerierSession) IsGeneralQueryAdded () error {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession IsGeneralQueryAdded
 return nil
}

func(np *IgmpQuerierSession) SetGeneralQueryInterval () error {
 //parameters: SessionHandle QueryInterval
 //AgtIgmpQuerierSession SetGeneralQueryInterval
 return nil
}

func(np *IgmpQuerierSession) GetGeneralQueryInterval ()(string, error) {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession GetGeneralQueryInterval
 return "", nil
}

func(np *IgmpQuerierSession) AddGroupQueryMessage () error {
 //parameters: SessionHandle GroupPoolHandle
 //AgtIgmpQuerierSession AddGroupQueryMessage
 return nil
}

func(np *IgmpQuerierSession) RemoveGroupQueryMessage () error {
 //parameters: SessionHandle GroupPoolHandle
 //AgtIgmpQuerierSession RemoveGroupQueryMessage
 return nil
}

func(np *IgmpQuerierSession) RemoveAllQueryMessages () error {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession RemoveAllQueryMessages
 return nil
}

func(np *IgmpQuerierSession) ListGroupQueryMessages ()(string, error) {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession ListGroupQueryMessages
 return "", nil
}

func(np *IgmpQuerierSession) AddSourcePools () error {
 //parameters: SessionHandle GroupPoolHandle SourcePoolHandles
 //AgtIgmpQuerierSession AddSourcePools
 return nil
}

func(np *IgmpQuerierSession) RemoveSourcePools () error {
 //parameters: SessionHandle GroupPoolHandle SourcePoolHandles
 //AgtIgmpQuerierSession RemoveSourcePools
 return nil
}

func(np *IgmpQuerierSession) RemoveAllSourcePools () error {
 //parameters: SessionHandle GroupPoolHandle
 //AgtIgmpQuerierSession RemoveAllSourcePools
 return nil
}

func(np *IgmpQuerierSession) ListSourcePools ()(string, error) {
 //parameters: SessionHandle GroupPoolHandle
 //AgtIgmpQuerierSession ListSourcePools
 return "", nil
}

func(np *IgmpQuerierSession) GetNumberOfSourcesInGroupPool ()(string, error) {
 //parameters: SessionHandle GroupPoolHandle
 //AgtIgmpQuerierSession GetNumberOfSourcesInGroupPool
 return "", nil
}

func(np *IgmpQuerierSession) ListGroupPoolsInvalid ()(string, error) {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession ListGroupPoolsInvalid
 return "", nil
}

func(np *IgmpQuerierSession) IsInvalid () error {
 //parameters: SessionHandle GroupPoolHandle
 //AgtIgmpQuerierSession IsInvalid
 return nil
}

func(np *IgmpQuerierSession) SetGroupQueryParameters () error {
 //parameters: SessionHandle GroupPoolHandle Lmqi Lmqc
 //AgtIgmpQuerierSession SetGroupQueryParameters
 return nil
}

func(np *IgmpQuerierSession) GetGroupQueryParameters ()(string, error) {
 //parameters: SessionHandle GroupPoolHandle
 //AgtIgmpQuerierSession GetGroupQueryParameters
 return "", nil
}

func(np *IgmpQuerierSession) StartGeneralQueryMessage () error {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession StartGeneralQueryMessage
 return nil
}

func(np *IgmpQuerierSession) StopGeneralQueryMessage () error {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession StopGeneralQueryMessage
 return nil
}

func(np *IgmpQuerierSession) IsGeneralQueryRunning () error {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession IsGeneralQueryRunning
 return nil
}

func(np *IgmpQuerierSession) StartGroupQueryMessages () error {
 //parameters: SessionHandle GroupPoolHandles
 //AgtIgmpQuerierSession StartGroupQueryMessages
 return nil
}

func(np *IgmpQuerierSession) StartAllQueryMessages () error {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession StartAllQueryMessages
 return nil
}

func(np *IgmpQuerierSession) StopGroupQueryMessages () error {
 //parameters: SessionHandle GroupPoolHandles
 //AgtIgmpQuerierSession StopGroupQueryMessages
 return nil
}

func(np *IgmpQuerierSession) StopAllQueryMessages () error {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession StopAllQueryMessages
 return nil
}

func(np *IgmpQuerierSession) ListGroupQueryMessagesRunning ()(string, error) {
 //parameters: SessionHandle
 //AgtIgmpQuerierSession ListGroupQueryMessagesRunning
 return "", nil
}

func(np *IgmpQuerierSession) IsGroupQueryRunning () error {
 //parameters: SessionHandle GroupPoolHandle
 //AgtIgmpQuerierSession IsGroupQueryRunning
 return nil
}

func(np *IgmpQuerierSession) SetInterfaceIpAddress () error {
 //parameters: SessionHandle TesterIpAddress
 //AgtIgmpQuerierSession SetInterfaceIpAddress
 return nil
}

func(np *IgmpQuerierSession) SetParameter () error {
 //parameters: 
 //AgtIgmpQuerierSession SetParameter
 return nil
}

