package codetect

import (
	"log"
	"os"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"gopkg.in/vmarkovtsev/CodeNeuron.v1/assets"
)

// CodeBoundary represents a start or an end of a detected code block.
type CodeBoundary struct {
	// PositionInRunes is the index of the boundary in the parsed *runes array*.
	// This is not a position in the byte stream.
	// The boundary goes *after* the corresponding rune index.
	PositionInRunes int
	// Start is true if the boundary is a start, otherwise, it is false for an end.
	Start bool
}

var instance = loadModel()

type model struct {
	graph *tf.Graph
	input1 tf.Output
	input2 tf.Output
	output tf.Output
	batchSize int
	sequenceLength int
}

func loadModel() *model {
	os.Setenv("TF_CPP_MIN_LOG_LEVEL", "3")
	modelBytes, err := assets.Asset("model.pb")
	if err != nil {
		log.Fatalf("failed to load model.pb from the assets: %v", err)
	}
	graph := tf.NewGraph()
	err = graph.Import(modelBytes, "")
	if err != nil {
		log.Fatalf("importing the model: %v", err)
	}
	input1 := graph.Operation("input_1_1").Output(0)
	input2 := graph.Operation("input_2_1").Output(0)
	output := graph.Operation("output").Output(0)
	inputShape, err := input1.Shape().ToSlice()
	if err != nil {
		log.Fatalf("Getting the input shape: %v", err)
	}
	batchSize := int(inputShape[0])
	sequenceLength := int(inputShape[1])
	return &model {
		graph: graph,
		input1: input1,
		input2: input2,
		output: output,
		batchSize: batchSize,
		sequenceLength: sequenceLength,
	}
}

func bakeText(text []rune) ([]*tf.Tensor, []*tf.Tensor, error) {
	batchSize, sequenceLength := instance.batchSize, instance.sequenceLength
	realSize := len(text) - sequenceLength
	inputSize := realSize
	if inputSize % batchSize != 0 {
		inputSize = (inputSize / batchSize + 1) * batchSize
	}
	batches1 := make([][][]uint8, inputSize / batchSize)
	batches2 := make([][][]uint8, inputSize / batchSize)
	for i := range batches1 {
		batches1[i] = make([][]uint8, batchSize)
		batches2[i] = make([][]uint8, batchSize)
	}
	for i := realSize; i < inputSize; i++ {
		arr := make([]uint8, sequenceLength)
		batches1[i / batchSize][i % batchSize] = arr
		batches2[i / batchSize][i % batchSize] = arr
	}
	pos := 0
	for x := range text {
		if x < sequenceLength / 2 {
			continue
		}
		if x >= len(text) - sequenceLength / 2 {
			break
		}
		arr1 := make([]uint8, sequenceLength)
		batches1[pos / batchSize][pos % batchSize] = arr1
		arr2 := make([]uint8, sequenceLength)
		batches2[pos / batchSize][pos % batchSize] = arr2
		for i := 0; i < sequenceLength; i++ {
			bi := x - sequenceLength + i + 1
			if bi >= 0 {
				val, exists := CHARS[text[bi]]
				if !exists {
					val = uint8(len(CHARS))
				}
				arr1[i] = val
			}
			bi = x + sequenceLength - i
			if bi < len(text) {
				val, exists := CHARS[text[bi]]
				if !exists {
					val = uint8(len(CHARS))
				}
				arr2[i] = val
			}
		}
		pos++
	}
	tensors1 := make([]*tf.Tensor, len(batches1))
	tensors2 := make([]*tf.Tensor, len(batches2))
	for i := range batches1 {
		batch1 := batches1[i]
		batch2 := batches2[i]
		input1, err := tf.NewTensor(batch1)
		if err != nil {
			return nil, nil, err
		}
		input2, err := tf.NewTensor(batch2)
		if err != nil {
			return nil, nil, err
		}
		tensors1[i] = input1
		tensors2[i] = input2
	}
	return tensors1, tensors2, nil
}

// Run detects the code block boundaries using CodeNeuron network.
// See GetSequenceLength() for the details which portion of the text is analyzed.
func Run(text string, session *tf.Session) ([]CodeBoundary, error) {
	if text == "" {
		return nil, nil
	}
	runes := []rune(text)
	tensors1, tensors2, err := bakeText(runes)
	if err != nil {
		return nil, err
	}
	pos := 0
	boundaries := []CodeBoundary{}
	for i := range tensors1 {
		result, err := session.Run(map[tf.Output]*tf.Tensor{
			instance.input1: tensors1[i], instance.input2: tensors2[i],
		}, []tf.Output{instance.output}, nil)
		if err != nil {
			return nil, err
		}
		probs := result[0].Value().([][]float32)
		for _, prob := range probs {
			offsetPos := pos + instance.sequenceLength / 2
			maxi := 2
			maxval := prob[2]
			if prob[0] > maxval {
				maxi = 0
				maxval = prob[0]
			}
			if prob[1] > maxval {
				maxi = 1
				maxval = prob[1]
			}
			if maxi == 0 {
				boundaries = append(boundaries, CodeBoundary{
					PositionInRunes: offsetPos,
					Start: true,
				})
			} else if maxi == 1 {
				boundaries = append(boundaries, CodeBoundary{
					PositionInRunes: offsetPos,
					Start: false,
				})
			}
			pos++
			if pos >= len(text) - instance.sequenceLength {
				break
			}
		}
	}
	return boundaries, nil
}

// OpenSession initializes a new Tensorflow session.
// Remember to defer session.Close()
func OpenSession() (*tf.Session, error) {
	return tf.NewSession(instance.graph, &tf.SessionOptions{})
}

// GetSequenceLength returns the sequence length of the RNN model.
// text[:length / 2] and text[-length / 2:] are not analyzed because the network has too little
// context. You can workaround this by appending and prepending some constant strings.
func GetSequenceLength() int {
	return instance.sequenceLength
}
