package sticker

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
)

// LabelNearParameters is the parameters for LabelNear.
type LabelNearParameters struct {
	// K is the number of the hash tables.
	K uint
	// L is the bit-width of bucket indices in each hash table.
	L uint
	// R is the size of a reservoir of each bucket.
	R uint
}

// NewLabelNearParameters returns an new default LabelNearParameters.
func NewLabelNearParameters() *LabelNearParameters {
	return &LabelNearParameters{
		K: 64,
		L: 16,
		R: 64,
	}
}

// LabelNear is a faster implementation of LabelNearest which uses the optimal Densified One Permutation Hashing (DOPH) and the reservoir sampling (Wang+ 2017).
//
// References:
//
// (Wang+ 2017) Y. Wang, A. Shrivastava, and J. Ryu. "FLASH: Randomized Algorithms Accelerated over CPU-GPU for Ultra-High Dimensional Similarity Search." arXiv preprint arXiv:1709.01190, 2017.
type LabelNear struct {
	// Dataset is the training dataset.
	Dataset *Dataset
	// Hashing is the Jaccard hashing.
	Hashing *JaccardHashing
}

// TrainLabelNear returns an trained LabelNear on the given training dataset ds.
//
// Currently, this function returns no error.
func TrainLabelNear(ds *Dataset, params *LabelNearParameters, debug *log.Logger) (*LabelNear, error) {
	newds := &Dataset{
		X: make(FeatureVectors, len(ds.X)),
		Y: make(LabelVectors, len(ds.Y)),
	}
	if debug != nil {
		debug.Printf("constructing JaccardHashing ...")
	}
	hashing := NewJaccardHashing(params.K, params.L, params.R)
	for i, xi := range ds.X {
		yi := ds.Y[i]
		lenxi := float32(0.0)
		for _, xipair := range xi {
			lenxi += xipair.Value * xipair.Value
		}
		lenxi = Sqrt32(lenxi)
		newxi := make(FeatureVector, len(xi))
		for j, xipair := range xi {
			newxi[j] = KeyValue32{xipair.Key, xipair.Value / lenxi}
		}
		newyi := make(LabelVector, len(yi))
		copy(newyi, yi)
		newds.X[i], newds.Y[i] = newxi, newyi
		hashing.Add(xi, uint32(i))
	}
	if debug != nil {
		bucketUsage, bucketSizeHist := hashing.Summary()
		debug.Printf("JaccardHashing(K=%d,L=%d,R=%d): bucketUsage=%d", hashing.K(), hashing.L(), hashing.R(), bucketUsage)
		debug.Printf("JaccardHashing(K=%d,L=%d,R=%d): bucketSizeHist=%d", hashing.K(), hashing.L(), hashing.R(), bucketSizeHist)
	}
	return &LabelNear{
		Dataset: newds,
		Hashing: hashing,
	}, nil
}

// DecodeLabelNearWithGobDecoder decodes LabelNear using decoder.
//
// This function returns an error in decoding.
func DecodeLabelNearWithGobDecoder(model *LabelNear, decoder *gob.Decoder) error {
	model.Dataset = &Dataset{}
	if err := DecodeDatasetWithGobDecoder(model.Dataset, decoder); err != nil {
		return fmt.Errorf("DecodeLabelNear: Dataset: %s", err)
	}
	model.Hashing = &JaccardHashing{}
	if err := DecodeJaccardHashingWithGobDecoder(model.Hashing, decoder); err != nil {
		return fmt.Errorf("DecodeLabelNear: Hashing: %s", err)
	}
	return nil
}

// DecodeLabelNear decodes LabelNear from r.
// Directly passing *os.File used by a gob.Decoder to this function causes mysterious errors.
// Thus, if users use gob.Decoder, then they should call DecodeLabelNearWithGobDecoder.
//
// This function returns an error in decoding.
func DecodeLabelNear(model *LabelNear, r io.Reader) error {
	return DecodeLabelNearWithGobDecoder(model, gob.NewDecoder(r))
}

// EncodeLabelNearWithGobEncoder decodes LabelNear using encoder.
//
// This function returns an error in decoding.
func EncodeLabelNearWithGobEncoder(model *LabelNear, encoder *gob.Encoder) error {
	if err := EncodeDatasetWithGobEncoder(model.Dataset, encoder); err != nil {
		return fmt.Errorf("EncodeLabelNear: Dataset: %s", err)
	}
	if err := EncodeJaccardHashingWithGobEncoder(model.Hashing, encoder); err != nil {
		return fmt.Errorf("EncodeLabelNear: Hashing: %s", err)
	}
	return nil
}

// EncodeLabelNear encodes LabelNear to w.
// Directly passing *os.File used by a gob.Encoder to this function causes mysterious errors.
// Thus, if users use gob.Encoder, then they should call EncodeLabelNearWithGobEncoder.
//
// This function returns an error in encoding.
func EncodeLabelNear(model *LabelNear, w io.Writer) error {
	return EncodeLabelNearWithGobEncoder(model, gob.NewEncoder(w))
}

// GobEncode returns the error always, because users should encode large LabelNear objects with EncodeLabelNear.
func (model *LabelNear) GobEncode() ([]byte, error) {
	return nil, fmt.Errorf("LabelNear should be encoded with EncodeLabelNear")
}

// FindNears returns the S near entries with each similarity for the given entry.
// The quantity c*S is used for sieving the candidates by hashing.
// See Predict for hyper-parameter details.
func (model *LabelNear) FindNears(x FeatureVector, c, S uint, beta float32) KeyValues32 {
	indexSimsTopS := make(KeyValues32, 0, S)
	nears := KeyCounts32(model.Hashing.FindNears(x))
	nears = nears.SortLargestCountsWithHeap(c * S)
	for _, nearPair := range nears {
		protoidx, count := nearPair.Key, nearPair.Count
		if count == 0 {
			break
		}
		xp := model.Dataset.X[protoidx]
		if sim, count := DotCount(x, xp); sim > 0.0 {
			// For efficiency, call Pow32 as short-cut style.
			jaccard := float32(count) / float32(len(x)+len(xp)-count)
			if beta == 0 {
			} else if beta == 1 {
				sim *= jaccard
			} else {
				sim *= Pow32(jaccard, beta)
			}
			if len(indexSimsTopS) == 0 {
				indexSimsTopS = append(indexSimsTopS, KeyValue32{uint32(protoidx), sim})
			} else if indexSimsTopS[len(indexSimsTopS)-1].Value > sim {
				if len(indexSimsTopS) < cap(indexSimsTopS) {
					indexSimsTopS = append(indexSimsTopS, KeyValue32{uint32(protoidx), sim})
				}
			} else {
				for rank := 0; rank < len(indexSimsTopS); rank++ {
					if sim >= indexSimsTopS[rank].Value {
						if len(indexSimsTopS) < cap(indexSimsTopS) {
							indexSimsTopS = append(indexSimsTopS, KeyValue32{0, 0})
						}
						copy(indexSimsTopS[rank+1:], indexSimsTopS[rank:])
						indexSimsTopS[rank] = KeyValue32{uint32(protoidx), sim}
						break
					}
				}
			}
		}
	}
	return indexSimsTopS
}

// Predict returns the results for the given data entry x with the sparse S-near neighborhood.
// The returned results are the top-K labels, the label histogram, and the slice of the data entry index and its similarity.
//
// alpha is the smoothing parameter for weighting the votes by each neighbor.
// beta is the smoothing parameter for balancing the Jaccard similarity and the cosine similarity.
func (model *LabelNear) Predict(x FeatureVector, K, c, S uint, alpha, beta float32) (LabelVector, map[uint32]float32, KeyValues32) {
	indexSimsTopS := model.FindNears(x, c, S, beta)
	labelHist := make(map[uint32]float32)
	xlen := float32(0.0)
	for _, xpair := range x {
		xlen += xpair.Value * xpair.Value
	}
	xlen = Sqrt32(xlen)
	for _, indexSim := range indexSimsTopS {
		value := Pow32(indexSim.Value/xlen, alpha)
		for _, label := range model.Dataset.Y[indexSim.Key] {
			labelHist[label] += value
		}
	}
	return RankTopK(labelHist, K), labelHist, indexSimsTopS
}

// PredictAll returns the top-K labels for each data entry in X with the sparse S-near neighborhood.
// See Predict for hyper-parameter details.
func (model *LabelNear) PredictAll(X FeatureVectors, K, c, S uint, alpha, beta float32) LabelVectors {
	Yhat := make(LabelVectors, 0, len(X))
	for _, xi := range X {
		yihat, _, _ := model.Predict(xi, K, c, S, alpha, beta)
		Yhat = append(Yhat, yihat)
	}
	return Yhat
}
