package main

import (
	"runtime"
	"fmt"
	"sort"
	"os"
	"bufio"
	"sync"
	"math"
	"math/rand"
	"strconv"

)

// https://stackoverflow.com/questions/29693708/sort-pair-in-golang
type pair struct {
	A int
	B int
}

type Pairs []pair

func (s Pairs) Len() int {
	return len(s)
}
func (s Pairs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Pairs) Less(i, j int) bool {
	return s[i].A > s[j].A
}

func degrees(gr [][] int, sz int, filename string, wg sync.WaitGroup) {
	res := make([]int, sz)
	for i := 0; i < sz; i++ {
		res[i] = len(gr[i])
	}
	ch := make(chan pair, 10)
	defer close(ch)
	go func() {
		saveToFile(ch, filename)
		wg.Done()
	}()
	for i, x := range res {
		ch <- pair{i, x}
	}
}

func queuePut(q[] int, u[] bool, e int, x int) ([] int, [] bool, int) {
	q[e] = x;
	e++;
	u[x] = true;
	return q, u, e
}

func queuePop(q[] int, b int) (int, int) {
	b++;
	var res int = q[b - 1];
	return res, b
}

func queueIsEmpty(b int, e int) bool {
	return b == e
}

// by Cormen, bfsTree
func bfsTree(gr[][] int, s int) [] int {
	n := len(gr)
	d := make([]int, n)
	q := make([]int, n)
	used := make([]bool, n)
	b, e := 0, 0

	put := func(x int) {
		q[e] = x;
		e++;
		used[x] = true
	}
	pop := func() int {
		b++;
		return q[b - 1]
	}
	empty := func() bool {
		return b == e
	}
	for i := 0; i < n; i++ {
		d[i] = -1;
	}
	put(s)
	d[s] = 0
	for !empty() {
		v := pop()
		for _, u := range gr[v] {
			if used[u] {
				continue
			}
			put(u)
			d[u] = d[v] + 1
		}
	}
	return d
}

// Calculate without workers
func shortPath(numberOfWorkers int, gr[][] int, n int, filename string) {
	result := make([][]int, n)
	for w := 0; w < numberOfWorkers; w++ {
		for v := 0; v < n; v++ {
			res := bfsTree(gr, v)
			for i := 0; i < len(res); i++ {
				result[v] = append(result[v], res[i])
			}
		}
	}
	var max int = 0
	var sum float64 = 0.0
	var cnt int = 0
	for i := 0; i < n; i++ {
		for j := 0; j < len(result); j++ {
			if max < result[i][j] {
				max = result[i][j]
			}
			if result[i][j] != -1 {
				sum += float64(result[i][j])
				cnt++;
			}

		}
	}

	fmt.Printf("[%s] diam=%d, avg dist=%.1f\n", filename, max, sum / float64(cnt))
}

func SccMax(gr[][] int, isRemove []bool) int {
	n := len(gr)
	buf := make([]int, n)
	used := make([]bool, n)
	b, e := 0, 0
	var max, curr int
	put := func(x int) {
		buf[e] = x
		e++
		used[x] = true
		curr++
		if curr > max {
			max = curr
		}
	}
	pop := func() int {
		b++; return buf[b - 1]
	}
	empty := func() bool {
		return b == e
	}
	for i := 0; i < n; i++ {
		if !used[i] && !isRemove[i] {
			curr = 0
			put(i)
			for !empty() {
				v := pop()
				for _, u := range gr[v] {
					if used[u] || isRemove[u] {
						continue
					}
					put(u)
				}
			}
		}
	}
	return max
}

func randomNodeDel(isRemove [] bool, sz int) [] bool {
	for {
		idx := rand.Intn(sz)
		if !isRemove[idx] {
			isRemove[idx] = true
			return isRemove
		}
	}
}

func attack(gr [][] int, n int, random int, wg sync.WaitGroup) {
	var attCnt [] int
	var attOk int = 0
	chAttack := make(chan pair, 10)
	go func() {
		saveToFile(chAttack, "attackNormal.dat")
		wg.Done()
	}()
	for i := 0; i < random; i++ {
		remove := make([]bool, n)
		var rcnt int = 0
		for n - rcnt > 10 {
			remove = randomNodeDel(remove, n)
			rcnt++
			max := SccMax(gr, remove)
			chAttack <- pair{rcnt, max}
			if float64(max) < math.Sqrt(float64(n - rcnt)) * 3 {
				attCnt = append(attCnt, rcnt)
				attOk++
				break
			}
		}
	}
	close(chAttack)
	fmt.Printf("%d out of %d trys of random brakedown succeeded\n", attOk, random)
	if attOk > 0 {
		fmt.Println("number of random removals to significantly damage (out of", n, "vertices:", attCnt)
	}
}

func attackMaxNode(gr [][] int, n int, wg sync.WaitGroup) {
	undGrByDeg := make([]pair, n)
	for i := 0; i < n; i++ {
		undGrByDeg[i] = pair{len(gr[i]), i}
	}
	sort.Sort(Pairs(undGrByDeg))
	ch := make(chan pair, 10)
	go func() {
		saveToFile(ch, "attackMaxNode.dat")
		wg.Done()
	}()
	var isOk bool
	var rcnt int = 0
	var pos int = 0
	remove := make([] bool, n)
	for n - rcnt > 10 {
		remove[pos] = true
		pos++
		rcnt++
		max := SccMax(gr, remove)
		ch <- pair{rcnt, max}
		if float64(max) < math.Sqrt(float64(n - rcnt)) * 3 {
			fmt.Printf("Attack is successfull, %d nodes removed\n", rcnt)
			isOk = true
			break
		}
	}
	close(ch)
	if !isOk {
		fmt.Println("Attack is failed")
	}
}

func saveToFile(ch chan pair, filename string) {
	file, error := os.Create(filename)
	if error != nil {
		panic(error)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	defer w.Flush()

	for it := range ch {
		_, err := fmt.Fprintf(w, "%d %d\n", it.A, it.B)
		if err != nil {
			panic(err)
		}
	}
}

func clusteringCC(gr [][] int, adj [][] bool, n int, wg sync.WaitGroup) {
	var sum float64 = 0.0
	for i := 0; i < n; i++ {
		var adjCount int = len(gr[i])
		var denom int = adjCount * (adjCount - 1)
		var sumNum int = 0;

		if (denom != 0) {
			for j := 0; j < len(gr[i]); j++ {
				for k := 0; k < len(gr[i]); k++ {
					if (j != k && adj[gr[i][j]][gr[i][k]]) {
						sumNum++
					}
				}
			}
			sum += (2 * (float64(sumNum)) / float64(denom))
			//if(sumNum !=0) {
			//	sum += (2 * float64(sumNum)) / float64(denom * (denom - 1) / 2)
			//}
		}
	}
	var all float64 = float64(1) / float64(n)
	var CC = all * sum;
	//var pEst float64 = (6.7 * math.Pow(10.0, 6.0)) / float64(n * (n - 1) / 2)
	fmt.Println("CCallG ", CC/*, " CCAv_esP ", pEst*/)

	ch := make(chan pairf, 10)
	go func() {
		savef(ch, "clustering.dat")
		wg.Done()
	}()
	for it := 0; it < n; it++ {
		sum := 0.0
		var h int = 0
		for i := 0; i < n; i++ {
			var adjCount int = len(gr[i])
			if (len(gr[i]) == it) {
				h++;
				var denom int = adjCount * (adjCount - 1)
				var sumNum int = 0;
				if (denom != 0) {
					for j := 0; j < len(gr[i]); j++ {
						for k := 0; k < len(gr[i]); k++ {
							if (j != k && adj[gr[i][j]][gr[i][k]]) {
								sumNum++
							}
						}
					}
					sum += ( 2 * (float64(sumNum)) / float64(denom))
					/*if(sumNum !=0) {
						sum += (2 * float64(sumNum)) / float64(denom * (denom - 1) / 2)
					}*/
				}
			}
		}
		if (h != 0) {
			all = float64(1) / float64(n)
			CC = all * sum;
			//pEst = (6.7 * math.Pow(10.0, 6.0)) / float64(h * (h - 1) / 2)
			ch <- pairf{it, CC}
			fmt.Println(it, " ", "CCallG ", CC /*, " CCAv_esP ", pEst*/)
		}
	}
	close(ch)

}

const periodlen = 5
const NoConv = 10000

func PR(g [][]int, damp float64) (result []float64, steps int) {
	n := len(g)
	prev, curr := make([]float64, n), make([]float64, n)
	bonus := damp / float64(n)
	popular := make([]int, n)

	//threshold := int(math.Sqrt(float64(n)))
	for i := 0; i < n; i++ {
		for _, u := range g[i] {
			popular[u]++

		}
	}
	p := 0
	const zz int = 0000000000000001
	for i := 0; i < n; i++ {
		if popular[i] >= zz {
			popular[p] = i
			p++
		}
	}
	popular = popular[:p]
	for i, initval := 0, 1 / float64(n); i < n; i++ {
		prev[i] = initval
	}

	calc := func(v int) {
		t := (1 - damp) * prev[v]
		if len(g[v]) == 0 {
			t /= float64(len(popular))
			for _, p := range popular {
				curr[p] += t
			}
			return
		}
		t /= float64(len(g[v]))
		for _, u := range g[v] {
			curr[u] += t
		}
	}

	var norms []float64

	for step := 1; step < NoConv; step++ {
		if step > 1 {
			sum := 0.
			prev, curr = curr, prev
			for i := 0; i < n; i++ {
				sum += curr[i]
				curr[i] = bonus
			}
			add := (1 - sum) / float64(n)
			for i := 0; i < n; i++ {
				curr[i] += add
			}
		}
		for i := 0; i < n; i++ {
			calc(i)
		}
		norm := 0.0
		for i := 0; i < n; i++ {
			norm += math.Abs(curr[i] - prev[i])
		}
		norms = append(norms, norm)
		if check(norms) && step > 10 {
			return curr, step
		}
		if len(norms) == 2 * periodlen {
			norms = norms[1:]
		}
	}
	return nil, -1
}

func check(norms []float64) bool {
	n := len(norms)
	if n < 2 {
		return false
	}
	min, max := minmax(norms)
	if max - min < 10e-5 {
		return true // all similar enough
	}
	if n < 4 {
		return false // lets do at least 5 steps to check for cycle
	}
	// lets cmp pairwise, check if there is period of len $periodlen or shorter
	for plen := 1; plen <= periodlen && plen * 2 < n; plen++ {
		t := norms[n - (2 * plen):]
		var mismatch bool
		for i := 0; i < plen; i++ {
			if math.Abs(t[i] - t[i + plen]) > 10e-5 {
				mismatch = true
				break
			}
		}
		if !mismatch {
			return true
		}
	}
	return false
}

func minmax(t []float64) (float64, float64) {
	min, max := t[0], t[0]
	for i := 1; i < len(t); i++ {
		if t[i] < min {
			min = t[i]
		} else if t[i] > max {
			max = t[i]
		}
	}
	return min, max
}

type pairf struct {
	a int
	b float64
}

func savef(ch chan pairf, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	for p := range ch {
		_, err := fmt.Fprintf(w, "%d %v\n", p.a, p.b)
		if err != nil {
			panic(err)
		}
	}
}

func PageRank(gr[][] int, damps []float64, wg sync.WaitGroup) {

	var xx int = 0
	for _, x := range damps {
		pr, steps := PR(gr, x)
		desc := "[page.Rank] "
		if x == 0 {
			desc = "without dampening"
		} else {
			desc = fmt.Sprintf("with damp factor=%.2f", x)
		}
		if pr == nil {
			desc += ": no convergence"
		} else {
			desc += fmt.Sprintf(": convergence reached after %d steps", steps)
		}
		fmt.Println(desc)

		ch := make(chan pairf, 10)
		defer close(ch)
		desc = fmt.Sprintf("%02.0f.pg", 100 * x)
		go func() {
			savef(ch, "page/desc" + strconv.Itoa(xx) +".dat")
			wg.Done()
		}()
		sort.Float64s(pr)
		for i, p := range pr {
			ch <- pairf{i + 1, p}
		}
		xx++
	}

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var n, m int
	fmt.Scanf("%d %d", &n, &m)
	fmt.Println("x = ", n, " y = ", m)

	grIn := make([][]int, n)
	grOut := make([][]int, n)
	grSum := make([][]int, n)
	AdjMatrix := make([][]bool, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			AdjMatrix[i] = append(AdjMatrix[i], false);
		}
	}

	for i := 0; i < m; i++ {
		var x, y int
		fmt.Scanf("%d %d", &x, &y)
		grOut[x] = append(grOut[x], y)
		grIn[y] = append(grIn[y], x)
		grSum[x] = append(grSum[x], y)
		grSum[y] = append(grSum[y], x)
		AdjMatrix[y][x] = true
		AdjMatrix[x][y] = true
	}
	//_ = grIn
	//_ = grOut
	//_ = grSum

	var wg sync.WaitGroup

	damps := []float64{0}
	for f := 0.1; f <= 1; f += .01 {
		damps = append(damps, f)
	}

	wg.Add(len(damps))
	wg.Add(3)
	wg.Add(2)
	degrees(grIn, n, "in.deg", wg)
	degrees(grOut, n, "out.deg", wg)
	degrees(grSum, n, "sum.deg", wg)
	shortPath(10, grIn, n, "inPaths.deg")//, wg)
	shortPath(1, grOut, n, "outPaths.deg")//, wg)
	shortPath(1, grSum, n, "sumPaths.deg")//, wg)
	attack(grSum, n, 10, wg)
	attackMaxNode(grSum, n, wg)
	clusteringCC(grSum, AdjMatrix, n, wg)
	PageRank(grOut, damps, wg);
}
