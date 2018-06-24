package sticker

import (
	"bytes"
	"encoding/gob"
	"log"
	"testing"

	"github.com/hiro4bbh/go-assert"
)

func TestTrainLabelConst(t *testing.T) {
	// Test usual cases.
	n := 100
	ds := &Dataset{
		X: make(FeatureVectors, 2*n),
		Y: make(LabelVectors, 2*n),
	}
	for i := 0; i < n; i++ {
		ds.X[n*0+i], ds.X[n*1+i] = FeatureVector{KeyValue32{0, 0.0}}, FeatureVector{KeyValue32{0, 0.0}}
		ds.Y[n*0+i], ds.Y[n*1+i] = LabelVector{0, 1}, LabelVector{1, 2}
	}
	model := goassert.New(t).SucceedNew(TrainLabelConst(ds, nil)).(*LabelConst)
	goassert.New(t, LabelVector{1, 0, 2}).Equal(model.LabelList)
	goassert.New(t, []float32{200, 100, 100}).Equal(model.LabelFreqList)
	Yhat := make(LabelVectors, 2*n)
	for i := 0; i < n; i++ {
		Yhat[2*i], Yhat[2*i+1] = LabelVector{1, 0, 2, ^uint32(0)}, LabelVector{1, 0, 2, ^uint32(0)}
	}
	goassert.New(t, Yhat).Equal(model.PredictAll(ds.X, 4))
	// Test encoder/decoder.
	var buf bytes.Buffer
	goassert.New(t, "LabelConst should be encoded with EncodeLabelConst").ExpectError(gob.NewEncoder(&buf).Encode(model))
	buf.Reset()
	goassert.New(t).SucceedWithoutError(EncodeLabelConst(model, &buf))
	var decodedModel LabelConst
	// gob.Decoder.Decode won't call LabelConst.GobDecode, because the encoder did not encode LabelConst.
	goassert.New(t).SucceedWithoutError(DecodeLabelConst(&decodedModel, &buf))
	goassert.New(t, model).Equal(&decodedModel)
	// Test debug log.
	// Currently, TrainLabelConst outputs nothing on debug log.
	var debugBuf bytes.Buffer
	goassert.New(t).SucceedNew(TrainLabelConst(ds, log.New(&debugBuf, "", 0)))
	goassert.New(t, true).Equal(debugBuf.String() == "")
}
