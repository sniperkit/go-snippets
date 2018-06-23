package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

func main() {
	fmt.Println("Make sure you launch this benchmark via ./run_bench.sh")
	runtime.GOMAXPROCS(runtime.NumCPU())
	const N = 100000000
	c := make([]float64, N)
	d := make([]float64, N)
	h := math.Pi / N

	beg1 := time.Now()
	for i := 0; i < N; i++ {
		c[i] = math.Exp(math.Sin(float64(i)*h) + math.Cos(math.Pi+float64(i)*h))
	}
	time1 := time.Since(beg1)
	fmt.Println("Sequential execution took: ", time1)

	beg2 := time.Now()
	//gomp
	for i := 0; i < N; i++ {
		d[i] = math.Exp(math.Sin(float64(i)*h) + math.Cos(math.Pi+float64(i)*h))
	}
	time2 := time.Since(beg2)
	fmt.Println("Parallel execution took: ", time2)
	fmt.Println("Speedup is ", float64(time1.Nanoseconds())/float64(time2.Nanoseconds()))
	for i := 0; i < N; i++ {
		if math.Abs(c[i]-d[i]) > 1e-6 {
			fmt.Println("Results of computations are not close enough: ", c[i], ", ", d[i])
			panic("Problems with parallelization")
		}
	}
	fmt.Println("Results of computations are close enough")
}
