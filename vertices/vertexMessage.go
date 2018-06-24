package vertices

// VertexMessage is a struct for passing messages from vertex to vertex in
// a worker. It is converted to a Msg to be sent across the network.
type VertexMessage struct {
	FromID    int
	Value     float64
	ToID      int
	Superstep int
}

// ActiveMessage is a struct to indicate whether a vertex is active or not
// after a superstep has completed.
type ActiveMessage struct {
	VertexID  int
	Active    bool
	Superstep int
}
