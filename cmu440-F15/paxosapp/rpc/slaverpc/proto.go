// This file contains constants and arguments used to perform RPCs between
// Slave nodes and Master nodes.

package slaverpc

type GetArgs struct {
	Key string
}

type GetReply struct {
	Value  []string
	Status int
}

type AppendArgs struct {
	Key   string
	Value []string
}

type AppendReply struct {
	//Nothing to reply here
}

type GetRankArgs struct {
	Key string
}

type GetRankReply struct {
	Value float64
}

type PutRankArgs struct {
	Key   string
	Value float64
}

type PutRankReply struct {
	//nothing here
}
