package n2xsdk

type AnalysisSource struct {
 Handler string
}

func(np *nalysisSource) SaveCaptureData () error {
 //parameters: PortHandle StartPacket EndPacket FilterHandle FileName Description
 //AgtAnalysisSource SaveCaptureData
 return nil
}

func(np *nalysisSource) IsCaptureSaveInProgress () error {
 //parameters: SourceHandle
 //AgtAnalysisSource IsCaptureSaveInProgress
 return nil
}

func(np *nalysisSource) WaitForCaptureSave () error {
 //parameters: SourceHandle
 //AgtAnalysisSource WaitForCaptureSave
 return nil
}

func(np *nalysisSource) OpenCaptureFile () error {
 //parameters: FileName
 //AgtAnalysisSource OpenCaptureFile
 return nil
}

func(np *nalysisSource) CloseCaptureFile () error {
 //parameters: SourceHandle
 //AgtAnalysisSource CloseCaptureFile
 return nil
}

func(np *nalysisSource) ListOpenCaptureFiles ()(string, error) {
 //parameters: 
 //AgtAnalysisSource ListOpenCaptureFiles
 return "", nil
}

func(np *nalysisSource) MergeCaptureData () error {
 //parameters: NumSources psaSourceHandles
 //AgtAnalysisSource MergeCaptureData
 return nil
}

func(np *nalysisSource) IsMergeCaptureDataInProgress () error {
 //parameters: MergeRequestID
 //AgtAnalysisSource IsMergeCaptureDataInProgress
 return nil
}

func(np *nalysisSource) WaitForMergeCaptureData () error {
 //parameters: MergeRequestID
 //AgtAnalysisSource WaitForMergeCaptureData
 return nil
}

func(np *nalysisSource) GetMergeCaptureSourceHandleByRequestID ()(string, error) {
 //parameters: MergeRequestID
 //AgtAnalysisSource GetMergeCaptureSourceHandleByRequestID
 return "", nil
}

func(np *nalysisSource) RemoveMergeCaptureSource () error {
 //parameters: SourceHandle
 //AgtAnalysisSource RemoveMergeCaptureSource
 return nil
}

func(np *nalysisSource) ListMergeCaptureSources ()(string, error) {
 //parameters: 
 //AgtAnalysisSource ListMergeCaptureSources
 return "", nil
}

func(np *nalysisSource) ListMergeCaptureSourceInputSources ()(string, error) {
 //parameters: SourceHandle
 //AgtAnalysisSource ListMergeCaptureSourceInputSources
 return "", nil
}

func(np *nalysisSource) GetMergeCaptureSourceTimes ()(string, error) {
 //parameters: SourceHandle
 //AgtAnalysisSource GetMergeCaptureSourceTimes
 return "", nil
}

func(np *nalysisSource) GetServiceRequestCountBySource ()(string, error) {
 //parameters: 
 //AgtAnalysisSource GetServiceRequestCountBySource
 return "", nil
}

func(np *nalysisSource) GetSourceType ()(string, error) {
 //parameters: SourceHandle
 //AgtAnalysisSource GetSourceType
 return "", nil
}

func(np *nalysisSource) WaitForSourceData () error {
 //parameters: SourceHandle
 //AgtAnalysisSource WaitForSourceData
 return nil
}

func(np *nalysisSource) GetSourceStats ()(string, error) {
 //parameters: SourceHandle
 //AgtAnalysisSource GetSourceStats
 return "", nil
}

func(np *nalysisSource) GetSourceTimes ()(string, error) {
 //parameters: SourceHandle
 //AgtAnalysisSource GetSourceTimes
 return "", nil
}

func(np *nalysisSource) GetCaptureFileTimes ()(string, error) {
 //parameters: SourceHandle
 //AgtAnalysisSource GetCaptureFileTimes
 return "", nil
}

func(np *nalysisSource) ListCapturePortInfo ()(string, error) {
 //parameters: SourceHandle
 //AgtAnalysisSource ListCapturePortInfo
 return "", nil
}

func(np *nalysisSource) RequestPacketsByRange () error {
 //parameters: SourceHandle RangeStart RangeEnd SearchBackwards FilterHandle MaxPackets MaxOctets
 //AgtAnalysisSource RequestPacketsByRange
 return nil
}

func(np *nalysisSource) RequestPacketsByTime () error {
 //parameters: NumSources psaSourceHandles StartTimestampSeconds StartTimestampNanoseconds FilterHandle MaxPackets MaxOctets
 //AgtAnalysisSource RequestPacketsByTime
 return nil
}

func(np *nalysisSource) RequestAnalysis () error {
 //parameters: SourceHandle FilterHandle Width
 //AgtAnalysisSource RequestAnalysis
 return nil
}

func(np *nalysisSource) CancelAnalysis () error {
 //parameters: SourceHandle
 //AgtAnalysisSource CancelAnalysis
 return nil
}

