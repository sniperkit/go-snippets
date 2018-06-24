// Implemntation of a concrete performance manager
//
// TODO incorporate historical data in some way, in case the last run was
// uncharacteristic. Perhaps a historical average
// Would require changing the "setElapsedTime" method and the "IsOptimal" method
//
///////////////////////////////////////////////////////////////////////////////

package manager

import (
	"math"
	"sort"
	"time"

	"project_c9f7_i5l8_o0p4_p0j8/msg"
)

// the same as the id of the first vertex in the partition
type partitionId msg.VertexId

type workerStats struct {
	id          msg.WorkerId
	latestTime  time.Duration
	speed       float64 // vertices per second
	partitions  []msg.Partition
	numVertices int
}

///////////////////////// performanceManager ////////////////////////

type performanceManager struct {
	workers      map[msg.WorkerId]*workerStats
	numVertices  int
	totalTime    time.Duration
	fastest      float64
	slowest      float64
	distribution *distribution
}

func (pm *performanceManager) AddWorker(worker msg.WorkerId) {
	if _, exists := pm.workers[worker]; !exists {
		w := new(workerStats)
		w.id = worker
		pm.workers[worker] = w
	}
}

func (pm *performanceManager) RemoveWorker(worker msg.WorkerId) {
	delete(pm.workers, worker)
}

func (pm *performanceManager) NumWorkers() int {
	return len(pm.workers)
}

func (pm *performanceManager) Workers() []msg.WorkerId {
	var ws []msg.WorkerId
	for w := range pm.workers {
		ws = append(ws, w)
	}
	return ws
}

func (pm *performanceManager) SetNumVertices(num int) {
	pm.numVertices = num
}

func (pm *performanceManager) updateSpeed(w *workerStats) {
	w.speed = float64(w.numVertices) / w.latestTime.Seconds()

	// Update extremes (if slowest is 0, it has never been set)
	if pm.slowest == 0 || w.speed < pm.slowest {
		pm.slowest = w.speed
	}
	if w.speed > pm.fastest {
		pm.fastest = w.speed
	}
}

func (pm *performanceManager) SetElapsedTime(worker msg.WorkerId, dur time.Duration) {
	// Ignore if an unrealistic time given
	if w, exists := pm.workers[worker]; dur > 0 && exists {
		pm.totalTime = pm.totalTime - w.latestTime + dur
		w.latestTime = dur

		if w.numVertices > 0 {
			pm.updateSpeed(w)
		}
	}
}

func (pm *performanceManager) GetSpeed(wid msg.WorkerId) float64 {
	if w, exists := pm.workers[wid]; exists {
		return w.speed
	}
	return 0
}

func (pm *performanceManager) FastestToSlowest() float64 {
	if pm.slowest == 0 {
		// Nothing has been set, so fastest and slowest are equal
		return 1
	}
	return pm.fastest / pm.slowest
}

func (pm *performanceManager) IsOptimal() bool {
	r := pm.FastestToSlowest()
	return r <= MaxOptimalRatio
}

func (pm *performanceManager) ResetSpeeds() {
	pm.totalTime = 0
	pm.fastest = 0
	pm.slowest = 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Also sets to average speed all new workers
func (pm *performanceManager) calculateTotalSpeed() (totalSpeed float64) {
	avgSpeed := 1.0 // temporarily initial value
	newWorkers := []*workerStats{}
	for _, w := range pm.workers {
		if w.speed == 0 {
			// must be new, so temporarily set speed to 1
			newWorkers = append(newWorkers, w)
			w.speed = avgSpeed
		} else {
			totalSpeed += w.speed
		}
	}

	if len(pm.workers) == len(newWorkers) {
		return avgSpeed * float64(len(newWorkers))
	}

	// set new worker speed to ACTUAL average speed
	avgSpeed = totalSpeed / float64(len(pm.workers)-len(newWorkers))
	for _, w := range newWorkers {
		w.speed = avgSpeed
		totalSpeed += avgSpeed
	}

	return totalSpeed

}

func (pm *performanceManager) sortedWorkers() []*workerStats {
	workers := make(bySpeed, len(pm.workers))
	i := 0
	for _, w := range pm.workers {
		workers[i] = w
		i++
	}
	sort.Sort(workers)
	return workers
}

func (pm *performanceManager) Redistribute() (assignsByWorker []Assignment) {
	totalNumParts := calculateNumParts(len(pm.workers))
	partitionSize := int(math.Ceil(float64(pm.numVertices) / float64(totalNumParts)))
	totalSpeed := pm.calculateTotalSpeed() // also sets avg speed for new workers
	pm.distribution.init(totalNumParts, partitionSize)

	// Since we round up, to ensure no vertices are lost, start with the
	// fastest workers, any empty partitions are on slowest worker(s)
	workers := pm.sortedWorkers()

	startVertex := 0
	for _, w := range workers {
		numParts := int(math.Ceil((w.speed / totalSpeed) * float64(totalNumParts)))
		w.numVertices = 0
		w.partitions = make([]msg.Partition, numParts)

		for i := range w.partitions {
			nextStartV := min(startVertex+pm.PartitionSize(), pm.numVertices)
			if nextStartV == startVertex {
				// We've assigned all the vertices earlier than expected
				w.partitions = w.partitions[:i]
				break
			}
			p := msg.NewPartition(startVertex, nextStartV)
			w.partitions[i] = p
			w.numVertices += p.Size()

			pm.distribution.setWorker(w.id, p)

			startVertex = nextStartV
		}
		assignsByWorker = append(assignsByWorker, Assignment{w.id, w.partitions})
	}
	return assignsByWorker
}

func (pm *performanceManager) PartitionSize() int {
	return pm.distribution.partitionSize
}

func (pm *performanceManager) GetWorker(vid msg.VertexId) msg.WorkerId {
	if vid.Int() >= pm.numVertices {
		return msg.WorkerId("")
	}
	return pm.distribution.getVertexWorker(vid)

}

func (pm *performanceManager) LoadAssignments(assigns []Assignment) int {
	pm.workers = make(map[msg.WorkerId]*workerStats)
	pm.distribution = new(distribution)

	// In case the number of workers has changed (there's been an undiscovered death)
	// we must findPartitionSize based on previous partitions
	totalNumParts := calculateNumParts(len(assigns))
	partitionSize := findPartitionSize(assigns)
	pm.distribution.init(totalNumParts, partitionSize)

	numVertices := 0

	for _, a := range assigns {
		for _, p := range a.Partitions {
			pm.distribution.setWorker(a.Worker, p)
			numVertices += p.Size()
		}
	}
	return numVertices
}

func findPartitionSize(assigns []Assignment) int {
	if len(assigns) == 0 {
		return 0
	}
	// In case one of the worker's had the tail end of the assignments,
	// we need to compare two partition sizes
	// (This doesn't work if dummy partition sizes of 0 are allowed)

	size1 := 0
	size2 := 0

	p := assigns[0].Partitions
	if len(p) > 0 {
		size1 = p[0].Size()
	}
	if len(p) > 1 {
		size2 = p[1].Size()
	} else if len(assigns) > 1 {
		p := assigns[1].Partitions
		if len(p) > 0 {
			size2 = p[0].Size()
		}
	}

	if size1 > size2 {
		return size1
	}
	return size2
}

func calculateNumParts(numWorkers int) (totalNumParts int) {
	return numWorkers * PartitionsPerWorkerAvg
}
