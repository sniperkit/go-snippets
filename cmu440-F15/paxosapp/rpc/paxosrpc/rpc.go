package paxosrpc

// STAFF USE ONLY! Students should not use this interface in their code.
type RemotePaxosNode interface {
	// Propose initializes proposing a value for a key, and replies with the
	// value that was committed for that key. Propose should not return until
	// a value has been committed, or 15 seconds have passed.
	Propose(args *ProposeArgs, reply *ProposeReply) error

	// GetValue looks up the value for a key, and replies with the value or with
	// the Status KeyNotFound.
	//GetValue(args *GetValueArgs, reply *GetValueReply) error

	// GetNextProposalNumber generates a proposal number which will be passed to
	// Propose. Proposal numbers should not repeat, and for a particular node
	// they should be strictly increasing.
	GetNextProposalNumber(args *ProposalNumberArgs, reply *ProposalNumberReply) error

	// Receive a Prepare message from another Paxos Node. The message contains
	// the key whose value is being proposed by the node sending the prepare
	// message. This function should respond with Status OK if the prepare is
	// accepted and Reject otherwise.
	RecvPrepare(args *PrepareArgs, reply *PrepareReply) error

	// Receive an Accept message from another Paxos Node. The message contains
	// the key whose value is being proposed by the node sending the accept
	// message. This function should respond with Status OK if the prepare is
	// accepted and Reject otherwise.
	RecvAccept(args *AcceptArgs, reply *AcceptReply) error

	// Receive a Commit message from another Paxos Node. The message contains
	// the key whose value was proposed by the node sending the commit
	// message. This function should respond with Status OK if the prepare is
	// accepted and Reject otherwise.
	RecvCommit(args *CommitArgs, reply *CommitReply) error

	// Notify another node of a replacement server which has started up. The
	// message contains the Server ID of the node being replaced, and the
	// hostport of the replacement node
	RecvReplaceServer(args *ReplaceServerArgs, reply *ReplaceServerReply) error

	// Request the value that was agreed upon for a particular round. A node
	// receiving this message should reply with the data (as an array of bytes)
	// needed to make the replacement server aware of the keys and values
	// committed so far.
	RecvReplaceCatchup(args *ReplaceCatchupArgs, reply *ReplaceCatchupReply) error

	// Get all crawled data from the database
	GetAllLinks(args *GetAllLinksArgs, reply *GetAllLinksReply) error

	// Get links for a specific URL
	GetLinks(args *GetLinksArgs, reply *GetLinksReply) error

	// Append some links
	Append(args *AppendArgs, reply *AppendReply) error

	// Get Rank
	GetRank(args *GetRankArgs, reply *GetRankReply) error

	// Put Rank
	PutRank(args *PutRankArgs, reply *PutRankReply) error
}

type PaxosNode struct {
	// Embed all methods into the struct. See the Effective Go section about
	// embedding for more details: golang.org/doc/effective_go.html#embedding
	RemotePaxosNode
}

// Wrap wraps t in a type-safe wrapper struct to ensure that only the desired
// methods are exported to receive RPCs.
func Wrap(t RemotePaxosNode) RemotePaxosNode {
	return &PaxosNode{t}
}
