package sticker

import (
	"bytes"
	"encoding/gob"
	"log"
	"testing"

	"github.com/hiro4bbh/go-assert"
)

func TestLabelOnePrune(t *testing.T) {
	model3 := &LabelOne{
		Params: NewLabelOneParameters(),
		Biases: []float32{1.0, 2.0, 3.0},
		WeightLists: map[uint32]KeyValues32{
			1: {KeyValue32{2, 1.0}, KeyValue32{3, 2.0}},
			3: {KeyValue32{1, 1.0}, KeyValue32{3, 2.0}},
		},
		Labels: LabelVector{1, 3, 2},
		Summaries: []map[string]interface{}{
			{}, {}, {},
		},
	}
	model1 := &LabelOne{
		Params: NewLabelOneParameters(),
		Biases: []float32{1.0},
		WeightLists: map[uint32]KeyValues32{
			3: {KeyValue32{1, 1.0}},
		},
		Labels: LabelVector{1},
		Summaries: []map[string]interface{}{
			{},
		},
	}
	goassert.New(t, model1).Equal(model3.Prune(1))
	goassert.New(t, model3).Equal(model3.Prune(3))
	goassert.New(t, model3).Equal(model3.Prune(5))
}

func TestTrainLabelOne(t *testing.T) {
	n := 10
	ds := &Dataset{
		X: make(FeatureVectors, 2*n),
		Y: make(LabelVectors, 2*n),
	}
	for i := 0; i < n; i++ {
		ds.X[2*i], ds.Y[2*i] = FeatureVector{KeyValue32{0, 1.0}}, LabelVector{0, 1}
		ds.X[2*i+1], ds.Y[2*i+1] = FeatureVector{KeyValue32{1, 1.0}}, LabelVector{0, 2}
	}
	params := NewLabelOneParameters()
	params.T = 5
	model := goassert.New(t).SucceedNew(TrainLabelOne(ds, params, nil)).(*LabelOne)
	goassert.New(t, uint(3)).Equal(model.Nrounds())
	YhatK2T5, YhatK5T5, YhatK2T2 := make(LabelVectors, 2*n), make(LabelVectors, 2*n), make(LabelVectors, 2*n)
	for i := 0; i < n; i++ {
		YhatK2T5[2*i], YhatK2T5[2*i+1] = LabelVector{0, 1}, LabelVector{0, 2}
		YhatK5T5[2*i], YhatK5T5[2*i+1] = LabelVector{0, 1, 2, ^uint32(0), ^uint32(0)}, LabelVector{0, 2, 1, ^uint32(0), ^uint32(0)}
		YhatK2T2[2*i], YhatK2T2[2*i+1] = LabelVector{0, 1}, LabelVector{0, 1}
	}
	goassert.New(t, YhatK2T5).Equal(model.PredictAll(ds.X, 2, 5))
	goassert.New(t, YhatK5T5).Equal(model.PredictAll(ds.X, 5, 5))
	goassert.New(t, YhatK2T2).Equal(model.PredictAll(ds.X, 2, 2))
	// debug logger is tested in TestDecodeEncodeLabelOne.
}

func TestDecodeEncodeLabelOne(t *testing.T) {
	// Test a default use-case with debug logger.
	// This test assumes that TrainLabelOne is already tested.
	n := 10
	ds := &Dataset{
		X: make(FeatureVectors, 2*n),
		Y: make(LabelVectors, 2*n),
	}
	for i := 0; i < n; i++ {
		ds.X[2*i], ds.Y[2*i] = FeatureVector{KeyValue32{0, 1.0}}, LabelVector{0, 1}
		ds.X[2*i+1], ds.Y[2*i+1] = FeatureVector{KeyValue32{1, 1.0}}, LabelVector{0, 2}
	}
	var debugBuf bytes.Buffer
	params := NewLabelOneParameters()
	params.T = 3
	model := goassert.New(t).SucceedNew(TrainLabelOne(ds, params, log.New(&debugBuf, "", 0))).(*LabelOne)
	goassert.New(t, true).Equal(debugBuf.String() != "")
	var buf bytes.Buffer
	goassert.New(t).SucceedWithoutError(EncodeLabelOne(model, &buf))
	var decodedModel LabelOne
	goassert.New(t).SucceedWithoutError(DecodeLabelOne(&decodedModel, &buf))
	goassert.New(t, model).Equal(&decodedModel)
	// gob.Decoder.Decode won't call LabelOne.GobDecode, because the encoder did not encode LabelOne.
	goassert.New(t, "LabelOne should be encoded with EncodeLabelOne").ExpectError(gob.NewEncoder(&buf).Encode(&decodedModel))
}
