package mapreduce

import "testing"

func mapper(s interface{}, output chan interface{}) {
	result := s.(string) + " (changed by mapper)"
	output <- result
}

func reducer(input chan interface{}, output chan interface{}) {
	result := make(map[string]int)

	for inputResult := range input {
		result[inputResult.(string)] = result[inputResult.(string)] + 1
	}

	output <- result
}

func sendInput() chan interface{} {
	output := make(chan interface{})

	go func() {

		output <- "three"
		output <- "one"
		output <- "two"
		output <- "three"
		output <- "two"
		output <- "three"

		close(output)
	}()

	return output
}

func TestProcess(t *testing.T) {
	inputChan := sendInput()
	result := Process(mapper, reducer, inputChan, 5).(map[string]int)

	if result["one (changed by mapper)"] != 1 ||
		result["two (changed by mapper)"] != 2 ||
		result["three (changed by mapper)"] != 3 {
		t.Error("Incorrect mapreduce result:", result)
	}
}
