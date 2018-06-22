package main

type Device interface {
	GetJsonPath() string
	GetName() string
	SetName(name string)
	RegisterCmd(sCmdName string, Bytes []byte)
	GetNumCommands() int
	DoCmd(sCmdName string)
	GetCommandsList() map[string][]byte
}
