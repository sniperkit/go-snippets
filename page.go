package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	zmq "github.com/pebbe/zmq4"
	"encoding/json"

	// "github.com/chrislusf/gleam/distributed"
	"github.com/chrislusf/gleam/flow"
	"github.com/chrislusf/gleam/gio"
	// "github.com/chrislusf/gleam/gio/mapper"
	"github.com/chrislusf/gleam/gio/reducer"
	"github.com/chrislusf/gleam/util"
)

var (
	voter = gio.RegisterMapper(makevotes)
	splitter = gio.RegisterMapper(split)
	isDistributed        = flag.Bool("distributed", false, "run in distributed or not")
)

func main() {
	gio.Init()
	flag.Parse()
	data := make(chan interface{})
	output := make(chan []interface{}, 500)
	var wg sync.WaitGroup

	//the zeromq stuff goes here
	//  Socket to talk to server
	subscriber, _ := zmq.NewSocket(zmq.SUB)

	subscriber.Connect("tcp://0.0.0.0:5555")

	//  Subscribe to zipcode, default is NYC, 10001

	subscriber.SetSubscribe("")

	//  Process 100 updates
	go func() {
		for {
			msg, _ := subscriber.Recv(0)

			if (msg == "done") {
				close(output)
				break
			}
			
			type Message struct {
				Node string
				Outlinks string
				Keywords string
				Score float32
			}
			
			message := Message{}
			err := json.Unmarshal([]byte(msg), &message)
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
			var row []interface{}
			
			row = append(row, message.Node, message.Outlinks, message.Score, message.Keywords)
			
			// fmt.Println(len(row))
			
			data <- row



		}
	}()

	//start consumers before pagerank starts
	for i := 0; i < 100; i++ {
		wg.Add(1)
		
		go func(input chan []interface{}) {
			//send chan data to cassandra here
			for x := range input {
				fmt.Println(x)
			}
			defer wg.Done()

		}(output)
	}

	start := flow.New("indexing").Channel(data)

	start = start.Map("split", splitter)

	keywords := start.Select("divert_keywords", flow.Field(1,4))
	// keywords.Printlnf("%s - %s")
	// word_doc_one := start.Select(1,4).
	// 	PartitionByKey("partition", 7).
	// 	Map("read content", registeredReadConent)

	// docFreq := word_doc_one.
	// 	ReduceBy("df", reducer.SumInt64, flow.Field(1,2)).
	// 	SelectKV("reorder", flow.Field(2), flow.Field(1,3))

	// termFreq :=
	// 	docFreq.
	// 	Select("discard", flow.Field(2,3)).
	// 	ReduceBy("word_doc_tf", reducer.SumInt64, flow.Field(1))

	// tdidf = docFreq.
	// 	Join("byWord", termFreq, flow.Field(1)).
	// 	Map("tfidf", registeredTfIdf)
	// 	// Sort("sort by tf/df", flow.Field(5)).
	// 	// OutputRow(func(row *util.Row) error {
	// 	// 	fmt.Printf("%s: %s tf=%d df=%d tf-idf=%f\n",
	// 	// 		row.K[0],
	// 	// 		row.V[0],
	// 	// 		row.V[1].(uint16),
	// 	// 		row.V[2].(uint16),
	// 	// 		row.V[3].(float32),
	// 	// 	)
	// 	// 	return nil
	// 	// })


	runs := 10
	for i := 0; i < runs; i++ {
		start = graph(start)
	}

	// index := start.join("byPage", tdidf, flow.Field(1))

	start = start.Join("add", keywords, flow.Field(1))

	start.OutputRow(func(row *util.Row) error {
			// fmt.Printf("%s: %f\n",
			// 	row.K[0],
			// 	row.V[1].(float32),

			// )
			var outputRow []interface{}
			outputRow = append(outputRow, row.K[0], row.V[0], row.V[1], row.V[2])
			output <- outputRow
			return nil
		})
	
	start.Run()

	wg.Wait()


}

func makevotes(x []interface{}) error {
	fmt.Print(len(x))

	rawurls := gio.ToString(x[1])
	score := gio.ToFloat64(x[2])
	urls := strings.Split(rawurls, ",")

	for url := range urls {
		gio.Emit(url, score/float64(len(urls)))
	}
	return nil
}

func graph(start *flow.Dataset) *flow.Dataset{
	page := start.
		Map("votes", voter).
		Printlnf("%s - %d").
		ReduceBy("count", reducer.SumFloat64, flow.Field(1))

	initial :=
		start.
		Select("access_urls", flow.Field(1,2))


	return initial.Join("byNode", page, flow.Field(1))
}
func split(x []interface{}) error {
	row := x[0].([]interface{})
	gio.Emit(row...)
	return nil
}

// func readContent(x []interface{}) error {

// 	page := gio.ToString(x[0])
// 	text := gio.ToString(x[1])

// 	//split text into list of words, E.g. tokenize here

// 	for word := range words {
// 		gio.Emit(word, page, 1)
// 	}
// 	return nil
// }

// func tfidf(x []interface{}) error {
// 	word := gio.ToString(x[0])
// 	page := gio.ToString(x[1])
// 	df := uint16(gio.ToInt64(x[2]))
// 	tf := uint16(gio.ToInt64(x[3]))

// 	gio.Emit(page, word, tf, df, float32(df)/float32(tf))
// 	return nil
// }

