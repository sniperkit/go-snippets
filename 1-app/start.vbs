strComputer = "."
Set objWMIService = GetObject("winmgmts:\\" & strComputer & "\root\cimv2")
Set colProcessList=objWMIService.ExecQuery ("select * from Win32_Process where Name='1-app.exe' ")
For Each objProcess in colProcessList
    objProcess.Terminate()
Next

Set oShell = CreateObject("WScript.Shell")
'Æô¶¯bash
oShell.Run "1-app.exe", 0, false