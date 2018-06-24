package plugin

import (
	"container/heap"
	"fmt"
	"log"

	"github.com/hiro4bbh/sticker"
)

// QueuePrioritizedByFloat32Item is the item structure of QueuePrioritizedByFloat32.
type QueuePrioritizedByFloat32Item struct {
	priority float32
	item     interface{}
}

// NewQueuePrioritizedByFloat32Item returns a new QueuePrioritizedByFloat32Item.
func NewQueuePrioritizedByFloat32Item(priority float32, item interface{}) QueuePrioritizedByFloat32Item {
	return QueuePrioritizedByFloat32Item{
		priority: priority,
		item:     item,
	}
}

// Item returns item.
func (item QueuePrioritizedByFloat32Item) Item() interface{} {
	return item.item
}

// Priority returns priority.
func (item QueuePrioritizedByFloat32Item) Priority() float32 {
	return item.priority
}

// QueuePrioritizedByFloat32 is the queue prioritized by float32.
// This implements interface heap.Interface.
type QueuePrioritizedByFloat32 []QueuePrioritizedByFloat32Item

// Len is for interface heap.Interface.
func (q QueuePrioritizedByFloat32) Len() int {
	return len(q)
}

// Less is for interface heap.Interface.
func (q QueuePrioritizedByFloat32) Less(i, j int) bool {
	return q[i].priority > q[j].priority
}

// Push is for interface heap.Interface.
func (q *QueuePrioritizedByFloat32) Push(item interface{}) {
	*q = append(*q, item.(QueuePrioritizedByFloat32Item))
}

// Pop is for interface heap.Interface.
func (q *QueuePrioritizedByFloat32) Pop() interface{} {
	item := (*q)[len(*q)-1]
	*q = (*q)[:len(*q)-1]
	return item
}

// Swap is for interface heap.Interface.
func (q QueuePrioritizedByFloat32) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

// Painter_TopLabels is the painter returns the most frequent top-K mis-classified labels.
//
// This is registered to Painters.
func Painter_TopLabels(ds *sticker.Dataset, Z []sticker.KeyValues32, K uint, debug *log.Logger) []uint32 {
	labelFreqs := make(map[uint32]float32)
	for i, zi := range Z {
		yi := ds.Y[i]
		maxNegZi := float32(0.0)
		for j := len(yi); j < len(zi); j++ {
			if zil := zi[j].Value; maxNegZi < zil {
				maxNegZi = zil
			}
		}
		for j, label := range yi {
			zil := zi[j].Value
			if zil <= maxNegZi {
				labelFreqs[label]++
			}
		}
	}
	filterK := 4 * K
	if filterK > uint(len(labelFreqs)) {
		filterK = uint(len(labelFreqs))
	}
	labelsTopK := sticker.RankTopK(labelFreqs, filterK)
	if debug != nil {
		msg := "Painter(topLabels): 1st order label frequencies:"
		for k := uint(0); k < filterK; k++ {
			msg += fmt.Sprintf(" %d:%g", labelsTopK[k], labelFreqs[labelsTopK[k]])
		}
		debug.Print(msg)
	}
	if K > uint(len(labelFreqs)) {
		K = uint(len(labelFreqs))
	}
	return labelsTopK[:K]
}

// Painter_TopLabelSubSet is the painter returns the most frequent top-K mis-classified co-occurring labels.
//
// This is registered to Painters.
func Painter_TopLabelSubSet(ds *sticker.Dataset, Z []sticker.KeyValues32, K uint, debug *log.Logger) []uint32 {
	if K == 0 {
		return []uint32{}
	}
	maxNegZ := make([]float32, len(Z))
	labelFreqs := make(map[uint32]float32)
	for i, zi := range Z {
		yi := ds.Y[i]
		maxNegZi := float32(0.0)
		for j := len(yi); j < len(zi); j++ {
			if zil := zi[j].Value; maxNegZi < zil {
				maxNegZi = zil
			}
		}
		maxNegZ[i] = maxNegZi
		for j, label := range yi {
			zil := zi[j].Value
			if zil <= maxNegZi {
				labelFreqs[label]++
			}
		}
	}
	filterK := 4 * K
	if filterK > uint(len(labelFreqs)) {
		filterK = uint(len(labelFreqs))
	}
	labelsTopK := sticker.RankTopK(labelFreqs, filterK)
	if debug != nil {
		msg := "Painter(topLabelSubSet): 1st order label frequencies:"
		for k := uint(0); k < filterK; k++ {
			msg += fmt.Sprintf(" %d:%g", labelsTopK[k], labelFreqs[labelsTopK[k]])
		}
		debug.Print(msg)
	}
	labelsTopK = labelsTopK[:filterK]
	queue := make(QueuePrioritizedByFloat32, filterK)
	for r, label := range labelsTopK {
		queue[r] = NewQueuePrioritizedByFloat32Item(labelFreqs[label], map[uint32]bool{label: true})
	}
	heap.Init(&queue)
	var labelMap map[uint32]bool
	for {
		item := heap.Pop(&queue).(QueuePrioritizedByFloat32Item)
		labelMap = item.Item().(map[uint32]bool)
		labelFreqs = make(map[uint32]float32)
		if uint(len(labelMap)) == K {
			break
		}
		for i, zi := range Z {
			yi := ds.Y[i]
			ncontaineds := 0
			for _, label := range yi {
				if _, ok := labelMap[label]; ok {
					ncontaineds++
				}
			}
			if ncontaineds < len(labelMap) {
				continue
			}
			maxNegZi := maxNegZ[i]
			for j, label := range yi {
				if _, ok := labelMap[label]; !ok {
					zil := zi[j].Value
					if zil <= maxNegZi {
						labelFreqs[label]++
					}
				}
			}
		}
		if len(labelFreqs) == 0 {
			break
		}
		subK := K
		if subK > uint(len(labelFreqs)) {
			subK = uint(len(labelFreqs))
		}
		labelsTopK := sticker.RankTopK(labelFreqs, subK)
		for _, label := range labelsTopK {
			labelMapNew := make(map[uint32]bool, len(labelMap))
			for l := range labelMap {
				labelMapNew[l] = true
			}
			labelMapNew[label] = true
			heap.Push(&queue, NewQueuePrioritizedByFloat32Item(labelFreqs[label], labelMapNew))
		}
		if filterK < uint(len(queue)) {
			queueOld := queue
			queue = make(QueuePrioritizedByFloat32, filterK)
			for r := uint(0); r < filterK; r++ {
				queue[r] = heap.Pop(&queueOld).(QueuePrioritizedByFloat32Item)
			}
			heap.Init(&queue)
		}
	}
	labels := make([]uint32, 0, len(labelMap))
	for label := range labelMap {
		labels = append(labels, label)
	}
	return labels
}

// Painters is the map from the painter's name to the painter function.
// A painter takes the dataset, the corresponding margin matrix, and the maximum number of requested labels, then returns the slice of selected labels.
// debug is used for debug logs.
var Painters = map[string]func(ds *sticker.Dataset, Z []sticker.KeyValues32, K uint, debug *log.Logger) []uint32{
	"topLabels":      Painter_TopLabels,
	"topLabelSubSet": Painter_TopLabelSubSet,
}
