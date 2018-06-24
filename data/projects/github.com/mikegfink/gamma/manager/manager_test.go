// To compile for use with gdb (works best in go >1.5)
// -c: compiles to a .test executable file
// -gcflags "-N -l": removes optimizations (inlining and registerizations)
// go test -c -gcflags "-N -l"

package manager

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"sort"
	"testing"
	"time"

	"project_c9f7_i5l8_o0p4_p0j8/msg"
)

var t *testing.T

///// Test parameters
const numV int = 17443453

var workers []msg.WorkerId = []msg.WorkerId{"w0", "w1", "w2", "w3", "w4", "w5", "w6", "w7", "w8", "w9"}

var numW int = len(workers)

const fastToSlow float64 = 1.5

/////////////////

func TestInitialization(tee *testing.T) {
	t = tee

	m := New(numV)
	for _, w := range workers {
		m.AddWorker(w)
	}

	m.AddWorker("w10") // to be deleted
	test("NumWorkers", 11, m.NumWorkers())
	m.RemoveWorker("w10")
	test("RemoveWorker", 10, m.NumWorkers())

	// Workers returns the correct slice
	ws := m.Workers()
	test("len(Workers())", 10, len(ws))
	for _, w := range ws {
		test("Workers() contains", true, contains(ws, w))
	}

	// Initialization tests
	test("FastestToSlowest", 1.0, m.FastestToSlowest())
	test("IsOptimal", true, m.IsOptimal())

	test("PartitionSize", 0, m.PartitionSize())
	test("GetWorker", msg.WorkerId(""), m.GetWorker(0))
}

func TestEqualDistribution(tee *testing.T) {
	t = tee

	nV := numV * PartitionsPerWorkerAvg * numW

	m := New(nV)
	for _, w := range workers {
		m.AddWorker(w)
	}

	///// ASSIGNING VERTICES /////////////////////
	assigns := m.Redistribute()
	test("len([]Assignment)", numW, len(assigns))
	// Equally-divided
	vCount := 0
	realAvg := int(math.Ceil(float64(nV) / float64(numW*PartitionsPerWorkerAvg)))
	for _, a := range assigns {
		test("Equal num parts", PartitionsPerWorkerAvg, len(a.Partitions))

		for _, p := range a.Partitions {
			vCount += p.Size()
			test("equal num vert per part", realAvg, p.Size())
		}
	}
	test("equal div, all V accounted for", nV, vCount)

	v1 := 0
	wid1 := m.GetWorker(msg.VertexId(0))
	test(fmt.Sprintf("worker %v, vertex 0: contained?", wid1), true, contains(workers, wid1))

	v2 := 0 + m.PartitionSize() - 1
	wid2 := m.GetWorker(msg.VertexId(v2))
	test(fmt.Sprintf("wid for v %v == wid for v %v", v1, v2), wid1, wid2)

	//// Due to the implementation, this should be true, but there is no
	// rule that that consecutive partitions need to appear on the same worker
	// so if this fails, it may not be a problem
	v3 := 0 + m.PartitionSize()
	wid3 := m.GetWorker(msg.VertexId(v3))
	test(fmt.Sprintf("wid for v %v == wid for v %v", v1, v3), wid1, wid3)

	v4 := nV - 1
	wid4 := m.GetWorker(msg.VertexId(v4))
	test(fmt.Sprintf("worker %v, vertex %v: contained?", wid4, v4), true, contains(workers, wid4))

	v5 := nV
	wid5 := m.GetWorker(msg.VertexId(v5))
	test(fmt.Sprintf("worker %v, vertex %v = nV: contained?", wid5, v5), false, contains(workers, wid5))

	// Test LoadAssignments
	loadNumV := m.LoadAssignments(assigns)
	test("LoadAssignments: load size == num vertices", nV, loadNumV)
}

func TestUnequalDistribution(tee *testing.T) {
	t = tee

	m := New(numV)
	m.AddWorker("w0")
	m.AddWorker("w9")
	m.Redistribute()

	for _, w := range workers {
		// Re-adds are ignored
		m.AddWorker(w)
	}

	fast := seconds(1)
	slow := seconds(fastToSlow)

	////////////////// Set times ////////////////
	m.SetElapsedTime("w0", fast)
	// Should set both fastest and slowest to fast speed
	testF(0.001, "SetElapsedTimeFast", 1.0, m.FastestToSlowest())
	test("optimal", true, m.IsOptimal())
	// Should now set different slow speed
	m.SetElapsedTime("w9", slow)
	testF(0.1, "FasttoSlow", float64(slow)/float64(fast), m.FastestToSlowest())
	test("is optimal?", false, m.IsOptimal())

	// Now check uneven redistribution /////////////////
	assigns := m.Redistribute()
	test("len([]Assignment)", numW, len(assigns))
	sizes := []int{}

	fmt.Println("WorkerID (speed): partition.start() partition.end()")
	vCount := 0
	maxP := 0
	minP := math.MaxInt32
	for _, a := range assigns {
		numP := len(a.Partitions)
		for _, p := range a.Partitions {
			sizes = append(sizes, int(p.Size()))
			fmt.Printf("%v (%v): %v %v\n", a.Worker, m.GetSpeed(a.Worker), p.Start(), p.End())
			vCount += p.Size()
		}
		if numP > maxP {
			maxP = numP
		}
		if numP < minP {
			minP = numP
		}
	}
	totalP := 0
	for _, a := range assigns {
		numP := len(a.Partitions)
		if minP < numP && numP < maxP {
			totalP += numP
		}
	}
	midP := totalP / (numW - 2)

	test("uneven div, all V accounted for", numV, vCount)
	testF(1, fmt.Sprintf("MinP %v, MidP %v, MaxP %d, fast %d, slow %d", minP, midP, maxP, fast, slow),
		m.FastestToSlowest(), float64(maxP)/float64(minP))
	test(fmt.Sprintf("minP < midP < maxP: %d < %d < %d", minP, midP, maxP), true, minP < midP && midP < maxP)

	sort.Ints(sizes)
	fmt.Printf("Sizes:\n%v\n", sizes)
	fmt.Printf("MinP %v, MidP %v, MaxP %d, fast %d, slow %d, maxP/minP %v, slow/fast %v\n",
		minP, midP, maxP, fast, slow, float64(maxP)/float64(minP), float64(slow)/float64(fast))
}

func test(summary string, expect, actual interface{}) {
	if !reflect.DeepEqual(expect, actual) {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("Line %d:: %s: Expected %v, Actual %v", line, summary, expect, actual)
	}

}
func testF(precision float64, summary string, expect, actual float64) {
	diff := math.Abs(expect - actual)
	if diff > precision {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("Line %d:: %s: Expected %v, Actual %v", line, summary, expect, actual)
	}
}

func seconds(s float64) time.Duration {
	d := float64(time.Second) * s
	return time.Duration(d)
}

func contains(ws []msg.WorkerId, id msg.WorkerId) bool {
	for _, w := range ws {
		if w == id {
			return true
		}
	}
	return false
}
