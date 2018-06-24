package vertices

// PageRankVertex implements the Vertex interface and allows the page rank
// algorithm to be performed on the vertex.
type PageRankVertex struct {
	BaseVertex
	NumVertices int
}

const vertValue = 0.15
const edgeWeight = 0.85

// Update runs one superstep on the PageRankVertex. It sends it's outgoing
// messages back to the engine via the engineChan.
func (prv *PageRankVertex) Update(step int, engineChan chan VertexMessage) bool {
	prv.Superstep = step
	sum := 0.0
	for _, msg := range prv.IncMsgs {
		sum += msg.Value
	}
	prv.Value = vertValue/float64(prv.NumVertices) + (float64(sum) * edgeWeight)

	outgoingPageRank := prv.Value / float64(len(prv.OutVertices))
	//log.Println("Vertex: Value: ", prv.Value, "OutValue:", outgoingPageRank, "Superstep: ", prv.Superstep, "Inc Sum:", sum)
	for _, id := range prv.OutVertices {
		// if id > 10000 || id == 0 {
		// 	log.Println("PageRank val: ", id, " ID: ", prv.ID)
		// }

		outMsg := VertexMessage{
			FromID:    prv.ID,
			Value:     outgoingPageRank,
			ToID:      id,
			Superstep: prv.Superstep,
		}
		engineChan <- outMsg
	}
	//log.Println("Sent messages from vertex id: ", prv.Id)

	// PageRank always has active vertices
	return true
}

// ReceiveMessages accepts messages for the next Superstep
func (prv *PageRankVertex) ReceiveMessages(msgs []VertexMessage) {
	// Could validate, but shouldn't be necessary.
	//log.Println("Received msgs:", prv.Id)
	// TODO: SHould be append?
	prv.IncMsgs = msgs
}
