// Provides management of vertices on workers, and determines when a
// redistribution of vertices would increase efficiency,
// based on performance in previous supersteps
//
// To instantiate, simply call
//		manager.New(currentNumberOfVerticesInGraph)
//

// TODO incorporate historical data in some way, in case the last run was
// uncharacteristic. Perhaps a historical average

//////////////////////////////////////////////////////////////////////////////

package manager

import (
	"time"

	"project_c9f7_i5l8_o0p4_p0j8/msg"
)

// TBD determine best values
const MaxOptimalRatio float64 = 100
const PartitionsPerWorkerAvg int = 10

// ===========================================================================

// The partitions assigned to a given worker
type Assignment struct {
	Worker     msg.WorkerId
	Partitions []msg.Partition
}

type Manager interface {
	// These are the workers whose performance metrics are being analyzed
	// Re-adding a worker isignored
	AddWorker(worker msg.WorkerId)
	RemoveWorker(worker msg.WorkerId)

	// Return the number of workers known
	NumWorkers() int

	// Returns a slice of all known worker Ids
	Workers() []msg.WorkerId

	// The number of vertices in the Graph.
	// If this changes, the vertices will have to be redistributed
	SetNumVertices(num int)

	// Returns the calculated size of each partition (except, perhaps, the last)
	PartitionSize() int

	// the time it took for worker to complete the last superstep
	// Also sets worker speed if vertices have been assigned
	SetElapsedTime(worker msg.WorkerId, dur time.Duration)

	// This worker's speed, in vertices per second, based on current
	// assignment and elapsed times
	GetSpeed(worker msg.WorkerId) float64

	// The ratio of the fastest worker to the slowest
	FastestToSlowest() float64

	// The max Fastest to Slowest ratio allowed for optimal distributions.
	// If this is exceeded, isOptimal returns false
	IsOptimal() bool

	// Set all speeds to 0, to ready for next run
	ResetSpeeds()

	// Partitions the vertices, assigns the partitions based on the latest
	// performance, and updates assignments. Returns the updated assignment info.
	Redistribute() []Assignment

	// Find the worker for the give vertex. If the vertex has not yet been
	// assigned, the result is unspecified
	GetWorker(vertex msg.VertexId) msg.WorkerId

	// Replaces current worker and parition info with these
	// Returns the number of loaded vertices
	LoadAssignments(assigns []Assignment) int
}

func New(numVertices int) Manager {
	pm := new(performanceManager)
	pm.workers = make(map[msg.WorkerId]*workerStats)
	pm.numVertices = numVertices
	pm.distribution = new(distribution)

	return pm
}
