package n2xsdk

type StreamPatternMatcher struct {
 Handler string
}

func(np *StreamPatternMatcher) ListAvailableMatchers ()(string, error) {
 //parameters: PortHandle
 //AgtStreamPatternMatcher ListAvailableMatchers
 return "", nil
}

func(np *StreamPatternMatcher) ListUsedMatchers ()(string, error) {
 //parameters: PortHandle
 //AgtStreamPatternMatcher ListUsedMatchers
 return "", nil
}

func(np *StreamPatternMatcher) ListUnusedMatchers ()(string, error) {
 //parameters: PortHandle
 //AgtStreamPatternMatcher ListUnusedMatchers
 return "", nil
}

func(np *StreamPatternMatcher) ResetMatcherField () error {
 //parameters: PortHandle Matcher
 //AgtStreamPatternMatcher ResetMatcherField
 return nil
}

func(np *StreamPatternMatcher) SetMatcherField () error {
 //parameters: PortHandle Matcher Field
 //AgtStreamPatternMatcher SetMatcherField
 return nil
}

func(np *StreamPatternMatcher) GetMatcherField ()(string, error) {
 //parameters: PortHandle Matcher
 //AgtStreamPatternMatcher GetMatcherField
 return "", nil
}

func(np *StreamPatternMatcher) ListAvailableMatcherFields ()(string, error) {
 //parameters: PortHandle
 //AgtStreamPatternMatcher ListAvailableMatcherFields
 return "", nil
}

func(np *StreamPatternMatcher) SetMatcherFieldByOffset () error {
 //parameters: PortHandle Matcher OffsetType ByteOffset BitMask
 //AgtStreamPatternMatcher SetMatcherFieldByOffset
 return nil
}

func(np *StreamPatternMatcher) GetMatcherFieldByOffset ()(string, error) {
 //parameters: PortHandle Matcher
 //AgtStreamPatternMatcher GetMatcherFieldByOffset
 return "", nil
}

func(np *StreamPatternMatcher) SetMatcherValues () error {
 //parameters: PortHandle StreamGroupHandle Matcher ValueType Count psaValues
 //AgtStreamPatternMatcher SetMatcherValues
 return nil
}

func(np *StreamPatternMatcher) GetMatcherValues ()(string, error) {
 //parameters: PortHandle StreamGroupHandle Matcher
 //AgtStreamPatternMatcher GetMatcherValues
 return "", nil
}

func(np *StreamPatternMatcher) ClearMatcherValues () error {
 //parameters: PortHandle StreamGroupHandle Matcher
 //AgtStreamPatternMatcher ClearMatcherValues
 return nil
}

func(np *StreamPatternMatcher) ListMatcherStreamGroups ()(string, error) {
 //parameters: PortHandle Matcher
 //AgtStreamPatternMatcher ListMatcherStreamGroups
 return "", nil
}

