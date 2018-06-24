// This file contains constants and arguments used to perform RPCs between
// two Paxos nodes. DO NOT MODIFY!

package paxosrpc

// Status represents the status of a RPC's reply.
type Status int
type Lookup int

const (
	NumCopies = 5
)

const (
	OK     Status = iota + 1 // Paxos replied OK
	Reject                   // Paxos rejected the message
	OtherValueCommitted
)

const (
	KeyFound    Lookup = iota + 1 // GetValue key found
	KeyNotFound                   // GetValue key not found
)

type ProposalNumberArgs struct {
	Key string
}

type ProposalNumberReply struct {
	N int
}

type ProposeArgs struct {
	N   int // Proposal number
	Key string
	V   [NumCopies]int
}

type ProposeReply struct {
	Status Status
	V      [NumCopies]int
}

type PrepareArgs struct {
	Key string
	N   int
}

type PrepareReply struct {
	Status Status
	N_a    int            // Highest proposal number accepted
	V_a    [NumCopies]int // Corresponding value
}

type AcceptArgs struct {
	Key string
	N   int
	V   [NumCopies]int
}

type AcceptReply struct {
	Status Status
}

type CommitArgs struct {
	Key string
	V   [NumCopies]int
}

type CommitReply struct {
	// No content, no reply necessary
}

type ReplaceServerArgs struct {
	SrvID    int // Server being replaced
	Hostport string
}

type ReplaceServerReply struct {
	// No content necessary
}

type ReplaceCatchupArgs struct {
	// No content necessary
}

type ReplaceCatchupReply struct {
	Data []byte
}

type GetAllLinksArgs struct {
	// No content necessary, just get the whole damn thing lol
}

type GetAllLinksReply struct {
	LinksMap map[string][]string
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
	//nothing to be returned
}

type GetLinksArgs struct {
	Key string
}

type GetLinksReply struct {
	Value []string
}

type AppendArgs struct {
	Key   string
	Value []string
}

type AppendReply struct {
	Status Status
}
