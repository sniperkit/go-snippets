package plugin

import (
	"math/rand"
	"testing"

	"github.com/hiro4bbh/go-assert"
	"github.com/hiro4bbh/sticker"
)

func TestDeterministicDatasetEntrySubSampler(t *testing.T) {
	ds := &sticker.Dataset{
		X: sticker.FeatureVectors{
			sticker.FeatureVector{sticker.KeyValue32{0, 1.0}}, sticker.FeatureVector{sticker.KeyValue32{1, 1.0}},
			sticker.FeatureVector{sticker.KeyValue32{2, 1.0}}, sticker.FeatureVector{sticker.KeyValue32{3, 1.0}},
		},
		Y: sticker.LabelVectors{
			sticker.LabelVector{0}, sticker.LabelVector{1}, sticker.LabelVector{2}, sticker.LabelVector{3},
		},
	}
	subSampler := NewDeterministicDatasetEntrySubSampler(3)
	goassert.New(t, []int{0, 1, 2}).Equal(subSampler.SubSample(ds, 0))
	goassert.New(t, []int{3}).Equal(subSampler.SubSample(ds, 1))
	goassert.New(t, []int{0, 1, 2}).Equal(subSampler.SubSample(ds, 2))
	goassert.New(t, []int{0, 1, 2}).Equal(subSampler.SubSample(ds, 0))
}

func TestRandomDatasetEntrySubSampler(t *testing.T) {
	ds := &sticker.Dataset{
		X: sticker.FeatureVectors{
			sticker.FeatureVector{sticker.KeyValue32{0, 1.0}}, sticker.FeatureVector{sticker.KeyValue32{1, 1.0}},
			sticker.FeatureVector{sticker.KeyValue32{2, 1.0}}, sticker.FeatureVector{sticker.KeyValue32{3, 1.0}},
		},
		Y: sticker.LabelVectors{
			sticker.LabelVector{0}, sticker.LabelVector{1}, sticker.LabelVector{2}, sticker.LabelVector{3},
		},
	}
	subSampler := NewRandomDatasetEntrySubSampler(3)
	var rng *rand.Rand
	rng = rand.New(rand.NewSource(0))
	goassert.New(t, []int{rng.Intn(ds.Size()), rng.Intn(ds.Size()), rng.Intn(ds.Size())}).Equal(subSampler.SubSample(ds, 0))
	rng = rand.New(rand.NewSource(1))
	goassert.New(t, []int{rng.Intn(ds.Size()), rng.Intn(ds.Size()), rng.Intn(ds.Size())}).Equal(subSampler.SubSample(ds, 1))
	rng = rand.New(rand.NewSource(2))
	goassert.New(t, []int{rng.Intn(ds.Size()), rng.Intn(ds.Size()), rng.Intn(ds.Size())}).Equal(subSampler.SubSample(ds, 2))
	rng = rand.New(rand.NewSource(0))
	goassert.New(t, []int{rng.Intn(ds.Size()), rng.Intn(ds.Size()), rng.Intn(ds.Size())}).Equal(subSampler.SubSample(ds, 0))
}

func TestDatasetNoneFeatureSubSampler(t *testing.T) {
	ds := &sticker.Dataset{
		X: sticker.FeatureVectors{
			sticker.FeatureVector{sticker.KeyValue32{0, 1.0}},
			sticker.FeatureVector{sticker.KeyValue32{0, 2.0}, sticker.KeyValue32{1, 3.0}},
			sticker.FeatureVector{sticker.KeyValue32{0, 4.0}, sticker.KeyValue32{2, 5.0}, sticker.KeyValue32{9, 6.0}},
			sticker.FeatureVector{sticker.KeyValue32{10, 7.0}},
		},
		Y: sticker.LabelVectors{
			sticker.LabelVector{0}, sticker.LabelVector{0, 1}, sticker.LabelVector{0, 2, 9}, sticker.LabelVector{10},
		},
	}
	goassert.New(t, ds).EqualWithoutError(DatasetFeatureSubSamplers["none"](ds, 0))
	goassert.New(t, ds).EqualWithoutError(DatasetFeatureSubSamplers["none"](ds, 1))
}

func TestDatasetSqrtFeatureSubSampler(t *testing.T) {
	ds := &sticker.Dataset{
		X: sticker.FeatureVectors{
			sticker.FeatureVector{sticker.KeyValue32{0, 1.0}},
			sticker.FeatureVector{sticker.KeyValue32{0, 2.0}, sticker.KeyValue32{1, 3.0}},
			sticker.FeatureVector{sticker.KeyValue32{0, 4.0}, sticker.KeyValue32{2, 5.0}, sticker.KeyValue32{9, 6.0}},
			sticker.FeatureVector{sticker.KeyValue32{10, 7.0}},
		},
		Y: sticker.LabelVectors{
			sticker.LabelVector{0}, sticker.LabelVector{0, 1}, sticker.LabelVector{0, 2, 9}, sticker.LabelVector{10},
		},
	}
	goassert.New(t, &sticker.Dataset{
		X: sticker.FeatureVectors{
			sticker.FeatureVector{},
			sticker.FeatureVector{},
			sticker.FeatureVector{sticker.KeyValue32{9, 6.0}},
			sticker.FeatureVector{sticker.KeyValue32{10, 7.0}},
		},
		Y: sticker.LabelVectors{
			sticker.LabelVector{0}, sticker.LabelVector{0, 1}, sticker.LabelVector{0, 2, 9}, sticker.LabelVector{10},
		},
	}).EqualWithoutError(DatasetFeatureSubSamplers["sqrt"](ds, 0))
	goassert.New(t, &sticker.Dataset{
		X: sticker.FeatureVectors{
			sticker.FeatureVector{},
			sticker.FeatureVector{sticker.KeyValue32{1, 3.0}},
			sticker.FeatureVector{},
			sticker.FeatureVector{sticker.KeyValue32{10, 7.0}},
		},
		Y: sticker.LabelVectors{
			sticker.LabelVector{0}, sticker.LabelVector{0, 1}, sticker.LabelVector{0, 2, 9}, sticker.LabelVector{10},
		},
	}).EqualWithoutError(DatasetFeatureSubSamplers["sqrt"](ds, 1))
}
