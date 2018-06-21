package n2xsdk

type LdpLspEgressPoolList struct {
 Handler string
}

func(np *LdpLspEgressPoolLis) AddLspPool () error {
 //parameters: SessionHandle FecType NumLsps PoolType
 //AgtLdpLspEgressPoolList AddLspPool
 return nil
}

func(np *LdpLspEgressPoolLis) RemoveLspPool () error {
 //parameters: SessionHandle PoolHandle
 //AgtLdpLspEgressPoolList RemoveLspPool
 return nil
}

func(np *LdpLspEgressPoolLis) ListLspPools ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpLspEgressPoolList ListLspPools
 return "", nil
}

func(np *LdpLspEgressPoolLis) ClearLspPools () error {
 //parameters: SessionHandle
 //AgtLdpLspEgressPoolList ClearLspPools
 return nil
}

func(np *LdpLspEgressPoolLis) EnablePoolList ()(string, error) {
 //parameters: SessionHandle PoolType
 //AgtLdpLspEgressPoolList EnablePoolList
 return "", nil
}

func(np *LdpLspEgressPoolLis) DisablePoolList ()(string, error) {
 //parameters: SessionHandle PoolType
 //AgtLdpLspEgressPoolList DisablePoolList
 return "", nil
}

func(np *LdpLspEgressPoolLis) IsEnabledPoolList ()(string, error) {
 //parameters: SessionHandle PoolType
 //AgtLdpLspEgressPoolList IsEnabledPoolList
 return "", nil
}

func(np *LdpLspEgressPoolLis) OpenAllLdpLspPools () error {
 //parameters: 
 //AgtLdpLspEgressPoolList OpenAllLdpLspPools
 return nil
}

func(np *LdpLspEgressPoolLis) CloseAllLdpLspPools () error {
 //parameters: 
 //AgtLdpLspEgressPoolList CloseAllLdpLspPools
 return nil
}

func(np *LdpLspEgressPoolLis) GetIncomingEgressPoolHandle ()(string, error) {
 //parameters: SessionHandle
 //AgtLdpLspEgressPoolList GetIncomingEgressPoolHandle
 return "", nil
}

