package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func niceAngles(a float64) string {
	return fmt.Sprintf("%d.%02d'%02d\"", int(a),
		int((a-float64(int(a)))*60),
		int((a*60-float64(int(a*60)))*60))
}

func main() {
	var a float64
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscan(scanner.Text(), &a)
		fmt.Println(niceAngles(a))
	}
}
