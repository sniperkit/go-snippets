package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	// "project_c9f7_i5l8_o0p4_p0j8/pregel"
	"strconv"
	"strings"

	"../pregel"
)

const numEngines = 6
const batch = 100

var messages = make(map[int]map[int][]pregel.VertexMessage)
var msgDistributionChan chan []pregel.VertexMessage

//var mutex = &sync.Mutex{}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func runPregel(numSteps int, engines []*pregel.Engine) {
	c := make(chan bool)
	go receiver()
	for i := 0; i < numSteps; i++ {
		currMsgs := make([]map[int][]pregel.VertexMessage, len(engines), len(engines))
		for _, engine := range engines {
			currMsgs[engine.Id] = make(map[int][]pregel.VertexMessage)
			for id, msg := range messages[engine.Id] {
				currMsgs[engine.Id][id] = msg
			}
			// log.Println("Messages for engine: ", engine.Id, len(currMsgs[engine.Id]))
		}

		clearMessages()
		for _, engine := range engines {
			go engine.Superstep(currMsgs[engine.Id], i, c)
		}
		for _ = range engines {
			<-c
		}
		// log.Println("Finished superstep.")
	}

}

func clearMessages() {
	for k := range messages {
		delete(messages, k)
	}
	//log.Println("Cleared messages.")
}

func receiver() {
	for {
		newMsgs := <-msgDistributionChan
		for _, msg := range newMsgs {
			if msg.ToId == 0 {
				log.Println("Receive val: ", msg.ToId, " ID: ", msg.FromId)
			}

			engID := (msg.ToId % numEngines)
			if _, ok := messages[engID]; ok {
				if messages[engID][msg.ToId] != nil {
					messages[engID][msg.ToId] = append(messages[engID][msg.ToId], msg)
				} else {
					messages[engID][msg.ToId] = []pregel.VertexMessage{msg}
				}
			} else {
				messages[engID] = make(map[int][]pregel.VertexMessage)
				messages[engID][msg.ToId] = []pregel.VertexMessage{msg}
			}
		}
	}
}

func main() {
	lines, err := readLines("../sampleData/Wiki-Vote.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}
	log.Println(len(lines))
	vertexMap := make(map[int][]int)
	for _, line := range lines {
		if line[0] == '#' {
			continue
		} else {
			entries := strings.Fields(line)
			//log.Println(entries)
			if len(entries) == 2 {
				id, err := strconv.Atoi(entries[0])
				if err != nil {
					log.Fatal(err)
				}
				val, err := strconv.Atoi(entries[1])
				if err != nil {
					log.Fatal(err)
				}
				if _, ok := vertexMap[val]; !ok {
					vertexMap[val] = []int{}
				}
				vertexMap[id] = append(vertexMap[id], val)
			}
		}
	}

	msgDistributionChan = make(chan []pregel.VertexMessage) //, numEngines)

	vertices := [numEngines]map[int][]int{}
	for id := 0; id < numEngines; id++ {
		vertices[id] = make(map[int][]int)
	}
	for vid, edges := range vertexMap {
		worker := vid % numEngines
		vertices[worker][vid] = edges
	}

	engines := make([]*pregel.Engine, numEngines, numEngines)
	for eid := 0; eid < numEngines; eid++ {
		//vertexMsgChan := make(chan pregel.VertexMessage)
		prvs := pregel.GetPageRankVertices(len(vertexMap), vertices[eid])
		engine := pregel.NewEngine(prvs, eid, msgDistributionChan, batch)

		engines[eid] = engine
	}

	runPregel(20, engines)

	output := []string{"Page Rank Values: \r\n"}
	for _, engine := range engines {
		vertices := engine.GetVertices()
		for _, vertex := range vertices {
			id := strconv.Itoa(vertex.GetId())
			val := strconv.FormatFloat(vertex.GetValue(), 'E', -1, 64)
			output = append(output, id+" "+val+"\r\n")
		}
	}
	_ = writeLines(output, "../sampleData/Wiki-Vote-Output.txt")

	n := len(vertexMap)
	log.Println("Number of vertices: ", n)

}
