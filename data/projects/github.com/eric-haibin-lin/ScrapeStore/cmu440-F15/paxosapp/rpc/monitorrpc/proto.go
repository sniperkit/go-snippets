// This file contains constants and arguments used to perform RPCs between
// Slave nodes and Master nodes.

package monitorrpc

type NodeType int

const (
	Master     NodeType = iota + 1 
	Slave                 
)

type HeartBeatArgs struct {
	Id int
	Type NodeType
}

type HeartBeatReply struct {
	//Nothing to reply 
}

