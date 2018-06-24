package worker

import (
	"log"
	"runtime"

	"project_c9f7_i5l8_o0p4_p0j8/db"
	"project_c9f7_i5l8_o0p4_p0j8/vertices"
)

// Worker struct holds the engines and communication between them
type Worker struct {
	messages            map[int][]vertices.VertexMessage
	engines             []*Engine
	msgDistributionChan chan vertices.VertexMessage
	vertexMap           map[int]vertices.Vertex
	stepDoneChan        chan int
	stopChan            chan bool
	stopReceiver        chan bool
	localMsgChan        chan vertices.VertexMessage
	inactiveVertexChan  chan vertices.ActiveMessage
	serverVtoVChan      chan vertices.VertexMessage
	hasStarted          bool
}

// NewWorker allows a new worker to be constructed with a particular batch
// size for network communication
func NewWorker(serverVtoVChan chan vertices.VertexMessage, inactiveVertexChan chan vertices.ActiveMessage, stepDoneChan chan int) *Worker {
	numEngines := runtime.NumCPU()
	worker := &Worker{
		messages:            make(map[int][]vertices.VertexMessage),
		engines:             make([]*Engine, numEngines, numEngines),
		msgDistributionChan: make(chan vertices.VertexMessage), //, numEngines)
		inactiveVertexChan:  inactiveVertexChan,
		serverVtoVChan:      serverVtoVChan,
		stepDoneChan:        stepDoneChan,
		stopChan:            make(chan bool),
		stopReceiver:        make(chan bool),
	}
	return worker
}

// PrepareSuperstep is to be run sequentially before the superstep is run in a
// goroutine. This ensures that only messages for this superstep are used.
func (w *Worker) PrepareSuperstep() {
	log.Println("Worker: Preparing for superstep, distributing messages to vertices and clearing for next round")
	w.distributeMessages()
	w.clearMessages()
}

func (w *Worker) distributeMessages() {
	for vid, vertex := range w.vertexMap {
		vertex.ReceiveMessages(w.messages[vid])
	}

}

// Superstep runs a single superstep on the vertices assigned to this worker
func (w *Worker) Superstep(stepNum int) {
	log.Println("Worker: Running superstep #:", stepNum)
	engDoneChan := make(chan bool)
	w.localMsgChan = make(chan vertices.VertexMessage)

	go w.receiveLocalMsgs(stepNum)

	for _, engine := range w.engines {
		go engine.Superstep(w.localMsgChan, stepNum, engDoneChan)
	}
	for _ = range w.engines {
		<-engDoneChan
	}
	close(w.localMsgChan)

}

func (w *Worker) clearMessages() {
	for k := range w.messages {
		delete(w.messages, k)
	}
}

func (w *Worker) receiveLocalMsgs(stepNum int) {
	for {
		select {
		case msg, more := <-w.localMsgChan:
			if more {
				if _, vok := w.vertexMap[msg.ToID]; vok {
					log.Println("Worker: Local message to ID: ", msg.ToID)
					w.msgDistributionChan <- msg
				} else {
					log.Println("Worker: Sending out id: ", msg.ToID)
					w.serverVtoVChan <- msg
				}
			} else {
				log.Println("Worker: Finished superstep: ", stepNum)
				w.stepDoneChan <- stepNum
				w.localMsgChan = nil
				return
			}
			break
		case _ = <-w.stopChan:
			log.Println("Worker: Aborting superstep.")
			return
		}
	}
}

// StopReceiver allows us to ensure that the receiver is turned off until
// the messages have been transferred to the engines and the worker is ready
// for messages in the next superstep.
func (w *Worker) StopReceiver() {
	log.Println("Worker: Stopping the Receiver")

	if !w.hasStarted {
		log.Println("Worker: Receiver was never started")
		return
	}

	if w.localMsgChan != nil {
		close(w.localMsgChan)
		<-w.stopChan
	}

	if w.msgDistributionChan != nil {
		close(w.msgDistributionChan)
		<-w.stopReceiver
	}

	w.hasStarted = false
	log.Println("Worker: Receivers stopped.")
}

// RunReceiver starts the receiver for accepting messages for a superstep
// Messages can be accepted from local or network vertices - for local or
// network vertices
func (w *Worker) RunReceiver() {
	log.Println("Worker: Starting the Receiver up.")
	w.hasStarted = true
	w.msgDistributionChan = make(chan vertices.VertexMessage)
	for {
		select {
		case msg, more := <-w.msgDistributionChan:
			if more {
				if _, ok := w.messages[msg.ToID]; ok {
					w.messages[msg.ToID] = append(w.messages[msg.ToID], msg)
				} else {
					w.messages[msg.ToID] = []vertices.VertexMessage{msg}
				}
			} else {
				log.Println("Worker: All messages received.")
				w.stopReceiver <- true
				return
			}
		}
	}
}

// ReceiveNetVertexMessage receives a VertexMessage and drops it in the message
// channel.
func (w *Worker) ReceiveNetVertexMessage(msg vertices.VertexMessage) {
	log.Print("Worker: Received outside message, id:", msg.ToID)
	if w.hasStarted {
		if _, vok := w.vertexMap[msg.ToID]; vok {
			w.msgDistributionChan <- msg
		} else {
			log.Println("Worker: That message doesn't belong here: ", msg.ToID)
			//w.serverVtoVChan <- msg
		}
	} else {
		log.Println("Worker: Message out of order, this worker is closed.", msg)
	}
}

// LoadVertices loads vertices into the engines.
// TODO: Add type of job from server
func (w *Worker) LoadVertices(jobName string, partitions []struct {
	min int
	max int
}) error {
	allBaseVertices := make(map[int]vertices.BaseVertex)
	log.Println("Worker: Loading vertices with name: ", jobName, "Partitions:", partitions)
	for _, partition := range partitions {
		minID := partition.min
		maxID := partition.max

		partitionBaseVertices, err := db.BatchGet(jobName, minID, maxID)
		if err != nil {
			return err
		}
		for bvid, bv := range partitionBaseVertices {
			allBaseVertices[bvid] = bv
			log.Print("Node Id:", bvid)
		}
	}
	log.Println("Worker: Loaded", len(allBaseVertices), "vertices.")

	numVertices, err := db.NumVertices(jobName)
	if err != nil {
		return err
	}
	w.StopReceiver()
	w.clearMessages()
	w.vertexMap = w.getVertices(numVertices, allBaseVertices, vertices.PageRank)

	engineMaps := make([]map[int]vertices.Vertex, len(w.engines))
	for id := 0; id < len(w.engines); id++ {
		engineMaps[id] = make(map[int]vertices.Vertex)
	}

	for vid, vertex := range w.vertexMap {
		engID := vid % len(w.engines)
		engineMaps[engID][vid] = vertex
		w.messages[vid] = vertex.GetMessages()
	}

	for eid := 0; eid < len(w.engines); eid++ {

		engine := NewEngine(engineMaps[eid], eid, w.inactiveVertexChan)

		w.engines[eid] = engine
	}
	go w.RunReceiver()
	return nil
}

// SaveVertices saves the current vertex information of this worker in the
// database under the provided jobName.
func (w *Worker) SaveVertices(jobName string) bool {
	log.Println("Worker: Saving vertices.")
	allVertices := make(map[int]vertices.Vertex)
	for _, engine := range w.engines {
		vertices := engine.GetVertices()
		for vid, vertex := range vertices {
			vertex.ReceiveMessages(w.messages[vid])
			allVertices[vid] = vertex
		}
	}

	err := db.BatchUpdate(jobName, allVertices)
	success := true
	if err != nil {
		success = false
	}
	return success
}

func (w *Worker) getVertices(numVertices int, vertexMap map[int]vertices.BaseVertex, jobType vertices.VertexType) map[int]vertices.Vertex {
	switch jobType {
	case vertices.PageRank:
		return vertices.GetPageRankVertices(numVertices, vertexMap)
	default:
		log.Panic("Worker: Unrecognized job type: ", jobType)
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
