package n2xsdk

type LspPingSession struct {
 Handler string
}

func(np *LspPingSession) Start () error {
 //parameters: SessionHandle
 //AgtLspPingSession Start
 return nil
}

func(np *LspPingSession) Reset () error {
 //parameters: SessionHandle
 //AgtLspPingSession Reset
 return nil
}

func(np *LspPingSession) Renew () error {
 //parameters: SessionHandle
 //AgtLspPingSession Renew
 return nil
}

