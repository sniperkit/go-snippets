package n2xsdk

type SessionManager struct {
 Handler string
}

func(np *SessionManager) GetInterfaceByName ()(string, error) {
 //parameters: InterfaceName
 //AgtSessionManager GetInterfaceByName
 return "", nil
}

func(np *SessionManager) ListInterfaces ()(string, error) {
 //parameters: 
 //AgtSessionManager ListInterfaces
 return "", nil
}

func(np *SessionManager) GetInterfaceDetails ()(string, error) {
 //parameters: InterfaceName
 //AgtSessionManager GetInterfaceDetails
 return "", nil
}

func(np *SessionManager) OpenSession () error {
 //parameters: SessionType SessionMode
 //AgtSessionManager OpenSession
 return nil
}

func(np *SessionManager) OpenSessionByVersion () error {
 //parameters: SessionType Version SessionMode
 //AgtSessionManager OpenSessionByVersion
 return nil
}

func(np *SessionManager) CloseSession () error {
 //parameters: SessionHandle
 //AgtSessionManager CloseSession
 return nil
}

func(np *SessionManager) ListSessionTypes ()(string, error) {
 //parameters: 
 //AgtSessionManager ListSessionTypes
 return "", nil
}

func(np *SessionManager) ListSessionVersions ()(string, error) {
 //parameters: SessionType
 //AgtSessionManager ListSessionVersions
 return "", nil
}

func(np *SessionManager) ListOpenSessions ()(string, error) {
 //parameters: 
 //AgtSessionManager ListOpenSessions
 return "", nil
}

func(np *SessionManager) ListOpenSessionsByType ()(string, error) {
 //parameters: SessionType
 //AgtSessionManager ListOpenSessionsByType
 return "", nil
}

func(np *SessionManager) ListOpenSessionsByVersion ()(string, error) {
 //parameters: SessionType Version
 //AgtSessionManager ListOpenSessionsByVersion
 return "", nil
}

func(np *SessionManager) GetSessionVersion ()(string, error) {
 //parameters: SessionHandle
 //AgtSessionManager GetSessionVersion
 return "", nil
}

func(np *SessionManager) GetSessionType ()(string, error) {
 //parameters: SessionHandle
 //AgtSessionManager GetSessionType
 return "", nil
}

func(np *SessionManager) GetSessionPort ()(string, error) {
 //parameters: SessionHandle
 //AgtSessionManager GetSessionPort
 return "", nil
}

func(np *SessionManager) GetSessionContext ()(string, error) {
 //parameters: SessionHandle
 //AgtSessionManager GetSessionContext
 return "", nil
}

func(np *SessionManager) SetSessionLabel () error {
 //parameters: SessionHandle SessionLabel
 //AgtSessionManager SetSessionLabel
 return nil
}

func(np *SessionManager) GetSessionLabel ()(string, error) {
 //parameters: SessionHandle
 //AgtSessionManager GetSessionLabel
 return "", nil
}

func(np *SessionManager) GetSessionPid ()(string, error) {
 //parameters: SessionHandle
 //AgtSessionManager GetSessionPid
 return "", nil
}

func(np *SessionManager) GetNumGuiConnections ()(string, error) {
 //parameters: SessionHandle GuiType
 //AgtSessionManager GetNumGuiConnections
 return "", nil
}

func(np *SessionManager) GetMaxGuiConnectionsByVersion ()(string, error) {
 //parameters: SessionType Version GuiType
 //AgtSessionManager GetMaxGuiConnectionsByVersion
 return "", nil
}

func(np *SessionManager) SetMaxGuiConnectionsByVersion () error {
 //parameters: SessionType Version GuiType MaxGuiSessions
 //AgtSessionManager SetMaxGuiConnectionsByVersion
 return nil
}

func(np *SessionManager) GetMaxGuiConnections ()(string, error) {
 //parameters: SessionType GuiType
 //AgtSessionManager GetMaxGuiConnections
 return "", nil
}

func(np *SessionManager) SetMaxGuiConnections () error {
 //parameters: SessionType GuiType MaxGuiSessions
 //AgtSessionManager SetMaxGuiConnections
 return nil
}

func(np *SessionManager) GetDefaultSessionVersion ()(string, error) {
 //parameters: SessionType
 //AgtSessionManager GetDefaultSessionVersion
 return "", nil
}

func(np *SessionManager) IsMaxGuiConnectionsUpdateable () error {
 //parameters: 
 //AgtSessionManager IsMaxGuiConnectionsUpdateable
 return nil
}

func(np *SessionManager) GetGlobalMaxGuiConnections ()(string, error) {
 //parameters: 
 //AgtSessionManager GetGlobalMaxGuiConnections
 return "", nil
}

