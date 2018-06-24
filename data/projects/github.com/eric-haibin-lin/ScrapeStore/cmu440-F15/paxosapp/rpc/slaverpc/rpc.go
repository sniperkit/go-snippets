package slaverpc

type RemoteSlaveNode interface {
	Append(args *AppendArgs, reply *AppendReply) error
	Get(args *GetArgs, reply *GetReply) error
	PutRank(args *PutRankArgs, reply *PutRankReply) error
	GetRank(args *GetRankArgs, reply *GetRankReply) error
}

type SlaveNode struct {
	// Embed all methods into the struct. See the Effective Go section about
	// embedding for more details: golang.org/doc/effective_go.html#embedding
	RemoteSlaveNode
}

// Wrap wraps t in a type-safe wrapper struct to ensure that only the desired
// methods are exported to receive RPCs.
func Wrap(t RemoteSlaveNode) RemoteSlaveNode {
	return &SlaveNode{t}
}
