package n2xsdk

type FieldStatisticsControl struct {
 Handler string
}

func(np *FieldStatisticsControl) SetIndexField () error {
 //parameters: NumPorts pPortHandles Field
 //AgtFieldStatisticsControl SetIndexField
 return nil
}

func(np *FieldStatisticsControl) GetIndexField ()(string, error) {
 //parameters: NumPorts pPortHandles
 //AgtFieldStatisticsControl GetIndexField
 return "", nil
}

func(np *FieldStatisticsControl) ListAvailableFieldsOnPort ()(string, error) {
 //parameters: PortHandle
 //AgtFieldStatisticsControl ListAvailableFieldsOnPort
 return "", nil
}

func(np *FieldStatisticsControl) SetFrameFilter () error {
 //parameters: NumPorts pPortHandles FrameFilter
 //AgtFieldStatisticsControl SetFrameFilter
 return nil
}

func(np *FieldStatisticsControl) GetFrameFilter ()(string, error) {
 //parameters: NumPorts pPortHandles
 //AgtFieldStatisticsControl GetFrameFilter
 return "", nil
}

func(np *FieldStatisticsControl) SetFieldPageOffset () error {
 //parameters: NumPorts pPortHandles NumPageOffsets pFieldPageOffsets
 //AgtFieldStatisticsControl SetFieldPageOffset
 return nil
}

func(np *FieldStatisticsControl) GetFieldPageOffset ()(string, error) {
 //parameters: NumPorts pPortHandles
 //AgtFieldStatisticsControl GetFieldPageOffset
 return "", nil
}

func(np *FieldStatisticsControl) ListSelectedFields ()(string, error) {
 //parameters: PortHandle
 //AgtFieldStatisticsControl ListSelectedFields
 return "", nil
}

