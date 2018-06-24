package main

import "fmt"

func bfs(gr[][] int, s int) [] int {
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

type pair struct{ A, B int }

func main() {
	var x, y int
	fmt.Scanf("%d %d", &x, &y)
	gr := make([][]int, x)
	for i := 0; i < y; i++ {
		var a, b int
		fmt.Scanf("%d %d", &a, &b)
		gr[a] = append(gr[a], b)
		//gr[b] = append(gr[a], b)
	}
	ans := make([]map[int]int, 10)
	d := make([]int, 10)
	ress := make([][]int, x)
	d[0] = d[0] + 1
	for w := 0; w < 1; w++ {
		ans[w] = make(map[int]int)
		for v := 0; v < x; v++ {
			res := bfs(gr, v)
			fmt.Println(v, " ", res)
			for i := 0; i < len(res); i++ {
				ress[v] = append(ress[v], res[i])
			}
		}
		fmt.Println("\n")
	}

	for i := 0; i < x; i++ {
		for j := 0; j < len(ress); j++ {
			fmt.Print(ress[i][j], " ")
		}
		fmt.Println()
	}

	var max int = 0
	var sum float64 = 0.0
	var cnt int = 0
	// find max
	for i := 0; i < x; i++ {
		for j := 0; j < len(ress); j++ {
			if max < ress[i][j] {
				max = ress[i][j]
			}
			if ress[i][j] != -1 {
				sum += float64(ress[i][j])
				cnt++;
			}

		}
	}
	fmt.Printf("diam=%d, avg dist=%.1f\n", max, sum / float64(cnt))

}
