package n2xsdk

type LiveCapture struct {
 Handler string
}

func(np *LiveCapture) SetMode () error {
 //parameters: PortHandle Mode
 //AgtLiveCapture SetMode
 return nil
}

func(np *LiveCapture) GetMode ()(string, error) {
 //parameters: PortHandle
 //AgtLiveCapture GetMode
 return "", nil
}

func(np *LiveCapture) Enable () error {
 //parameters: PortHandle
 //AgtLiveCapture Enable
 return nil
}

func(np *LiveCapture) Disable () error {
 //parameters: PortHandle
 //AgtLiveCapture Disable
 return nil
}

func(np *LiveCapture) Start () error {
 //parameters: PortHandle
 //AgtLiveCapture Start
 return nil
}

func(np *LiveCapture) Stop () error {
 //parameters: PortHandle
 //AgtLiveCapture Stop
 return nil
}

func(np *LiveCapture) GetState ()(string, error) {
 //parameters: PortHandle
 //AgtLiveCapture GetState
 return "", nil
}

func(np *LiveCapture) EnableAutomaticStopCaptureOnBufferFull () error {
 //parameters: PortHandle
 //AgtLiveCapture EnableAutomaticStopCaptureOnBufferFull
 return nil
}

func(np *LiveCapture) DisableAutomaticStopCaptureOnBufferFull () error {
 //parameters: PortHandle
 //AgtLiveCapture DisableAutomaticStopCaptureOnBufferFull
 return nil
}

func(np *LiveCapture) IsAutomaticStopCaptureOnBufferFull () error {
 //parameters: PortHandle
 //AgtLiveCapture IsAutomaticStopCaptureOnBufferFull
 return nil
}

func(np *LiveCapture) SetBufferSizeMultiplier () error {
 //parameters: PortHandle MultiplierValue
 //AgtLiveCapture SetBufferSizeMultiplier
 return nil
}

func(np *LiveCapture) GetBufferSizeMultiplier ()(string, error) {
 //parameters: PortHandle
 //AgtLiveCapture GetBufferSizeMultiplier
 return "", nil
}

func(np *LiveCapture) GetMostRecentCaptureFileInSession ()(string, error) {
 //parameters: PortHandle
 //AgtLiveCapture GetMostRecentCaptureFileInSession
 return "", nil
}

func(np *LiveCapture) GetPacketsCapturedCount ()(string, error) {
 //parameters: PortHandle
 //AgtLiveCapture GetPacketsCapturedCount
 return "", nil
}

func(np *LiveCapture) GetPacketsDroppedCount ()(string, error) {
 //parameters: PortHandle
 //AgtLiveCapture GetPacketsDroppedCount
 return "", nil
}

