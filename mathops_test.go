package tfcgo

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"testing"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	op "github.com/tensorflow/tensorflow/tensorflow/go/op"
)

func TestSquare(test *testing.T) {

	d, _, err := loadTestingData("numbers.csv")
	if err != nil {
		log.Fatal(err)
	}

	t1, err := tf.NewTensor(d)
	if err != nil {
		test.Fatalf("%q: %v", "TestSquare", err)
	}
	t2, e := Square(t1)
	if e != nil {
		test.Fatalf("%q: %v", "TestSquare", e)

	}

	tensor1Data := t1.Value().([][]float64)
	tensor2Data := t2.Value().([][]float64)
	for x, f1 := range tensor1Data {
		f2 := tensor2Data[x]
		if f1[0]*f1[0] != f2[0] {
			test.Fatalf("TestSquare Failed ", f1[0], f2[0])
		}

	}

}

func convertToFloat64(recs []string) []float64 {
	var numbers []float64
	for _, i := range recs {
		if n, err := strconv.ParseFloat(i, 64); err == nil {
			numbers = append(numbers, n)
		}
	}
	return numbers
}

func loadTestingData(fileName string) (records [][]float64, labels []string, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)

	var recs [][]float64
	var lbls []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		recs = append(recs, convertToFloat64(record[1:2]))
		//lbls = append(lbls, record[5])
	}
	return recs, lbls, err
}

func makeConst(s *op.Scope, value interface{}) (tf.Output, error) {
	t, ok := value.(*tf.Tensor)
	if !ok {
		var err error
		if t, err = tf.NewTensor(value); err != nil {
			return tf.Output{}, err
		}
	}
	op := s.AddOperation(tf.OpSpec{
		Type: "Const",
		Name: "Name",
		Attrs: map[string]interface{}{
			"dtype": t.DataType(),
			"value": t,
		},
	})
	return op.Output(0), nil
}
