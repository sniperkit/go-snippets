package n2xsdk

type AnalysisFilter struct {
 Handler string
}

func(np *nalysisFilter) GetType ()(string, error) {
 //parameters: Handle
 //AgtAnalysisFilter GetType
 return "", nil
}

func(np *nalysisFilter) GetName ()(string, error) {
 //parameters: Handle
 //AgtAnalysisFilter GetName
 return "", nil
}

func(np *nalysisFilter) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtAnalysisFilter GetLockCount
 return "", nil
}

func(np *nalysisFilter) ClearFilter () error {
 //parameters: FilterHandle
 //AgtAnalysisFilter ClearFilter
 return nil
}

func(np *nalysisFilter) AddPatternCriterion () error {
 //parameters: FilterHandle PatternHandle
 //AgtAnalysisFilter AddPatternCriterion
 return nil
}

func(np *nalysisFilter) RemovePatternCriterion () error {
 //parameters: FilterHandle PatternHandle
 //AgtAnalysisFilter RemovePatternCriterion
 return nil
}

func(np *nalysisFilter) AddValueCriterion () error {
 //parameters: FilterHandle ValueCriterion Value
 //AgtAnalysisFilter AddValueCriterion
 return nil
}

func(np *nalysisFilter) RemoveValueCriterion () error {
 //parameters: FilterHandle ValueCriterion
 //AgtAnalysisFilter RemoveValueCriterion
 return nil
}

func(np *nalysisFilter) AddStatusCriterion () error {
 //parameters: FilterHandle StatusCriterion IsSet
 //AgtAnalysisFilter AddStatusCriterion
 return nil
}

func(np *nalysisFilter) RemoveStatusCriterion () error {
 //parameters: FilterHandle StatusCriterion
 //AgtAnalysisFilter RemoveStatusCriterion
 return nil
}

func(np *nalysisFilter) ListPatternCriteria ()(string, error) {
 //parameters: FilterHandle
 //AgtAnalysisFilter ListPatternCriteria
 return "", nil
}

func(np *nalysisFilter) ListValueCriteria ()(string, error) {
 //parameters: FilterHandle
 //AgtAnalysisFilter ListValueCriteria
 return "", nil
}

func(np *nalysisFilter) ListStatusCriteria ()(string, error) {
 //parameters: FilterHandle
 //AgtAnalysisFilter ListStatusCriteria
 return "", nil
}

