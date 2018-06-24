//////////////////////////////////////////////////////////////////////////////
// Processes messages from and replies to workers, given particular request ID.
//////////////////////////////////////////////////////////////////////////////

package job

import (
	"fmt"
	"log"
	"project_c9f7_i5l8_o0p4_p0j8/db"
	"project_c9f7_i5l8_o0p4_p0j8/manager"
	"project_c9f7_i5l8_o0p4_p0j8/msg"
	"runtime"
	"time"
)

// Timeout if the server is waiting for a reply from a worker
const timeout time.Duration = 5 * time.Second

// Save Checkpoint after this many checkpoints, and repeat
const checkpoint_rate = 10

// If all nodes aren't yet inactive (i.e., having voted to halt), Pregel will stop after this supserstep
const max_supersteps = 20

func Run(
	request msg.Request,
	workers []msg.WorkerId,
	cIn chan msg.FromWorker,
	cOut chan msg.FromServer,
	cDone chan msg.Result) {
	log.Printf("Request started: %v", request)

	num_vertices, err := db.NumVertices(request.DBAccess.Key())

	if len(workers) == 0 || err != nil {
		result := msg.Result{msg.Failure, request}
		cDone <- result
		return
	}

	mgr := manager.New(num_vertices)
	addWorkers(mgr, workers)

	// If something fails while communicating with workers, there will be a panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%v: Requeuening Incomplete request %v", r, request)
			// reset superstep to last checkpointstep
			request.Superstep = request.CheckpointStep
			cDone <- msg.Result{msg.Incomplete, request}
		}
	}()

	for {

		assigns := mgr.Redistribute()
		sendAssignments(assigns, request.DBAccess.Key(), cIn, cOut)
		done := iterateSupersteps(&request, mgr, cIn, cOut) // Won't be done if needs redistribution to rebalance work
		if done || request.Superstep > max_supersteps {
			log.Printf("Job COMPLETE: %v", request)
			result := msg.Result{msg.Success, request}
			cDone <- result
			return
		}
	}

}

func addWorkers(mgr manager.Manager, workers []msg.WorkerId) {
	for _, worker := range workers {
		mgr.AddWorker(worker)
	}
}

func logMessageFW(fw msg.FromWorker) {
	_, _, line, _ := runtime.Caller(1)
	log.Printf("Line %v -- RECEIVED msg Type %v; Data %v", line, msg.TypeStr(fw.Type), fw)
}

func logMessageFS(fs msg.FromServer) {
	_, _, line, _ := runtime.Caller(1)
	log.Printf("Line %v -- SENDING msg Type %v; Data %v", line, msg.TypeStr(fs.Type), fs)
}

func sendAssignments(assigns []manager.Assignment, dbKey string, cIn chan msg.FromWorker, cOut chan msg.FromServer) {
	log.Printf("Sending %v Assignments to %v", len(assigns), dbKey)
	for _, a := range assigns {
		cOut <- msg.NewAssign(dbKey, a.Partitions, a.Worker)
	}

	acks := make(map[msg.WorkerId]struct{})
	for {
		select {
		case fw := <-cIn:
			logMessageFW(fw)
			if fw.Type == msg.PartitionAck {
				if fw.Msg == 0 {
					//db problems
					panic("Workers Unable to load vertices from db.")
				} else {
					// no db problems
					acks[fw.SrcWorker] = struct{}{}
				}
			} else {
				log.Printf("Error: Msg type expected PartitionAck, Received %v - %v\n", msg.TypeStr(fw.Type), fw)
			}
			if len(acks) == len(assigns) {
				return
			}
		case <-time.After(timeout):
			panic("Timed out waiting for PartitionAck")
		}
	}

}

// Returns True if completed Pregel
// Returns False if needs redistribution of workers before continuing
// Panics if incompmlete for some reason (worker died, or joined, so this needs to be restarted at checkpoint)
func iterateSupersteps(request *msg.Request, mgr manager.Manager,
	cIn chan msg.FromWorker, cOut chan msg.FromServer) bool {
	log.Printf("Beginning SUPERSTEP %v", request.Superstep)

	// Used for calculation elapsed times
	startTimes := map[msg.WorkerId]time.Time{}

	for _, w := range mgr.Workers() {
		fs := msg.NewSuperstep(request.Superstep, w)
		startTimes[w] = time.Now()
		logMessageFS(fs)
		cOut <- fs

	}

	// true means existance means finished; Type is Done or Inactive
	dones := make(map[msg.WorkerId]msg.Type)
	for {
		halt := true
		select {
		case fw := <-cIn:
			logMessageFW(fw)
			switch fw.Type {
			case msg.Inactive:
				// Since it did no work, don't update elapsed time
				dones[fw.SrcWorker] = msg.Inactive
			case msg.Done:
				// It has complete its superstep, so update its elapsed time
				wid := fw.SrcWorker
				elapsedTime := time.Since(startTimes[wid])
				log.Printf("Superstep %v, Worker %v, time %v", request.Superstep, wid, elapsedTime)
				mgr.SetElapsedTime(wid, elapsedTime)
				halt = false
				dones[fw.SrcWorker] = msg.Done
			case msg.V2V:
				halt = false
				dstWorker := mgr.GetWorker(fw.DstVertex)
				fs := msg.NewV2VServer(fw.DstVertex, fw.Msg, dstWorker, fw.StepNum, fw.SrcVertex, fw.SrcWorker)
				logMessageFS(fs)
				cOut <- fs
			default:
				log.Printf("Error: Unexpected message during Superstep %v: %v",
					request.Superstep, fw)
			}

			if len(dones) == mgr.NumWorkers() {
				log.Printf("Completed Superstep %v", request.Superstep)
				log.Printf("\tFastest to Slowest ratio: %v", mgr.FastestToSlowest())
				log.Printf("\tIs in optimal range?: %v", mgr.IsOptimal())

				// Superstep is complete
				request.Superstep++
				if request.Superstep >= max_supersteps {
					halt = true
				}
				if (request.Superstep%checkpoint_rate == 0) || halt || !mgr.IsOptimal() {
					saveCheckpoint(request.DBAccess.OtherKey(), mgr.Workers(), cIn, cOut)
					(&request.DBAccess).SwapKeys()
					request.CheckpointStep = request.Superstep

					if halt {
						if request.DBAccess.PrimaryKey() != request.DBAccess.Key() {
							// It most recently saved into the Secondary key
							// so we have to copy it to the primary key for the client
							saveCheckpoint(request.DBAccess.OtherKey(), mgr.Workers(), cIn, cOut)
						}
						return true
					} else if !mgr.IsOptimal() {
						return false
					} else {
						return false
						//panic("No problem. Just checking if there are any new workers.")
					}
				}
				mgr.ResetSpeeds()
				return iterateSupersteps(request, mgr, cIn, cOut)
			}
		case <-time.After(timeout):
			panic(fmt.Sprintf("Timed out during Superstep %v", request.Superstep))
		}
	}
}

func saveCheckpoint(dbKey string, workers []msg.WorkerId, cIn chan msg.FromWorker, cOut chan msg.FromServer) {
	log.Printf("Saving CHECKPOINTS in %v", dbKey)
	for _, w := range workers {
		fs := msg.NewSaveCheckpoint(dbKey, w)
		logMessageFS(fs)
		cOut <- fs
	}

	acks := make(map[msg.WorkerId]struct{})
	for {
		select {
		case fw := <-cIn:
			logMessageFW(fw)
			if fw.Type == msg.SaveCheckpointAck {
				if fw.Msg == 0 {
					//db problems
					panic("Workers unable to SaveCheckpoint")
				} else {
					// no db problems
					acks[fw.SrcWorker] = struct{}{}
				}
			} else {
				log.Printf("Error: Msg type expected SaveCheckpointAck, Received %v", fw)
			}
			if len(acks) == len(workers) {
				return
			}
		case <-time.After(timeout):
			panic("Timed out waiting for SaveCheckpointAck")
		}
	}
}
