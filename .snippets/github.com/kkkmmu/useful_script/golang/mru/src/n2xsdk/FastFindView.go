package n2xsdk

type FastFindView struct {
 Handler string
}

func(np *FastFindView) GetFastFindDataSet ()(string, error) {
 //parameters: ViewHandle
 //AgtFastFindView GetFastFindDataSet
 return "", nil
}

func(np *FastFindView) SetFastFindDataSet () error {
 //parameters: ViewHandle DataSetHandle
 //AgtFastFindView SetFastFindDataSet
 return nil
}

func(np *FastFindView) GetMeasurementType ()(string, error) {
 //parameters: ViewHandle
 //AgtFastFindView GetMeasurementType
 return "", nil
}

func(np *FastFindView) SetMeasurementType () error {
 //parameters: ViewHandle MeasurementType
 //AgtFastFindView SetMeasurementType
 return nil
}

func(np *FastFindView) IsDataSourceValid () error {
 //parameters: ViewHandle
 //AgtFastFindView IsDataSourceValid
 return nil
}

func(np *FastFindView) GetColumns ()(string, error) {
 //parameters: ViewHandle
 //AgtFastFindView GetColumns
 return "", nil
}

func(np *FastFindView) GetSortOrder ()(string, error) {
 //parameters: ViewHandle
 //AgtFastFindView GetSortOrder
 return "", nil
}

func(np *FastFindView) SetSortOrder () error {
 //parameters: ViewHandle Columns Directions
 //AgtFastFindView SetSortOrder
 return nil
}

func(np *FastFindView) ClearSort () error {
 //parameters: ViewHandle
 //AgtFastFindView ClearSort
 return nil
}

func(np *FastFindView) GetFilter ()(string, error) {
 //parameters: ViewHandle
 //AgtFastFindView GetFilter
 return "", nil
}

func(np *FastFindView) SetFilter () error {
 //parameters: ViewHandle Filter
 //AgtFastFindView SetFilter
 return nil
}

func(np *FastFindView) ClearFilter () error {
 //parameters: ViewHandle
 //AgtFastFindView ClearFilter
 return nil
}

func(np *FastFindView) GetUnfilteredNumberOfRows ()(string, error) {
 //parameters: ViewHandle
 //AgtFastFindView GetUnfilteredNumberOfRows
 return "", nil
}

func(np *FastFindView) GetNumberOfRows ()(string, error) {
 //parameters: ViewHandle
 //AgtFastFindView GetNumberOfRows
 return "", nil
}

func(np *FastFindView) GetRows ()(string, error) {
 //parameters: ViewHandle firstRowNumber numberOfRows
 //AgtFastFindView GetRows
 return "", nil
}

func(np *FastFindView) Export () error {
 //parameters: ViewHandle Format Filename
 //AgtFastFindView Export
 return nil
}

func(np *FastFindView) GetDefaultDirectory ()(string, error) {
 //parameters: 
 //AgtFastFindView GetDefaultDirectory
 return "", nil
}

