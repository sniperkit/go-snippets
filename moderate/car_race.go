package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type segment struct {
	length, angle float64
}

type car struct {
	number  int
	laptime float64
}
type ranking []car

func (r ranking) Len() int { return len(r) }
func (r ranking) Less(i, j int) bool {
	return r[i].laptime < r[j].laptime
}
func (r ranking) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	scanner.Scan()
	s := strings.Fields(scanner.Text())
	track := make([]segment, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		var v, w float64
		fmt.Sscan(s[i], &v)
		fmt.Sscan(s[i+1], &w)
		track[i/2] = segment{v, w}
	}

	var cars []car
	for scanner.Scan() {
		var (
			num                                int
			vmax, accel, brake, laptime, vinit float64
		)
		fmt.Sscan(scanner.Text(), &num, &vmax, &accel, &brake)
		for _, i := range track {
			away := accel * (1 - (vinit / vmax)) * (vinit + vmax) / (2 * 3600)
			vend := vmax * (1 - (i.angle / 180))
			bway := brake * (1 - (vend / vmax)) * (vend + vmax) / (2 * 3600)
			laptime += 3600 * (2*away/(vinit+vmax) + (i.length-away-bway)/vmax + 2*bway/(vend+vmax))
			vinit = vend
		}
		cars = append(cars, car{num, laptime})
	}
	sort.Sort(ranking(cars))

	for _, i := range cars {
		fmt.Printf("%d %.2f\n", i.number, i.laptime)
	}

}
