package tfcgo

import (
	"fmt"
	"reflect"

	"github.com/gonum/floats"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

//SumOfSquares - accepts a value (expecting [][]float64) and returns sum of squares
func SumOfSquares(s *op.Scope, value interface{}) (*tf.Tensor, tf.Output, error) {

	t, err := tf.NewTensor(value)
	if err != nil {
		return nil, tf.Output{}, err
	}

	var squaredData []float64
	switch reflect.TypeOf(t.Value()).Kind() {
	case reflect.Slice:
		tensorData := t.Value().([][]float64)
		for _, f := range tensorData {
			norm := floats.Norm(f, 2)
			squaredData = append(squaredData, norm*norm)
		}
	}

	return MakeTensorAndOutput(s, squaredData)
}

//Log - accepts a tensor of [][]float64 and squares []float64[0]
func Log(s *op.Scope, value interface{}) (*tf.Tensor, tf.Output, error) {
	t, err := tf.NewTensor(value)
	if err != nil {
		return nil, tf.Output{}, err
	}

	var sData []float64
	switch reflect.TypeOf(t.Value()).Kind() {
	case reflect.Slice:
		tensorData := t.Value().([][]float64)
		for _, f := range tensorData {
			log := floats.LogSumExp(f)
			sData = append(sData, log)
		}
	}

	return MakeTensorAndOutput(s, sData)
}

//Square - accepts a tensor of [][]float64 and squares []float64[0]
func Square(s *op.Scope, value interface{}) (*tf.Tensor, tf.Output, error) {
	t, err := tf.NewTensor(value)
	if err != nil {
		return nil, tf.Output{}, err
	}
	var squaredData [][]float64
	switch reflect.TypeOf(t.Value()).Kind() {
	case reflect.Slice:
		tensorData := t.Value().([][]float64)
		for _, f := range tensorData {
			var square []float64
			square = append(square, f[0]*f[0])
			squaredData = append(squaredData, square)
		}
	}

	return MakeTensorAndOutput(s, squaredData)
}

//Square - accepts a tensor of [][]float64 and squares []float64[0]
func SquareT(t *tf.Tensor) (*tf.Tensor, error) {
	var squaredData [][]float64
	switch reflect.TypeOf(t.Value()).Kind() {
	case reflect.Slice:
		tensorData := t.Value().([][]float64)
		for _, f := range tensorData {
			var square []float64
			square = append(square, f[0]*f[0])
			squaredData = append(squaredData, square)
		}
	}

	tensor, e := tf.NewTensor(squaredData)
	if e != nil {
		return nil, e
	}
	return tensor, nil
}

//ReduceMean - mean accross elements
func ReduceMean(t *tf.Tensor) (*tf.Tensor, error) {
	var meanData []float64
	switch reflect.TypeOf(t.Value()).Kind() {
	case reflect.Slice:
		tensorData := t.Value().([][]float64)

		for i, f := range tensorData {
			var mean float64
			//mean += float64(i) * f
			fmt.Println(i, f)
			meanData = append(meanData, mean)
		}
	}

	tensor, e := tf.NewTensor(meanData)
	if e != nil {
		return nil, e
	}
	return tensor, nil
}

func RandomSample(t *tf.Tensor) (*tf.Tensor, error) {

	//var randomData [][]float64
	tensorData := t.Value().([][]float64)

	// for i, f := range tensorData {
	// 	var mean float64
	// 	//mean += float64(i) * f
	// 	fmt.Println(i, f)
	// 	meanData = append(meanData, mean)
	// }

	tensor, e := tf.NewTensor(tensorData)
	if e != nil {
		return nil, e
	}
	return tensor, nil
}
