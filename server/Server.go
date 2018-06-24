package main

import (
	"log"
	"os"
	"project_c9f7_i5l8_o0p4_p0j8/job"
	"project_c9f7_i5l8_o0p4_p0j8/msg"
	"time"
	"github.com/arcaneiceman/GoVector/govec"
)

// Assumptions Made:

// Global TODOs:
// TODO: garbage collect/close connections? worker/clients/main()
// TODO: consistent naming schemes
// TODO: Documentation
// TODO: usage

var Logger *govec.GoLog

//============================================================
func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

//============================================================
func init() {
	Logger = govec.Initialize("server", "serverlogfile")
}

//============================================================
// This function keeps pulling from PendingRequests and completes the jobs
func work(clientManager *ClientManager, workerManager *WorkerManager) {
	for {
		// Assign workers for this task.
		selectedWorkers := workerManager.SelectWorkers()
		if len(selectedWorkers) == 0 {
			log.Println("work(): No available workers.")
			time.Sleep(time.Second)
			continue
		}

		// Get a request.
		request, ok := clientManager.GetRequest()
		if ok == false {
			log.Println("work(): No pending requests.")
			time.Sleep(time.Second)
		} else {
			log.Printf("work(): Handling request: %v with workers: %v\n", request, selectedWorkers)

			// Create the message channels assigned for this request.
			// TODO: check the size of these parameters.
			cIn := make(chan msg.FromWorker, 10000)
			cOut := make(chan msg.FromServer, 10000)
			cResult := make(chan msg.Result)

			workerManager.PrepareWorkers(selectedWorkers, cIn)

			// Run the job.
			jobResult := msg.Result{msg.Nil, msg.Request{}}
			go job.Run(request, selectedWorkers, cIn, cOut, cResult)
			for {
				select {
				case jobResult = <-cResult:
					if jobResult.Val == msg.Success {
						log.Printf("work(): Success - %v", msg.ResultStr(jobResult))
						break
					} else {
						log.Printf("work(): Not Success yet: %v", msg.ResultStr(jobResult))
					}

				case msgOut := <-cOut:
					log.Printf("work(): sending out message - %v\n", msgOut)
					workerManager.SendMessageToWorker(msgOut)
				}

				if jobResult.Val != msg.Nil {
					break // We only want to break when jobResult has been set.
				}
			}

			// Cleanup.
			workerManager.CloseWorkers(selectedWorkers)
			clientManager.CompletedRequest(request, jobResult)
			close(cIn)
			close(cOut)
			close(cResult)

			log.Printf("work(): Finished handling request %v\n", request)
		}
	}
}

//============================================================
// Entry point to the server.
func main() {
	log.SetFlags(log.Lshortfile)

	clientServiceAddr := os.Args[1]
	clientManager := new(ClientManager)
	clientManager.Initialize(clientServiceAddr)

	workerServiceAddr := os.Args[2]
	workerManager := new(WorkerManager)
	workerManager.Initialize(workerServiceAddr)

	work(clientManager, workerManager)
}
