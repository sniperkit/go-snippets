package plugin

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sort"

	"github.com/hiro4bbh/sticker"
)

// LabelBoostParameters is the parameters for LabelBoost.
type LabelBoostParameters struct {
	// RankerTrainerName is the used BinaryRankerTrainer name.
	RankerTrainerName string
	// C is the penalty parameter for BinaryRankerTrainer.
	C float32
	// Epsilon is the tolerance parameter for BinaryClassifierTrainer.
	Epsilon float32
	// NegativeSampleSize is the size of each negative sample for Multi-Label Ranking Hinge Boosting.
	// Specify 0 for Multi-Label Hinge Boosting.
	NegativeSampleSize uint
	// PainterK is the maximum number of the painted target label.
	PainterK uint
	// PainterName is the used Painter name.
	PainterName string
	// T is the maxinum number of boosting rounds.
	T uint
}

// NewLabelBoostParameters returns an LabelBoostParameters initialized with the default values.
func NewLabelBoostParameters() *LabelBoostParameters {
	return &LabelBoostParameters{
		RankerTrainerName:  "L1SVC_PrimalSGD",
		C:                  float32(1.0),
		Epsilon:            float32(0.01),
		NegativeSampleSize: uint(10),
		PainterName:        "topLabelSubSet",
		PainterK:           uint(1),
		T:                  uint(100),
	}
}

// LabelBoost is the multi-label boosting model.
type LabelBoost struct {
	// Params is the used LabelBoostParameters.
	Params *LabelBoostParameters
	// Biases is the bias slice used by splitters on each boosting round.
	Biases []float32
	// Weights is the weight sparse matrix used by splitters on each boosting round.
	// Weights is the map from the feature key to the (roundID, the weight on the feature of #roundID splitter) slice.
	// This data structure reduces the number of times that the classifier accesses the golang's map a lot.
	WeightLists map[uint32]sticker.KeyValues32
	// LabelLists is the label list slice used in each boosting round.
	// Each label list has the labels stickered to the entry if the classifier at the round returns positive score on the entry.
	LabelLists []sticker.LabelVector
	// The following members are not required.
	//
	// Summaries is the summary object slice for each boosting round.
	// The entries in this summary is considered to provide compact and useful information in best-effort, so this specification would be loose and rapidly changing.
	Summaries []map[string]interface{}
}

// TrainLabelBoost returns an trained LabelBoost on the given dataset ds.
func TrainLabelBoost(ds *sticker.Dataset, params *LabelBoostParameters, debug *log.Logger) (*LabelBoost, error) {
	rng := rand.New(rand.NewSource(0))
	painter, ok := Painters[params.PainterName]
	if !ok {
		return nil, fmt.Errorf("unknown PainterName: %s", params.PainterName)
	}
	rankerTrainer, ok := BinaryRankerTrainers[params.RankerTrainerName]
	if !ok {
		return nil, fmt.Errorf("unknown RankerTrainerName: %s", params.RankerTrainerName)
	}
	algoName := "MLRHB"
	if params.NegativeSampleSize == 0 {
		algoName = "MLHB"
	}
	n := ds.Size()
	// Initialize the margin matrix.
	Z := make([]sticker.KeyValues32, n)
	for i, yi := range ds.Y {
		zi := make(sticker.KeyValues32, 0, len(yi))
		for _, label := range yi {
			zi = append(zi, sticker.KeyValue32{label, 0.0})
		}
		Z[i] = zi
	}
	biases, weightLists := []float32{}, make(map[uint32]sticker.KeyValues32)
	labelLists := []sticker.LabelVector{}
	summaries := []map[string]interface{}{}
	for t := uint(1); t <= params.T; t++ {
		labelList := painter(ds, Z, params.PainterK, debug)
		if debug != nil {
			debug.Printf("TrainLabelBoost(%s): t=%d: Painter(%s,K=%d): selected labels: %v", algoName, t, params.PainterName, params.PainterK, labelList)
		}
		labelSet := make(map[uint32]bool)
		for _, label := range labelList {
			labelSet[label] = true
		}
		positiveLists := make(map[uint32][]int, len(labelSet))
		for i, yi := range ds.Y {
			for _, label := range yi {
				if _, ok := labelSet[label]; ok {
					positiveLists[label] = append(positiveLists[label], i)
				}
			}
		}
		var deltas []bool
		var pairIndices [][2]int
		var pairMargins, pairCs []float32
		if params.NegativeSampleSize == 0 {
			deltas = make([]bool, 0, n)
			pairIndices, pairMargins, pairCs = make([][2]int, 0, n), make([]float32, 0, n), make([]float32, 0, n)
			// Multi-Label Hinge Boosting.
			// Extract the positive and negative entries.
			for i, yi := range ds.Y {
				zi := Z[i]
				nremains, minPosZi, maxNegZi := len(labelSet), sticker.Inf32(+1.0), float32(0.0)
				for j, label := range yi {
					if _, ok := labelSet[label]; ok {
						nremains--
					}
					if zij := zi[j].Value; minPosZi > zij {
						minPosZi = zij
					}
				}
				for j := len(yi); j < len(zi); j++ {
					if zij := zi[j].Value; maxNegZi < zij {
						maxNegZi = zij
					}
				}
				deltai, weighti := false, params.C
				if nremains == 0 {
					deltai = true
				} else if nremains < len(labelSet) {
					weighti = 0.0
				}
				deltas = append(deltas, deltai)
				pairIndex := [2]int{-1, i}
				if deltai {
					pairIndex = [2]int{i, -1}
				}
				pairIndices, pairMargins, pairCs = append(pairIndices, pairIndex), append(pairMargins, minPosZi-maxNegZi), append(pairCs, weighti)
			}
			// Reweighting for balancing the positive and negative weights.
			nnegs, nposs := 0, 0
			for i, deltai := range deltas {
				if pairCs[i] > 0.0 {
					if deltai {
						nposs++
					} else {
						nnegs++
					}
				}
			}
			if nposs < nnegs {
				for i := range pairCs {
					if deltas[i] {
						pairCs[i] *= float32(nnegs) / float32(nposs)
					}
				}
			}
			if debug != nil {
				debug.Printf("TrainLabelBoost(%s): t=%d: extracted %d negative(s) and %d positive(s)", algoName, t, nnegs, nposs)
			}
		} else {
			// Multi-Label Ranking Hinge Boosting.
			npairs := 0
			for _, positiveList := range positiveLists {
				npairs += int(params.NegativeSampleSize) * len(positiveList)
			}
			pairIndices, pairMargins, pairCs = make([][2]int, 0, npairs), make([]float32, 0, npairs), make([]float32, 0, npairs)
			// Sample the positive/negative pairs.
			negativeList := make([]int, n)
			for _, label := range labelList {
				positiveList := positiveLists[label]
				for i, pi, ni := 0, 0, 0; i < n; i++ {
					if pi < len(positiveList) && positiveList[pi] == i {
						pi++
						continue
					}
					negativeList[ni] = i
					ni++
				}
				for _, i := range positiveList {
					zi, zil := Z[i], float32(0.0)
					for _, zipair := range zi {
						if zipair.Key == label {
							zil = zipair.Value
							break
						}
					}
					for nn := uint(0); nn < params.NegativeSampleSize; nn++ {
						j := negativeList[rng.Intn(n-len(positiveList))]
						zj, zjl := Z[j], float32(0.0)
						for _, zjpair := range zj {
							if zjpair.Key == label {
								zjl = zjpair.Value
								break
							}
						}
						pairIndices, pairMargins, pairCs = append(pairIndices, [2]int{i, j}), append(pairMargins, zil-zjl), append(pairCs, params.C/float32(len(labelList)*int(params.NegativeSampleSize)*len(positiveList)))
					}
				}
			}
			if debug != nil {
				debug.Printf("TrainLabelBoost(%s): t=%d: extracted positive/negative %d pair(s)", algoName, t, len(pairIndices))
			}
		}
		// Training the splitter.
		splitter, err := rankerTrainer(ds.X, pairIndices, pairMargins, pairCs, params.Epsilon, nil)
		if err != nil {
			return nil, fmt.Errorf("BinaryRankerTrainer(%s): %s", params.RankerTrainerName, err)
		}
		var Zt []float32
		if params.NegativeSampleSize == 0 {
			var tn, fn, fp, tp uint
			tn, fn, fp, tp, Zt, _ = splitter.ReportPerformance(ds.X, deltas)
			if debug != nil {
				debug.Printf("TrainLabelBoost(%s): t=%d: trained the splitter: tn=%d, fn=%d, fp=%d, tp=%d", algoName, t, tn, fn, fp, tp)
			}
		} else {
			Zt = splitter.PredictAll(ds.X)
		}
		// Update the margin matrix.
		for i, zi := range Z {
			zti := Zt[i]
			updatedLabelSet := make(map[uint32]bool)
			for j, zij := range zi {
				label := zij.Key
				if _, ok := labelSet[label]; ok {
					zi[j].Value += zti
					updatedLabelSet[label] = true
				}
			}
			for label := range labelSet {
				if _, ok := updatedLabelSet[label]; !ok {
					zi = append(zi, sticker.KeyValue32{label, 0.0 + zti})
				}
			}
			Z[i] = zi
		}
		// Append the round information.
		biases = append(biases, splitter.Bias)
		for feature, value := range splitter.Weight {
			weightLists[feature] = append(weightLists[feature], sticker.KeyValue32{uint32(t - 1), value})
		}
		labelLists = append(labelLists, labelList)
		summary := make(map[string]interface{})
		summaries = append(summaries, summary)
	}
	return &LabelBoost{
		Params:      params,
		Biases:      biases,
		WeightLists: weightLists,
		LabelLists:  labelLists,
		Summaries:   summaries,
	}, nil
}

// DecodeLabelBoostWithGobDecoder decodes LabelBoost using decoder.
//
// This function returns an error in decoding.
func DecodeLabelBoostWithGobDecoder(model *LabelBoost, decoder *gob.Decoder) error {
	model.Params = &LabelBoostParameters{}
	if err := decoder.Decode(model.Params); err != nil {
		return fmt.Errorf("DecodeLabelBoost: Params: %s", err)
	}
	if err := decoder.Decode(&model.Biases); err != nil {
		return fmt.Errorf("DecodeLabelBoost: Biases: %s", err)
	}
	var lenWeightLists int
	if err := decoder.Decode(&lenWeightLists); err != nil {
		return fmt.Errorf("DecodeLabelBoost: len(WeightLists): %s", err)
	}
	model.WeightLists = make(map[uint32]sticker.KeyValues32)
	for i := 0; i < lenWeightLists; i++ {
		var feature uint32
		if err := decoder.Decode(&feature); err != nil {
			return fmt.Errorf("DecodeLabelBoost: #%d WeightList: %s", i, err)
		}
		var weightList sticker.KeyValues32
		if err := decoder.Decode(&weightList); err != nil {
			return fmt.Errorf("DecodeLabelBoost: #%d WeightList: %s", i, err)
		}
		model.WeightLists[feature] = weightList
	}
	if err := decoder.Decode(&model.LabelLists); err != nil {
		return fmt.Errorf("DecodeLabelBoost: LabelLists: %s", err)
	}
	if err := decoder.Decode(&model.Summaries); err != nil {
		return fmt.Errorf("DecodeLabelBoost: Summaries: %s", err)
	}
	return nil
}

// DecodeLabelBoost decodes LabelBoost from r.
// Directly passing *os.File used by a gob.Decoder to this function causes mysterious errors.
// Thus, if users use gob.Decoder, then they should call DecodeLabelBoostWithGobDecoder.
//
// This function returns an error in decoding.
func DecodeLabelBoost(model *LabelBoost, r io.Reader) error {
	return DecodeLabelBoostWithGobDecoder(model, gob.NewDecoder(r))
}

// EncodeLabelBoostWithGobEncoder decodes LabelBoost using encoder.
//
// This function returns an error in decoding.
func EncodeLabelBoostWithGobEncoder(model *LabelBoost, encoder *gob.Encoder) error {
	if err := encoder.Encode(model.Params); err != nil {
		return fmt.Errorf("EncodeLabelBoost: Params: %s", err)
	}
	if err := encoder.Encode(model.Biases); err != nil {
		return fmt.Errorf("EncodeLabelBoost: Biases: %s", err)
	}
	features := make([]int, 0, len(model.WeightLists))
	for feature := range model.WeightLists {
		features = append(features, int(feature))
	}
	sort.Ints(features)
	if err := encoder.Encode(len(model.WeightLists)); err != nil {
		return fmt.Errorf("EncodeLabelBoost: len(WeightLists): %s", err)
	}
	for i, feature := range features {
		if err := encoder.Encode(uint32(feature)); err != nil {
			return fmt.Errorf("EncodeLabelBoost: #%d WeightList: %s", i, err)
		}
		if err := encoder.Encode(model.WeightLists[uint32(feature)]); err != nil {
			return fmt.Errorf("EncodeLabelBoost: #%d WeightList: %s", i, err)
		}
	}
	if err := encoder.Encode(model.LabelLists); err != nil {
		return fmt.Errorf("EncodeLabelBoost: LabelLists: %s", err)
	}
	if err := encoder.Encode(model.Summaries); err != nil {
		return fmt.Errorf("EncodeLabelBoost: Summaries: %s", err)
	}
	return nil
}

// EncodeLabelBoost encodes LabelForest to w.
// Directly passing *os.File used by a gob.Encoder to this function causes mysterious errors.
// Thus, if users use gob.Encoder, then they should call EncodeLabelBoostWithGobEncoder.
//
// This function returns an error in encoding.
func EncodeLabelBoost(model *LabelBoost, w io.Writer) error {
	return EncodeLabelBoostWithGobEncoder(model, gob.NewEncoder(w))
}

// GobEncode returns the error always, because users should encode large LabelBoost objects with EncodeLabelBoost.
func (model *LabelBoost) GobEncode() ([]byte, error) {
	return nil, fmt.Errorf("LabelBoost should be encoded with EncodeLabelBoost")
}

// Nrounds return the number of the rounds.
func (model *LabelBoost) Nrounds() uint {
	return uint(len(model.Biases))
}

// Predict returns the top-K predicted labels for the given data point x with the first T rounds.
func (model *LabelBoost) Predict(x sticker.FeatureVector, K uint, T uint) sticker.LabelVector {
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
		for _, label := range model.LabelLists[t] {
			y[label] += zt
		}
	}
	return sticker.RankTopK(y, K)
}

// PredictAll returns the slice of the top-K predicted labels for each data point in X with the first T rounds.
func (model *LabelBoost) PredictAll(X sticker.FeatureVectors, K uint, T uint) sticker.LabelVectors {
	Y := make(sticker.LabelVectors, 0, len(X))
	for _, xi := range X {
		Y = append(Y, model.Predict(xi, K, T))
	}
	return Y
}
