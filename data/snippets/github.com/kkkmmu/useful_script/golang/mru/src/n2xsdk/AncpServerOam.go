package n2xsdk

type AncpServerOam struct {
 Handler string
}

func(np *ncpServerOam) SetLoopbackCount () error {
 //parameters: SessionPoolHandle LoopbackCount
 //AgtAncpServerOam SetLoopbackCount
 return nil
}

func(np *ncpServerOam) GetLoopbackCount ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtAncpServerOam GetLoopbackCount
 return "", nil
}

func(np *ncpServerOam) SetLoopbackTimeout () error {
 //parameters: SessionPoolHandle LoopbackTimeout
 //AgtAncpServerOam SetLoopbackTimeout
 return nil
}

func(np *ncpServerOam) GetLoopbackTimeout ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtAncpServerOam GetLoopbackTimeout
 return "", nil
}

func(np *ncpServerOam) SetOpaqueData () error {
 //parameters: SessionPoolHandle FirstToken SecondToken
 //AgtAncpServerOam SetOpaqueData
 return nil
}

func(np *ncpServerOam) GetOpaqueData ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtAncpServerOam GetOpaqueData
 return "", nil
}

func(np *ncpServerOam) AddTlv () error {
 //parameters: SessionPoolHandle TlvType TlvLength TlvValue
 //AgtAncpServerOam AddTlv
 return nil
}

func(np *ncpServerOam) ListTlvs ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtAncpServerOam ListTlvs
 return "", nil
}

func(np *ncpServerOam) RemoveTlv () error {
 //parameters: SessionPoolHandle TlvType
 //AgtAncpServerOam RemoveTlv
 return nil
}

func(np *ncpServerOam) RemoveAllTlvs () error {
 //parameters: SessionPoolHandle
 //AgtAncpServerOam RemoveAllTlvs
 return nil
}

func(np *ncpServerOam) SetTlv () error {
 //parameters: SessionPoolHandle TlvType TlvLength TlvValue
 //AgtAncpServerOam SetTlv
 return nil
}

func(np *ncpServerOam) GetTlv ()(string, error) {
 //parameters: SessionPoolHandle TlvType
 //AgtAncpServerOam GetTlv
 return "", nil
}

func(np *ncpServerOam) AddSubTlv () error {
 //parameters: SessionPoolHandle TlvType SubTlvType SubTlvLength SubTlvValue
 //AgtAncpServerOam AddSubTlv
 return nil
}

func(np *ncpServerOam) ListSubTlvs ()(string, error) {
 //parameters: SessionPoolHandle TlvType
 //AgtAncpServerOam ListSubTlvs
 return "", nil
}

func(np *ncpServerOam) RemoveSubTlv () error {
 //parameters: SessionPoolHandle TlvType SubTlvType
 //AgtAncpServerOam RemoveSubTlv
 return nil
}

func(np *ncpServerOam) RemoveAllSubTlvs () error {
 //parameters: SessionPoolHandle TlvType
 //AgtAncpServerOam RemoveAllSubTlvs
 return nil
}

func(np *ncpServerOam) SetSubTlv () error {
 //parameters: SessionPoolHandle TlvType SubTlvType SubTlvLength SubTlvValue
 //AgtAncpServerOam SetSubTlv
 return nil
}

func(np *ncpServerOam) GetSubTlv ()(string, error) {
 //parameters: SessionPoolHandle TlvType SubTlvType
 //AgtAncpServerOam GetSubTlv
 return "", nil
}

