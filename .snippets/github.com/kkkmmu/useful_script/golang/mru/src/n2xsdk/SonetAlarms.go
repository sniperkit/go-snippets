package n2xsdk

type SonetAlarms struct {
 Handler string
}

func(np *SonetAlarms) SetTransmitAlarm () error {
 //parameters: PortHandle Alarm
 //AgtSonetAlarms SetTransmitAlarm
 return nil
}

func(np *SonetAlarms) GetTransmitAlarm ()(string, error) {
 //parameters: PortHandle
 //AgtSonetAlarms GetTransmitAlarm
 return "", nil
}

func(np *SonetAlarms) GetAlarmStartTimestamp ()(string, error) {
 //parameters: PortHandle
 //AgtSonetAlarms GetAlarmStartTimestamp
 return "", nil
}

func(np *SonetAlarms) GetAlarmStopTimestamp ()(string, error) {
 //parameters: PortHandle
 //AgtSonetAlarms GetAlarmStopTimestamp
 return "", nil
}

func(np *SonetAlarms) SetAlarmValue () error {
 //parameters: PortHandle Alarm Value
 //AgtSonetAlarms SetAlarmValue
 return nil
}

func(np *SonetAlarms) GetAlarmValue ()(string, error) {
 //parameters: PortHandle Alarm
 //AgtSonetAlarms GetAlarmValue
 return "", nil
}

