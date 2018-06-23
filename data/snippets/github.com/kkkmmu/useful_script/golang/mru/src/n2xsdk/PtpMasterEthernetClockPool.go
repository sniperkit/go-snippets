package n2xsdk

type PtpMasterEthernetClockPool struct {
 Handler string
}

func(np *PtpMasterEthernetClockPool) SetClockQuality () error {
 //parameters: SessionPoolHandle ClockClass ClockAccuracy OffsetScaledLogVariance Priority1 Priority2 AlternateMasterFlag StepsRemoved TimeSource
 //AgtPtpMasterEthernetClockPool SetClockQuality
 return nil
}

func(np *PtpMasterEthernetClockPool) GetClockQuality ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool GetClockQuality
 return "", nil
}

func(np *PtpMasterEthernetClockPool) EnableAnnounce () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool EnableAnnounce
 return nil
}

func(np *PtpMasterEthernetClockPool) DisableAnnounce () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool DisableAnnounce
 return nil
}

func(np *PtpMasterEthernetClockPool) IsEnabledAnnounce () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool IsEnabledAnnounce
 return nil
}

func(np *PtpMasterEthernetClockPool) SetAnnounceAttributes () error {
 //parameters: SessionPoolHandle LogMessageInterval ReceiptTimeout
 //AgtPtpMasterEthernetClockPool SetAnnounceAttributes
 return nil
}

func(np *PtpMasterEthernetClockPool) GetAnnounceAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool GetAnnounceAttributes
 return "", nil
}

func(np *PtpMasterEthernetClockPool) EnableSync () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool EnableSync
 return nil
}

func(np *PtpMasterEthernetClockPool) DisableSync () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool DisableSync
 return nil
}

func(np *PtpMasterEthernetClockPool) IsEnabledSync () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool IsEnabledSync
 return nil
}

func(np *PtpMasterEthernetClockPool) SetSyncAttributes () error {
 //parameters: SessionPoolHandle TwoStepFlagMode LogMessageInterval
 //AgtPtpMasterEthernetClockPool SetSyncAttributes
 return nil
}

func(np *PtpMasterEthernetClockPool) GetSyncAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool GetSyncAttributes
 return "", nil
}

func(np *PtpMasterEthernetClockPool) EnableFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool EnableFollowUp
 return nil
}

func(np *PtpMasterEthernetClockPool) DisableFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool DisableFollowUp
 return nil
}

func(np *PtpMasterEthernetClockPool) IsEnabledFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool IsEnabledFollowUp
 return nil
}

func(np *PtpMasterEthernetClockPool) SetFollowUpAttributes () error {
 //parameters: SessionPoolHandle FollowUpDelay
 //AgtPtpMasterEthernetClockPool SetFollowUpAttributes
 return nil
}

func(np *PtpMasterEthernetClockPool) GetFollowUpAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool GetFollowUpAttributes
 return "", nil
}

func(np *PtpMasterEthernetClockPool) EnableDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool EnableDelayResponse
 return nil
}

func(np *PtpMasterEthernetClockPool) DisableDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool DisableDelayResponse
 return nil
}

func(np *PtpMasterEthernetClockPool) IsEnabledDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool IsEnabledDelayResponse
 return nil
}

func(np *PtpMasterEthernetClockPool) SetDelayResponseAttributes () error {
 //parameters: SessionPoolHandle LogMessageInterval
 //AgtPtpMasterEthernetClockPool SetDelayResponseAttributes
 return nil
}

func(np *PtpMasterEthernetClockPool) GetDelayResponseAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool GetDelayResponseAttributes
 return "", nil
}

func(np *PtpMasterEthernetClockPool) SetDomainNumber () error {
 //parameters: SessionPoolHandle DomainNumber
 //AgtPtpMasterEthernetClockPool SetDomainNumber
 return nil
}

func(np *PtpMasterEthernetClockPool) GetDomainNumber ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool GetDomainNumber
 return "", nil
}

func(np *PtpMasterEthernetClockPool) SetPortNumber () error {
 //parameters: SessionPoolHandle PortNumber
 //AgtPtpMasterEthernetClockPool SetPortNumber
 return nil
}

func(np *PtpMasterEthernetClockPool) GetPortNumber ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool GetPortNumber
 return "", nil
}

func(np *PtpMasterEthernetClockPool) SetClockId () error {
 //parameters: SessionPoolHandle ClockId
 //AgtPtpMasterEthernetClockPool SetClockId
 return nil
}

func(np *PtpMasterEthernetClockPool) GetClockId ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool GetClockId
 return "", nil
}

func(np *PtpMasterEthernetClockPool) SetClockIdIncrementingRange () error {
 //parameters: SessionPoolHandle FirstClockId Repeat Increment
 //AgtPtpMasterEthernetClockPool SetClockIdIncrementingRange
 return nil
}

func(np *PtpMasterEthernetClockPool) GetClockIdIncrementingRange ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool GetClockIdIncrementingRange
 return "", nil
}

func(np *PtpMasterEthernetClockPool) SetClockIdList ()(string, error) {
 //parameters: SessionPoolHandle ClockIdList
 //AgtPtpMasterEthernetClockPool SetClockIdList
 return "", nil
}

func(np *PtpMasterEthernetClockPool) GetClockIdList ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool GetClockIdList
 return "", nil
}

func(np *PtpMasterEthernetClockPool) SetPathDelayMode () error {
 //parameters: SessionPoolHandle PathDelayMode
 //AgtPtpMasterEthernetClockPool SetPathDelayMode
 return nil
}

func(np *PtpMasterEthernetClockPool) GetPathDelayMode ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool GetPathDelayMode
 return "", nil
}

func(np *PtpMasterEthernetClockPool) EnablePeerDelayRequest () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool EnablePeerDelayRequest
 return nil
}

func(np *PtpMasterEthernetClockPool) DisablePeerDelayRequest () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool DisablePeerDelayRequest
 return nil
}

func(np *PtpMasterEthernetClockPool) IsEnabledPeerDelayRequest () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool IsEnabledPeerDelayRequest
 return nil
}

func(np *PtpMasterEthernetClockPool) EnablePeerDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool EnablePeerDelayResponse
 return nil
}

func(np *PtpMasterEthernetClockPool) DisablePeerDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool DisablePeerDelayResponse
 return nil
}

func(np *PtpMasterEthernetClockPool) IsEnabledPeerDelayResponse () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool IsEnabledPeerDelayResponse
 return nil
}

func(np *PtpMasterEthernetClockPool) EnablePeerDelayResponseFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool EnablePeerDelayResponseFollowUp
 return nil
}

func(np *PtpMasterEthernetClockPool) DisablePeerDelayResponseFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool DisablePeerDelayResponseFollowUp
 return nil
}

func(np *PtpMasterEthernetClockPool) IsEnabledPeerDelayResponseFollowUp () error {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool IsEnabledPeerDelayResponseFollowUp
 return nil
}

func(np *PtpMasterEthernetClockPool) SetPeerToPeerPathDelayAttributes () error {
 //parameters: SessionPoolHandle LogMessageInterval TwoStepFlagMode FollowUpDelay
 //AgtPtpMasterEthernetClockPool SetPeerToPeerPathDelayAttributes
 return nil
}

func(np *PtpMasterEthernetClockPool) GetPeerToPeerPathDelayAttributes ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool GetPeerToPeerPathDelayAttributes
 return "", nil
}

func(np *PtpMasterEthernetClockPool) SetTimestampCalibration () error {
 //parameters: SessionPoolHandle TxOffset RxOffset
 //AgtPtpMasterEthernetClockPool SetTimestampCalibration
 return nil
}

func(np *PtpMasterEthernetClockPool) GetTimestampCalibration ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool GetTimestampCalibration
 return "", nil
}

func(np *PtpMasterEthernetClockPool) GetVlanPriority ()(string, error) {
 //parameters: DeviceHandle VlanTagIndex
 //AgtPtpMasterEthernetClockPool GetVlanPriority
 return "", nil
}

func(np *PtpMasterEthernetClockPool) SetVlanPriority () error {
 //parameters: DeviceHandle VlanTagIndex VlanPriority
 //AgtPtpMasterEthernetClockPool SetVlanPriority
 return nil
}

func(np *PtpMasterEthernetClockPool) SetDestinationMacAddress () error {
 //parameters: SessionPoolHandle DestinationMac
 //AgtPtpMasterEthernetClockPool SetDestinationMacAddress
 return nil
}

func(np *PtpMasterEthernetClockPool) GetDestinationMacAddress ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool GetDestinationMacAddress
 return "", nil
}

func(np *PtpMasterEthernetClockPool) SetPeerMacAddress () error {
 //parameters: SessionPoolHandle PeerMac
 //AgtPtpMasterEthernetClockPool SetPeerMacAddress
 return nil
}

func(np *PtpMasterEthernetClockPool) GetPeerMacAddress ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPtpMasterEthernetClockPool GetPeerMacAddress
 return "", nil
}

