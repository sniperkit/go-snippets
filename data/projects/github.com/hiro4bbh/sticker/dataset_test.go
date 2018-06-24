package sticker

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"github.com/hiro4bbh/go-assert"
)

func TestFloat32Slice(t *testing.T) {
	x := []float32{5.0, 3.0, 1.0, 4.0, 2.0}
	sort.Sort(Float32Slice(x))
	goassert.New(t, []float32{1.0, 2.0, 3.0, 4.0, 5.0}).Equal(x)
}

func TestSummarizeFloat32Slice(t *testing.T) {
	min0, q25_0, med0, q75_0, max0, mean0 := SummarizeFloat32Slice([]float32{})
	goassert.New(t, "NaN,NaN,NaN,NaN,NaN,NaN").Equal(fmt.Sprintf("%f,%f,%f,%f,%f,%f", min0, q25_0, med0, q75_0, max0, mean0))
	goassert.New(t, float32(1.0), float32(1.0), float32(1.0), float32(1.0), float32(1.0), float32(1.0)).Equal(SummarizeFloat32Slice([]float32{1.0}))
	goassert.New(t, float32(1.0), float32(1.25), float32(1.5), float32(1.75), float32(2.0), float32(1.5)).Equal(SummarizeFloat32Slice([]float32{1.0, 2.0}))
	goassert.New(t, float32(1.0), float32(2.0), float32(3.0), float32(4.0), float32(5.0), float32(3.0)).Equal(SummarizeFloat32Slice([]float32{5.0, 3.0, 1.0, 4.0, 2.0}))
}

func TestAvgTotalVariationAmongSparseVectors(t *testing.T) {
	goassert.New(t, float32(0.0)).Equal(AvgTotalVariationAmongSparseVectors(SparseVectors{}))
	goassert.New(t, "NaN").Equal(fmt.Sprintf("%g", AvgTotalVariationAmongSparseVectors(SparseVectors{
		SparseVector{0: 1.0, 1: 1.0, 2: 1.0},
		SparseVector{},
	})))
	goassert.New(t, float32(0.0)).Equal(AvgTotalVariationAmongSparseVectors(SparseVectors{
		SparseVector{0: 1.0, 1: 1.0, 2: 1.0},
	}))
	goassert.New(t, float32(0.0)).Equal(AvgTotalVariationAmongSparseVectors(SparseVectors{
		SparseVector{0: 1.0, 1: 1.0, 2: 1.0},
		SparseVector{0: 1.0, 1: 1.0, 2: 1.0},
	}))
	goassert.New(t, float32(0.25)).Equal(AvgTotalVariationAmongSparseVectors(SparseVectors{
		SparseVector{0: 1.0, 1: 1.0, 2: 1.0, 4: 1.0},
		SparseVector{0: 1.0, 1: 1.0, 2: 1.0, 3: 1.0},
	}))
	goassert.New(t, float32(0.0)).Equal(AvgTotalVariationAmongSparseVectors(createBenchmarkAvgTotalVariationAmongSparseVectors()))
}

func TestSparsifyVector(t *testing.T) {
	goassert.New(t, SparseVector{
		1: 1.0, 3: -1.0,
	}).Equal(SparsifyVector([]float32{0.0, 1.0, 0.0, -1.0, 0.0}))
}

func TestKeyCountMap32(t *testing.T) {
	m := NewKeyCountMap32(27)
	goassert.New(t, 32).Equal(len(m))
	goassert.New(t, map[uint32]uint32{}).Equal(m.Map())
	m.Inc(1)
	m.Inc(2)
	m.Inc(1)
	m.Inc(31)
	goassert.New(t, map[uint32]uint32{1: 2, 2: 1, 31: 1}).Equal(m.Map())
	goassert.New(t, KeyCount32{1, 2}).Equal(m.Get(1))
	goassert.New(t, KeyCount32{2, 1}).Equal(m.Get(2))
	goassert.New(t, KeyCount32{31, 1}).Equal(m.Get(31))
	goassert.New(t, KeyCount32{0, 0}).Equal(m.Get(0))
	goassert.New(t, KeyCount32{0, 0}).Equal(m.Get(63))
}

func TestKeyCounts32ExtractLargestCountsByInsert(t *testing.T) {
	// The keys are ignored.
	kcs := KeyCounts32{
		KeyCount32{0, 4}, KeyCount32{0, 3}, KeyCount32{0, 1}, KeyCount32{0, 5}, KeyCount32{0, 2},
	}
	expected := KeyCounts32{
		KeyCount32{0, 5}, KeyCount32{0, 4}, KeyCount32{0, 3}, KeyCount32{0, 2}, KeyCount32{0, 1},
	}
	goassert.New(t, expected).Equal(kcs.ExtractLargestCountsByInsert(5))
	goassert.New(t, expected[:3]).Equal(kcs.ExtractLargestCountsByInsert(3))
	goassert.New(t, expected).Equal(kcs.ExtractLargestCountsByInsert(10))
}

func TestKeyCounts32SortLargestCountsWithHeap(t *testing.T) {
	clone := func(kcs KeyCounts32) KeyCounts32 {
		kcs2 := make(KeyCounts32, len(kcs))
		copy(kcs2, kcs)
		return kcs2
	}
	// The keys are ignored.
	kcs := KeyCounts32{
		KeyCount32{0, 4}, KeyCount32{0, 3}, KeyCount32{0, 1}, KeyCount32{0, 5}, KeyCount32{0, 2},
	}
	expected := KeyCounts32{
		KeyCount32{0, 5}, KeyCount32{0, 4}, KeyCount32{0, 3}, KeyCount32{0, 2}, KeyCount32{0, 1},
	}
	goassert.New(t, expected).Equal(clone(kcs).SortLargestCountsWithHeap(5))
	goassert.New(t, expected[:3]).Equal(clone(kcs).SortLargestCountsWithHeap(3))
	goassert.New(t, expected[:4]).Equal(clone(kcs).SortLargestCountsWithHeap(4))
	goassert.New(t, expected).Equal(clone(kcs).SortLargestCountsWithHeap(10))
}

func TestKeyValues32OrderedByKey(t *testing.T) {
	data := KeyValues32OrderedByKey(KeyValues32{
		KeyValue32{4, 2.0}, KeyValue32{1, 1.0}, KeyValue32{1, 2.0}, KeyValue32{2, 3.0}, KeyValue32{3, 3.0},
	})
	sort.Sort(data)
	goassert.New(t, KeyValues32{
		KeyValue32{1, 1.0}, KeyValue32{1, 2.0}, KeyValue32{2, 3.0}, KeyValue32{3, 3.0}, KeyValue32{4, 2.0},
	}).Equal(KeyValues32(data))
}

func TestKeyValues32OrderedByValue(t *testing.T) {
	data := KeyValues32OrderedByValue(KeyValues32{
		KeyValue32{5, 2.0}, KeyValue32{1, 1.0}, KeyValue32{2, 1.0}, KeyValue32{3, 1.0}, KeyValue32{4, 3.0},
	})
	sort.Sort(data)
	goassert.New(t, KeyValues32{
		KeyValue32{3, 1.0}, KeyValue32{2, 1.0}, KeyValue32{1, 1.0}, KeyValue32{5, 2.0}, KeyValue32{4, 3.0},
	}).Equal(KeyValues32(data))
}

func TestDotCount(t *testing.T) {
	goassert.New(t, float32(1.0), 1).Equal(DotCount(FeatureVector{KeyValue32{1, 1.0}}, FeatureVector{KeyValue32{1, 1.0}}))
	goassert.New(t, float32(1.0), 1).Equal(DotCount(FeatureVector{KeyValue32{2, 1.0}}, FeatureVector{KeyValue32{1, 1.0}, KeyValue32{2, 1.0}}))
	goassert.New(t, float32(1.0), 1).Equal(DotCount(FeatureVector{KeyValue32{1, 1.0}, KeyValue32{2, 1.0}}, FeatureVector{KeyValue32{2, 1.0}}))
}

func TestFeatureVectors(t *testing.T) {
	X := FeatureVectors{
		FeatureVector{KeyValue32{0, 1.0}},
		FeatureVector{KeyValue32{0, 2.0}, KeyValue32{1, 3.0}},
		FeatureVector{KeyValue32{0, 4.0}, KeyValue32{2, 5.0}, KeyValue32{9, 6.0}},
	}
	goassert.New(t, 10).Equal(X.Dim())
}

func TestLabelVector(t *testing.T) {
	y := LabelVector{3, 2, 1}
	sort.Sort(y)
	goassert.New(t, LabelVector{1, 2, 3}).Equal(y)
}

func TestLabelVectors(t *testing.T) {
	Y := LabelVectors{
		LabelVector{0}, LabelVector{0, 1}, LabelVector{0, 2, 7},
	}
	goassert.New(t, 8).Equal(Y.Dim())
}

func TestDecodeEncodeDataset(t *testing.T) {
	ds := &Dataset{
		X: FeatureVectors{
			FeatureVector{KeyValue32{0, 1.0}},
			FeatureVector{KeyValue32{0, 2.0}, KeyValue32{1, 3.0}},
			FeatureVector{KeyValue32{0, 4.0}, KeyValue32{2, 5.0}, KeyValue32{9, 6.0}},
			FeatureVector{KeyValue32{10, 7.0}},
		},
		Y: LabelVectors{
			LabelVector{0}, LabelVector{0, 1}, LabelVector{0, 2, 9}, LabelVector{10},
		},
	}
	var buf bytes.Buffer
	goassert.New(t).SucceedWithoutError(EncodeDataset(ds, &buf))
	var decodedDs Dataset
	goassert.New(t).SucceedWithoutError(DecodeDataset(&decodedDs, &buf))
	goassert.New(t, ds).Equal(&decodedDs)
	// gob.Decoder.Decode won't call Dataset.GobDecode, because the encoder did not encode Dataset.
	goassert.New(t, "Dataset should be encoded with EncodeDataset").ExpectError(gob.NewEncoder(&buf).Encode(&decodedDs))
}

func TestDatasetFeatureSubSet(t *testing.T) {
	ds := &Dataset{
		X: FeatureVectors{
			FeatureVector{KeyValue32{0, 1.0}},
			FeatureVector{KeyValue32{0, 2.0}, KeyValue32{1, 3.0}},
			FeatureVector{KeyValue32{0, 4.0}, KeyValue32{2, 5.0}, KeyValue32{9, 6.0}},
			FeatureVector{KeyValue32{10, 7.0}},
		},
		Y: LabelVectors{
			LabelVector{0}, LabelVector{0, 1}, LabelVector{0, 2, 9}, LabelVector{10},
		},
	}
	subds := &Dataset{
		X: FeatureVectors{
			FeatureVector{KeyValue32{0, 1.0}},
			FeatureVector{KeyValue32{0, 2.0}},
			FeatureVector{KeyValue32{0, 4.0}, KeyValue32{9, 6.0}},
			FeatureVector{},
		},
		Y: LabelVectors{
			LabelVector{0}, LabelVector{0, 1}, LabelVector{0, 2, 9}, LabelVector{10},
		},
	}
	goassert.New(t, subds).Equal(ds.FeatureSubSet(map[uint32]struct{}{0: {}, 5: {}, 9: {}}))
}

func TestDatasetSize(t *testing.T) {
	ds := &Dataset{
		X: FeatureVectors{
			FeatureVector{KeyValue32{0, 1.0}},
			FeatureVector{KeyValue32{0, 2.0}, KeyValue32{1, 3.0}},
			FeatureVector{KeyValue32{0, 4.0}, KeyValue32{2, 5.0}, KeyValue32{9, 6.0}},
		},
		Y: LabelVectors{
			LabelVector{0}, LabelVector{0, 1}, LabelVector{0, 2, 9}, LabelVector{10},
		},
	}
	goassert.New(t, 3).Equal(ds.Size())
}

func TestDatasetSubSet(t *testing.T) {
	ds := &Dataset{
		X: FeatureVectors{
			FeatureVector{KeyValue32{0, 1.0}},
			FeatureVector{KeyValue32{0, 2.0}, KeyValue32{1, 3.0}},
			FeatureVector{KeyValue32{0, 4.0}, KeyValue32{2, 5.0}, KeyValue32{9, 6.0}},
		},
		Y: LabelVectors{
			LabelVector{0}, LabelVector{0, 1}, LabelVector{0, 2, 9}, LabelVector{10},
		},
	}
	goassert.New(t, &Dataset{
		X: FeatureVectors{ds.X[1], ds.X[1], ds.X[0]},
		Y: LabelVectors{ds.Y[1], ds.Y[1], ds.Y[0]},
	}).Equal(ds.SubSet([]int{1, 1, 0}))
}

func TestDatasetWriteToText(t *testing.T) {
	ds := &Dataset{
		X: FeatureVectors{
			FeatureVector{KeyValue32{0, 1.0}},
			FeatureVector{KeyValue32{0, 2.0}, KeyValue32{1, 3.0}},
			FeatureVector{KeyValue32{0, 4.0}, KeyValue32{2, 5.0}, KeyValue32{9, 6.0}},
			FeatureVector{KeyValue32{10, 7.0}},
		},
		Y: LabelVectors{
			LabelVector{0}, LabelVector{0, 1}, LabelVector{0, 2, 9}, LabelVector{10},
		},
	}
	var buf bytes.Buffer
	goassert.New(t).SucceedWithoutError(ds.WriteTextDataset(&buf))
	ds2 := goassert.New(t).SucceedNew(ReadTextDataset(&buf)).(*Dataset)
	goassert.New(t, ds).Equal(ds2)
}

func TestReadTextDataset(t *testing.T) {
	ds := &Dataset{
		X: FeatureVectors{
			FeatureVector{KeyValue32{0, 1.0}},
			FeatureVector{KeyValue32{0, 2.0}, KeyValue32{1, 3.0}},
			FeatureVector{KeyValue32{0, 4.0}, KeyValue32{2, 5.0}, KeyValue32{9, 6.0}},
		},
		Y: LabelVectors{
			LabelVector{0}, LabelVector{0, 1}, LabelVector{0, 2, 9},
		},
	}
	goassert.New(t, ds).EqualWithoutError(ReadTextDataset(strings.NewReader("3 10 10\r\n0 0:1\n0,1 0:2 1:3\r\n0,2,9 0:4 2:5 9:6\n")))
	goassert.New(t, "cannot read first line").ExpectError(ReadTextDataset(strings.NewReader("")))
	goassert.New(t, "illegal first line").ExpectError(ReadTextDataset(strings.NewReader("3\n")))
	goassert.New(t, "illegal nentries in first line").ExpectError(ReadTextDataset(strings.NewReader("x y z\n")))
	goassert.New(t, "illegal nfeatures in first line").ExpectError(ReadTextDataset(strings.NewReader("3 y z\n")))
	goassert.New(t, "illegal nlabels in first line").ExpectError(ReadTextDataset(strings.NewReader("3 10 z\n")))
	goassert.New(t, "L1: cannot read line").ExpectError(ReadTextDataset(strings.NewReader("3 10 10\n")))
	goassert.New(t, "L1: illegal #1 label ID").ExpectError(ReadTextDataset(strings.NewReader("3 10 10\nx\n")))
	goassert.New(t, "L1: too large #1 label ID \\(>= 10\\)").ExpectError(ReadTextDataset(strings.NewReader("3 10 10\n10\n")))
	goassert.New(t, "L1: illegal #1 featureID:value pair").ExpectError(ReadTextDataset(strings.NewReader("3 10 10\n0 x\n")))
	goassert.New(t, "L1: illegal #1 feature ID").ExpectError(ReadTextDataset(strings.NewReader("3 10 10\n0 x:y\n")))
	goassert.New(t, "L1: too large #1 feature ID \\(>= 10\\)").ExpectError(ReadTextDataset(strings.NewReader("3 10 10\n0 10:y\n")))
	goassert.New(t, "L1: illegal #1 feature value").ExpectError(ReadTextDataset(strings.NewReader("3 10 10\n0 0:y\n")))
}

func createBenchmarkAvgTotalVariationAmongSparseVectors() SparseVectors {
	n, m := 50, 20
	svs := make(SparseVectors, n)
	for i := range svs {
		sv := make(SparseVector, m)
		for j := 0; j < m; j++ {
			sv[uint32(j)] = float32(j)
		}
		svs[i] = sv
	}
	return svs
}

func BenchmarkAvgTotalVariationAmongSparseVectors(b *testing.B) {
	svs := createBenchmarkAvgTotalVariationAmongSparseVectors()
	b.ResetTimer()
	for t := 0; t < b.N; t++ {
		AvgTotalVariationAmongSparseVectors(svs)
	}
}

func BenchmarkKeyCounts32ExtractLargestCountsByInsert(b *testing.B) {
	rng := rand.New(rand.NewSource(0))
	for t := 0; t < b.N; t++ {
		// The keys are ignored.
		kcs := make(KeyCounts32, 64*4096)
		for i := 0; i < len(kcs)/2; i++ {
			kcs[rng.Intn(len(kcs))].Count = rng.Uint32()
		}
		kcs.ExtractLargestCountsByInsert(2 * 75)
	}
}

func BenchmarkKeyCounts32SortLargestCountsWithHeap(b *testing.B) {
	rng := rand.New(rand.NewSource(0))
	for t := 0; t < b.N; t++ {
		// The keys are ignored.
		kcs := make(KeyCounts32, 64*4096)
		for i := 0; i < len(kcs)/2; i++ {
			kcs[rng.Intn(len(kcs))].Count = rng.Uint32()
		}
		kcs.SortLargestCountsWithHeap(2 * 75)
	}
}

func BenchmarkKeyCountMap32(b *testing.B) {
	rng := rand.New(rand.NewSource(0))
	for t := 0; t < b.N; t++ {
		m := NewKeyCountMap32(64 * 1024)
		for i := 0; i < len(m)/2; i++ {
			m.Inc(rng.Uint32())
		}
	}
}

func BenchmarkMapUint32Uint32(b *testing.B) {
	rng := rand.New(rand.NewSource(0))
	for t := 0; t < b.N; t++ {
		m := make(map[uint32]uint32)
		for i := 0; i < 64*1024/2; i++ {
			m[rng.Uint32()] += 1
		}
	}
}
