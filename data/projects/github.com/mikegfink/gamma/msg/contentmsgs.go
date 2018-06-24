// A number of different types used for messages and by messages between
// workers and server

package msg

import (
	"fmt"
	"project_c9f7_i5l8_o0p4_p0j8/db"
)

//import "../db"

//import "../db"

// ===========================================================================
// ===========================================================================
// Node->Node messages

type Type int

// The comments on each line indicate which special fields are expected to be filled
const (
	// Server -> Worker message types
	NilType        Type = iota // can't start with 0 or else message Decoding fails
	Assign                     // DBKey, Partition, DstWorker
	Superstep                  // StepNum, DstWorker
	SaveCheckpoint             // DBKey, DstWorker
	LoadCheckpoint             // DBKey, DstWorker

	// Worker->Server messsage types
	PartitionAck      // SrcWorker
	Done              // SrcWorker
	Inactive          // SrcWorker
	SaveCheckpointAck // SrcWorker
	LoadCheckpointAck // SrcWorker

	// Both Server->Worker and Worker->Server
	V2V // DstVertex, Msg, [DstWorker (FromServer only)], (SrcVertex, SrcWorker, StepNum)

)

func TypeStr(t Type) string {
	switch t {
	case Assign:
		return "Assign"
	case Superstep:
		return "Superstep"
	case SaveCheckpoint:
		return "SaveCheckpoint"
	case LoadCheckpoint:
		return "LoadCheckpoint"
	case PartitionAck:
		return "PartitionAck"
	case Done:
		return "Done"
	case Inactive:
		return "Inactive"
	case SaveCheckpointAck:
		return "SaveCheckpointAck"
	case LoadCheckpointAck:
		return "LoadCheckpointAck"
	case V2V:
		return "V2V"
	default:
		return fmt.Sprintf("Illegal msg.State: %v", t)
	}

}

// ===========================================================================

type FromServer struct {
	Type Type

	SrcWorker WorkerId
	DstWorker WorkerId

	SrcVertex VertexId
	DstVertex VertexId

	Msg     float64 // TBD: Should this really be []byte to allow for diff job types?
	StepNum int     // Current superstep number

	Partitions []Partition
	// HACK
	// PartitionsStart []int
	// PartitionsEnd   []int
	DBKey string // This will be one of the two db.Access keys

	Mid int //message id (for debugging)
}

var mcounter int = 0 // counter ti give Message Ids; only for hacky debugging

func NewAssign(dbKey string, partitions []Partition, dstWorker WorkerId) FromServer {
	var fs FromServer
	fs.Type = Assign
	fs.DBKey = dbKey
	fs.Partitions = partitions

	fs.DstWorker = dstWorker

	fs.Mid = mcounter
	mcounter++
	return fs
}

func NewSuperstep(stepNum int, dstWorker WorkerId) FromServer {
	var fs FromServer
	fs.Type = Superstep
	fs.StepNum = stepNum
	fs.DstWorker = dstWorker

	fs.Mid = mcounter
	mcounter++
	return fs
}

func NewV2VServer(dst VertexId, msg float64, dstWorker WorkerId, stepNum int,
	src VertexId, srcWorker WorkerId) FromServer {
	var fs FromServer
	fs.Type = V2V
	fs.DstVertex = dst
	fs.Msg = msg
	fs.DstWorker = dstWorker

	fs.StepNum = stepNum
	fs.SrcVertex = src
	fs.SrcWorker = srcWorker

	fs.Mid = mcounter
	mcounter++
	return fs
}

func NewSaveCheckpoint(dbKey string, dstWorker WorkerId) FromServer {
	var fs FromServer
	fs.Type = SaveCheckpoint
	fs.DBKey = dbKey
	fs.DstWorker = dstWorker

	fs.Mid = mcounter
	mcounter++
	return fs
}

// ===========================================================================

type FromWorker struct {
	Type Type

	SrcWorker WorkerId
	DstWorker WorkerId

	DstVertex VertexId
	SrcVertex VertexId

	Msg float64 // TBD: Should this really be []byte to allow for diff job types?

	// The rest are only for debugging purposes
	StepNum int
}

// ===========================================================================
// Client->Server messages

// The client sends to request that a job is run.
type ClientRequestMsg struct {
	ClientId  int
	RequestId int
	DBAccess  db.Access
	// TODO: Parameters needed for passing and processing the graph data
	// For example:
	// GraphBinary []byte (or other format)
	// MaxSuperStepCount int
	// RequestType <- if we want
}

// ===========================================================================
// ========================Common Types=======================================

type WorkerId string

type VertexId int

func (vid VertexId) Int() int {
	return int(vid)
}

/////////////////////////// Partition ///////////////////

/// As in slices, startIdx is first index, and the last index is one before endIdx
type Partition struct {
	StartIdx int
	EndIdx   int
}

func NewPartition(startIdx, endIdx int) Partition {
	return Partition{startIdx, endIdx}
}

func (p Partition) Size() int {
	return p.EndIdx - p.StartIdx
}

// The first index in p
func (p Partition) Start() VertexId {
	return VertexId(p.StartIdx)
}

// The last index + 1 in p
func (p Partition) End() VertexId {
	return VertexId(p.EndIdx)
}

////////////////////////////////////////////////////////////
