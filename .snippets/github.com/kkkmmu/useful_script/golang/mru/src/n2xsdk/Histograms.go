package n2xsdk

type Histograms struct {
 Handler string
}

func(np *Histograms) IsHistogrammingSupported () error {
 //parameters: PortHandle
 //AgtHistograms IsHistogrammingSupported
 return nil
}

func(np *Histograms) GetMaximumNumberOfClasses ()(string, error) {
 //parameters: HistogramsHandle
 //AgtHistograms GetMaximumNumberOfClasses
 return "", nil
}

func(np *Histograms) GetMaximumNumberOfBuckets ()(string, error) {
 //parameters: HistogramsHandle
 //AgtHistograms GetMaximumNumberOfBuckets
 return "", nil
}

func(np *Histograms) AddPorts () error {
 //parameters: HistogramsHandle NumPorts pPortHandles
 //AgtHistograms AddPorts
 return nil
}

func(np *Histograms) RemovePorts () error {
 //parameters: HistogramsHandle NumPorts pPortHandles
 //AgtHistograms RemovePorts
 return nil
}

func(np *Histograms) ListPorts ()(string, error) {
 //parameters: HistogramsHandle
 //AgtHistograms ListPorts
 return "", nil
}

func(np *Histograms) ListAvailableClassFields ()(string, error) {
 //parameters: HistogramsHandle
 //AgtHistograms ListAvailableClassFields
 return "", nil
}

func(np *Histograms) SetClassField () error {
 //parameters: HistogramsHandle Field
 //AgtHistograms SetClassField
 return nil
}

func(np *Histograms) ClearClassField () error {
 //parameters: HistogramsHandle
 //AgtHistograms ClearClassField
 return nil
}

func(np *Histograms) GetClassField ()(string, error) {
 //parameters: HistogramsHandle
 //AgtHistograms GetClassField
 return "", nil
}

func(np *Histograms) AddClasses () error {
 //parameters: HistogramsHandle NumLowerBounds pLowerBounds NumUpperBounds pUpperBounds
 //AgtHistograms AddClasses
 return nil
}

func(np *Histograms) RemoveClasses () error {
 //parameters: HistogramsHandle NumClasses pClassHandles
 //AgtHistograms RemoveClasses
 return nil
}

func(np *Histograms) ListClasses ()(string, error) {
 //parameters: HistogramsHandle
 //AgtHistograms ListClasses
 return "", nil
}

func(np *Histograms) SetBucketBoundaries () error {
 //parameters: HistogramsHandle NumBucketBoundaries pBucketBoundaries
 //AgtHistograms SetBucketBoundaries
 return nil
}

func(np *Histograms) GetBucketBoundaries ()(string, error) {
 //parameters: HistogramsHandle
 //AgtHistograms GetBucketBoundaries
 return "", nil
}

func(np *Histograms) DeriveBucketBoundaries () error {
 //parameters: NumberOfBuckets MinimumBound MaximumBound
 //AgtHistograms DeriveBucketBoundaries
 return nil
}

func(np *Histograms) SetDerivedBucketBoundaries () error {
 //parameters: HistogramsHandle NumberOfBuckets MinimumBound MaximumBound
 //AgtHistograms SetDerivedBucketBoundaries
 return nil
}

func(np *Histograms) AreBucketBoundariesDerived () error {
 //parameters: HistogramsHandle
 //AgtHistograms AreBucketBoundariesDerived
 return nil
}

func(np *Histograms) GetHistograms ()(string, error) {
 //parameters: HistogramsHandle
 //AgtHistograms GetHistograms
 return "", nil
}

func(np *Histograms) GetStatistics ()(string, error) {
 //parameters: HistogramsHandle
 //AgtHistograms GetStatistics
 return "", nil
}

func(np *Histograms) GetMeasurementExtents ()(string, error) {
 //parameters: HistogramsHandle
 //AgtHistograms GetMeasurementExtents
 return "", nil
}

func(np *Histograms) ResetResults () error {
 //parameters: HistogramsHandle
 //AgtHistograms ResetResults
 return nil
}

