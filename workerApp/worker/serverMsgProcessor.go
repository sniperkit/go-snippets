package worker

import (
	"log"

	"project_c9f7_i5l8_o0p4_p0j8/msg"
	"project_c9f7_i5l8_o0p4_p0j8/vertices"
)

// ServerMsgProcessor is an adapter for receiving messages from the server
// and calling the appropriate methods in the worker
type ServerMsgProcessor struct {
	worker          *Worker
	outMsgChan      chan msg.FromWorker
	vToVMsgChan     chan vertices.VertexMessage
	inactiveMsgChan chan vertices.ActiveMessage
	wID             msg.WorkerId
	stepDoneChan    chan int
}

// NewServerMsgProcessor creates a new message processor and worker to
// send the messages to.
func NewServerMsgProcessor(wID msg.WorkerId, outMsgChan chan msg.FromWorker) *ServerMsgProcessor {
	vtoVMsgChan := make(chan vertices.VertexMessage)
	inactiveMsgChan := make(chan vertices.ActiveMessage)
	stepDoneChan := make(chan int)
	worker := NewWorker(vtoVMsgChan, inactiveMsgChan, stepDoneChan)

	smp := &ServerMsgProcessor{
		worker:          worker,
		outMsgChan:      outMsgChan,
		vToVMsgChan:     vtoVMsgChan,
		inactiveMsgChan: inactiveMsgChan,
		stepDoneChan:    stepDoneChan,
		wID:             wID,
	}
	return smp
}

// Process accepts a pregel.Message and processes it into
func (smp *ServerMsgProcessor) Process(serverMsg msg.FromServer) {
	switch serverMsg.Type {
	case msg.V2V:
		log.Println("SMP: Received a V2V message.")
		vertexMsg := vertices.VertexMessage{
			FromID:    serverMsg.SrcVertex.Int(),
			Value:     serverMsg.Msg,
			ToID:      serverMsg.DstVertex.Int(),
			Superstep: serverMsg.StepNum,
		}
		smp.worker.ReceiveNetVertexMessage(vertexMsg)
		break
	case msg.Assign:
		log.Println("SMP: Received an assign/load message.")
		jobName := serverMsg.DBKey

		var workerPartitions = []struct {
			min int
			max int
		}{}
		for _, partition := range serverMsg.Partitions {
			minID := partition.Start().Int()
			maxID := partition.End().Int()
			workerPartition := struct {
				min int
				max int
			}{minID, maxID}
			workerPartitions = append(workerPartitions, workerPartition)
		}
		err := smp.worker.LoadVertices(jobName, workerPartitions)
		ackMsg := msg.FromWorker{
			Type:      msg.PartitionAck,
			SrcWorker: smp.wID,
		}
		if err == nil {
			log.Println("SMP: Successfully loaded vertices.")
			ackMsg.Msg = 1.0
		} else {
			log.Println("SMP: Failed to load vertices.")
			ackMsg.Msg = 0.0
		}
		smp.outMsgChan <- ackMsg
		break
	case msg.SaveCheckpoint:
		log.Println("SMP: Received a save vertices message.")
		jobName := serverMsg.DBKey
		success := smp.worker.SaveVertices(jobName)
		ackMsg := msg.FromWorker{
			Type:      msg.SaveCheckpointAck,
			SrcWorker: smp.wID,
		}
		if success {
			log.Println("SMP: Successfully saved vertices.")
			ackMsg.Msg = 1.0
		} else {
			log.Println("SMP: Failed to save vertices.")
			ackMsg.Msg = 0.0
		}
		smp.outMsgChan <- ackMsg
		break

	//// Server Never sends this. It only sends msg.Assign
	//case msg.LoadCheckpoint:
	//jobName := serverMsg.DBKey
	//minID := serverMsg.Partition.Start().Int()
	//maxID := serverMsg.Partition.End().Int()
	//smp.worker.StopReceiver()
	//err := smp.worker.LoadVertices(jobName, minID, maxID)
	//ackMsg := msg.FromWorker{
	//Type:      msg.LoadCheckpointAck,
	//SrcWorker: smp.wID,
	//}
	//if err == nil {
	//ackMsg.Msg = 1.0
	//} else {
	//ackMsg.Msg = 0.0
	//}
	//smp.outMsgChan <- ackMsg
	//break

	case msg.Superstep:
		log.Println("SMP: Received a start superstep message.")
		smp.worker.PrepareSuperstep()
		go smp.worker.Superstep(serverMsg.StepNum)
		go func() {
			for {
				select {
				case step := <-smp.stepDoneChan:
					ackMsg := msg.FromWorker{
						Type:      msg.Done,
						SrcWorker: smp.wID,
						StepNum:   step,
					}
					smp.outMsgChan <- ackMsg
					return
				case vToVMsg := <-smp.vToVMsgChan:
					ackMsg := msg.FromWorker{
						Type:      msg.V2V,
						SrcWorker: smp.wID,
						StepNum:   vToVMsg.Superstep,
						SrcVertex: msg.VertexId(vToVMsg.FromID),
						DstVertex: msg.VertexId(vToVMsg.ToID),
						Msg:       vToVMsg.Value,
					}
					smp.outMsgChan <- ackMsg
					break
				case inactiveMsg := <-smp.inactiveMsgChan:
					ackMsg := msg.FromWorker{
						Type:      msg.Inactive,
						SrcWorker: smp.wID,
						SrcVertex: msg.VertexId(inactiveMsg.VertexID),
						StepNum:   inactiveMsg.Superstep,
					}
					smp.outMsgChan <- ackMsg
					break
				}
			}
		}()
		break
	default:
		log.Panic("SMP: Unrecognized message type: ", serverMsg.Type)

	}
}
