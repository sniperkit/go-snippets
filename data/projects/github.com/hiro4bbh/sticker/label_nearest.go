package sticker

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
)

// SimCountPair is the data structure for float32 similarity and uint32 count.
type SimCountPair struct {
	Sim   float32
	Count uint32
}

// LabelNearestContext is a context used in inference.
// This is not protected by any mutex, so this should not be accessed by multiple goroutines.
type LabelNearestContext []SimCountPair

// LabelNearest is the sparse weighted nearest neighborhood.
// Sparse means in the following 3 reasons:
//   (i) The similarity used in constructing the nearest neighborhood defined by the inner-product on the features activated in the given entry.
//   (ii) The positive labels assigned to the entries in the training dataset are only used for averaging.
//   (iii) The entries in the training dataset whose inner-product is not positive are not used.
//
// LabelNearest is only the optimized data structure of the training dataset for searching nearest neighborhood.
type LabelNearest struct {
	// NfeaturesList is the slice of the number of features contained in each training data entry.
	NfeaturesList []uint32
	// FeatureIndexList is the map from the feature to the feature index.
	// The feature index contains the list of the pair of the entry index and the corresponding feature value.
	// The feature vector for the each entry is normalized for computing cos similarity effectively in the inference.
	FeatureIndexList map[uint32]KeyValues32
	// LabelVectors is the label vectors in the training dataset.
	LabelVectors LabelVectors
}

// TrainLabelNearest returns an trained LabelNearest on the given training dataset ds.
//
// Currently, this function returns no error.
func TrainLabelNearest(ds *Dataset, debug *log.Logger) (*LabelNearest, error) {
	nfeaturesList := make([]uint32, ds.Size())
	lenX := make([]float32, ds.Size())
	for i, xi := range ds.X {
		nfeaturesList[i] = uint32(len(xi))
		l := float32(0.0)
		for _, xipair := range xi {
			l += xipair.Value * xipair.Value
		}
		lenX[i] = Sqrt32(l)
	}
	featureIndexList := make(map[uint32]KeyValues32)
	featureActs := make([]float32, 0, len(ds.X))
	for i, xi := range ds.X {
		lenxi := lenX[i]
		for _, xipair := range xi {
			featureIndexList[xipair.Key] = append(featureIndexList[xipair.Key], KeyValue32{uint32(i), xipair.Value / lenxi})
		}
		featureActs = append(featureActs, float32(len(xi)))
	}
	featureActsMin, featureActsQ25, featureActsMed, featureActsQ75, featureActsMax, featureActsAvg := SummarizeFloat32Slice(featureActs)
	if debug != nil {
		debug.Printf("TrainLabelNearest: feature occurrence summary: min=%.0f, 1st-qu=%.0f, median=%.0f, 3rd-qu=%.0f, max=%.0f, avg=%.2f", featureActsMin, featureActsQ25, featureActsMed, featureActsQ75, featureActsMax, featureActsAvg)
	}
	labelVectors := make(LabelVectors, 0, len(ds.Y))
	for _, yi := range ds.Y {
		yihat := make(LabelVector, len(yi))
		copy(yihat, yi)
		labelVectors = append(labelVectors, yihat)
	}
	return &LabelNearest{
		NfeaturesList:    nfeaturesList,
		FeatureIndexList: featureIndexList,
		LabelVectors:     labelVectors,
	}, nil
}

// DecodeLabelNearestWithGobDecoder decodes LabelNearest using decoder.
//
// This function returns an error in decoding.
func DecodeLabelNearestWithGobDecoder(model *LabelNearest, decoder *gob.Decoder) error {
	model.NfeaturesList = []uint32{}
	if err := decoder.Decode(&model.NfeaturesList); err != nil {
		return fmt.Errorf("DecodeLabelNearest: NfeaturesList: %s", err)
	}
	var nfeatures uint32
	if err := decoder.Decode(&nfeatures); err != nil {
		return fmt.Errorf("DecodeLabelNearest: nfeatures: %s", err)
	}
	model.FeatureIndexList = make(map[uint32]KeyValues32, nfeatures)
	for j := uint32(0); j < nfeatures; j++ {
		var feature uint32
		if err := decoder.Decode(&feature); err != nil {
			return fmt.Errorf("DecodeLabelNearest: #%d feature: %s", j, err)
		}
		featureIndex := KeyValues32{}
		if err := decoder.Decode(&featureIndex); err != nil {
			return fmt.Errorf("DecodeLabelNearest: FeatureIndexList[%d]: %s", feature, err)
		}
		model.FeatureIndexList[feature] = featureIndex
	}
	model.LabelVectors = make(LabelVectors, 0, len(model.NfeaturesList))
	for i := 0; i < len(model.NfeaturesList); i++ {
		labelVector := LabelVector{}
		if err := decoder.Decode(&labelVector); err != nil {
			return fmt.Errorf("DecodeLabelNearest: LabelVectors[%d]: %s", i, err)
		}
		model.LabelVectors = append(model.LabelVectors, labelVector)
	}
	return nil
}

// DecodeLabelNearest decodes LabelNearest from r.
// Directly passing *os.File used by a gob.Decoder to this function causes mysterious errors.
// Thus, if users use gob.Decoder, then they should call DecodeLabelNearestWithGobDecoder.
//
// This function returns an error in decoding.
func DecodeLabelNearest(model *LabelNearest, r io.Reader) error {
	return DecodeLabelNearestWithGobDecoder(model, gob.NewDecoder(r))
}

// EncodeLabelNearestWithGobEncoder decodes LabelNearest using encoder.
//
// This function returns an error in decoding.
func EncodeLabelNearestWithGobEncoder(model *LabelNearest, encoder *gob.Encoder) error {
	if err := encoder.Encode(model.NfeaturesList); err != nil {
		return fmt.Errorf("EncodeLabelNearest: NfeatureList: %s", err)
	}
	nfeatures := uint32(len(model.FeatureIndexList))
	if err := encoder.Encode(nfeatures); err != nil {
		return fmt.Errorf("EncodeLabelNearest: nfeatures: %s", err)
	}
	j := 0
	for feature, featureIndex := range model.FeatureIndexList {
		if err := encoder.Encode(feature); err != nil {
			return fmt.Errorf("EncodeLabelNearest: #%d feature: %s", j, err)
		}
		if err := encoder.Encode(featureIndex); err != nil {
			return fmt.Errorf("EncodeLabelNearest: FeatureIndexList[%d]: %s", feature, err)
		}
		j++
	}
	for i, labelVector := range model.LabelVectors {
		if err := encoder.Encode(labelVector); err != nil {
			return fmt.Errorf("EncodeLabelNearest: LabelVectors[%d]: %s", i, err)
		}
	}
	return nil
}

// EncodeLabelNearest encodes LabelNearest to w.
// Directly passing *os.File used by a gob.Encoder to this function causes mysterious errors.
// Thus, if users use gob.Encoder, then they should call EncodeLabelNearestWithGobEncoder.
//
// This function returns an error in encoding.
func EncodeLabelNearest(model *LabelNearest, w io.Writer) error {
	return EncodeLabelNearestWithGobEncoder(model, gob.NewEncoder(w))
}

// GobEncode returns the error always, because users should encode large LabelNearest objects with EncodeLabelNearest.
func (model *LabelNearest) GobEncode() ([]byte, error) {
	return nil, fmt.Errorf("LabelNearest should be encoded with EncodeLabelNearest")
}

// FindNearests returns the S nearest entries with each similarity for the given entry.
// See Predict for hyper-parameter details.
func (model *LabelNearest) FindNearests(x FeatureVector, S uint, beta float32) KeyValues32 {
	return model.FindNearestsWithContext(x, S, beta, model.NewContext())
}

// FindNearestsWithContext is FindNearests with the specified LabelNearestContext.
func (model *LabelNearest) FindNearestsWithContext(x FeatureVector, S uint, beta float32, ctx LabelNearestContext) KeyValues32 {
	simCounts := []SimCountPair(ctx)
	for _, xpair := range x {
		featureIndex := model.FeatureIndexList[xpair.Key]
		for _, indexValue := range featureIndex {
			simCounts[indexValue.Key].Sim += indexValue.Value * xpair.Value
			simCounts[indexValue.Key].Count++
		}
	}
	indexSimsTopS := make(KeyValues32, 0, S)
	for i, simCount := range simCounts {
		if sim, count := simCount.Sim, simCount.Count; count > 0 {
			if sim > 0.0 {
				// For efficiency, call Pow32 as short-cut style.
				jaccard := float32(count) / float32(uint32(len(x))+model.NfeaturesList[i]-count)
				if beta == 0 {
				} else if beta == 1 {
					sim *= jaccard
				} else {
					sim *= Pow32(jaccard, beta)
				}
				if len(indexSimsTopS) == 0 {
					indexSimsTopS = append(indexSimsTopS, KeyValue32{uint32(i), sim})
				} else if indexSimsTopS[len(indexSimsTopS)-1].Value > sim {
					if len(indexSimsTopS) < cap(indexSimsTopS) {
						indexSimsTopS = append(indexSimsTopS, KeyValue32{uint32(i), sim})
					}
				} else {
					for rank := 0; rank < len(indexSimsTopS); rank++ {
						if sim >= indexSimsTopS[rank].Value {
							if len(indexSimsTopS) < cap(indexSimsTopS) {
								indexSimsTopS = append(indexSimsTopS, KeyValue32{0, 0})
							}
							copy(indexSimsTopS[rank+1:], indexSimsTopS[rank:])
							indexSimsTopS[rank] = KeyValue32{uint32(i), sim}
							break
						}
					}
				}
			}
			simCounts[i] = SimCountPair{0.0, 0}
		}
	}
	return indexSimsTopS
}

// NewContext returns a new context for some inference memory.
func (model *LabelNearest) NewContext() LabelNearestContext {
	return make([]SimCountPair, len(model.LabelVectors))
}

// Predict returns the results for the given data entry x with the sparse S-nearest neighborhood.
// The returned results are the top-K labels, the label histogram, and the slice of the data entry index and its similarity.
//
// alpha is the smoothing parameter for weighting the votes by each neighbor.
// beta is the smoothing parameter for balancing the Jaccard similarity and the cosine similarity.
func (model *LabelNearest) Predict(x FeatureVector, K, S uint, alpha, beta float32) (LabelVector, map[uint32]float32, KeyValues32) {
	return model.PredictWithContext(x, K, S, alpha, beta, model.NewContext())
}

// PredictWithContext is Predict with the specified LabelNearestContext.
func (model *LabelNearest) PredictWithContext(x FeatureVector, K, S uint, alpha, beta float32, ctx LabelNearestContext) (LabelVector, map[uint32]float32, KeyValues32) {
	indexSimsTopS := model.FindNearestsWithContext(x, S, beta, ctx)
	labelHist := make(map[uint32]float32)
	xlen := float32(0.0)
	for _, xpair := range x {
		xlen += xpair.Value * xpair.Value
	}
	xlen = Sqrt32(xlen)
	for _, indexSim := range indexSimsTopS {
		value := Pow32(indexSim.Value/xlen, alpha)
		for _, label := range model.LabelVectors[indexSim.Key] {
			labelHist[label] += value
		}
	}
	return RankTopK(labelHist, K), labelHist, indexSimsTopS
}

// PredictAll returns the top-K labels for each data entry in X with the sparse S-nearest neighborhood.
// See Predict for hyper-parameter details.
func (model *LabelNearest) PredictAll(X FeatureVectors, K, S uint, alpha, beta float32) LabelVectors {
	Yhat := make(LabelVectors, 0, len(X))
	ctx := model.NewContext()
	for _, xi := range X {
		yihat, _, _ := model.PredictWithContext(xi, K, S, alpha, beta, ctx)
		Yhat = append(Yhat, yihat)
	}
	return Yhat
}
