package vertices

// Vertex interface that vertices should employ to be
// used in a Pregel system
type Vertex interface {
	Update(step int, engineChan chan VertexMessage) bool
	GetID() int
	GetValue() float64
	GetOutVertices() []int
	GetActive() bool
	GetSuperstep() int
	GetMessages() []VertexMessage
	ReceiveMessages(msgs []VertexMessage)
}

// BaseVertex struct that can be embedded in particular
// vertex types like a PageRankVertex
type BaseVertex struct {
	ID          int
	Value       float64
	OutVertices []int
	IncMsgs     []VertexMessage
	Active      bool
	Superstep   int
}

// GetID returns the ID of the vertex
func (bv *BaseVertex) GetID() int {
	return bv.ID
}

// GetValue returns the value of the vertex
func (bv *BaseVertex) GetValue() float64 {
	return bv.Value
}

// GetOutVertices returns the value of the vertex
func (bv *BaseVertex) GetOutVertices() []int {
	return bv.OutVertices
}

// GetSuperstep returns the current superstep
func (bv *BaseVertex) GetSuperstep() int {
	return bv.Superstep
}

// GetMessages returns the incoming messages of the vertex for this superstep
func (bv *BaseVertex) GetMessages() []VertexMessage {
	return bv.IncMsgs
}

// GetActive returns the active status of the vertex
func (bv *BaseVertex) GetActive() bool {
	return bv.Active
}

// GetPageRankVertices creates PageRankVertices from BaseVertices
func GetPageRankVertices(numVertices int, baseVertices map[int]BaseVertex) map[int]Vertex {
	prvMap := make(map[int]Vertex, len(baseVertices))

	for id, baseVertex := range baseVertices {
		// baseVertex := BaseVertex{
		// 	Id:          id,
		// 	Value:       0.0,
		// 	OutVertices: edges,
		// 	IncMsgs:     []VertexMessage{},
		// 	Active:      true,
		// 	Superstep:   0,
		// }
		prvMap[id] = &PageRankVertex{
			BaseVertex:  baseVertex,
			NumVertices: numVertices,
		}
	}
	return prvMap
}
