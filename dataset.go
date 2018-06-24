package sticker

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"math/bits"
	"sort"
	"strconv"
	"strings"
)

// Float32Slice implements the interface sort.Interface.
type Float32Slice []float32

func (x Float32Slice) Len() int {
	return len(x)
}

func (x Float32Slice) Less(i, j int) bool {
	return x[i] < x[j]
}

func (x Float32Slice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

// SummarizeFloat32Slice returns the 5-summary(minimum, 1st quantile, median, 3rd quantile and maximum) and the average of the given float32 slice.
func SummarizeFloat32Slice(x []float32) (min, q25, med, q75, max, avg float32) {
	if len(x) == 0 {
		return NaN32(), NaN32(), NaN32(), NaN32(), NaN32(), NaN32()
	} else if len(x) == 1 {
		return x[0], x[0], x[0], x[0], x[0], x[0]
	}
	xs := make([]float32, len(x))
	copy(xs, x)
	sort.Sort(Float32Slice(xs))
	min, max = xs[0], xs[len(xs)-1]
	if len(xs)%2 == 1 {
		med = xs[len(xs)/2]
	} else {
		med = 0.5*xs[len(xs)/2-1] + 0.5*xs[len(xs)/2]
	}
	r25 := (1.0-0.25)*0.0 + 0.25*float32(len(xs)-1)
	if _, r25f := Modf32(r25); r25f == 0.0 {
		q25 = xs[int(Floor32(r25))]
	} else {
		q25 = (Ceil32(r25)-r25)*xs[int(Floor32(r25))] + (r25-Floor32(r25))*xs[int(Ceil32(r25))]
	}
	r75 := (1.0-0.75)*0.0 + 0.75*float32(len(xs)-1)
	if _, r75f := Modf32(r75); r75f == 0.0 {
		q75 = xs[int(Floor32(r75))]
	} else {
		q75 = (Ceil32(r75)-r75)*xs[int(Floor32(r75))] + (r75-Floor32(r75))*xs[int(Ceil32(r75))]
	}
	for _, x := range xs {
		avg += x
	}
	avg /= float32(len(xs))
	return
}

// SparseVector is the map from uint32 key to float32 value.
type SparseVector map[uint32]float32

// SparseVectors is the slice of SparseVector.
type SparseVectors []SparseVector

// AvgTotalVariationAmongSparseVectors returns the average total-variation distance among the sparse vectors.
// This function returns 0.0 if there is at most one sparse vector, and returns NaN if some sparse vectors are empty.
func AvgTotalVariationAmongSparseVectors(svs SparseVectors) float32 {
	n := len(svs)
	if n <= 1 {
		return 0.0
	}
	keySet := make(map[uint32]bool)
	Zs := make([]float32, n)
	for i, svi := range svs {
		Z := float32(0.0)
		for key, value := range svi {
			keySet[key] = true
			Z += value
		}
		if Z == 0 {
			return NaN32()
		}
		Zs[i] = Z
	}
	keys := make([]int, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)
	sum := float32(0.0)
	for i := 0; i < n; i++ {
		svi, Zi := svs[i], Zs[i]
		for j := i + 1; j < n; j++ {
			svj, Zj := svs[j], Zs[j]
			for _, key := range keys {
				sum += Abs32(svi[uint32(key)]/Zi - svj[uint32(key)]/Zj)
			}
		}
	}
	return 0.5 * sum / (float32(n) * float32(n-1) / 2)
}

// SparsifyVector returns a SparseVector converted from v.
func SparsifyVector(v []float32) SparseVector {
	sv := make(map[uint32]float32)
	for i, vi := range v {
		if vi != 0.0 {
			sv[uint32(i)] = vi
		}
	}
	return sv
}

// HashUint32 returns the hashed value of the given uint32 x.
// This comes from MurmurHash3 fmix32 (https://github.com/aappleby/smhasher/blob/master/src/MurmurHash3.cpp).
func HashUint32(x uint32) uint32 {
	h := uint64(x + 1)
	h ^= h >> 16
	h = (h * 0x85ebca6b) & 0xffffffff
	h ^= h >> 13
	h = (h * 0xc2b2ae35) & 0xffffffff
	return uint32(h ^ (h >> 16))
}

// KeyCount32 is the pair of uint32 feature key and its uint32 value.
type KeyCount32 struct {
	Key, Count uint32
}

// KeyCounts32 is the slice of KeyCount32.
type KeyCounts32 []KeyCount32

// ExtractLargestCountsByInsert returns the only K largest entries.
func (kcs KeyCounts32) ExtractLargestCountsByInsert(K uint) KeyCounts32 {
	if K > uint(len(kcs)) {
		K = uint(len(kcs))
	}
	kcs2 := make(KeyCounts32, 0, K)
	for _, keyCount := range kcs {
		if count := keyCount.Count; count > 0 {
			if len(kcs2) == 0 {
				kcs2 = append(kcs2, keyCount)
			} else if kcs2[len(kcs2)-1].Count > count {
				if len(kcs2) < cap(kcs2) {
					kcs2 = append(kcs2, keyCount)
				}
			} else {
				for rank := 0; rank < len(kcs2); rank++ {
					if count > kcs2[rank].Count {
						if len(kcs2) < cap(kcs2) {
							kcs2 = append(kcs2, KeyCount32{0, 0})
						}
						copy(kcs2[rank+1:], kcs2[rank:])
						kcs2[rank] = keyCount
						break
					}
				}
			}
		}
	}
	return kcs2
}

// SortLargestCountsWithHeap sorts the only K largest entries at the first as maintaining the heap, and returns the shrunk slice to the self.
func (kcs KeyCounts32) SortLargestCountsWithHeap(K uint) KeyCounts32 {
	// This implementation comes from http://web.archive.org/web/20140807181610/http://fallabs.com/blog-ja/promenade.cgi?id=104.
	// Strategy:
	//   Basically, this is an in-place heap sort.
	//   Retain the heap at the first whose size is at most K for efficiently sorting the K largest counts.
	if K > uint(len(kcs)) {
		K = uint(len(kcs))
	}
	// First, the only first entry is automatically in the heap.
	cur := 1
	// The first K entries are inserted into the heap.
	// Any entry in the heap satisfies that the entry is "NOT LARGER" than any descendants of it.
	// The order in the heap are reversed at the third step.
	// The heap structure is retained as follows:
	//   [0] -> [1] -> [3] -> ...
	//              -> [4] -> ...
	//       -> [2] -> [5] -> ...
	//              -> [6] -> ...
	// Hence, prev(k) = (k - 1)/2, left(k) = 2*k + 1, and right(k) = 2*k + 2.
	for ; cur < int(K); cur++ {
		// Insert the cur-th entry into the heap.
		cidx := cur
		for cidx > 0 {
			pidx := (cidx - 1) / 2 // prev(cidx)
			if !(kcs[pidx].Count > kcs[cidx].Count) {
				break
			}
			// Swap the current entry with the parent one, because the current one is smaller.
			kcs[cidx], kcs[pidx] = kcs[pidx], kcs[cidx]
			// Perform this recursively at the parent entry.
			cidx = pidx
		}
	}
	// Second, the remain entries are inserted into the heap as keeping the heap size is k.
	for cur < len(kcs) {
		// Insert the current entry if it is larger than the smallest one in the heap.
		if kcs[cur].Count > kcs[0].Count {
			// Procedure A: Insert the current entry into the size K heap.
			// Swap the current entry and the smallest one in the heap.
			kcs[0], kcs[cur] = kcs[cur], kcs[0]
			// Insert the current entry in the heap.
			pidx, bot := 0, int(K)/2
			for pidx < bot {
				// Take the smaller child as the current entry.
				cidx := 2*pidx + 1 // left(cidx)
				if cidx < int(K)-1 && kcs[cidx].Count > kcs[cidx+1].Count {
					cidx++
				}
				if kcs[cidx].Count > kcs[pidx].Count {
					break
				}
				// Swap the current entry with the selected child, because the current one is larger.
				kcs[pidx], kcs[cidx] = kcs[cidx], kcs[pidx]
				// Perform this recursively at the child entry.
				pidx = cidx
			}
		}
		cur++
	}
	// Third, the entries in the heap are reversed.
	// This is achieved by shrinking the heap one by one:
	//   Taking the largest entry as the current one, insert it to the heap.
	// Hence, the largest entry is being selected as the current one as pushing down the smaller entries.
	cur = int(K) - 1
	for cur > 0 {
		// Apply procedure A in the size cur heap.
		kcs[0], kcs[cur] = kcs[cur], kcs[0]
		pidx, bot := 0, cur/2
		for pidx < bot {
			cidx := 2*pidx + 1
			if cidx < cur-1 && kcs[cidx].Count > kcs[cidx+1].Count {
				cidx++
			}
			if kcs[cidx].Count > kcs[pidx].Count {
				break
			}
			kcs[pidx], kcs[cidx] = kcs[cidx], kcs[pidx]
			pidx = cidx
		}
		cur--
	}
	return kcs[:K]
}

// KeyCountMap32 is the faster map of KeyCount32s.
// This cannot has entries with value 0.
//
// Currently, this does not support expansion at insertion.
// Users can iterate the entries with raw access to the internal, so this is for expert use in order to achieve faster counting.
type KeyCountMap32 KeyCounts32

// NewKeyCountMap32 returns a new KeyCountMap32.
func NewKeyCountMap32(capacity uint) KeyCountMap32 {
	return make(KeyCountMap32, 1<<uint(bits.Len(capacity)))
}

// Get returns the entry with the given key.
func (m KeyCountMap32) Get(key uint32) KeyCount32 {
	for k := HashUint32(key) & uint32(len(m)-1); m[k].Count > 0; k = (k + 1) & uint32(len(m)-1) {
		if m[k].Key == key {
			return KeyCount32{m[k].Key, m[k].Count}
		}
	}
	return KeyCount32{0, 0}
}

// Inc increments the entry's value with the given key, and returns the entry.
func (m KeyCountMap32) Inc(key uint32) KeyCount32 {
	k := HashUint32(key) & uint32(len(m)-1)
	for ; m[k].Count > 0; k = (k + 1) & uint32(len(m)-1) {
		if m[k].Key == key {
			m[k].Count++
			return KeyCount32{m[k].Key, m[k].Count}
		}
	}
	m[k] = KeyCount32{key, 1}
	return KeyCount32{m[k].Key, m[k].Count}
}

// Map returns the map version of self.
func (m KeyCountMap32) Map() map[uint32]uint32 {
	ma := make(map[uint32]uint32)
	for _, mpair := range m {
		if mpair.Count > 0 {
			ma[mpair.Key] = mpair.Count
		}
	}
	return ma
}

// KeyValue32 is the pair of uint32 feature key and its float32 value.
type KeyValue32 struct {
	Key   uint32
	Value float32
}

// KeyValues32 is the KeyValue32 slice.
type KeyValues32 []KeyValue32

// KeyValues32OrderedByKey is KeyValues32 implementing sort.Interface for sorting in the key order.
// If the keys are same, then these key-values are sorted in increasing value order.
type KeyValues32OrderedByKey KeyValues32

func (kvs KeyValues32OrderedByKey) Len() int {
	return len(kvs)
}

func (kvs KeyValues32OrderedByKey) Less(i, j int) bool {
	return kvs[i].Key < kvs[j].Key || (kvs[i].Key == kvs[j].Key && kvs[i].Value < kvs[j].Value)
}

func (kvs KeyValues32OrderedByKey) Swap(i, j int) {
	kvs[i], kvs[j] = kvs[j], kvs[i]
}

// KeyValues32OrderedByValue is KeyValues32 implementing sort.Interface for sorting in the value order.
// If the values are same, then these key-values are sorted in decreasing key order, because this is intended for sorting in decreasing key/value order in reverse mode.
type KeyValues32OrderedByValue KeyValues32

func (kvs KeyValues32OrderedByValue) Len() int {
	return len(kvs)
}

func (kvs KeyValues32OrderedByValue) Less(i, j int) bool {
	return kvs[i].Value < kvs[j].Value || (kvs[i].Value == kvs[j].Value && kvs[i].Key >= kvs[j].Key)
}

func (kvs KeyValues32OrderedByValue) Swap(i, j int) {
	kvs[i], kvs[j] = kvs[j], kvs[i]
}

// FeatureVector is the sparse and static feature vector.
// The elements should be ordered by feature ID (key).
type FeatureVector = KeyValues32OrderedByKey

// DotCount returns the inner product and the size of the intersect of the supports between x and y.
func DotCount(x, y FeatureVector) (float32, int) {
	xj, yj, d, count := 0, 0, float32(0.0), 0
	for xj < len(x) {
		for yj < len(y) && x[xj].Key > y[yj].Key {
			yj++
		}
		if yj >= len(y) {
			break
		}
		if x[xj].Key == y[yj].Key {
			d += x[xj].Value * y[yj].Value
			yj++
			count++
		}
		xj++
	}
	return d, count
}

// FeatureVectors is the FeatureVector slice.
type FeatureVectors []FeatureVector

// Dim returns the calculated dimension of FeatureVectors.
// This is the maximum feature key plus 1.
func (X FeatureVectors) Dim() (d int) {
	for _, xi := range X {
		for _, xipair := range xi {
			if d <= int(xipair.Key) {
				d = int(xipair.Key) + 1
			}
		}
	}
	return
}

// LabelVector is the sparse label vector which is the slice of label key.
type LabelVector []uint32

func (labels LabelVector) Len() int {
	return len(labels)
}

func (labels LabelVector) Less(i, j int) bool {
	return labels[i] < labels[j]
}

func (labels LabelVector) Swap(i, j int) {
	labels[i], labels[j] = labels[j], labels[i]
}

// LabelVectors is the LabelVector slice.
type LabelVectors []LabelVector

// Dim returns the calculated dimension of label vectors.
// This is the maximum label ID plus 1.
func (Y LabelVectors) Dim() (d int) {
	for _, yi := range Y {
		for _, label := range yi {
			if d <= int(label) {
				d = int(label) + 1
			}
		}
	}
	return
}

// Dataset is a collection of the pair of one feature vector and one label vector.
type Dataset struct {
	X FeatureVectors
	Y LabelVectors
}

// ReadTextDataset returns a new Dataset from reader.
//
// The data from reader is assumed to be formatted like used in LIBLINEAR and LIBSVM.
// It is a plain text whose cells are separated by single space.
// The first line has the number of entries, features, and labels.
// The remaining lines have entries of the dataset.
// Each entry is encoded in one line like (comma-separated label) (feature:value)*.
//
// This function returns an error in reading the dataset.
func ReadTextDataset(reader io.Reader) (*Dataset, error) {
	br := bufio.NewReader(reader)
	line, err := br.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("cannot read first line")
	}
	line = strings.TrimSpace(line)
	nentriesNfeaturesNlabels := strings.Split(line, " ")
	if len(nentriesNfeaturesNlabels) != 3 {
		return nil, fmt.Errorf("illegal first line")
	}
	nentries, err := strconv.ParseUint(nentriesNfeaturesNlabels[0], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("illegal nentries in first line")
	}
	nfeatures, err := strconv.ParseUint(nentriesNfeaturesNlabels[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("illegal nfeatures in first line")
	}
	nlabels, err := strconv.ParseUint(nentriesNfeaturesNlabels[2], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("illegal nlabels in first line")
	}
	X := make(FeatureVectors, nentries)
	Y := make(LabelVectors, nentries)
	for i := uint64(0); i < nentries; i++ {
		line, err := br.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("L%d: cannot read line", i+1)
		}
		line = strings.TrimSpace(line)
		entry := strings.Split(line, " ")
		labelStrs := strings.Split(entry[0], ",")
		entry = entry[1:]
		y := make(LabelVector, 0, len(labelStrs))
		for j, labelStr := range labelStrs {
			label, err := strconv.ParseUint(labelStr, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("L%d: illegal #%d label ID", i+1, j+1)
			}
			if label >= nlabels {
				return nil, fmt.Errorf("L%d: too large #%d label ID (>= %d)", i+1, j+1, nlabels)
			}
			y = append(y, uint32(label))
		}
		x := make(FeatureVector, 0, len(entry))
		for j, featureValueStr := range entry {
			featureValue := strings.Split(featureValueStr, ":")
			if len(featureValue) != 2 {
				return nil, fmt.Errorf("L%d: illegal #%d featureID:value pair", i+1, j+1)
			}
			feature, err := strconv.ParseUint(featureValue[0], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("L%d: illegal #%d feature ID", i+1, j+1)
			}
			if feature >= nfeatures {
				return nil, fmt.Errorf("L%d: too large #%d feature ID (>= %d)", i+1, j+1, nfeatures)
			}
			value, err := strconv.ParseFloat(featureValue[1], 32)
			if err != nil {
				return nil, fmt.Errorf("L%d: illegal #%d feature value", i+1, j+1)
			}
			x = append(x, KeyValue32{uint32(feature), float32(value)})
		}
		sort.Sort(x)
		X[i], Y[i] = x, y
	}
	return &Dataset{
		X: X,
		Y: Y,
	}, nil
}

// FeatureSubSet returns the sub-set of the dataset whose entry has only features in the given set of features.
// For efficiency, the label vectors of sub-dataset has references to the one of the dataset.
func (ds *Dataset) FeatureSubSet(features map[uint32]struct{}) *Dataset {
	subds := &Dataset{
		X: FeatureVectors{},
		Y: LabelVectors{},
	}
	for i, xi := range ds.X {
		xinew := FeatureVector{}
		for _, xipair := range xi {
			if _, ok := features[xipair.Key]; ok {
				xinew = append(xinew, xipair)
			}
		}
		subds.X, subds.Y = append(subds.X, xinew), append(subds.Y, ds.Y[i])
	}
	return subds
}

// DecodeDatasetWithGobDecoder decodes Dataset using decoder.
//
// This function returns an error in decoding.
func DecodeDatasetWithGobDecoder(ds *Dataset, decoder *gob.Decoder) error {
	var size int
	if err := decoder.Decode(&size); err != nil {
		return fmt.Errorf("DecodeDataset: size: %s", err)
	}
	ds.X = make(FeatureVectors, size)
	for i := range ds.X {
		if err := decoder.Decode(&ds.X[i]); err != nil {
			return fmt.Errorf("DecodeDataset: #%d feature vector: %s", i, err)
		}
	}
	ds.Y = make(LabelVectors, size)
	for i := range ds.Y {
		if err := decoder.Decode(&ds.Y[i]); err != nil {
			return fmt.Errorf("DecodeDataset: #%d label vector: %s", i, err)
		}
	}
	return nil
}

// DecodeDataset decodes Dataset from r.
// Directly passing *os.File used by a gob.Decoder to this function causes mysterious errors.
// Thus, if users use gob.Decoder, then they should call DecodeDatasetWithGobDecoder.
//
// This function returns an error in decoding.
func DecodeDataset(ds *Dataset, r io.Reader) error {
	return DecodeDatasetWithGobDecoder(ds, gob.NewDecoder(r))
}

// EncodeDatasetWithGobEncoder decodes Dataset using encoder.
//
// This function returns an error in decoding.
func EncodeDatasetWithGobEncoder(ds *Dataset, encoder *gob.Encoder) error {
	if err := encoder.Encode(ds.Size()); err != nil {
		return fmt.Errorf("EncodeDataset: size: %s", err)
	}
	for i, xi := range ds.X {
		if err := encoder.Encode(xi); err != nil {
			return fmt.Errorf("EncodeDataset: #%d feature vector: %s", i, err)
		}
	}
	for i, yi := range ds.Y {
		if err := encoder.Encode(yi); err != nil {
			return fmt.Errorf("EncodeDataset: #%d label vector: %s", i, err)
		}
	}
	return nil
}

// EncodeDataset encodes Dataset to w.
// Directly passing *os.File used by a gob.Encoder to this function causes mysterious errors.
// Thus, if users use gob.Encoder, then they should call EncodeDatasetWithGobEncoder.
//
// This function returns an error in encoding.
func EncodeDataset(ds *Dataset, w io.Writer) error {
	return EncodeDatasetWithGobEncoder(ds, gob.NewEncoder(w))
}

// GobEncode returns the error always, because users should encode large Dataset objects with EncodeDataset.
func (ds *Dataset) GobEncode() ([]byte, error) {
	return nil, fmt.Errorf("Dataset should be encoded with EncodeDataset")
}

// Size returns the size of the dataset.
func (ds *Dataset) Size() int {
	return len(ds.X)
}

// SubSet returns the sub-set of the dataset selected by the given index slice.
// For efficiency, the all vectors of sub-dataset has references to the ones of the dataset.
//
// This function panics if the index is out of range.
func (ds *Dataset) SubSet(indices []int) *Dataset {
	subds := &Dataset{
		X: make(FeatureVectors, len(indices)),
		Y: make(LabelVectors, len(indices)),
	}
	for ii, i := range indices {
		subds.X[ii], subds.Y[ii] = ds.X[i], ds.Y[i]
	}
	return subds
}

// WriteTextDataset writes the dataset to the given writer as ReadTextDataset supports.
//
// This function returns an error in writing.
func (ds *Dataset) WriteTextDataset(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%d %d %d\n", ds.Size(), ds.X.Dim(), ds.Y.Dim()); err != nil {
		return err
	}
	for i, xi := range ds.X {
		for j, label := range ds.Y[i] {
			if j > 0 {
				if _, err := w.Write([]byte(",")); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintf(w, "%d", label); err != nil {
				return err
			}
		}
		for _, xipair := range xi {
			if _, err := fmt.Fprintf(w, " %d:%f", xipair.Key, xipair.Value); err != nil {
				return err
			}
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}
	return nil
}
