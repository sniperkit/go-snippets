package vertices

// VertexType is an enumeration to identify what type of operation
// a worker will be running on its vertices
type VertexType int

const (
	// PageRank indicates that the page rank algorithm will be run on these
	// vertices
	PageRank VertexType = iota
)
