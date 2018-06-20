package lib

import (
	"sync"
)

const (
	maxCommandLen = 200
)

//CommandList keep the parsed commands
type CommandList struct {
	commands []string
	lock     *sync.RWMutex
}

//NewCommandList ...
func NewCommandList() *CommandList {
	return &CommandList{
		commands: make([]string, 0),
		lock:     new(sync.RWMutex),
	}
}

//Log the executed command
func (cl *CommandList) Log(command string) {
	cl.lock.Lock()
	defer cl.lock.Unlock()

	if len(command) > 0 {
		if len(cl.commands) < maxCommandLen {
			cl.commands = append(cl.commands, command)
		} else {
			newList := make([]string, 0)
			newList = append(newList, cl.commands[1:]...)
			newList = append(newList, command)
			cl.commands = newList
		}
	}
}

//Commands return all logged commands
func (cl *CommandList) Commands() []string {
	cl.lock.RLock()
	defer cl.lock.RUnlock()

	commands := make([]string, 0)
	if len(cl.commands) > 0 {
		commands = append(commands, cl.commands...)
	}

	return commands
}
