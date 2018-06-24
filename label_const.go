package sticker

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
)

// LabelConst is the multi-label constant model.
type LabelConst struct {
	// LabelList and LabelFreqList are the label and its frequency list in descending order in the training set occurrences.
	LabelList     LabelVector
	LabelFreqList []float32
}

// TrainLabelConst returns an trained LabelConst on the given dataset ds.
func TrainLabelConst(ds *Dataset, debug *log.Logger) (*LabelConst, error) {
	labelFreqs := make(map[uint32]float32)
	for _, yi := range ds.Y {
		for _, label := range yi {
			labelFreqs[label]++
		}
	}
	labelList := RankTopK(labelFreqs, uint(len(labelFreqs)))
	labelFreqList := make([]float32, uint(len(labelFreqs)))
	for rank, label := range labelList {
		labelFreqList[rank] = labelFreqs[label]
	}
	return &LabelConst{
		LabelList:     labelList,
		LabelFreqList: labelFreqList,
	}, nil
}

// DecodeLabelConstWithGobDecoder decodes LabelConst using decoder.
//
// This function returns an error in decoding.
func DecodeLabelConstWithGobDecoder(model *LabelConst, decoder *gob.Decoder) error {
	if err := decoder.Decode(&model.LabelList); err != nil {
		return fmt.Errorf("DecodeLabelBoost: LabelList: %s", err)
	}
	if err := decoder.Decode(&model.LabelFreqList); err != nil {
		return fmt.Errorf("DecodeLabelBoost: LabelFreqList: %s", err)
	}
	return nil
}

// DecodeLabelConst decodes LabelConst from r.
// Directly passing *os.File used by a gob.Decoder to this function causes mysterious errors.
// Thus, if users use gob.Decoder, then they should call DecodeLabelConstWithGobDecoder.
//
// This function returns an error in decoding.
func DecodeLabelConst(model *LabelConst, r io.Reader) error {
	return DecodeLabelConstWithGobDecoder(model, gob.NewDecoder(r))
}

// EncodeLabelConstWithGobEncoder decodes LabelConst using encoder.
//
// This function returns an error in decoding.
func EncodeLabelConstWithGobEncoder(model *LabelConst, encoder *gob.Encoder) error {
	if err := encoder.Encode(model.LabelList); err != nil {
		return fmt.Errorf("EncodeLabelBoost: LabelList: %s", err)
	}
	if err := encoder.Encode(model.LabelFreqList); err != nil {
		return fmt.Errorf("EncodeLabelBoost: LabelFreqList: %s", err)
	}
	return nil
}

// EncodeLabelConst encodes LabelConst to w.
// Directly passing *os.File used by a gob.Encoder to this function causes mysterious errors.
// Thus, if users use gob.Encoder, then they should call EncodeLabelBoostWithGobEncoder.
//
// This function returns an error in encoding.
func EncodeLabelConst(model *LabelConst, w io.Writer) error {
	return EncodeLabelConstWithGobEncoder(model, gob.NewEncoder(w))
}

// GobEncode returns the error always, because users should encode large LabelConst objects with EncodeLabelConst.
func (model *LabelConst) GobEncode() ([]byte, error) {
	return nil, fmt.Errorf("LabelConst should be encoded with EncodeLabelConst")
}

// PredictAll returns the top-K labels for each data entry in X.
func (model *LabelConst) PredictAll(X FeatureVectors, K uint) LabelVectors {
	predictedLabels := make(LabelVector, K)
	Kmax := K
	if Kmax > uint(len(model.LabelList)) {
		Kmax = uint(len(model.LabelList))
	}
	copy(predictedLabels, model.LabelList[:Kmax])
	for rank := uint(len(model.LabelList)); rank < K; rank++ {
		predictedLabels[rank] = ^uint32(0)
	}
	Y := make(LabelVectors, 0, len(X))
	for range X {
		yi := make(LabelVector, K)
		copy(yi, predictedLabels)
		Y = append(Y, yi)
	}
	return Y
}
