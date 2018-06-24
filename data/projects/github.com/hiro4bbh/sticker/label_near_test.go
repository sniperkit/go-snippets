package sticker

import (
	"bytes"
	"encoding/gob"
	"log"
	"testing"

	"github.com/hiro4bbh/go-assert"
)

func TestDecodeEncodeLabelNear(t *testing.T) {
	ds := &Dataset{
		X: FeatureVectors{
			FeatureVector{KeyValue32{2, 2.0}}, FeatureVector{KeyValue32{1, 1.0}},
			FeatureVector{KeyValue32{1, 1.0}, KeyValue32{3, 3.0}}, FeatureVector{KeyValue32{4, 4.0}},
			FeatureVector{KeyValue32{1, 1.0}, KeyValue32{3, 3.0}, KeyValue32{5, 5.0}},
			FeatureVector{KeyValue32{1, 1.0}, KeyValue32{3, 3.0}, KeyValue32{5, 5.0}, KeyValue32{6, 6.0}},
		},
		Y: LabelVectors{
			LabelVector{2}, LabelVector{1},
			LabelVector{3}, LabelVector{4},
			LabelVector{5},
			LabelVector{6},
		},
	}
	var debugBuf bytes.Buffer
	params := NewLabelNearParameters()
	model := goassert.New(t).SucceedNew(TrainLabelNear(ds, params, log.New(&debugBuf, "", 0))).(*LabelNear)
	goassert.New(t, true).Equal(debugBuf.String() != "")
	var buf bytes.Buffer
	goassert.New(t).SucceedWithoutError(EncodeLabelNear(model, &buf))
	var decodedModel LabelNear
	goassert.New(t).SucceedWithoutError(DecodeLabelNear(&decodedModel, &buf))
	goassert.New(t, model).Equal(&decodedModel)
	// gob.Decoder.Decode won't call LabelNear.GobDecode, because the encoder did not encode LabelNear.
	goassert.New(t, "LabelNear should be encoded with EncodeLabelNear").ExpectError(gob.NewEncoder(&buf).Encode(&decodedModel))
}

func TestLabelNearFindNears(t *testing.T) {
	// This only tests the efficient computation of Pow32.
	ds := &Dataset{
		X: FeatureVectors{FeatureVector{KeyValue32{1, 1.0}, KeyValue32{2, 1.0}, KeyValue32{3, 1.0}, KeyValue32{4, 1.0}}},
		Y: LabelVectors{LabelVector{1}},
	}
	params := NewLabelNearParameters()
	model := goassert.New(t).SucceedNew(TrainLabelNear(ds, params, nil)).(*LabelNear)
	goassert.New(t, KeyValues32{KeyValue32{0, 0.5}}).Equal(model.FindNears(FeatureVector{KeyValue32{1, 1.0}}, 5, 1, 0.0))
	goassert.New(t, KeyValues32{KeyValue32{0, 0.125}}).Equal(model.FindNears(FeatureVector{KeyValue32{1, 1.0}}, 5, 1, 1.0))
	goassert.New(t, KeyValues32{KeyValue32{0, 0.03125}}).Equal(model.FindNears(FeatureVector{KeyValue32{1, 1.0}}, 5, 1, 2.0))
}

func TestLabelNearPredictAll(t *testing.T) {
	ds := &Dataset{
		X: FeatureVectors{
			FeatureVector{KeyValue32{2, 2.0}}, FeatureVector{KeyValue32{1, 1.0}},
			FeatureVector{KeyValue32{1, 1.0}, KeyValue32{3, 3.0}}, FeatureVector{KeyValue32{4, 4.0}},
			FeatureVector{KeyValue32{1, 1.0}, KeyValue32{3, 3.0}, KeyValue32{5, 5.0}},
			FeatureVector{KeyValue32{1, 1.0}, KeyValue32{3, 3.0}, KeyValue32{5, 5.0}, KeyValue32{6, 6.0}},
		},
		Y: LabelVectors{
			LabelVector{2}, LabelVector{1},
			LabelVector{3}, LabelVector{4},
			LabelVector{5},
			LabelVector{6},
		},
	}
	params := NewLabelNearParameters()
	model := goassert.New(t).SucceedNew(TrainLabelNear(ds, params, nil)).(*LabelNear)
	goassert.New(t, LabelVectors{
		LabelVector{1, ^uint32(0), ^uint32(0)},
		LabelVector{4, ^uint32(0), ^uint32(0)},
		LabelVector{^uint32(0), ^uint32(0), ^uint32(0)},
		LabelVector{3, ^uint32(0), ^uint32(0)},
		LabelVector{^uint32(0), ^uint32(0), ^uint32(0)},
	}).Equal(model.PredictAll(FeatureVectors{
		FeatureVector{KeyValue32{1, 1.0}},
		FeatureVector{KeyValue32{4, 1.0}},
		FeatureVector{KeyValue32{10, 1.0}, KeyValue32{11, 1.0}},
		FeatureVector{KeyValue32{1, 1.0}, KeyValue32{3, 1.0}},
		FeatureVector{KeyValue32{1, -1.0}},
	}, 3, 5, 1, 1.0, 1.0))
	goassert.New(t, LabelVectors{
		LabelVector{1, 3, 5},
		LabelVector{4, ^uint32(0), ^uint32(0)},
		LabelVector{^uint32(0), ^uint32(0), ^uint32(0)},
		LabelVector{3, 1, 5},
		LabelVector{^uint32(0), ^uint32(0), ^uint32(0)},
	}).Equal(model.PredictAll(FeatureVectors{
		FeatureVector{KeyValue32{1, 1.0}},
		FeatureVector{KeyValue32{4, 1.0}},
		FeatureVector{KeyValue32{10, 1.0}, KeyValue32{11, 1.0}},
		FeatureVector{KeyValue32{1, 1.0}, KeyValue32{3, 1.0}},
		FeatureVector{KeyValue32{1, -1.0}},
	}, 3, 5, 3, 1.0, 1.0))
	// Test the sorted order.
	ds2 := &Dataset{
		X: FeatureVectors{
			FeatureVector{KeyValue32{1, 1.0}, KeyValue32{2, 1.0}, KeyValue32{3, 1.0}, KeyValue32{4, 1.0}, KeyValue32{5, 1.0}},
			FeatureVector{KeyValue32{1, 1.0}, KeyValue32{3, 1.0}},
			FeatureVector{KeyValue32{1, 1.0}, KeyValue32{2, 1.0}, KeyValue32{3, 1.0}, KeyValue32{5, 1.0}},
			FeatureVector{KeyValue32{1, 1.0}, KeyValue32{3, 1.0}, KeyValue32{5, 1.0}},
			FeatureVector{KeyValue32{3, 1.0}},
		},
		Y: LabelVectors{
			LabelVector{1}, LabelVector{2}, LabelVector{3}, LabelVector{4}, LabelVector{5},
		},
	}
	model2 := goassert.New(t).SucceedNew(TrainLabelNear(ds2, params, nil)).(*LabelNear)
	goassert.New(t, LabelVectors{
		LabelVector{1, 3, 4, 2, 5},
	}).Equal(model2.PredictAll(FeatureVectors{
		FeatureVector{KeyValue32{1, 1.0}, KeyValue32{2, 1.0}, KeyValue32{3, 1.0}, KeyValue32{4, 1.0}, KeyValue32{5, 1.0}},
	}, 5, 5, 5, 1.0, 1.0))
}
