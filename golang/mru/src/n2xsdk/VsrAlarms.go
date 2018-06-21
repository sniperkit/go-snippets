package n2xsdk

type VsrAlarms struct {
 Handler string
}

func(np *VsrAlarms) SetTransmitAlarm () error {
 //parameters: PortHandle Alarm
 //AgtVsrAlarms SetTransmitAlarm
 return nil
}

func(np *VsrAlarms) GetTransmitAlarm ()(string, error) {
 //parameters: PortHandle
 //AgtVsrAlarms GetTransmitAlarm
 return "", nil
}

