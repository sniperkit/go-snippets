package n2xsdk

type ControlModule struct {
 Handler string
}

func(np *ControlModule) RunScript () error {
 //parameters: filename
 //AgtControlModule RunScript
 return nil
}

func(np *ControlModule) SetLogLevel () error {
 //parameters: LogLevel
 //AgtControlModule SetLogLevel
 return nil
}

func(np *ControlModule) GetLogLevel ()(string, error) {
 //parameters: 
 //AgtControlModule GetLogLevel
 return "", nil
}

func(np *ControlModule) SetControllerLogFileName () error {
 //parameters: LogName
 //AgtControlModule SetControllerLogFileName
 return nil
}

func(np *ControlModule) SetGuiPort () error {
 //parameters: Port
 //AgtControlModule SetGuiPort
 return nil
}

func(np *ControlModule) SetGuiPortIpAddress () error {
 //parameters: Port IpAddress
 //AgtControlModule SetGuiPortIpAddress
 return nil
}

func(np *ControlModule) EnableSlaveLogging () error {
 //parameters: 
 //AgtControlModule EnableSlaveLogging
 return nil
}

func(np *ControlModule) DisableSlaveLogging () error {
 //parameters: 
 //AgtControlModule DisableSlaveLogging
 return nil
}

func(np *ControlModule) StopScript () error {
 //parameters: 
 //AgtControlModule StopScript
 return nil
}

func(np *ControlModule) GetVersion ()(string, error) {
 //parameters: 
 //AgtControlModule GetVersion
 return "", nil
}

func(np *ControlModule) GetScriptState ()(string, error) {
 //parameters: 
 //AgtControlModule GetScriptState
 return "", nil
}

func(np *ControlModule) IsSlaveLoggingEnabled () error {
 //parameters: 
 //AgtControlModule IsSlaveLoggingEnabled
 return nil
}

func(np *ControlModule) SetGlobalScriptVariable () error {
 //parameters: VarName Value
 //AgtControlModule SetGlobalScriptVariable
 return nil
}

func(np *ControlModule) SetLocalScriptVariable () error {
 //parameters: ssmId VarName Value
 //AgtControlModule SetLocalScriptVariable
 return nil
}

func(np *ControlModule) RunScriptProcedure () error {
 //parameters: ssmId ProcName
 //AgtControlModule RunScriptProcedure
 return nil
}

func(np *ControlModule) CompileScript () error {
 //parameters: ScriptName
 //AgtControlModule CompileScript
 return nil
}

func(np *ControlModule) CompileSuite () error {
 //parameters: ScriptName
 //AgtControlModule CompileSuite
 return nil
}

func(np *ControlModule) GetGlobalVariable ()(string, error) {
 //parameters: VarName
 //AgtControlModule GetGlobalVariable
 return "", nil
}

func(np *ControlModule) GetGlobalArrayMember ()(string, error) {
 //parameters: VarName Index
 //AgtControlModule GetGlobalArrayMember
 return "", nil
}

func(np *ControlModule) SetGlobalVariable () error {
 //parameters: VarName Value
 //AgtControlModule SetGlobalVariable
 return nil
}

func(np *ControlModule) SetGlobalArrayMember () error {
 //parameters: VarName Index Value
 //AgtControlModule SetGlobalArrayMember
 return nil
}

func(np *ControlModule) SetLocalVariable () error {
 //parameters: ssmId VarName Value
 //AgtControlModule SetLocalVariable
 return nil
}

func(np *ControlModule) SetLocalArrayMember () error {
 //parameters: ssmId VarName Index Value
 //AgtControlModule SetLocalArrayMember
 return nil
}

func(np *ControlModule) HaltScriptProcedure () error {
 //parameters: ssmId
 //AgtControlModule HaltScriptProcedure
 return nil
}

func(np *ControlModule) SetSsmStatsUpdateInterval () error {
 //parameters: ssmId Interval
 //AgtControlModule SetSsmStatsUpdateInterval
 return nil
}

func(np *ControlModule) ContinueFromSyncPoint () error {
 //parameters: 
 //AgtControlModule ContinueFromSyncPoint
 return nil
}

func(np *ControlModule) OpenRemoteGuiFile () error {
 //parameters: RemoteFile Session
 //AgtControlModule OpenRemoteGuiFile
 return nil
}

func(np *ControlModule) WriteRemoteGuiFile () error {
 //parameters: Line
 //AgtControlModule WriteRemoteGuiFile
 return nil
}

func(np *ControlModule) CloseRemoteGuiFile () error {
 //parameters: 
 //AgtControlModule CloseRemoteGuiFile
 return nil
}

func(np *ControlModule) RunRemoteExpectScript () error {
 //parameters: Session
 //AgtControlModule RunRemoteExpectScript
 return nil
}

func(np *ControlModule) GetMctsGuiState ()(string, error) {
 //parameters: 
 //AgtControlModule GetMctsGuiState
 return "", nil
}

func(np *ControlModule) SetMctsGuiRunning () error {
 //parameters: Port IpAddress
 //AgtControlModule SetMctsGuiRunning
 return nil
}

func(np *ControlModule) UnSetMctsGuiRunning () error {
 //parameters: 
 //AgtControlModule UnSetMctsGuiRunning
 return nil
}

