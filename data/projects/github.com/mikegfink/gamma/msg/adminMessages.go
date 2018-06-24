package msg

import (
	"fmt"
	"project_c9f7_i5l8_o0p4_p0j8/db"
)

// Represent a job request from a client. The OutChannel is used to indicate when this job has been handled and should be returned to the client (or stored if Client is unavailable) since Client<->Server is RPC.
type Request struct {
	ClientId       int
	RequestId      int
	DBAccess       db.Access
	Superstep      int
	OutChannel     chan Result
	CheckpointStep int
}

// The result of a job operation
type Result struct {
	Val     ResultVal // The value, as shown in const below
	Request Request   // In case value is not Success, here is the requet to be requeed
}

// status of the request
type ResultVal int

const (
	Nil        ResultVal = iota
	Success              // Completed, can return to client
	Incomplete           // Retry
	Failure              // Fatal error, unfinishable
)

func ResultStr(r Result) string {
	switch r.Val {
	case Nil:
		return "Nil"
	case Success:
		return "Success"
	case Incomplete:
		return "Incomplete"
	case Failure:
		return "Failure"
	default:
		return fmt.Sprintf("Illegal msg.ResultVal: %v", r.Val)

	}
}

// ===============================================

// ErrMsg Basic error message
type ErrMsg struct {
	Error string
}

// Client - Server Messages:

// ClientConnectionMsg The client sends this when opening a connection
type ClientConnectionMsg struct {
	ClientId int
}

// ServerConnectionResp A server's reply to a connection attempt.
// TODO: this should return the value of any job previously run/pending jobs.
type ServerConnectionResp struct {
	IsAccepted  bool
	PendingJobs []int
}

// ServerRequestResp A server's reply to a request.
// If the request was not completed, reply val will be the error string.
// If the request was successfully completed, reply val will be the return data.
type ServerRequestResp struct {
	RequestId int
	Success   bool
	ReplyVal  string
}

// Server <-> Worker Messages:

// WorkerConnectionMsg Sent by worker when connecting to the server.
type WorkerConnectionMsg struct {
	WorkerId      WorkerId
	WorkerAddress string
}

// WorkerConnectionResp Response from server to a Worker Connection Requirest
type WorkerConnectionResp struct {
	WorkerId   WorkerId
	IsAccepted bool
}
