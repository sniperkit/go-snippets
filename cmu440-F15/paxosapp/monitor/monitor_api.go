package monitor

import (
	"github.com/cmu440-F15/paxosapp/rpc/monitorrpc"
)

type MonitorNode interface {
	// HeartBeat RPC is for master nodes to contact monitor server. The monitor server 
	// mark down the master's id for liveness check 
	HeartBeat(args *monitorrpc.HeartBeatArgs, reply *monitorrpc.HeartBeatReply) error
}
