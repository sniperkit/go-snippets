package sticker

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sort"
)

// LabelOneParameters is the parameters for LabelOne.
type LabelOneParameters struct {
	// ClassifierTrainerName is the used BinaryClassifierTrainer name.
	ClassifierTrainerName string
	// C is the penalty parameter for BinaryClassifierTrainer.
	C float32
	// Epsilon is the tolerance parameter for BinaryClassifierTrainer.
	Epsilon float32
	// T is the maximum number of the rounds, which is equal to the maximum number of the target labels.
	T uint
}

// NewLabelOneParameters returns an LabelOneParameters initialized with the default values.
func NewLabelOneParameters() *LabelOneParameters {
	return &LabelOneParameters{
		ClassifierTrainerName: "L1Logistic_PrimalSGD",
		C:       float32(1.0),
		Epsilon: float32(1.0e-05),
		T:       uint(100),
	}
}

// LabelOne is the One-versus-Rest classifier for multi-label ranking.
// The t-th classifier (t = 1, ..., T) is the classifier for the top-t frequently occurring label.
type LabelOne struct {
	// Params is the used LabelOneParameters.
	Params *LabelOneParameters
	// Biases is the bias slice used by splitters on each classifier.
	Biases []float32
	// Weights is the weight sparse matrix used by each classifier.
	// This is the map from the feature key to the (roundID, the weight on the feature of #roundID splitter) slice.
	// This data structure reduces the number of times that the classifier accesses the golang's map a lot.
	WeightLists map[uint32]KeyValues32
	// Labels is the label slice used in each classifier.
	// The t-th label is the target label of the t-th classifier.
	Labels LabelVector
	// The following members are not required.
	//
	// Summaries is the summary object slice for each boosting round.
	// The entries in this summary is considered to provide compact and useful information in best-effort, so this specification would be loose and rapidly changing.
	Summaries []map[string]interface{}
}

// TrainLabelOne returns an trained LabelOne on the given dataset ds.
func TrainLabelOne(ds *Dataset, params *LabelOneParameters, debug *log.Logger) (*LabelOne, error) {
	classifierTrainer, ok := BinaryClassifierTrainers[params.ClassifierTrainerName]
	if !ok {
		return nil, fmt.Errorf("unknown ClassifierTrainerName: %s", params.ClassifierTrainerName)
	}
	n := ds.Size()
	biases, weightLists := []float32{}, make(map[uint32]KeyValues32)
	labels := []uint32{}
	summaries := []map[string]interface{}{}
	// Collect the label frequencies.
	labelFreqs := make(map[uint32]float32)
	for _, yi := range ds.Y {
		for _, label := range yi {
			labelFreqs[label]++
		}
	}
	T := params.T
	if T > uint(len(labelFreqs)) {
		T = uint(len(labelFreqs))
	}
	labelTopT := RankTopK(labelFreqs, T)
	for t := uint(1); t <= T; t++ {
		targetLabel := labelTopT[t-1]
		// Assign each data point to the positive class i.f.f. it has the target label.
		deltas := make([]bool, n)
		deltaFreq := make(map[bool]int)
		for i, yi := range ds.Y {
			delta := false
			for _, label := range yi {
				if targetLabel == label {
					delta = true
					break
				}
			}
			deltas[i] = delta
			deltaFreq[delta]++
		}
		// Training the splitter.
		if debug != nil {
			debug.Printf("TrainLabelOne: t=%d: training the splitter on %d negative(s) and %d positive(s) ...", t, deltaFreq[false], deltaFreq[true])
		}
		splitter, err := classifierTrainer(ds.X, deltas, params.C, params.Epsilon, debug)
		if err != nil {
			return nil, fmt.Errorf("BinaryClassifierTrainer(%s): %s", params.ClassifierTrainerName, err)
		}
		tn, fn, fp, tp, _, _ := splitter.ReportPerformance(ds.X, deltas)
		if debug != nil {
			debug.Printf("TrainLabelOne: t=%d: trained the splitter (tn=%d, fn=%d, fp=%d, tp=%d) ...", t, tn, fn, fp, tp)
		}
		// Append the round information.
		biases = append(biases, splitter.Bias)
		for feature, value := range splitter.Weight {
			weightLists[feature] = append(weightLists[feature], KeyValue32{uint32(t - 1), value})
		}
		labels = append(labels, targetLabel)
		summary := make(map[string]interface{})
		summary["splitPerf"] = map[string]int{"tn": int(tn), "fn": int(fn), "fp": int(fp), "tp": int(tp)}
		summaries = append(summaries, summary)
	}
	return &LabelOne{
		Params:      params,
		Biases:      biases,
		WeightLists: weightLists,
		Labels:      labels,
		Summaries:   summaries,
	}, nil
}

// DecodeLabelOneWithGobDecoder decodes LabelOne using decoder.
//
// This function returns an error in decoding.
func DecodeLabelOneWithGobDecoder(model *LabelOne, decoder *gob.Decoder) error {
	model.Params = &LabelOneParameters{}
	if err := decoder.Decode(model.Params); err != nil {
		return fmt.Errorf("DecodeLabelOne: Params: %s", err)
	}
	if err := decoder.Decode(&model.Biases); err != nil {
		return fmt.Errorf("DecodeLabelOne: Biases: %s", err)
	}
	var lenWeightLists int
	if err := decoder.Decode(&lenWeightLists); err != nil {
		return fmt.Errorf("DecodeLabelOne: len(WeightLists): %s", err)
	}
	model.WeightLists = make(map[uint32]KeyValues32)
	for i := 0; i < lenWeightLists; i++ {
		var feature uint32
		if err := decoder.Decode(&feature); err != nil {
			return fmt.Errorf("DecodeLabelOne: #%d WeightList: %s", i, err)
		}
		var weightList KeyValues32
		if err := decoder.Decode(&weightList); err != nil {
			return fmt.Errorf("DecodeLabelOne: #%d WeightList: %s", i, err)
		}
		model.WeightLists[feature] = weightList
	}
	if err := decoder.Decode(&model.Labels); err != nil {
		return fmt.Errorf("DecodeLabelOne: Labels: %s", err)
	}
	if err := decoder.Decode(&model.Summaries); err != nil {
		return fmt.Errorf("DecodeLabelOne: Summaries: %s", err)
	}
	return nil
}

// DecodeLabelOne decodes LabelOne from r.
// Directly passing *os.File used by a gob.Decoder to this function causes mysterious errors.
// Thus, if users use gob.Decoder, then they should call DecodeLabelOneWithGobDecoder.
//
// This function returns an error in decoding.
func DecodeLabelOne(model *LabelOne, r io.Reader) error {
	return DecodeLabelOneWithGobDecoder(model, gob.NewDecoder(r))
}

// EncodeLabelOneWithGobEncoder decodes LabelOne using encoder.
//
// This function returns an error in decoding.
func EncodeLabelOneWithGobEncoder(model *LabelOne, encoder *gob.Encoder) error {
	if err := encoder.Encode(model.Params); err != nil {
		return fmt.Errorf("EncodeLabelOne: Params: %s", err)
	}
	if err := encoder.Encode(model.Biases); err != nil {
		return fmt.Errorf("EncodeLabelOne: Biases: %s", err)
	}
	features := make([]int, 0, len(model.WeightLists))
	for feature := range model.WeightLists {
		features = append(features, int(feature))
	}
	sort.Ints(features)
	if err := encoder.Encode(len(model.WeightLists)); err != nil {
		return fmt.Errorf("EncodeLabelOne: len(WeightLists): %s", err)
	}
	for i, feature := range features {
		if err := encoder.Encode(uint32(feature)); err != nil {
			return fmt.Errorf("EncodeLabelOne: #%d WeightList: %s", i, err)
		}
		if err := encoder.Encode(model.WeightLists[uint32(feature)]); err != nil {
			return fmt.Errorf("EncodeLabelOne: #%d WeightList: %s", i, err)
		}
	}
	if err := encoder.Encode(model.Labels); err != nil {
		return fmt.Errorf("EncodeLabelOne: Labels: %s", err)
	}
	if err := encoder.Encode(model.Summaries); err != nil {
		return fmt.Errorf("EncodeLabelOne: Summaries: %s", err)
	}
	return nil
}

// EncodeLabelOne encodes LabelOne to w.
// Directly passing *os.File used by a gob.Encoder to this function causes mysterious errors.
// Thus, if users use gob.Encoder, then they should call EncodeLabelOneWithGobEncoder.
//
// This function returns an error in encoding.
func EncodeLabelOne(model *LabelOne, w io.Writer) error {
	return EncodeLabelOneWithGobEncoder(model, gob.NewEncoder(w))
}

// GobEncode returns the error always, because users should encode large LabelOne objects with EncodeLabelOne.
func (model *LabelOne) GobEncode() ([]byte, error) {
	return nil, fmt.Errorf("LabelOne should be encoded with EncodeLabelOne")
}

// Nrounds return the number of the rounds.
func (model *LabelOne) Nrounds() uint {
	return uint(len(model.Labels))
}

// Predict returns the top-K predicted labels for the given data entry x with the first T rounds.
func (model *LabelOne) Predict(x FeatureVector, K uint, T uint) LabelVector {
	if T > model.Nrounds() {
		T = model.Nrounds()
	}
	z := make([]float32, T)
	for _, xpair := range x {
		for _, weightpair := range model.WeightLists[xpair.Key] {
			if uint(weightpair.Key) >= T {
				break
			}
			z[weightpair.Key] += weightpair.Value * xpair.Value
		}
	}
	y := make(map[uint32]float32)
	for t, zt := range z {
		zt += model.Biases[t]
		y[model.Labels[t]] += zt
	}
	return RankTopK(y, K)
}

// PredictAll returns the slice of the top-K predicted labels for each data entry in X with the first T rounds.
func (model *LabelOne) PredictAll(X FeatureVectors, K uint, T uint) LabelVectors {
	Y := make(LabelVectors, 0, len(X))
	for _, xi := range X {
		Y = append(Y, model.Predict(xi, K, T))
	}
	return Y
}

// Prune returns the pruned LabelOne which has at most T rounds.
func (model *LabelOne) Prune(T uint) *LabelOne {
	if T > model.Nrounds() {
		T = model.Nrounds()
	}
	newModel := &LabelOne{
		Params:      model.Params,
		Biases:      model.Biases[:T],
		WeightLists: make(map[uint32]KeyValues32),
		Labels:      model.Labels[:T],
		Summaries:   model.Summaries[:T],
	}
	for feature, weightList := range model.WeightLists {
		l := 0
		for _, weightpair := range weightList {
			if weightpair.Key > uint32(T) {
				break
			}
			l++
		}
		if l > 0 {
			newModel.WeightLists[feature] = weightList[:l]
		}
	}
	return newModel
}
