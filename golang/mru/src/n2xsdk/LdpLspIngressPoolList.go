package n2xsdk

type LdpLspIngressPoolList struct {
 Handler string
}

func(np *LdpLspIngressPoolLis) AddLspPool () error {
 //parameters: SessionHandle FecType NumLsps PoolType
 //AgtLdpLspIngressPoolList AddLspPool
 return nil
}

func(np *LdpLspIngressPoolLis) RemoveLspPool () error {
 //parameters: SessionHandle PoolHandle
 //AgtLdpLspIngressPoolList RemoveLspPool
 return nil
}

func(np *LdpLspIngressPoolLis) ListLspPools ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpLspIngressPoolList ListLspPools
 return "", nil
}

func(np *LdpLspIngressPoolLis) ClearLspPools () error {
 //parameters: SessionHandle
 //AgtLdpLspIngressPoolList ClearLspPools
 return nil
}

func(np *LdpLspIngressPoolLis) EnablePoolList ()(string, error) {
 //parameters: SessionHandle PoolType
 //AgtLdpLspIngressPoolList EnablePoolList
 return "", nil
}

func(np *LdpLspIngressPoolLis) DisablePoolList ()(string, error) {
 //parameters: SessionHandle PoolType
 //AgtLdpLspIngressPoolList DisablePoolList
 return "", nil
}

func(np *LdpLspIngressPoolLis) IsEnabledPoolList ()(string, error) {
 //parameters: SessionHandle PoolType
 //AgtLdpLspIngressPoolList IsEnabledPoolList
 return "", nil
}

func(np *LdpLspIngressPoolLis) OpenAllLdpLspPools () error {
 //parameters: 
 //AgtLdpLspIngressPoolList OpenAllLdpLspPools
 return nil
}

func(np *LdpLspIngressPoolLis) CloseAllLdpLspPools () error {
 //parameters: 
 //AgtLdpLspIngressPoolList CloseAllLdpLspPools
 return nil
}

func(np *LdpLspIngressPoolLis) GetIncomingIngressPoolHandle ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpLspIngressPoolList GetIncomingIngressPoolHandle
 return "", nil
}

