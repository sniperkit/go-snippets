package main

import (
	"container/list"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"project_c9f7_i5l8_o0p4_p0j8/db"
	"project_c9f7_i5l8_o0p4_p0j8/msg"
	// "github.com/arcaneiceman/GoVector/govec"
)

//====================================================================
// Data Structures

// All the Client Requests that are not yet completed. Used as a FIFO queue.
// TODO: this could probably be changed to a channel.
var PendingRequests list.List

// TODO: add constraints on the RequestId to be non-negative.
// The Requests that have been accepted for each client. Indexed by ClientId. If the entry exists and is a positive integer, there is an accepted request.
var ClientRequests map[int]int

// Stores the last completed request for a ClientId.
// TODO: make this the appropriate data type for the result. Maybe we can cache the result on a temp file.
var CompletedRequests map[int]msg.Result

//====================================================================
//====================================================================
type ClientManager int

func (cm *ClientManager) Initialize(serviceAddr string) {
	// Initialize the data structures.
	ClientRequests = make(map[int]int)
	CompletedRequests = make(map[int]msg.Result)

	// Initialize the RPC Service.
	clientRPCService := new(ClientService)
	rpc.Register(clientRPCService)
	clientListener, err := net.Listen("tcp", serviceAddr)
	checkErr(err)
	log.Println("ClientService: listening for clients at %v", serviceAddr)
	go handleRPC(clientListener.(*net.TCPListener))
}

func (cm *ClientManager) GetRequest() (msg.Request, bool) {
	if PendingRequests.Len() == 0 {
		defer captureDroppedClient()

		c := make(chan msg.Result)
		var d db.Access
		return msg.Request{0, 0, d, 0, c, 0}, false
	} else {
		// TODO: update ClientRequests?
		front := PendingRequests.Front()
		request := front.Value.(msg.Request)
		PendingRequests.Remove(front)
		return request, true
	}
}

func (cm *ClientManager) CompletedRequest(request msg.Request, result msg.Result) {
	defer captureDroppedClient()

	log.Printf("Completed Request %v, result: %v\n", request, result)

	if result.Val == msg.Incomplete {
		PendingRequests.PushBack(result.Request)
	} else {
		// Store the request result
		CompletedRequests[request.ClientId] = result

		// Indicate that we are done processing the request for this client
		ClientRequests[request.ClientId] = 0

		// Continue the RPC call.
		request.OutChannel <- msg.Result{}
	}
}

func captureDroppedClient() {
	if r := recover(); r != nil {
		log.Printf("Recovered Dropped RPC %v\n", r)
	}
}

//====================================================================
//====================================================================
type ClientService int

// TODO: check valid args
func (cs *ClientService) Connect(args *msg.ClientConnectionMsg, reply *msg.ServerConnectionResp) error {
	validArgs := true
	if validArgs {
		Logger.LogLocalEvent(fmt.Sprintf("Client-%v-Connecting", args.ClientId))
		log.Printf("ClientService accepting connection from: %v\n", args.ClientId)
		// TODO: update the response with any previously completed requests.
		// completed, ok = CompletedRequests[args.ClientId]
		// if (ok) {
		// } else {
		// }
		reply.IsAccepted = true
	} else {
		reply.IsAccepted = false
	}

	return nil
}

func (cs *ClientService) Request(args *msg.ClientRequestMsg, reply *msg.ServerRequestResp) error {
	Logger.LogLocalEvent(fmt.Sprintf("ClientRequest-%v-Start", args.RequestId))
	// Check args.
	if args.RequestId <= 0 {
		log.Println("Client requested with non-positive RequestId.")
		reply.Success = false
		return nil
	}

	// Check whether we are currently handling a request.
	currentlyHandling, ok := ClientRequests[args.ClientId]
	if ok && currentlyHandling > 0 {
		log.Println("Client requested a new job while we are processing one already.")
		reply.Success = false
	} else {
		// Create and store the request.
		c := make(chan msg.Result)
		request := msg.Request{args.ClientId, args.RequestId, args.DBAccess, 0, c, 0}
		PendingRequests.PushBack(request)
		ClientRequests[args.ClientId] = args.RequestId

		// Wait for the request to be completed.
		<-c

		result := CompletedRequests[args.ClientId]
		log.Printf("Request completed with result: %v\n", result)

		// TODO: add result as a field of reply
		reply.Success = true
	}

	Logger.LogLocalEvent(fmt.Sprintf("ClientRequest-%v-Done", args.RequestId))

	return nil
}

func handleRPC(l *net.TCPListener) {
	for {
		conn, err := l.Accept()
		checkErr(err)
		go rpc.ServeConn(conn)
	}
}
