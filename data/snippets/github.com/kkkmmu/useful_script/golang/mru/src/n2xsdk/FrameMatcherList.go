package n2xsdk

type FrameMatcherList struct {
 Handler string
}

func(np *FrameMatcherLis) AddFrameMatcher () error {
 //parameters: PortHandle
 //AgtFrameMatcherList AddFrameMatcher
 return nil
}

func(np *FrameMatcherLis) RemoveFrameMatcher () error {
 //parameters: PortHandle pFrameMatcherHandle
 //AgtFrameMatcherList RemoveFrameMatcher
 return nil
}

func(np *FrameMatcherLis) ListFrameMatchers ()(string, error) {
 //parameters: PortHandle
 //AgtFrameMatcherList ListFrameMatchers
 return "", nil
}

func(np *FrameMatcherLis) GetFreeFrameMatcherCount ()(string, error) {
 //parameters: PortHandle
 //AgtFrameMatcherList GetFreeFrameMatcherCount
 return "", nil
}

func(np *FrameMatcherLis) SetName () error {
 //parameters: FrameMatcherHandle Name
 //AgtFrameMatcherList SetName
 return nil
}

func(np *FrameMatcherLis) GetName ()(string, error) {
 //parameters: FrameMatcherHandle
 //AgtFrameMatcherList GetName
 return "", nil
}

