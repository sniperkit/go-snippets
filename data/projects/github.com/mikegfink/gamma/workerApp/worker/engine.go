package worker

import (
	"log"

	"project_c9f7_i5l8_o0p4_p0j8/vertices"
)

// Engine is a struct that represents a vertex calculation engine that is
// responsible for a group of vertices in a worker.
type Engine struct {
	ID                 int
	vertexMap          map[int]vertices.Vertex
	inactiveVertexChan chan vertices.ActiveMessage
}

// NewEngine creates a new Engine with the given values.
func NewEngine(vertexMap map[int]vertices.Vertex, ID int, inactiveVertexChan chan vertices.ActiveMessage) *Engine {
	engine := Engine{
		ID:                 ID,
		vertexMap:          vertexMap,
		inactiveVertexChan: inactiveVertexChan,
	}
	return &engine
}

// GetVertices returns a map of the vertices this engine is responsible for
func (e *Engine) GetVertices() map[int]vertices.Vertex {
	return e.vertexMap
}

// Superstep runs a superstep on each vertex of the engine.
func (e *Engine) Superstep(workerMsgChan chan vertices.VertexMessage, step int, done chan bool) {
	//log.Println("Starting superstep for engine: ", e.ID, "with vertices: ", len(e.vertexMap))

	for _, vertex := range e.vertexMap {
		active := vertex.Update(step, workerMsgChan)
		if !active {
			inactiveMsg := vertices.ActiveMessage{
				VertexID:  vertex.GetID(),
				Active:    false,
				Superstep: vertex.GetSuperstep(),
			}
			e.inactiveVertexChan <- inactiveMsg
		}
	}
	done <- true

}

func (e *Engine) distributeMessages(msgs map[int][]vertices.VertexMessage) {

	for id, vertexMsgs := range msgs {
		if _, ok := e.vertexMap[id]; !ok {
			log.Println("That message shouldn't have been sent here, id: ", id, e.ID, vertexMsgs[0].ToID)
		} else {
			e.vertexMap[id].ReceiveMessages(vertexMsgs)
		}

	}
}

// func (e *Engine) runSender(stopChan chan bool) {
// 	//log.Println("Starting run sender")
// 	for {
// 		select {
// 		case outMessage := <-e.vertexMsgChan:
// 			if outMessage.ToId == 0 {
// 				log.Println("Engine: That's not right, no node with id 0:", outMessage)
// 			}
// 			if e.outMsgs == nil {
// 				e.outMsgs = []vertices.VertexMessage(nil)
// 			}
//
// 			e.outMsgs = append(e.outMsgs, outMessage)
//
// 			if len(e.outMsgs) > e.batchSize {
// 				for _, msg := range e.outMsgs {
// 					if msg.ToId > 10000 || msg.ToId == 0 {
// 						log.Println("Engine val: ", msg.ToId, " ID: ", msg.FromId, len(e.outMsgs), e.ID)
// 					}
// 				}
// 				e.workerMsgChan <- e.outMsgs
// 				e.outMsgs = []vertices.VertexMessage(nil)
// 			}
// 			break
// 		case _ = <-stopChan:
// 			return
// 		}
//
// 	}
// }
