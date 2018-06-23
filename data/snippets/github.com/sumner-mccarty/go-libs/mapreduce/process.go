package mapreduce

// MapperCollector is a channel that collects the output from mapper tasks
type MapperCollector chan chan interface{}

// MapperFunc is a function that performs the mapping part of the MapReduce job
type MapperFunc func(interface{}, chan interface{})

// ReducerFunc is a function that performs the reduce part of the MapReduce job
type ReducerFunc func(chan interface{}, chan interface{})

// Run accepts the user's mapper and reducer fn along with the input channel and how many worker channels it should use
func Process(mapper MapperFunc, reducer ReducerFunc, input chan interface{}, numWorkers int) interface{} {

	reducerInput := make(chan interface{})
	reducerOutput := make(chan interface{})
	mapperCollector := make(MapperCollector, numWorkers)

	go reducer(reducerInput, reducerOutput)
	go reducerDispatcher(mapperCollector, reducerInput)
	go mapperDispatcher(mapper, input, mapperCollector)

	return <-reducerOutput
}

func mapperDispatcher(mapper MapperFunc, input chan interface{}, collector MapperCollector) {
	for item := range input {
		taskOutput := make(chan interface{})
		go mapper(item, taskOutput)
		collector <- taskOutput
	}
	close(collector)
}

func reducerDispatcher(collector MapperCollector, reducerInput chan interface{}) {
	for output := range collector {
		reducerInput <- <-output
	}
	close(reducerInput)
}
