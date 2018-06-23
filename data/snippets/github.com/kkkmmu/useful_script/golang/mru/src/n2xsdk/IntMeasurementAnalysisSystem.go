package n2xsdk

type IntMeasurementAnalysisSystem struct {
 Handler string
}

func(np *IntMeasurementAnalysisSystem) DataSourceAvailable () error {
 //parameters: ViewHandle
 //IntMeasurementAnalysisSystem DataSourceAvailable
 return nil
}

func(np *IntMeasurementAnalysisSystem) GetHelpDirectory ()(string, error) {
 //parameters: 
 //IntMeasurementAnalysisSystem GetHelpDirectory
 return "", nil
}

func(np *IntMeasurementAnalysisSystem) FindKeyByBeginWith () error {
 //parameters: ViewHandle columnName value
 //IntMeasurementAnalysisSystem FindKeyByBeginWith
 return nil
}

func(np *IntMeasurementAnalysisSystem) FindKeyByValue () error {
 //parameters: ViewHandle columnName value
 //IntMeasurementAnalysisSystem FindKeyByValue
 return nil
}

func(np *IntMeasurementAnalysisSystem) GetRowIndicesByKey ()(string, error) {
 //parameters: ViewHandle key count
 //IntMeasurementAnalysisSystem GetRowIndicesByKey
 return "", nil
}

func(np *IntMeasurementAnalysisSystem) GetRowKeys ()(string, error) {
 //parameters: ViewHandle index count
 //IntMeasurementAnalysisSystem GetRowKeys
 return "", nil
}

func(np *IntMeasurementAnalysisSystem) GetUniqueColumnValues ()(string, error) {
 //parameters: ViewHandle columnName maxCount includeFilteredOut roundDataTime
 //IntMeasurementAnalysisSystem GetUniqueColumnValues
 return "", nil
}

func(np *IntMeasurementAnalysisSystem) GetRows ()(string, error) {
 //parameters: ViewHandle index count
 //IntMeasurementAnalysisSystem GetRows
 return "", nil
}

func(np *IntMeasurementAnalysisSystem) GetTableProperties ()(string, error) {
 //parameters: ViewHandle
 //IntMeasurementAnalysisSystem GetTableProperties
 return "", nil
}

func(np *IntMeasurementAnalysisSystem) SetSortOrder () error {
 //parameters: ViewHandle sortOrder
 //IntMeasurementAnalysisSystem SetSortOrder
 return nil
}

func(np *IntMeasurementAnalysisSystem) GetSortOrder ()(string, error) {
 //parameters: ViewHandle
 //IntMeasurementAnalysisSystem GetSortOrder
 return "", nil
}

func(np *IntMeasurementAnalysisSystem) GetNamedFilters ()(string, error) {
 //parameters: 
 //IntMeasurementAnalysisSystem GetNamedFilters
 return "", nil
}

func(np *IntMeasurementAnalysisSystem) AddNamedFilter () error {
 //parameters: name NamedFilterBytes
 //IntMeasurementAnalysisSystem AddNamedFilter
 return nil
}

func(np *IntMeasurementAnalysisSystem) Subscribe () error {
 //parameters: eventSink
 //IntMeasurementAnalysisSystem Subscribe
 return nil
}

func(np *IntMeasurementAnalysisSystem) Unsubscribe () error {
 //parameters: eventSink
 //IntMeasurementAnalysisSystem Unsubscribe
 return nil
}

