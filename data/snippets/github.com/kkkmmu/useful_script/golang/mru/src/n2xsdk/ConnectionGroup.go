package n2xsdk

type ConnectionGroup struct {
 Handler string
}

func(np *ConnectionGroup) SetApplicationType () error {
 //parameters: ConnectionGroupHandle ApplicationType
 //AgtConnectionGroup SetApplicationType
 return nil
}

func(np *ConnectionGroup) GetApplicationType ()(string, error) {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup GetApplicationType
 return "", nil
}

func(np *ConnectionGroup) GetClientPort ()(string, error) {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup GetClientPort
 return "", nil
}

func(np *ConnectionGroup) SetServerPort () error {
 //parameters: ConnectionGroupHandle ServerPortHandle
 //AgtConnectionGroup SetServerPort
 return nil
}

func(np *ConnectionGroup) GetServerPort ()(string, error) {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup GetServerPort
 return "", nil
}

func(np *ConnectionGroup) SetNumberOfClients () error {
 //parameters: ConnectionGroupHandle NumberOfClients
 //AgtConnectionGroup SetNumberOfClients
 return nil
}

func(np *ConnectionGroup) GetNumberOfClients ()(string, error) {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup GetNumberOfClients
 return "", nil
}

func(np *ConnectionGroup) SetNumberOfServers () error {
 //parameters: ConnectionGroupHandle NumberOfServers
 //AgtConnectionGroup SetNumberOfServers
 return nil
}

func(np *ConnectionGroup) GetNumberOfServers ()(string, error) {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup GetNumberOfServers
 return "", nil
}

func(np *ConnectionGroup) SetClientDistributionType () error {
 //parameters: ConnectionGroupHandle ClientDistributionType
 //AgtConnectionGroup SetClientDistributionType
 return nil
}

func(np *ConnectionGroup) GetClientDistributionType ()(string, error) {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup GetClientDistributionType
 return "", nil
}

func(np *ConnectionGroup) GetNumberOfConnections ()(string, error) {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup GetNumberOfConnections
 return "", nil
}

func(np *ConnectionGroup) GetNumberOfConnectionsPerServer ()(string, error) {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup GetNumberOfConnectionsPerServer
 return "", nil
}

func(np *ConnectionGroup) SetApplicationFileSize () error {
 //parameters: ConnectionGroupHandle ApplicationFileSize
 //AgtConnectionGroup SetApplicationFileSize
 return nil
}

func(np *ConnectionGroup) GetApplicationFileSize ()(string, error) {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup GetApplicationFileSize
 return "", nil
}

func(np *ConnectionGroup) SetConnectionCloseDelay () error {
 //parameters: ConnectionGroupHandle ConnectionCloseDelay
 //AgtConnectionGroup SetConnectionCloseDelay
 return nil
}

func(np *ConnectionGroup) GetConnectionCloseDelay ()(string, error) {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup GetConnectionCloseDelay
 return "", nil
}

func(np *ConnectionGroup) SetConnectionTeardownType () error {
 //parameters: ConnectionGroupHandle ConnectionTeardownType
 //AgtConnectionGroup SetConnectionTeardownType
 return nil
}

func(np *ConnectionGroup) GetConnectionTeardownType ()(string, error) {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup GetConnectionTeardownType
 return "", nil
}

func(np *ConnectionGroup) SetTargetLoad () error {
 //parameters: ConnectionGroupHandle IsClientSide LoadUnits Load
 //AgtConnectionGroup SetTargetLoad
 return nil
}

func(np *ConnectionGroup) GetTargetLoad ()(string, error) {
 //parameters: ConnectionGroupHandle IsClientSide LoadUnits
 //AgtConnectionGroup GetTargetLoad
 return "", nil
}

func(np *ConnectionGroup) ListSupportedTcpParameters ()(string, error) {
 //parameters: 
 //AgtConnectionGroup ListSupportedTcpParameters
 return "", nil
}

func(np *ConnectionGroup) SetTcpParameterValues () error {
 //parameters: ConnectionGroupHandle IsClientSide NumberOfTcpParameters psaTcpParameters NumberOfValues psaValues
 //AgtConnectionGroup SetTcpParameterValues
 return nil
}

func(np *ConnectionGroup) GetTcpParameterValues ()(string, error) {
 //parameters: ConnectionGroupHandle IsClientSide NumberOfTcpParameters psaTcpParameters
 //AgtConnectionGroup GetTcpParameterValues
 return "", nil
}

func(np *ConnectionGroup) ListTcpParameterValues ()(string, error) {
 //parameters: ConnectionGroupHandle IsClientSide
 //AgtConnectionGroup ListTcpParameterValues
 return "", nil
}

func(np *ConnectionGroup) GetPdu ()(string, error) {
 //parameters: ConnectionGroupHandle IsClientSide
 //AgtConnectionGroup GetPdu
 return "", nil
}

func(np *ConnectionGroup) ListClientServerIpAddressPairs ()(string, error) {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup ListClientServerIpAddressPairs
 return "", nil
}

func(np *ConnectionGroup) Enable () error {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup Enable
 return nil
}

func(np *ConnectionGroup) Disable () error {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup Disable
 return nil
}

func(np *ConnectionGroup) IsEnabled () error {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup IsEnabled
 return nil
}

func(np *ConnectionGroup) SetNumberOfConnectionsPerSecond () error {
 //parameters: ConnectionGroupHandle NumberOfConnectionsPerSecond
 //AgtConnectionGroup SetNumberOfConnectionsPerSecond
 return nil
}

func(np *ConnectionGroup) GetNumberOfConnectionsPerSecond ()(string, error) {
 //parameters: ConnectionGroupHandle
 //AgtConnectionGroup GetNumberOfConnectionsPerSecond
 return "", nil
}

