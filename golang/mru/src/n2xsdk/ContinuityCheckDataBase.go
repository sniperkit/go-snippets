package n2xsdk

type ContinuityCheckDataBase struct {
 Handler string
}

func(np *ContinuityCheckDataBase) Enable () error {
 //parameters: DeviceHandle
 //AgtContinuityCheckDataBase Enable
 return nil
}

func(np *ContinuityCheckDataBase) Disable () error {
 //parameters: DeviceHandle
 //AgtContinuityCheckDataBase Disable
 return nil
}

func(np *ContinuityCheckDataBase) IsEnabled () error {
 //parameters: DeviceHandle
 //AgtContinuityCheckDataBase IsEnabled
 return nil
}

func(np *ContinuityCheckDataBase) EnableAutoDiscovery () error {
 //parameters: DeviceHandle
 //AgtContinuityCheckDataBase EnableAutoDiscovery
 return nil
}

func(np *ContinuityCheckDataBase) DisableAutoDiscovery () error {
 //parameters: DeviceHandle
 //AgtContinuityCheckDataBase DisableAutoDiscovery
 return nil
}

func(np *ContinuityCheckDataBase) IsAutoDiscoveryEnabled () error {
 //parameters: DeviceHandle
 //AgtContinuityCheckDataBase IsAutoDiscoveryEnabled
 return nil
}

func(np *ContinuityCheckDataBase) GetCcdbRowId ()(string, error) {
 //parameters: DeviceHandle VlanId1 VlanId2 Level MacAddress MepId
 //AgtContinuityCheckDataBase GetCcdbRowId
 return "", nil
}

func(np *ContinuityCheckDataBase) GetCcdbRowIdList ()(string, error) {
 //parameters: DeviceHandle VlanId1 VlanId2 Level MacAddress MepId
 //AgtContinuityCheckDataBase GetCcdbRowIdList
 return "", nil
}

func(np *ContinuityCheckDataBase) AddRemoteMep () error {
 //parameters: DeviceHandle
 //AgtContinuityCheckDataBase AddRemoteMep
 return nil
}

func(np *ContinuityCheckDataBase) ListRemoteMeps ()(string, error) {
 //parameters: DeviceHandle
 //AgtContinuityCheckDataBase ListRemoteMeps
 return "", nil
}

func(np *ContinuityCheckDataBase) RemoveRemoteMep () error {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase RemoveRemoteMep
 return nil
}

func(np *ContinuityCheckDataBase) GetCcdbNumRows ()(string, error) {
 //parameters: DeviceHandle
 //AgtContinuityCheckDataBase GetCcdbNumRows
 return "", nil
}

func(np *ContinuityCheckDataBase) RemoveAllRemoteMeps () error {
 //parameters: DeviceHandle
 //AgtContinuityCheckDataBase RemoveAllRemoteMeps
 return nil
}

func(np *ContinuityCheckDataBase) SetRemoteMepId () error {
 //parameters: DeviceHandle RemoteMepKey MepId
 //AgtContinuityCheckDataBase SetRemoteMepId
 return nil
}

func(np *ContinuityCheckDataBase) GetRemoteMepId ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetRemoteMepId
 return "", nil
}

func(np *ContinuityCheckDataBase) SetRemoteMepMacAddress () error {
 //parameters: DeviceHandle RemoteMepKey MacAddress
 //AgtContinuityCheckDataBase SetRemoteMepMacAddress
 return nil
}

func(np *ContinuityCheckDataBase) GetRemoteMepMacAddress ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetRemoteMepMacAddress
 return "", nil
}

func(np *ContinuityCheckDataBase) SetRemoteMepMdLevel () error {
 //parameters: DeviceHandle RemoteMepKey Level
 //AgtContinuityCheckDataBase SetRemoteMepMdLevel
 return nil
}

func(np *ContinuityCheckDataBase) GetRemoteMepMdLevel ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetRemoteMepMdLevel
 return "", nil
}

func(np *ContinuityCheckDataBase) SetRemoteMepMegLevel () error {
 //parameters: DeviceHandle RemoteMepKey Level
 //AgtContinuityCheckDataBase SetRemoteMepMegLevel
 return nil
}

func(np *ContinuityCheckDataBase) GetRemoteMepMegLevel ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetRemoteMepMegLevel
 return "", nil
}

func(np *ContinuityCheckDataBase) GetRemoteMepState ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetRemoteMepState
 return "", nil
}

func(np *ContinuityCheckDataBase) GetRemoteMepLastStateChange ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetRemoteMepLastStateChange
 return "", nil
}

func(np *ContinuityCheckDataBase) GetRemoteMepRdiBit ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetRemoteMepRdiBit
 return "", nil
}

func(np *ContinuityCheckDataBase) EnableExpectedConnectivity () error {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase EnableExpectedConnectivity
 return nil
}

func(np *ContinuityCheckDataBase) DisableExpectedConnectivity () error {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase DisableExpectedConnectivity
 return nil
}

func(np *ContinuityCheckDataBase) IsExpectedConnectivityEnabled () error {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase IsExpectedConnectivityEnabled
 return nil
}

func(np *ContinuityCheckDataBase) SetExpectedConnectivity () error {
 //parameters: DeviceHandle
 //AgtContinuityCheckDataBase SetExpectedConnectivity
 return nil
}

func(np *ContinuityCheckDataBase) SetContinuityCheckInterval () error {
 //parameters: DeviceHandle RemoteMepKey CcInterval
 //AgtContinuityCheckDataBase SetContinuityCheckInterval
 return nil
}

func(np *ContinuityCheckDataBase) GetContinuityCheckInterval ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetContinuityCheckInterval
 return "", nil
}

func(np *ContinuityCheckDataBase) SetStandard () error {
 //parameters: DeviceHandle RemoteMepKey Standard
 //AgtContinuityCheckDataBase SetStandard
 return nil
}

func(np *ContinuityCheckDataBase) GetStandard ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetStandard
 return "", nil
}

func(np *ContinuityCheckDataBase) AddVlan () error {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase AddVlan
 return nil
}

func(np *ContinuityCheckDataBase) ListVlans ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase ListVlans
 return "", nil
}

func(np *ContinuityCheckDataBase) RemoveVlan () error {
 //parameters: DeviceHandle RemoteMepKey VlanStackIndex
 //AgtContinuityCheckDataBase RemoveVlan
 return nil
}

func(np *ContinuityCheckDataBase) SetVlanId () error {
 //parameters: DeviceHandle RemoteMepKey VlanStackIndex VlanId
 //AgtContinuityCheckDataBase SetVlanId
 return nil
}

func(np *ContinuityCheckDataBase) GetVlanId ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey VlanStackIndex
 //AgtContinuityCheckDataBase GetVlanId
 return "", nil
}

func(np *ContinuityCheckDataBase) SetMdNameFormat () error {
 //parameters: DeviceHandle RemoteMepKey MdNameFormat
 //AgtContinuityCheckDataBase SetMdNameFormat
 return nil
}

func(np *ContinuityCheckDataBase) SetMdNameString () error {
 //parameters: DeviceHandle RemoteMepKey MdName Length
 //AgtContinuityCheckDataBase SetMdNameString
 return nil
}

func(np *ContinuityCheckDataBase) SetMdNameMacAddrTwoOct () error {
 //parameters: DeviceHandle RemoteMepKey MacAddress Integer
 //AgtContinuityCheckDataBase SetMdNameMacAddrTwoOct
 return nil
}

func(np *ContinuityCheckDataBase) SetMdNameCharString () error {
 //parameters: DeviceHandle RemoteMepKey MdName Length
 //AgtContinuityCheckDataBase SetMdNameCharString
 return nil
}

func(np *ContinuityCheckDataBase) GetMdNameFormat ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetMdNameFormat
 return "", nil
}

func(np *ContinuityCheckDataBase) GetMdNameMacAddrTwoOct ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetMdNameMacAddrTwoOct
 return "", nil
}

func(np *ContinuityCheckDataBase) GetMdNameString ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetMdNameString
 return "", nil
}

func(np *ContinuityCheckDataBase) GetMdNameCharString ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetMdNameCharString
 return "", nil
}

func(np *ContinuityCheckDataBase) SetSmaNameFormat () error {
 //parameters: DeviceHandle RemoteMepKey SmaNameFormat
 //AgtContinuityCheckDataBase SetSmaNameFormat
 return nil
}

func(np *ContinuityCheckDataBase) SetSmaNamePriVid () error {
 //parameters: DeviceHandle RemoteMepKey VlanId
 //AgtContinuityCheckDataBase SetSmaNamePriVid
 return nil
}

func(np *ContinuityCheckDataBase) SetSmaNameCharString () error {
 //parameters: DeviceHandle RemoteMepKey SmaName Length
 //AgtContinuityCheckDataBase SetSmaNameCharString
 return nil
}

func(np *ContinuityCheckDataBase) SetSmaNameTwoOct () error {
 //parameters: DeviceHandle RemoteMepKey Integer
 //AgtContinuityCheckDataBase SetSmaNameTwoOct
 return nil
}

func(np *ContinuityCheckDataBase) SetSmaNameVpnId () error {
 //parameters: DeviceHandle RemoteMepKey VpnId
 //AgtContinuityCheckDataBase SetSmaNameVpnId
 return nil
}

func(np *ContinuityCheckDataBase) GetSmaNameFormat ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetSmaNameFormat
 return "", nil
}

func(np *ContinuityCheckDataBase) GetSmaNamePriVid ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetSmaNamePriVid
 return "", nil
}

func(np *ContinuityCheckDataBase) GetSmaNameCharString ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetSmaNameCharString
 return "", nil
}

func(np *ContinuityCheckDataBase) GetSmaNameTwoOct ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetSmaNameTwoOct
 return "", nil
}

func(np *ContinuityCheckDataBase) GetSmaNameVpnId ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetSmaNameVpnId
 return "", nil
}

func(np *ContinuityCheckDataBase) SetMegIdFormat () error {
 //parameters: DeviceHandle RemoteMepKey MegIdFormat
 //AgtContinuityCheckDataBase SetMegIdFormat
 return nil
}

func(np *ContinuityCheckDataBase) GetMegIdFormat ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetMegIdFormat
 return "", nil
}

func(np *ContinuityCheckDataBase) SetMegIdValue () error {
 //parameters: DeviceHandle RemoteMepKey MegId
 //AgtContinuityCheckDataBase SetMegIdValue
 return nil
}

func(np *ContinuityCheckDataBase) GetMegIdValue ()(string, error) {
 //parameters: DeviceHandle RemoteMepKey
 //AgtContinuityCheckDataBase GetMegIdValue
 return "", nil
}

