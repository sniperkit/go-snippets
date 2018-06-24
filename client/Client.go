/*
Usage:
$ go run Client.go [serverAddr TCP ip:port] [clientId] [GraphInfoPath] [VertexValue]

serverAddr: The address of the Server.
clientId: ID of the client
GraphInfoPath: the path to the text file containing the vertices
VertexValue: starting value for each vertex
*/

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"project_c9f7_i5l8_o0p4_p0j8/db"
	"project_c9f7_i5l8_o0p4_p0j8/msg"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: add arguments for graph and data processing. Textfile maybe?
func main() {
	// Parse Arguments:
	serverAddr := os.Args[1]
	client := os.Args[2]
	pathToGraph := os.Args[3]
	value := os.Args[4]

	log.SetFlags(log.Lshortfile)

	clientId, err := strconv.Atoi(client)
	checkErr(err)
	val, err := strconv.Atoi(value)
	checkErr(err)

	// Open connection to Server.
	service, err := rpc.Dial("tcp", serverAddr)
	checkErr(err)

	var connectArgs msg.ClientConnectionMsg
	var connectReply msg.ServerConnectionResp
	connectArgs.ClientId = clientId
	err = service.Call("ClientService.Connect", connectArgs, &connectReply)
	checkErr(err)
	fmt.Println("Server Connection Successful %b", connectReply.IsAccepted)

	// TODO also make a Secondary collection
	jobname := createJobName(clientId)
	access, err := db.CreateNewJob(jobname, pathToGraph, float64(val))
	checkErr(err)

	// TODO: Handle any pending jobs(?)

	fmt.Println("Sending Request.")
	// Send job request
	var requestArgs msg.ClientRequestMsg
	var requestReply msg.ServerRequestResp

	requestArgs.ClientId = clientId
	requestArgs.DBAccess = access
	// TODO set Secondary collection

	// TODO: this should come from the Server?
	requestArgs.RequestId = 123
	err = service.Call("ClientService.Request", requestArgs, &requestReply)
	checkErr(err)
	fmt.Println("Request success %b", requestReply.Success)

	outfile := "../sampleData/" + access.PrimaryKey() + "-out"
	err = db.PrintToFile(access.PrimaryKey(), outfile)
}

func createJobName(clientId int) string {
	seed := rand.NewSource(time.Now().UnixNano())
	randomNumber := rand.New(seed)
	id := randomNumber.Int63()
	primary := fmt.Sprintf("%d-%d", clientId, id)

	return primary
}

//PSEUDOCODE:

/*
=============================
Client:

Send ConnectionRequest
if (!reply.ConnectionAccepted)
  exit()


// Notes:
// We might not want to create a client process that just attempts to reconnect forever if a connection drops. That being said, we can be sensible and attempt to reconnect for some measure of time, but I don't think we should essentially hang forever.

if (reply.PendingResults)
  printoutresults for previous jobs

Send(Request)
if (!reply.RequestAccepted)
  exit()

Wait for response
=============================

Server:
Open ClientListener: go thread
Open WorkerListener: go thread
Run Scheduler()

ClientListener:
  Decode the connection request

  // Connection
  If the client has already connected before (clientID)
    check for job that previously completed, add those to reply

  // Request
  If the client already has a job running
    reply "already processing error"
  elseif no ressources/bad data
    reply err
  else
    add the job to PendingJobsQueues

WorkerListener
  Decode connection request
  Add worker to WorkerPool

Scheduler
  // Maintains and processes the pendingjobqueue
  Determine how many jobs to run
  get some jobs from the queue
  Split the workers on the jobs from WorkerPool
  Run for each job (separate processes?)
  add the jobs back on the queue when done

Run(jobid, workersid, )
  assign partitions using PartitionManager
  if first superstep run checkpoint
  run superstep some number of times
  run checkpoint
  // if a worker dies here, abort this iteration
  // if a worker fails (connection is live, workers say), redistribute its partitions to other workers.
  // if too many workers fail, abort this iteration
  // If we abort an iteration too many times -> return failure to client?
  // Update the PartitionManager as needed throughout this

Superstep


=============================
Worker

Send connection request
if (!accept)
  exit

wait for message:
partition:
  update the vertices
superstep:
  // This is pretty much already handled
  // just keep an open socket for incoming messages, send them to the Engine
checkpoint:
  run checkpoint on the vertices


// if the server dies here attempt to reconnect


=============================
Engine
// The Engine does not communicate with a network! It iterfaces with a Worker.

=============================

Assumptions!
- ClientIDs are unique and clients always reconnect with the same ID.
- A single client can only request one task at a time.
- Once a job is accepted, it will only remain in the database for some finite amount of time.
*/
