package monitorrpc

type RemoteMonitorNode interface {
	HeartBeat(args *HeartBeatArgs, reply *HeartBeatReply) error
}

type MonitorNode struct {
	// Embed all methods into the struct. See the Effective Go section about
	// embedding for more details: golang.org/doc/effective_go.html#embedding
	RemoteMonitorNode
}

// Wrap wraps t in a type-safe wrapper struct to ensure that only the desired
// methods are exported to receive RPCs.
func Wrap(t RemoteMonitorNode) RemoteMonitorNode {
	return &MonitorNode{t}
}