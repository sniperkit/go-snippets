package n2xsdk

type AnalysisSetSelection struct {
 Handler string
}

func(np *nalysisSetSelection) SelectAnalysisSets () error {
 //parameters: NumSets psaAnalysisSetHandles
 //AgtAnalysisSetSelection SelectAnalysisSets
 return nil
}

func(np *nalysisSetSelection) ListSelectedAnalysisSets ()(string, error) {
 //parameters: 
 //AgtAnalysisSetSelection ListSelectedAnalysisSets
 return "", nil
}

func(np *nalysisSetSelection) SetPercentiles () error {
 //parameters: NumPercentiles pPercentiles
 //AgtAnalysisSetSelection SetPercentiles
 return nil
}

func(np *nalysisSetSelection) RequestDistribution () error {
 //parameters: NumBuckets
 //AgtAnalysisSetSelection RequestDistribution
 return nil
}

func(np *nalysisSetSelection) CancelDistribution () error {
 //parameters: 
 //AgtAnalysisSetSelection CancelDistribution
 return nil
}

func(np *nalysisSetSelection) WaitForDistribution () error {
 //parameters: 
 //AgtAnalysisSetSelection WaitForDistribution
 return nil
}

func(np *nalysisSetSelection) WaitForExpansion () error {
 //parameters: 
 //AgtAnalysisSetSelection WaitForExpansion
 return nil
}

func(np *nalysisSetSelection) WaitForScale () error {
 //parameters: 
 //AgtAnalysisSetSelection WaitForScale
 return nil
}

func(np *nalysisSetSelection) GetScaleStatus ()(string, error) {
 //parameters: 
 //AgtAnalysisSetSelection GetScaleStatus
 return "", nil
}

func(np *nalysisSetSelection) GetExpansionStatus ()(string, error) {
 //parameters: 
 //AgtAnalysisSetSelection GetExpansionStatus
 return "", nil
}

func(np *nalysisSetSelection) IncreaseScale () error {
 //parameters: CenterTimeSeconds CenterTimeNanoseconds RelativeScaleFactor Width
 //AgtAnalysisSetSelection IncreaseScale
 return nil
}

func(np *nalysisSetSelection) DecreaseScale () error {
 //parameters: NumLevels
 //AgtAnalysisSetSelection DecreaseScale
 return nil
}

func(np *nalysisSetSelection) ExpandDistribution () error {
 //parameters: BucketNumber
 //AgtAnalysisSetSelection ExpandDistribution
 return nil
}

func(np *nalysisSetSelection) ContractDistribution () error {
 //parameters: NumLevels
 //AgtAnalysisSetSelection ContractDistribution
 return nil
}

