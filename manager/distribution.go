// Once partitioned, this keeps track of which partitions (and
// subsequently which vertices) are on each worker

package manager

import "project_c9f7_i5l8_o0p4_p0j8/msg"

/////////////////////////// distribution /////////////////////

type distribution struct {
	assignments   []msg.WorkerId
	partitionSize int
}

func (d *distribution) init(numPartitions int, partitionSize int) {
	d.assignments = make([]msg.WorkerId, numPartitions)
	d.partitionSize = partitionSize
}

func (d *distribution) setWorker(wid msg.WorkerId, part msg.Partition) {
	if d.partitionSize > 0 {
		pindex := int(part.Start()) / d.partitionSize
		d.assignments[pindex] = wid
	}
}

func (d *distribution) getPartitionWorker(part msg.Partition) msg.WorkerId {
	return d.getVertexWorker(part.Start())
}

func (d *distribution) getVertexWorker(vid msg.VertexId) msg.WorkerId {
	if d.partitionSize == 0 {
		return msg.WorkerId("")
	}
	pindex := vid.Int() / d.partitionSize
	return d.assignments[pindex]
}
