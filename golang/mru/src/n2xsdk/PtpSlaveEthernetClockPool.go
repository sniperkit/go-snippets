package n2xsdk

type PtpSlaveEthernetClockPool struct {
 Handler string
}

func(np *PtpSlaveEthernetClockPool) SetAnnounceReceiptTimeout () error {
 //parameters: SessionPoolHandle AnnounceTimeout
 //AgtPtpSlaveEthernetClockPool SetAnnounceReceiptTimeout
 return nil
}

func(np *PtpSlaveEthernetClockPool) GetAnnounceReceiptTimeout ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool GetAnnounceReceiptTimeout
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) EnableDelayRequest () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool EnableDelayRequest
 return nil
}

func(np *PtpSlaveEthernetClockPool) DisableDelayRequest () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool DisableDelayRequest
 return nil
}

func(np *PtpSlaveEthernetClockPool) IsEnabledDelayRequest () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool IsEnabledDelayRequest
 return nil
}

func(np *PtpSlaveEthernetClockPool) SetDelayRequestAttributes () error {
 //parameters: SessionPoolHandle IntervalMode LogMessageInterval RandomDistributionEnabled
 //AgtPtpSlaveEthernetClockPool SetDelayRequestAttributes
 return nil
}

func(np *PtpSlaveEthernetClockPool) GetDelayRequestAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool GetDelayRequestAttributes
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) EnableActivityTrace () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool EnableActivityTrace
 return nil
}

func(np *PtpSlaveEthernetClockPool) DisableActivityTrace () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool DisableActivityTrace
 return nil
}

func(np *PtpSlaveEthernetClockPool) IsEnabledActivityTrace () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool IsEnabledActivityTrace
 return nil
}

func(np *PtpSlaveEthernetClockPool) EnableActivityTraceFilter () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool EnableActivityTraceFilter
 return nil
}

func(np *PtpSlaveEthernetClockPool) DisableActivityTraceFilter () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool DisableActivityTraceFilter
 return nil
}

func(np *PtpSlaveEthernetClockPool) IsEnabledActivityTraceFilter () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool IsEnabledActivityTraceFilter
 return nil
}

func(np *PtpSlaveEthernetClockPool) SetActivityTraceAttributes () error {
 //parameters: SessionPoolHandle InstanceNumber
 //AgtPtpSlaveEthernetClockPool SetActivityTraceAttributes
 return nil
}

func(np *PtpSlaveEthernetClockPool) GetActivityTraceAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool GetActivityTraceAttributes
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) RefreshActivityTrace () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool RefreshActivityTrace
 return nil
}

func(np *PtpSlaveEthernetClockPool) ClearActivityTrace () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool ClearActivityTrace
 return nil
}

func(np *PtpSlaveEthernetClockPool) ListActivityTraceRows ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool ListActivityTraceRows
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) GetActivityTraceRow ()(string, error) {
 //parameters: SessionPoolHandle ActivityTraceRowId
 //AgtPtpSlaveEthernetClockPool GetActivityTraceRow
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) SetDomainNumber () error {
 //parameters: SessionPoolHandle DomainNumber
 //AgtPtpSlaveEthernetClockPool SetDomainNumber
 return nil
}

func(np *PtpSlaveEthernetClockPool) GetDomainNumber ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool GetDomainNumber
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) SetPortNumber () error {
 //parameters: SessionPoolHandle PortNumber
 //AgtPtpSlaveEthernetClockPool SetPortNumber
 return nil
}

func(np *PtpSlaveEthernetClockPool) GetPortNumber ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool GetPortNumber
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) SetClockId () error {
 //parameters: SessionPoolHandle ClockId
 //AgtPtpSlaveEthernetClockPool SetClockId
 return nil
}

func(np *PtpSlaveEthernetClockPool) GetClockId ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool GetClockId
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) SetClockIdIncrementingRange () error {
 //parameters: SessionPoolHandle FirstClockId Repeat Increment
 //AgtPtpSlaveEthernetClockPool SetClockIdIncrementingRange
 return nil
}

func(np *PtpSlaveEthernetClockPool) GetClockIdIncrementingRange ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool GetClockIdIncrementingRange
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) SetClockIdList ()(string, error) {
 //parameters: SessionPoolHandle ClockIdList
 //AgtPtpSlaveEthernetClockPool SetClockIdList
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) GetClockIdList ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool GetClockIdList
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) SetPathDelayMode () error {
 //parameters: SessionPoolHandle PathDelayMode
 //AgtPtpSlaveEthernetClockPool SetPathDelayMode
 return nil
}

func(np *PtpSlaveEthernetClockPool) GetPathDelayMode ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool GetPathDelayMode
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) EnablePeerDelayRequest () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool EnablePeerDelayRequest
 return nil
}

func(np *PtpSlaveEthernetClockPool) DisablePeerDelayRequest () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool DisablePeerDelayRequest
 return nil
}

func(np *PtpSlaveEthernetClockPool) IsEnabledPeerDelayRequest () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool IsEnabledPeerDelayRequest
 return nil
}

func(np *PtpSlaveEthernetClockPool) EnablePeerDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool EnablePeerDelayResponse
 return nil
}

func(np *PtpSlaveEthernetClockPool) DisablePeerDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool DisablePeerDelayResponse
 return nil
}

func(np *PtpSlaveEthernetClockPool) IsEnabledPeerDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool IsEnabledPeerDelayResponse
 return nil
}

func(np *PtpSlaveEthernetClockPool) EnablePeerDelayResponseFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool EnablePeerDelayResponseFollowUp
 return nil
}

func(np *PtpSlaveEthernetClockPool) DisablePeerDelayResponseFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool DisablePeerDelayResponseFollowUp
 return nil
}

func(np *PtpSlaveEthernetClockPool) IsEnabledPeerDelayResponseFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool IsEnabledPeerDelayResponseFollowUp
 return nil
}

func(np *PtpSlaveEthernetClockPool) SetPeerToPeerPathDelayAttributes () error {
 //parameters: SessionPoolHandle LogMessageInterval TwoStepFlagMode FollowUpDelay
 //AgtPtpSlaveEthernetClockPool SetPeerToPeerPathDelayAttributes
 return nil
}

func(np *PtpSlaveEthernetClockPool) GetPeerToPeerPathDelayAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool GetPeerToPeerPathDelayAttributes
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) SetTimestampCalibration () error {
 //parameters: SessionPoolHandle TxOffset RxOffset
 //AgtPtpSlaveEthernetClockPool SetTimestampCalibration
 return nil
}

func(np *PtpSlaveEthernetClockPool) GetTimestampCalibration ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool GetTimestampCalibration
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) GetVlanPriority ()(string, error) {
 //parameters: DeviceHandle VlanTagIndex
 //AgtPtpSlaveEthernetClockPool GetVlanPriority
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) SetVlanPriority () error {
 //parameters: DeviceHandle VlanTagIndex VlanPriority
 //AgtPtpSlaveEthernetClockPool SetVlanPriority
 return nil
}

func(np *PtpSlaveEthernetClockPool) SetDestinationMacAddress () error {
 //parameters: SessionPoolHandle DestinationMac
 //AgtPtpSlaveEthernetClockPool SetDestinationMacAddress
 return nil
}

func(np *PtpSlaveEthernetClockPool) GetDestinationMacAddress ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool GetDestinationMacAddress
 return "", nil
}

func(np *PtpSlaveEthernetClockPool) SetPeerMacAddress () error {
 //parameters: SessionPoolHandle PeerMac
 //AgtPtpSlaveEthernetClockPool SetPeerMacAddress
 return nil
}

func(np *PtpSlaveEthernetClockPool) GetPeerMacAddress ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpSlaveEthernetClockPool GetPeerMacAddress
 return "", nil
}

