package vptree

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	//	"sync"
	"testing"
)

const (
	D = 8        // number of point dimensions
	N = 1 * 1024 // number of sample points
)

var (
	isamples []Point
	iqueries []Point
)

// 1-D test samples
var (
	set1D_12 = []Point{
		&point{0.0},
		&point{0.025},
		&point{0.1},
		&point{0.125},
		&point{0.15},
		&point{0.2},
		&point{0.25},
		&point{0.5},
		&point{0.75},
		&point{0.9},
		&point{0.95},
		&point{1.0},
	}
	set1D_24 = []Point{
		&point{0.00},
		&point{0.02},
		&point{0.08},
		&point{0.10},
		&point{0.22},
		&point{0.23},
		&point{0.25},
		&point{0.28},
		&point{0.29},
		&point{0.30},
		&point{0.33},
		&point{0.36},
		&point{0.40},
		&point{0.50},
		&point{0.51},
		&point{0.52},
		&point{0.70},
		&point{0.75},
		&point{0.80},
		&point{0.85},
		&point{0.90},
		&point{0.94},
		&point{0.98},
		&point{1.0},
	}
)

//// point /////////////////////////////////////////////////////////////////////

type point []float32

func (p point) String() string {
	txt := fmt.Sprintf("point(")
	last := len(p) - 1
	for i := range p[:last] {
		txt += fmt.Sprintf("%.3f ", p[i])
	}
	txt += fmt.Sprintf("%.3f)", p[last])
	return txt
}

// euclidean
func (p *point) Distance(to Point) float32 {
	pto := to.(*point)
	distance := float32(0.0)
	for i := range *p {
		d := (*p)[i] - (*pto)[i]
		distance += d * d
	}
	return float32(math.Sqrt(float64(distance)))
}

/*
// manhattan
func (p *point) Distance(to Point) float32 {
	pto := to.(*point)
	distance := float32(0.0)
	for i := range *p {
		d := (*p)[i] - (*pto)[i]
		if d < 0.0 {
			d = -d
		}
		distance += d
	}
	return distance
}
*/

//// initialisation ////////////////////////////////////////////////////////////

func init() {
	rand.Seed(4)

	// create random samples
	fmt.Printf("generating %d random %d-D test samples...\n", N, D)
	samples := make([]point, N)
	for i := range samples {
		samples[i] = make([]float32, D)
		for d := 0; d < D; d++ {
			samples[i][d] = rand.Float32()
		}
	}
	// make interfaced samples
	isamples = make([]Point, len(samples))
	for i := range isamples {
		isamples[i] = &samples[i]
	}

	// create random queries
	fmt.Printf("generating %d random %d-D test queries...\n", N, D)
	queries := make([]point, N)
	for i := range queries {
		queries[i] = make([]float32, D)
		for d := 0; d < D; d++ {
			queries[i][d] = rand.Float32()
		}
	}
	// make interfaced queries
	iqueries = make([]Point, len(queries))
	for i := range iqueries {
		iqueries[i] = &queries[i]
	}
}

//// testing ///////////////////////////////////////////////////////////////////

/*
func TestDemo(t *testing.T) {
	for i := 0; i < 10; i++ {
		vp := NewVPTree(set1D_24, uint(i))
		info := vp.Info()

		fmt.Printf("\n%s\n", vp)
		fmt.Printf("%d nodes, %d leaves, max depth %d\n",
			info.NNodes, info.NLeaves, info.MaxDepth)

		query := Point(&point{0.666})
		r, d := vp.KNN(query, 1)
		//r, d := vp.ApproximateKNN(query, 1, 3)
		fmt.Printf("r %v, d %v\n", r, d)
	}
}
*/

/*
func TestDataset(t *testing.T) {
	// load sample points
	samples, err := ParsePoints("testpoints-16")
	if err != nil {
		t.Error("%s", err)
	}
	// make interfaced samples
	isamples = make([]Point, len(samples))
	for i := range isamples {
		isamples[i] = &samples[i]
	}

	// load query points
	queries, err := ParsePoints("testpoints-16")
	if err != nil {
		t.Error("%s", err)
	}
	// make interfaced queries
	iqueries = make([]Point, len(queries))
	for i := range iqueries {
		iqueries[i] = &queries[i]
	}

	vp := NewVPTree(isamples, 123)

	// nearest neighbor queries
	K := 3
	for i := range iqueries {
		fmt.Printf("query %s\n", iqueries[i])
		res, dists := vp.KNN(iqueries[i], K)
		for r := range res {
			p := res[r].(*point)
			fmt.Printf("near %s : %.3f\n", p, dists[r])
		}
		fmt.Println()
	}

	// range queries
	R := float32(1.0)
	for i := range iqueries {
		fmt.Printf("query %s\n", iqueries[i])
		res, dists := vp.RangeNN(iqueries[i], R)
		for r := range res {
			p := res[r].(*point)
			fmt.Printf("near %s : %.3f\n", p, dists[r])
		}
		fmt.Println()
	}
}
*/

// how does an empty tree behave:
// zero results should be returned
func TestEmpty(t *testing.T) {
	query := Point(&point{0.0})
	vp := NewVPTree(nil, 1)

	nn, dists := vp.KNN(query, 1)
	if len(nn) != 0 || len(dists) != 0 {
		t.Errorf("%d results, should have been 0", len(nn))
	}
}

// what happens if k is greater than the number of samples:
// k should be limited to the number of samples
func TestKgreaterthanN(t *testing.T) {
	query := Point(&point{0.0, 0.0, 0.0})
	samples := []Point{
		&point{0.0, 1.0, 2.0},
		&point{1.0, 2.0, 3.0},
		&point{2.0, 3.0, 4.0},
	}
	vp := NewVPTree(samples, 1)

	nn, dists := vp.KNN(query, 10)
	if len(nn) != len(samples) || len(dists) != len(samples) {
		t.Errorf("%d results, should have been %d", len(nn), len(samples))
	}
}

// what happens in the presence of ties in the kNN results:
// exactly k results should be returned
func TestTies(t *testing.T) {
	query := Point(&point{0.0, 0.0})
	samples := []Point{
		&point{1.0, 0.0},
		&point{0.0, 1.0},
	}
	vp := NewVPTree(samples, 1)

	nn, dists := vp.KNN(query, 1)
	if len(nn) != 1 || len(dists) != 1 {
		t.Errorf("%d results, should have been 1", len(nn))
	}
}

// what happens if there are duplicate sample points:
// they behave as if they were unique
func TestDuplicates(t *testing.T) {
	query := Point(&point{0.0})
	samples := []Point{
		&point{1.0}, &point{2.0}, &point{2.0}, &point{3.0},
	}
	vp := NewVPTree(samples, 1)

	k := len(samples)
	nn, dists := vp.KNN(query, k)
	if len(nn) != k || len(dists) != k {
		t.Errorf("%d results, should have been %d", len(nn), k)
	}
}

// test NN on random samples
func TestRandomNN(t *testing.T) {
	vp := NewVPTree(isamples, 1)

	for i := range iqueries {
		nnRef, distsRef := kNN_linear(isamples, iqueries[i], 1)
		nn, dist := vp.NN(iqueries[i])

		if !matchResults(nnRef, []Point{nn}, distsRef, []float32{dist}) {
			t.Error("VPtree.KNN() and reference results don't match")
		}
	}
}

// test kNN on random samples
func TestRandomKNN(t *testing.T) {
	vp := NewVPTree(isamples, 1)
	k := 10

	for i := range iqueries {
		nnRef, distsRef := kNN_linear(isamples, iqueries[i], k)
		nn, dists := vp.KNN(iqueries[i], k)

		if !matchResults(nnRef, nn, distsRef, dists) {
			t.Error("VPtree.KNN() and reference results don't match")
		}
	}
}

// test RangeNN on random samples
func TestRandomRangeNN(t *testing.T) {
	vp := NewVPTree(isamples, 1)
	r := float32(1.0)

	for i := range iqueries {
		nnRef, distsRef := rangeNN_linear(isamples, iqueries[i], r)
		nn, dists := vp.RangeNN(iqueries[i], r)

		if !matchResults(nnRef, nn, distsRef, dists) {
			t.Error("VPtree.RangeNN() and reference results don't match")
		}
	}
}

/*
// TODO: use hardcoded data
func TestKNN_linear(t *testing.T) {
	fmt.Printf("kNN_linear()...\n")
	knn, _ := kNN_linear(isamples, iqueries[0], 3)
	fmt.Printf("knn: %v\n", knn)
}

// TODO: use hardcoded data
func TestRangeNN_linear(t *testing.T) {
	fmt.Printf("rangeNN_linear()...\n")
	knn, _ := rangeNN_linear(isamples, iqueries[0], 1.0)
	fmt.Printf("nn: %v\n", knn)
}
*/

//// benchmarks ////////////////////////////////////////////////////////////////

func BenchmarkNothing(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}

func BenchmarkNewVPTree_0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewVPTree(isamples, 0)
	}
}

func BenchmarkNewVPTree_1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewVPTree(isamples, 1)
	}
}

func BenchmarkNewVPTree_4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewVPTree(isamples, 4)
	}
}

func BenchmarkNewVPTree_16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewVPTree(isamples, 16)
	}
}

func BenchmarkNewVPTree_64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewVPTree(isamples, 64)
	}
}

func Benchmark1NN(b *testing.B) {
	vp := NewVPTree(isamples, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for q := range iqueries {
			_, _ = vp.KNN(iqueries[q], 1)
		}
	}
}

func Benchmark2NN(b *testing.B) {
	vp := NewVPTree(isamples, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for q := range iqueries {
			_, _ = vp.KNN(iqueries[q], 2)
		}
	}
}

func Benchmark4NN(b *testing.B) {
	vp := NewVPTree(isamples, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for q := range iqueries {
			_, _ = vp.KNN(iqueries[q], 4)
		}
	}
}

func Benchmark8NN(b *testing.B) {
	vp := NewVPTree(isamples, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for q := range iqueries {
			_, _ = vp.KNN(iqueries[q], 8)
		}
	}
}

func BenchmarkNN(b *testing.B) {
	vp := NewVPTree(isamples, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for q := range iqueries {
			_, _ = vp.NN(iqueries[q])
		}
	}
}

func BenchmarkRange_01(b *testing.B) {
	vp := NewVPTree(isamples, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for q := range iqueries {
			_, _ = vp.RangeNN(iqueries[q], 0.1)
		}
	}
}

func BenchmarkRange_02(b *testing.B) {
	vp := NewVPTree(isamples, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for q := range iqueries {
			_, _ = vp.RangeNN(iqueries[q], 0.2)
		}
	}
}

func BenchmarkRange_04(b *testing.B) {
	vp := NewVPTree(isamples, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for q := range iqueries {
			_, _ = vp.RangeNN(iqueries[q], 0.4)
		}
	}
}

func BenchmarkRange_08(b *testing.B) {
	vp := NewVPTree(isamples, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for q := range iqueries {
			_, _ = vp.RangeNN(iqueries[q], 0.8)
		}
	}
}

func BenchmarkApproximate_8(b *testing.B) {
	vp := NewVPTree(isamples, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for q := range iqueries {
			_, _ = vp.ApproximateKNN(iqueries[q], 1, 8)
		}
	}
}

func BenchmarkApproximate_16(b *testing.B) {
	vp := NewVPTree(isamples, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for q := range iqueries {
			_, _ = vp.ApproximateKNN(iqueries[q], 1, 16)
		}
	}
}

func BenchmarkApproximate_32(b *testing.B) {
	vp := NewVPTree(isamples, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for q := range iqueries {
			_, _ = vp.ApproximateKNN(iqueries[q], 1, 32)
		}
	}
}

func BenchmarkApproximate_64(b *testing.B) {
	vp := NewVPTree(isamples, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for q := range iqueries {
			_, _ = vp.ApproximateKNN(iqueries[q], 1, 64)
		}
	}
}

//// utilities /////////////////////////////////////////////////////////////////

func matchResults(pA, pB []Point, dA, dB []float32) bool {
	// match points
	if len(pA) != len(pB) {
		return false
	}
	for i := range pA {
		ppA := pA[i].(*point)
		ppB := pB[i].(*point)
		if ppA != ppB {
			// Accept different points if they have identical distances, this
			// is a benign hint of sort order differences in the case of ties.
			if dA[i] == dB[i] {
				return true
			} else {
				return false
			}
		}
	}

	// match distances
	if len(dA) != len(dB) {
		return false
	}
	for i := range dA {
		if dA[i] != dB[i] {
			//fmt.Printf("d[%d] mismatch: %.3f - %.3f\n", i, dA[i], dB[i])
			return false
		}
	}

	return true
}

// ParsePoints returns a slice of points contained in the given input file.
// A point (with arbitrary dimensions) is formatted as a line with comma
// separated floating point values. All points in the file should have an
// equal number of dimensions.
func ParsePoints(name string) ([]point, error) {
	points := []point{}
	nDim := 0 // number of dimensions of the input points

	file, err := os.Open(name)
	if err != nil {
		return points, fmt.Errorf("os.Open(): %s", err)
	}
	defer file.Close()

	// scan lines
	nLines := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue // skip empty or comment lines
		}
		fields := strings.FieldsFunc(line,
			func(r rune) bool { return r == ',' })

		// check number of dimensions of point
		if nLines == 0 { // first point establishes the reference
			nDim = len(fields)
		} else {
			if len(fields) != nDim {
				return points,
					fmt.Errorf("line %d dimension mismatch", nLines)
			}
		}

		// parse point coordinates
		p := point{}
		for i := range fields {
			field := strings.TrimSpace(fields[i])
			v, err := strconv.ParseFloat(field, 32)
			if err != nil {
				return points,
					fmt.Errorf("strconv.ParseFloat(): %s (line %d)",
						err, nLines)
			}
			p = append(p, float32(v))
		}

		points = append(points, p)
		nLines++
	}
	if err := scanner.Err(); err != nil {
		return points, fmt.Errorf("scanner.Scan(): %s", err)
	}

	return points, nil
}

// reference: linear k nearest neighbor search
func kNN_linear(
	samples []Point, query Point, k int) ([]Point, []float32) {

	if k > len(samples) {
		fmt.Fprintf(os.Stderr, "reducing k from %d to %d\n", k, len(samples))
		k = len(samples)
	}

	pq := make(priorityQueue, 0, k)
	// push the k first samples onto the priority queue
	for i := 0; i < k; i++ {
		d := query.Distance(samples[i])
		heap.Push(&pq, &_Item{samples[i], d})
	}
	// find the k samples with the smallest distances
	for i := k; i < len(samples); i++ {
		d := query.Distance(samples[i])
		if d < pq.Top().(*_Item).d {
			heap.Pop(&pq)
			heap.Push(&pq, &_Item{samples[i], d})
		}
	}
	return extractReverse(pq)
}

// reference: linear range search
func rangeNN_linear(
	samples []Point, query Point, r float32) ([]Point, []float32) {

	pq := make(priorityQueue, 0, 8)
	// push all samples with distance <= r onto the priority queue
	for i := range samples {
		d := query.Distance(samples[i])
		if d <= r {
			heap.Push(&pq, &_Item{samples[i], d})
		}
	}
	return extractReverse(pq)
}
