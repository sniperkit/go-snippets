package n2xsdk

type VsrInterface struct {
 Handler string
}

func(np *VsrInterface) SetCodeViolationErrorMode () error {
 //parameters: PortHandle ErrorMode
 //AgtVsrInterface SetCodeViolationErrorMode
 return nil
}

func(np *VsrInterface) GetCodeViolationErrorMode ()(string, error) {
 //parameters: PortHandle
 //AgtVsrInterface GetCodeViolationErrorMode
 return "", nil
}

func(np *VsrInterface) EnableCrcCorrection () error {
 //parameters: PortHandle
 //AgtVsrInterface EnableCrcCorrection
 return nil
}

func(np *VsrInterface) DisableCrcCorrection () error {
 //parameters: PortHandle
 //AgtVsrInterface DisableCrcCorrection
 return nil
}

func(np *VsrInterface) IsCrcCorrectionEnabled () error {
 //parameters: PortHandle
 //AgtVsrInterface IsCrcCorrectionEnabled
 return nil
}

func(np *VsrInterface) EnableChannelProtection () error {
 //parameters: PortHandle
 //AgtVsrInterface EnableChannelProtection
 return nil
}

func(np *VsrInterface) DisableChannelProtection () error {
 //parameters: PortHandle
 //AgtVsrInterface DisableChannelProtection
 return nil
}

func(np *VsrInterface) IsChannelProtectionEnabled () error {
 //parameters: PortHandle
 //AgtVsrInterface IsChannelProtectionEnabled
 return nil
}

func(np *VsrInterface) EnableTxLoopback () error {
 //parameters: PortHandle
 //AgtVsrInterface EnableTxLoopback
 return nil
}

func(np *VsrInterface) DisableTxLoopback () error {
 //parameters: PortHandle
 //AgtVsrInterface DisableTxLoopback
 return nil
}

func(np *VsrInterface) IsTxLoopbackEnabled () error {
 //parameters: PortHandle
 //AgtVsrInterface IsTxLoopbackEnabled
 return nil
}

func(np *VsrInterface) IsVsrSupported () error {
 //parameters: PortHandle
 //AgtVsrInterface IsVsrSupported
 return nil
}

func(np *VsrInterface) SetErrorChannels () error {
 //parameters: PortHandle ErrorType Length psaChannels
 //AgtVsrInterface SetErrorChannels
 return nil
}

func(np *VsrInterface) GetErrorChannels ()(string, error) {
 //parameters: PortHandle ErrorType
 //AgtVsrInterface GetErrorChannels
 return "", nil
}

func(np *VsrInterface) PatchcordReversalOn () error {
 //parameters: PortHandle
 //AgtVsrInterface PatchcordReversalOn
 return nil
}

func(np *VsrInterface) PatchcordReversalOff () error {
 //parameters: PortHandle
 //AgtVsrInterface PatchcordReversalOff
 return nil
}

func(np *VsrInterface) IsPatchcordReversalOn () error {
 //parameters: PortHandle
 //AgtVsrInterface IsPatchcordReversalOn
 return nil
}

