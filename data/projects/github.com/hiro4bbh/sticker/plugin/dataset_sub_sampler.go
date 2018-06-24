package plugin

import (
	"math/rand"
	"sort"

	"github.com/hiro4bbh/sticker"
)

// DatasetEntrySubSampler is the interface for entry sub-sampler on the dataset.
type DatasetEntrySubSampler interface {
	// SubSample returns the index slice contained in the sub-sample with the given seed.
	SubSample(ds *sticker.Dataset, seed uint) []int
}

// DeterministicDatasetEntrySubSampler is a deterministic DatasetEntrySubSampler.
// The sub-sampler simply returns the sub-dataset with the given size in order of the given dataset.
// The seed is used as the sub-sample start index.
//
// This implements interface DatasetEntrySubSampler.
type DeterministicDatasetEntrySubSampler struct {
	n uint
}

// NewDeterministicDatasetEntrySubSampler returns an new DeterministicDatasetEntrySubSampler.
func NewDeterministicDatasetEntrySubSampler(n uint) DatasetEntrySubSampler {
	return &DeterministicDatasetEntrySubSampler{
		n: n,
	}
}

// SubSample is for interface DatasetEntrySubSampler.
func (sampler *DeterministicDatasetEntrySubSampler) SubSample(ds *sticker.Dataset, seed uint) []int {
	maxSeed := (uint(ds.Size()) + sampler.n - 1) / sampler.n
	start, end := (seed%maxSeed)*sampler.n, (seed%maxSeed+1)*sampler.n
	if end > uint(ds.Size()) {
		end = uint(ds.Size())
	}
	indices := make([]int, 0, end-start)
	for i := int(start); i < int(end); i++ {
		indices = append(indices, i)
	}
	return indices
}

// RandomDatasetEntrySubSampler is a random DatasetEntrySubSampler.
// This sub-sampler returns the sub-dataset with the given size with replacement from the given dataset.
// seed is used as the seed of the random number generator.
//
// This implements interface DatasetEntrySubSampler.
type RandomDatasetEntrySubSampler struct {
	n uint
}

// NewRandomDatasetEntrySubSampler returns an new RandomDatasetEntrySubSampler.
func NewRandomDatasetEntrySubSampler(n uint) DatasetEntrySubSampler {
	return &RandomDatasetEntrySubSampler{
		n: n,
	}
}

// SubSample is for interface DatasetEntrySubSampler.
func (sampler *RandomDatasetEntrySubSampler) SubSample(ds *sticker.Dataset, seed uint) []int {
	rng := rand.New(rand.NewSource(int64(seed)))
	indices := make([]int, sampler.n)
	for i := 0; i < int(sampler.n); i++ {
		indices[i] = rng.Intn(ds.Size())
	}
	return indices
}

// DatasetNoneFeatureSubSampler does not any feature sub-sampling, and returns the given dataset itself.
//
// This function returns no error currently.
func DatasetNoneFeatureSubSampler(ds *sticker.Dataset, seed int64) (*sticker.Dataset, error) {
	return ds, nil
}

// DatasetSqrtFeatureSubSampler sub-samples sqrt(J) features (J is the number of the used features), and returns the feature sub-sampled dataset.
// This is registered to sticker.DatasetFeatureSubSamplers.
//
// This function returns no error currently.
func DatasetSqrtFeatureSubSampler(ds *sticker.Dataset, seed int64) (*sticker.Dataset, error) {
	// Select sqrt(J) features
	featureSet := make(map[uint32]struct{})
	for _, xi := range ds.X {
		for _, xipair := range xi {
			featureSet[xipair.Key] = struct{}{}
		}
	}
	maxJ := int(sticker.Sqrt32(float32(len(featureSet))))
	features := make([]int, 0, len(featureSet))
	for feature := range featureSet {
		features = append(features, int(feature))
	}
	// Force determinism for reproducibility (the iteration order of golang's map is changed every time!).
	sort.Ints(features)
	rng := rand.New(rand.NewSource(seed))
	for j := 0; j < maxJ; j++ {
		j_ := j + rng.Intn(len(features)-j)
		features[j], features[j_] = features[j_], features[j]
	}
	features = features[:maxJ]
	featureSet = make(map[uint32]struct{})
	for _, feature := range features {
		featureSet[uint32(feature)] = struct{}{}
	}
	return ds.FeatureSubSet(featureSet), nil
}

// DatasetFeatureSubSampler is the type of feature sub-samplers.
// A sub-sampler returns the dataset whose features are sub-sampled.
type DatasetFeatureSubSampler func(ds *sticker.Dataset, seed int64) (*sticker.Dataset, error)

// DatasetFeatureSubSamplers is the map from the sub-sampler name to the corresponding sub-sampler.
var DatasetFeatureSubSamplers = map[string]DatasetFeatureSubSampler{
	"none": DatasetNoneFeatureSubSampler,
	"sqrt": DatasetSqrtFeatureSubSampler,
}

// DefaultDatasetFeatureSubSamplerName is the default DatasetFeatureSubSampler name.
const DefaultDatasetFeatureSubSamplerName = "none"
