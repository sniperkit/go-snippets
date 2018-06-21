package n2xsdk

type VsrError struct {
 Handler string
}

func(np *VsrError) VsrErrorOn () error {
 //parameters: PortHandle VsrError
 //AgtVsrError VsrErrorOn
 return nil
}

func(np *VsrError) VsrErrorOff () error {
 //parameters: PortHandle VsrError
 //AgtVsrError VsrErrorOff
 return nil
}

func(np *VsrError) IsVsrErrorOn () error {
 //parameters: PortHandle VsrError
 //AgtVsrError IsVsrErrorOn
 return nil
}

