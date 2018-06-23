package n2xsdk

type Bgp4UpdateMessageCaptureFilter struct {
 Handler string
}

func(np *Bgp4UpdateMessageCaptureFilter) Enable () error {
 //parameters: UpdateMessageCaptureFilterHandle
 //AgtBgp4UpdateMessageCaptureFilter Enable
 return nil
}

func(np *Bgp4UpdateMessageCaptureFilter) Disable () error {
 //parameters: UpdateMessageCaptureFilterHandle
 //AgtBgp4UpdateMessageCaptureFilter Disable
 return nil
}

func(np *Bgp4UpdateMessageCaptureFilter) Reset () error {
 //parameters: UpdateMessageCaptureFilterHandle
 //AgtBgp4UpdateMessageCaptureFilter Reset
 return nil
}

