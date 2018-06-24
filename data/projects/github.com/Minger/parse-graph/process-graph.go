package main

// image a microservice that takes a stream of raw graph data and formats it into nodes and edges
// for ingestion by a dynamic graph builder
// fan out the processing, then fan in the results to one channel
// we don't cache or look up hashes for previously seen strings: for simplicity we trade raw computation
// for memory and storage
// data csv comes from http://kgullikson88.github.io/blog/static/PyPiAnalyzer.html
// a follow up the spark graphx code:
// https://github.com/Minger/experiments/blob/master/spark-packages-pagerank.scala

import (
	"encoding/csv"
	"fmt"
	"github.com/dmiller/go-seq/murmur3"
	"io"
	"log"
	"os"
	"sync"
)

// Edge joins two nodes
type Edge struct {
	Source       uint32
	Destination  uint32
	Relationship string
}

// Node represents a package
type Node struct {
	VertexID uint32
	Name     string
}

// this func lifted from http://stackoverflow.com/questions/32027590/efficient-read-and-write-csv-in-go
func processCSV(rc io.Reader) (ch chan []string) {
	ch = make(chan []string, 10)
	go func() {
		r := csv.NewReader(rc)
		if _, err := r.Read(); err != nil { //read header
			log.Fatal(err)
		}
		defer close(ch)
		for {
			rec, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)

			}
			ch <- rec
		}
	}()
	return
}

func processRow(in <-chan []string) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		for {
			row := <-in
			destination := row[1]
			if len(destination) > 0 {
				destinationHash := murmur3.HashString(destination)
				node := Node{destinationHash, destination}
				out <- node

				source := row[2]
				if len(source) > 0 {
					sourceHash := murmur3.HashString(source)
					node := Node{sourceHash, source}
					out <- node
					edge := Edge{sourceHash, destinationHash, "requires"}
					out <- edge
				}
			}
		}
		close(out)
	}()
	return out
}

// lifted with slight modificaton from the docs
func merge(cs [4]<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	out := make(chan interface{})

	output := func(c <-chan interface{}) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	file, err := os.Open("requirements.csv")
	if err != nil {
		log.Fatal(err)
	}

	fileCh := processCSV(file)

	a := processRow(fileCh)
	b := processRow(fileCh)
	c := processRow(fileCh)
	d := processRow(fileCh)

	var chans = [4]<-chan interface{}{a, b, c, d}
	results := merge(chans)

	for {
		graph := <-results
		fmt.Printf("%+v\n", graph)
	}
}
