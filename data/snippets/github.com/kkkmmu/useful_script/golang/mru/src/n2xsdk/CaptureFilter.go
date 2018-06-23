package n2xsdk

type CaptureFilter struct {
 Handler string
}

func(np *CaptureFilter) ClearAllFilters () error {
 //parameters: PortHandle
 //AgtCaptureFilter ClearAllFilters
 return nil
}

func(np *CaptureFilter) AddStatusFilter () error {
 //parameters: PortHandle StatusFilter FilterAction
 //AgtCaptureFilter AddStatusFilter
 return nil
}

func(np *CaptureFilter) RemoveStatusFilter () error {
 //parameters: PortHandle StatusFilter
 //AgtCaptureFilter RemoveStatusFilter
 return nil
}

func(np *CaptureFilter) ListStatusFilters ()(string, error) {
 //parameters: PortHandle
 //AgtCaptureFilter ListStatusFilters
 return "", nil
}

func(np *CaptureFilter) AddFrameMatcherFilters () error {
 //parameters: PortHandle FrameMatcherHandleCount psaFrameMatcherHandles FilterAction
 //AgtCaptureFilter AddFrameMatcherFilters
 return nil
}

func(np *CaptureFilter) RemoveFrameMatcherFilters () error {
 //parameters: PortHandle FrameMatcherHandleCount psaFrameMatcherHandles
 //AgtCaptureFilter RemoveFrameMatcherFilters
 return nil
}

func(np *CaptureFilter) ListFrameMatcherFilters ()(string, error) {
 //parameters: PortHandle
 //AgtCaptureFilter ListFrameMatcherFilters
 return "", nil
}

