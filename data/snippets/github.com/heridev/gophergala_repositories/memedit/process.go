package main

type Process struct {
	PROCESSENTRY PROCESSENTRY32
	ModuleList []MODULEENTRY32
}

func (p *Process) getPid() uint32 {
	return p.PROCESSENTRY.ProcessID
}

func (p *Process) getName() string {

	return p.PROCESSENTRY.getName()
}