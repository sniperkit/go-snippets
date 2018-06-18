package tp

import (
	"fmt"
	"math"
	//"time"
)

const (
	// default EPSILON value for checking if a value is zero or not
	EPSILON = float64(1e-6)

	// default max iterations for optimizing the solution
	MAX_ITER = 100
)

type Problem struct {

	// 'static' variables
	epsilon, infinity float64
	maxIter           int

	// inputs, could be adjusted if supply/demand is unbalanced
	supply     []float64
	demand     []float64
	costMatrix [][]float64

	// balance flag
	//  -1: supply < demand
	//   0: supply == demand
	//   1: supply > demand
	balanced int

	// supply, demand size after adjustment if inputs are unbalanced
	sLen, dLen int

	// total amount to tranport from producers to consumers
	// including the dummy amount if inputs are unbalanced
	quatity float64

	// iteration count
	iterCnt int

	// u, v matrix
	u, v []float64

	// optimization starting cell
	row, col int

	// When computing u,v
	// 0: init state (not scanned)
	// 1: scan candidate
	// 2: already scanned
	// -------------------
	// When finding the loop
	// 0: row/col is available
	// 1: row/col is occupied
	rowFlags, colFlags []int

	// loop (link-list) head
	loop *cell

	// solution flow
	flow [][]*flowcell
}

// to solve degeneracy, use a struct to indicate
// if it is a basic variable with value 0
type flowcell struct {
	basic bool
	value float64
}

type cell struct {
	// cell location
	row, col int

	// flag to
	//  1) mark direction to next cell
	//  2) mark if this is the odd/even cell in chain/loop
	// value:
	//  true:  horizontal/odd
	//  false: vertical/even
	flag bool

	// minimum flow value of the even cell in the loop/chain
	loopEvenMinFlow float64

	// pointers to prev/next cell
	prev, next *cell
}

func (c *cell) String() string {
	direction, sign := "V", "-"
	if c.flag { // horizontal, odd-cell
		direction, sign = "H", "+"
	}
	return fmt.Sprintf("(%v,%v)/%v/%v", c.row, c.col, direction, sign)
}

func createProblem(s, d []float64, c [][]float64, maxIter int, epsilon float64) (*Problem, error) {
	sLen := len(s)
	if sLen < 1 {
		return nil, fmt.Errorf("not enough producers, need at least 1!")
	}
	dLen := len(d)
	if dLen < 1 {
		return nil, fmt.Errorf("not enough consumers, need at least 1! ")
	}
	if sLen != len(c) {
		return nil, fmt.Errorf("producer count doesn't match 1st dimension length of costMatrix!")
	}
	if dLen != len(c[0]) {
		return nil, fmt.Errorf("consumer count doesn't match 2nd dimension length of costMatrix!")
	}

	var sSum, dSum, quatity float64
	for i := 0; i < sLen; i++ {
		if s[i] < epsilon {
			return nil, fmt.Errorf("supply[%v]=%v is too small (<%v)!", i, s[i], epsilon)
		}
		sSum += s[i]
	}
	for i := 0; i < dLen; i++ {
		if d[i] < epsilon {
			return nil, fmt.Errorf("demand[%v]=%v is too small (<%v)!", i, d[i], epsilon)
		}
		dSum += d[i]
	}
	diff := float64(0)
	balanced := 0
	if sSum > dSum {
		diff = sSum - dSum
		quatity = sSum
		balanced = 1
	} else if sSum < dSum {
		diff = dSum - sSum
		quatity = dSum
		balanced = -1
	} else {
		quatity = sSum
		balanced = 0
	}
	var supply, demand []float64
	var costMatrix [][]float64
	if diff > epsilon { // unbalanced
		if sSum > dSum {
			// copy supply
			supply = make([]float64, sLen)
			for i := 0; i < sLen; i++ {
				supply[i] = s[i]
			}
			// copy demand and add one more (diff)
			demand = make([]float64, dLen+1)
			for i := 0; i < dLen; i++ {
				demand[i] = d[i]
			}
			demand[dLen] = diff
			// copy cost matrix, append one column (0s)
			costMatrix = make([][]float64, sLen)
			for i := 0; i < sLen; i++ {
				costMatrix[i] = make([]float64, dLen+1)
				for j := 0; j < dLen; j++ {
					costMatrix[i][j] = c[i][j]
				}
				costMatrix[i][dLen] = 0
			}
			// fix demand size
			dLen += 1
		} else { // sSum < dSum
			// copy supply and add one more (diff)
			supply = make([]float64, sLen+1)
			for i := 0; i < sLen; i++ {
				supply[i] = s[i]
			}
			supply[sLen] = diff
			// copy demand
			demand = make([]float64, dLen)
			for i := 0; i < dLen; i++ {
				demand[i] = d[i]
			}
			// copy cost matrix and append one row (0s)
			costMatrix = make([][]float64, sLen+1)
			for i := 0; i < sLen; i++ {
				costMatrix[i] = make([]float64, dLen)
				for j := 0; j < dLen; j++ {
					costMatrix[i][j] = c[i][j]
				}
			}
			costMatrix[sLen] = make([]float64, dLen)
			// fix supply size
			sLen += 1
		}
	} else { // balanced
		supply = make([]float64, sLen)
		for i := 0; i < sLen; i++ {
			supply[i] = s[i]
		}
		demand = make([]float64, dLen)
		for i := 0; i < dLen; i++ {
			demand[i] = d[i]
		}
		costMatrix = make([][]float64, sLen)
		for i := 0; i < sLen; i++ {
			costMatrix[i] = make([]float64, dLen)
			for j := 0; j < dLen; j++ {
				costMatrix[i][j] = c[i][j]
			}
		}
	}

	// create flow slice
	flow := make([][]*flowcell, sLen)
	for i := 0; i < sLen; i++ {
		flow[i] = make([]*flowcell, dLen)
		for j := 0; j < dLen; j++ {
			flow[i][j] = &flowcell{}
		}
	}
	// create Problem struct
	return &Problem{
		epsilon:  epsilon,
		infinity: math.Inf(1),
		maxIter:  maxIter,

		supply:     supply,
		demand:     demand,
		costMatrix: costMatrix,
		balanced:   balanced,

		sLen:     sLen,
		dLen:     dLen,
		quatity:  quatity,
		iterCnt:  0,
		u:        make([]float64, sLen),
		v:        make([]float64, dLen),
		row:      -1,
		col:      -1,
		rowFlags: make([]int, sLen),
		colFlags: make([]int, dLen),
		loop:     nil,

		flow: flow,
	}, nil
}

func (es *Problem) printSolution() {
	fmt.Println("[Solution]")
	cost := float64(0)
	for i := 0; i < es.sLen; i++ {
		for j := 0; j < es.dLen; j++ {
			if !es.flow[i][j].basic {
				continue
			}
			//fmt.Printf("i=%v,j=%v,supply=%v,demand=%v,transport=%v\n",
			//i, j, es.supply[i], es.demand[j], es.flow[i][j])
			c := es.costMatrix[i][j]
			fval := es.flow[i][j].value
			//if fval >= 0 {
			fmt.Printf(" (%v,%v),cost=%v,flow=%v\n", i, j, c, fval)
			cost += c * fval
			//}
		}
	}
	fmt.Printf("cost=%v\n", cost)
	fmt.Println("")
}

// find the initial solution with the "Minimal Cost" method
func (es *Problem) findFeasibleSolution() int {
	//fmt.Println("[Finding feasible solution ...]")
	//t1 := time.Now()

	sLen, dLen, epsilon, infinity := es.sLen, es.dLen, es.epsilon, es.infinity

	quatity := es.quatity
	//fmt.Printf("quatity=%v\n", quatity)
	flowCnt := 0

	for {
		// the least-cost (selected) row/column index
		si, sj := -1, -1
		var minCost = infinity
		// loop to find the least cost row/column
		for i := 0; i < sLen; i++ {
			// skip row if supply is "0"
			if es.supply[i] <= 0 {
				continue
			}
			for j := 0; j < dLen; j++ {
				// skip column if demand is "0"
				if es.demand[j] <= 0 {
					continue
				}
				cost := es.costMatrix[i][j]
				if cost < minCost {
					si, sj, minCost = i, j, cost
				} else if cost == minCost {
					// for same cost cell, choose the one which
					// transports more
					sq := math.Min(es.supply[si], es.demand[sj])
					q := math.Min(es.supply[i], es.demand[j])
					if q > sq {
						si, sj = i, j
					}
				}
			}
		}
		// substract the selected quatity from supply/demand
		s := es.supply[si]
		d := es.demand[sj]
		diff := s - d
		q := float64(0)
		if diff > epsilon { // s > d
			q = d
			es.supply[si] = diff
			es.demand[sj] = 0
		} else if diff < -epsilon { // s < d
			q = s
			es.supply[si] = 0
			es.demand[sj] = -diff
		} else { // s == d
			q = s
			es.supply[si] = 0
			es.demand[sj] = 0
		}
		// remove flow value from total quatity
		quatity = quatity - q
		//fmt.Printf("got min cost cell at (%v,%v)/%v, flow=%v, left=%v\n", si, sj, es.costMatrix[si][sj], q, quatity)
		// set basic variable
		fc := es.flow[si][sj]
		fc.basic = true
		fc.value = q
		flowCnt += 1

		if quatity <= epsilon {
			break
		}
	}

	//fmt.Printf("findFeasibleSolution() done in %v\n", time.Now().Sub(t1))
	//fmt.Println("")
	return flowCnt
}

func (es *Problem) computeUV() error {
	//fmt.Println("[Computing U,V ...]")
	//t1 := time.Now()

	sLen, dLen := es.sLen, es.dLen

	// reset row/col flags
	for i := 0; i < sLen; i++ {
		es.rowFlags[i] = 0
	}
	for i := 0; i < dLen; i++ {
		es.colFlags[i] = 0
	}
	// computed counts
	uComputedCnt, vComputedCnt := 0, 0
	// set u[0] = 0
	es.u[0] = float64(0)
	es.rowFlags[0] = 1
	uComputedCnt += 1
	//fmt.Println("set u[0]=0")
	//fmt.Println("")

	more2scan := false

	for {
		// scan rows
		more2scan = false
		for row := 0; row < sLen; row++ {
			if es.rowFlags[row] != 1 {
				continue
			}
			//fmt.Printf(">> scanning row #%v ...\n", row)
			for col := 0; col < dLen; col++ {
				if !es.flow[row][col].basic || es.colFlags[col] > 1 {
					//fmt.Printf(">>>> skipping col #%v.\n", col)
					continue
				}
				v := es.costMatrix[row][col] - es.u[row]
				es.v[col] = v
				vComputedCnt += 1
				//fmt.Printf(">>>> computed v[%v]=%v\n", col, v)
				if es.colFlags[col] == 0 {
					es.colFlags[col] = 1
					more2scan = true
					//fmt.Printf(">>>> set colFlags[%v]=1\n", col)
				}
			}
			es.rowFlags[row] = 2
			//fmt.Printf(">> set rowFlags[%v]=2\n", row)
			//fmt.Println("")
		}

		//fmt.Printf("rowFlags=%v, colFlags=%v, uComputedCnt=%v, vComputedCnt=%v\n", es.rowFlags, es.colFlags, uComputedCnt, vComputedCnt)
		//fmt.Println("")
		if !more2scan || (uComputedCnt == sLen && vComputedCnt == dLen) {
			//fmt.Println("break!")
			break
		}

		// scan columns
		more2scan = false
		for col := 0; col < dLen; col++ {
			if es.colFlags[col] != 1 {
				continue
			}
			//fmt.Printf(">> scanning col #%v ...\n", col)
			for row := 0; row < sLen; row++ {
				if !es.flow[row][col].basic || es.rowFlags[row] > 1 {
					//fmt.Printf(">>>> skipping row #%v.\n", row)
					continue
				}
				u := es.costMatrix[row][col] - es.v[col]
				es.u[row] = u
				uComputedCnt += 1
				//fmt.Printf(">>>> computed u[%v]=%v\n", row, u)
				if es.rowFlags[row] == 0 {
					es.rowFlags[row] = 1
					more2scan = true
					//fmt.Printf(">>>> set rowFlags[%v]=1\n", row)
				}
			}
			es.colFlags[col] = 2
			//fmt.Printf(">> set colFlags[%v]=2\n", col)
			//fmt.Println("")
		}

		//fmt.Printf("rowFlags=%v, colFlags=%v, uComputedCnt=%v, vComputedCnt=%v\n", es.rowFlags, es.colFlags, uComputedCnt, vComputedCnt)
		//fmt.Println("")
		if !more2scan || (uComputedCnt == sLen && vComputedCnt == dLen) {
			//fmt.Println("break!")
			break
		}
	}

	//fmt.Printf("computed u=%v\n", es.u)
	//fmt.Printf("computed v=%v\n", es.v)
	//fmt.Println("")

	if uComputedCnt != sLen || vComputedCnt != dLen {
		return fmt.Errorf("[computeUV()] U: %v/%v, V: %v/%v", uComputedCnt, sLen, vComputedCnt, dLen)
	} else {
		//fmt.Printf("computeUV() finished in %v\n", time.Now().Sub(t1))
		return nil
	}
}

func (es *Problem) isOptimal() bool {
	//fmt.Println("[Checking if current solution is optimal ...]")
	//t1 := time.Now()
	// find the base cell by computing the penalty for all no-flow cell
	sLen, dLen := es.sLen, es.dLen
	epsilon := es.epsilon

	es.row, es.col = -1, -1
	pMax := float64(-1)
	for i := 0; i < sLen; i++ {
		for j := 0; j < dLen; j++ {
			if es.flow[i][j].basic {
				continue
			}
			p := es.u[i] + es.v[j] - es.costMatrix[i][j]
			//fmt.Printf("i=%v,j=%v,p=%v\n", i, j, p)
			if p > epsilon && p > pMax {
				es.row, es.col, pMax = i, j, p
			}
		}
	}
	var optimal bool = (es.row == -1)
	/*if optimal {
		fmt.Println("current solution is optimal!")
	} else {
		fmt.Printf("current solution is NOT optimal! Optimization starting cell: (%v, %v)\n", es.row, es.col)
	}
	fmt.Println("")*/
	//fmt.Printf("isOptimal() finished in %v\n", time.Now().Sub(t1))
	return optimal
}

func (es *Problem) findLoop() error {
	//fmt.Println("[Finding a valid loop ...]")
	sLen, dLen := es.sLen, es.dLen
	infinity := es.infinity
	// reset row/col flags
	for i := 0; i < sLen; i++ {
		es.rowFlags[i] = 0
	}
	for i := 0; i < dLen; i++ {
		es.colFlags[i] = 0
	}
	// create the head cell
	// we only need to start with a horizontal cell (flag=true)
	// as the loop should end with a vertical cell, same as
	// starting it with a vertical cell (flag=false) which should
	// end with a horizontal cell
	head := &cell{
		row:             es.row,
		col:             es.col,
		flag:            true,
		loopEvenMinFlow: infinity,
		prev:            nil,
		next:            nil,
	}

	curr := head
	step := 0

	//fmt.Printf("starting from cell at %v\n", curr)
	for {

		nexti, nextj := -1, -1
		if curr.flag { // horizontal
			i := curr.row
			next := curr.next
			start := 0
			if next != nil {
				start = next.col + 1
			}
			if start < dLen {
				for j := start; j < dLen; j++ {
					if j == curr.col || !es.flow[i][j].basic || es.colFlags[j] == 1 {
						continue
					}
					nexti, nextj = i, j
					es.rowFlags[i] = 1
					break
				}
			}
		} else { // vertical
			j := curr.col
			next := curr.next
			start := 0
			if next != nil {
				start = next.row + 1
			}
			if start < sLen {
				for i := start; i < sLen; i++ {
					if i == curr.row || !es.flow[i][j].basic || es.rowFlags[i] == 1 {
						continue
					}
					nexti, nextj = i, j
					es.colFlags[j] = 1
					break
				}
			}
		}

		if nexti >= 0 { // found next cell
			nextFlag := !curr.flag
			loopEvenMinFlow := curr.loopEvenMinFlow
			// new cell is the even one in the chain,
			// need to calculate the min flow
			if !nextFlag { // next is the even cell
				nextFlow := es.flow[nexti][nextj].value
				if nextFlow < loopEvenMinFlow {
					loopEvenMinFlow = nextFlow
				}
			}

			next := &cell{
				row:             nexti,
				col:             nextj,
				flag:            nextFlag,
				loopEvenMinFlow: loopEvenMinFlow,
				prev:            curr,
				next:            nil,
			}
			curr.next = next
			curr = next
			step += 1

			//fmt.Printf("step #%v: found next cell %v\n", step, curr)

			if curr.col == head.col {
				// found a valid loop
				break
			}
		} else { // didn't find next cell, need to go back
			if curr == head {
				// cannot go back from head which means we couldn't
				// find a valid loop
				//fmt.Printf("cannot find a valid loop starting from %v\n", head)
				head = nil
				break
			} else {
				curr = curr.prev
				if curr.flag { // horizontal
					es.rowFlags[curr.row] = 0
				} else { // vertical
					es.colFlags[curr.col] = 0
				}
				step += 1

				//fmt.Printf("step #%v: go back to %v\n", step, curr)
			}
		}
	}

	if head != nil {
		// save loop even-cell min flow to head cell for easy access
		head.loopEvenMinFlow = curr.loopEvenMinFlow
		// save loop head cell
		es.loop = head

		/*fmt.Println("found a valid loop:")
		p := head
		for p != nil {
			if p != head {
				fmt.Print(" -> ")
			}
			fmt.Printf("%v", p)
			//fmt.Printf("%v/%v/%v", p, es.flow[p.row][p.col], p.loopEvenMinFlow)
			p = p.next
		}
		fmt.Println("")
		fmt.Printf("loop even-cell min flow: %v\n", head.loopEvenMinFlow)
		fmt.Println("")*/

		return nil
	} else {
		// clear loop head cell
		es.loop = nil
		//fmt.Println("")
		return fmt.Errorf("[findLoop()] cannot find a valid loop starting from (%v,%v).", es.row, es.col)
	}
}

func (es *Problem) fixDegeneracy(dCnt int) error {
	//fmt.Printf("[fixing degeneracy for %v variables ...]\n", dCnt)
	//t1 := time.Now()
	sLen, dLen := es.sLen, es.dLen
	left := dCnt
	for i := 0; i < sLen; i++ {
		for j := 0; j < dLen; j++ {
			fc := es.flow[i][j]
			if fc.basic {
				continue
			}
			// try to find an independent cell in another word,
			// starting from this cell it shouldn't form a loop
			es.row, es.col = i, j
			if err := es.findLoop(); err == nil {
				continue
			}
			// no loop found, set it as basic cell
			fc.basic = true
			fc.value = 0
			left -= 1
			//fmt.Printf("assigned (%v,%v) as 0-value basic cell.\n", i, j)
			if left == 0 {
				break
			}
		}
		if left == 0 {
			break
		}
	}
	if left > 0 {
		return fmt.Errorf("[fixDegeneracy(%v)] failed to find %v independent cells.", dCnt, left)
	}
	//fmt.Printf("fixDegenracy finished in %v\n", time.Now().Sub(t1))
	return nil
}

func (es *Problem) applyOptimization() {
	//fmt.Println("[Applying optimaztion ...]")
	//t1 := time.Now()
	p := es.loop
	q := p.loopEvenMinFlow
	epsilon := es.epsilon
	removed := false
	for p != nil {
		row, col := p.row, p.col
		fc := es.flow[row][col]
		if p.flag { // odd cell
			fc.basic = true
			fc.value = fc.value + q
			//fmt.Printf("added %v to (%v,%v)\n", q, row, col)
		} else { // even cell
			fval := fc.value - q
			fc.value = fval
			//fmt.Printf("substracted %v from (%v,%v)\n", q, row, col)
			// only remove the first one if multiply flow cells
			// reach value 0, this is to prevent degeneracy
			if !removed && fval <= epsilon {
				fc.basic = false
				removed = true
				//fmt.Printf("removed (%v,%v) from basic variables.\n", row, col)
			}
		}
		p = p.next
	}

	//fmt.Println()
	//fmt.Printf("applyOptimization finished in %v\n", time.Now().Sub(t1))
}

// Solve the transportation problem.
// Returns error if something goes wrong.
func (es *Problem) Solve() error {
	//fmt.Println("[Solving the problem ...]")
	//t1 := time.Now()
	flowCnt := es.findFeasibleSolution()
	//es.printSolution()
	//fmt.Printf("m=%v,n=%v,flowCnt=%v\n", es.sLen, es.dLen, flowCnt)
	//t2 := time.Now()
	//fmt.Printf("found feasible solution in %v\n", t2.Sub(t1))

	// fix degeneracy
	correctCnt := es.sLen + es.dLen - 1
	dCnt := correctCnt - flowCnt
	if dCnt > 0 {
		if err := es.fixDegeneracy(dCnt); err != nil {
			return err
		}
	} else if dCnt < 0 {
		return fmt.Errorf("feasible solution has %v basic cells, while it should be %v!", flowCnt, correctCnt)
	}
	//es.printSolution()
	//t3 := time.Now()
	//fmt.Printf("fixed degeneracy in %v\n", t3.Sub(t2))

	maxIter := es.maxIter
	//optimal := false
	for {
		//t4 := time.Now()
		//fmt.Printf("iteration #%v\n", es.iterCnt)
		if err := es.computeUV(); err != nil {
			return err
		}
		//fmt.Printf("computed u=%v v=%v\n", es.u, es.v)
		if es.isOptimal() {
			//fmt.Println("we got optimal solution!")
			//optimal = true
			break
		}
		//fmt.Printf("optimization start cell: (%v, %v)\n", es.row, es.col)
		if err := es.findLoop(); err != nil {
			return err
		}
		es.applyOptimization()
		//es.printSolution()
		es.iterCnt += 1
		//t5 := time.Now()
		//fmt.Printf("finished optimization iteration #%v in %v\n", es.iterCnt, t5.Sub(t4))
		if maxIter > 0 && es.iterCnt >= maxIter {
			break
		}
	}
	//t6 := time.Now()
	//fmt.Printf("finished %v optimization iterations in %v\n", es.iterCnt, t6.Sub(t3))

	//fmt.Printf("is optimal solution: %v\n", optimal)
	//fmt.Printf("iteration ran: %v\n", es.iterCnt)
	//fmt.Println()

	return nil
}

// Get the solution cost, should be called after calling Solve().
func (es *Problem) GetCost() float64 {
	sLen, dLen := es.sLen, es.dLen
	cost := float64(0)
	for i := 0; i < sLen; i++ {
		for j := 0; j < dLen; j++ {
			fc := es.flow[i][j]
			if !fc.basic || fc.value == 0 {
				continue
			}
			cost += fc.value * es.costMatrix[i][j]
		}
	}
	return cost
}

// Get the flow matrix, should be called after calling Solve().
func (es *Problem) GetFlow() [][]float64 {
	balanced, sLen, dLen := es.balanced, 0, 0
	if balanced < 0 {
		sLen, dLen = es.sLen-1, es.dLen
	} else if balanced > 0 {
		sLen, dLen = es.sLen, es.dLen-1
	} else {
		sLen, dLen = es.sLen, es.dLen
	}
	flow := make([][]float64, sLen)
	for i := 0; i < sLen; i++ {
		flow[i] = make([]float64, dLen)
		for j := 0; j < dLen; j++ {
			fc := es.flow[i][j]
			if !fc.basic || fc.value == 0 {
				continue
			}
			fval := fc.value
			flow[i][j] = fval
		}
	}
	return flow
}

// Get the solution (both the total cost and the flow matrix), should
// be called after calling Solve().
func (es *Problem) GetCostAndFlow() (float64, [][]float64) {
	balanced, sLen, dLen := es.balanced, 0, 0
	if balanced < 0 {
		sLen, dLen = es.sLen-1, es.dLen
	} else if balanced > 0 {
		sLen, dLen = es.sLen, es.dLen-1
	} else {
		sLen, dLen = es.sLen, es.dLen
	}
	cost := float64(0)
	flow := make([][]float64, sLen)
	for i := 0; i < sLen; i++ {
		flow[i] = make([]float64, dLen)
		for j := 0; j < dLen; j++ {
			fc := es.flow[i][j]
			if !fc.basic || fc.value == 0 {
				continue
			}
			fval := fc.value
			flow[i][j] = fval
			cost += fval * es.costMatrix[i][j]
		}
	}
	return cost, flow
}

// Create a transportation problem from the given args.
//
//  supply, demand: positive float64 array/slice.
//  costs: 2-D matrix, row size should match supply length, column
//         size should match demand length.
//
//  opts is for optional args:
//   opts[0]: MAX_ITER, max iterations to run when optimizing the
//            solution, default to 100.
//   opts[1]: EPSILON, used to tell if a float64 value is zero or not,
//            default to 1e-6.
//   if you need to use a non-default EPSILON (opt[1]), you must also
//   set MAX_ITER (opt[0]).
//
//  returns the Problem{} struct.
func CreateProblem(supply, demand []float64, costs [][]float64, opts ...float64) (*Problem, error) {
	maxIter, epsilon := MAX_ITER, EPSILON
	optsLen := len(opts)
	if optsLen > 0 {
		maxIter = int(opts[0])
		if optsLen > 1 {
			epsilon = opts[1]
			if epsilon > float64(1e-3) {
				return nil, fmt.Errorf("Given epsilon is too big (>1e-3): %v", opts[1])
			}
		}
	}

	//fmt.Printf("maxIter=%v, epsilon=%v\n", maxIter, epsilon)

	p, err := createProblem(supply, demand, costs, maxIter, epsilon)
	if err == nil {
		return p, nil
	} else {
		return nil, err
	}
}
