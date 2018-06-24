package sticker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hiro4bbh/go-assert"
)

func assertBinaryClassAssignmentEqual(t *testing.T, expected, got []bool) bool {
	t.Helper()
	if expected[0] == got[0] {
		goassert.New(t, expected).Equal(got)
		return false
	}
	expectedNeg := make([]bool, len(expected))
	for i, b := range expected {
		expectedNeg[i] = !b
	}
	goassert.New(t, expectedNeg).Equal(got)
	return true
}

func TestIdealDCG(t *testing.T) {
	IdealDCG0 := IdealDCG(0)
	goassert.New(t, "NaN (float32)").Equal(fmt.Sprintf("%#v (%T)", IdealDCG0, IdealDCG0))
	goassert.New(t, float32(1.0)).Equal(IdealDCG(1))
	goassert.New(t, float32(1.0+1.0/LogBinary32(1.0+2.0))).Equal(IdealDCG(2))
	goassert.New(t, float32(1.0+1.0/LogBinary32(1.0+2.0)+0.5)).Equal(IdealDCG(3))
}

func TestInvertRanks(t *testing.T) {
	goassert.New(t, map[uint32]int{0: 1, 1: 3, 2: 4, 3: 5, 9: 2}).Equal(InvertRanks(LabelVector{0, 9, 1, 2, 3, ^uint32(0)}))
}

func TestRankTopK(t *testing.T) {
	// Test on small distributions.
	goassert.New(t, LabelVector{}).Equal(RankTopK(SparseVector{0: 4, 1: 2, 2: 2, 3: 2, 9: 4}, 0))
	goassert.New(t, LabelVector{0}).Equal(RankTopK(SparseVector{0: 4, 1: 2, 2: 2, 3: 2, 9: 4}, 1))
	goassert.New(t, LabelVector{0, 9}).Equal(RankTopK(SparseVector{0: 4, 1: 2, 2: 2, 3: 2, 9: 4}, 2))
	goassert.New(t, LabelVector{0, 9, 1}).Equal(RankTopK(SparseVector{0: 4, 1: 2, 2: 2, 3: 2, 9: 4}, 3))
	goassert.New(t, LabelVector{0, 9, 1, 2, 3, ^uint32(0)}).Equal(RankTopK(SparseVector{0: 4, 1: 2, 2: 2, 3: 2, 9: 4}, 6))
	// Test on large distributions.
	largeLabelDist, n, m := SparseVector{0: 5, 1: 4, 2: 3, 3: 2, 4: 1}, uint(1000), uint(5)
	for i := len(largeLabelDist); i < int(n); i++ {
		largeLabelDist[uint32(i)] = 0
	}
	ranks := make(LabelVector, n+m)
	for i := 0; i < int(n); i++ {
		ranks[i] = uint32(i)
	}
	for i := int(n); i < len(ranks); i++ {
		ranks[i] = ^uint32(0)
	}
	goassert.New(t, ranks[:n/10]).Equal(RankTopK(largeLabelDist, n/10))
	goassert.New(t, ranks[:n]).Equal(RankTopK(largeLabelDist, n))
	goassert.New(t, ranks[:n+m]).Equal(RankTopK(largeLabelDist, n+m))
}

func TestReportNDCG(t *testing.T) {
	Y := LabelVectors{
		LabelVector{0, 9}, LabelVector{9}, LabelVector{8},
		LabelVector{0, 3, 9}, LabelVector{0, 8, 9}, LabelVector{0, 9},
		LabelVector{0, 2, 3, 9},
	}
	// Illegal case: K = 0
	goassert.New(t, "[]float32{NaN"+strings.Repeat(", NaN", len(Y)-1)+"}").Equal(fmt.Sprintf("%#v", ReportNDCG(Y, 0, LabelVectors{
		LabelVector{}, LabelVector{}, LabelVector{},
		LabelVector{}, LabelVector{}, LabelVector{},
		LabelVector{},
	})))
	goassert.New(t, []float32{
		1.0, 0.0, 0.0,
		1.0, 1.0, 1.0,
		1.0,
	}).Equal(ReportNDCG(Y, 1, LabelVectors{
		LabelVector{0}, LabelVector{0}, LabelVector{0},
		LabelVector{9}, LabelVector{9}, LabelVector{9},
		LabelVector{9},
	}))
	goassert.New(t, []float32{
		1.0, (1.0 / LogBinary32(1.0+2.0)) / IdealDCG(1), 0.0,
		1.0, 1.0, 1.0,
		1.0,
	}).Equal(ReportNDCG(Y, 2, LabelVectors{
		LabelVector{0, 9}, LabelVector{0, 9}, LabelVector{0, 9},
		LabelVector{9, 0}, LabelVector{9, 0}, LabelVector{9, 0},
		LabelVector{9, 0},
	}))
	goassert.New(t, []float32{
		1.0, (1.0 / LogBinary32(1.0+2.0)) / IdealDCG(1), 0.0,
		1.0, (1.0 + 1.0/LogBinary32(1.0+2.0)) / IdealDCG(3), 1.0,
		1.0,
	}).Equal(ReportNDCG(Y, 3, LabelVectors{
		LabelVector{0, 9, ^uint32(0)}, LabelVector{0, 9, ^uint32(0)}, LabelVector{0, 9, ^uint32(0)},
		LabelVector{9, 0, 3}, LabelVector{9, 0, 3}, LabelVector{9, 0, 3},
		LabelVector{9, 0, 2},
	}))
	goassert.New(t, []float32{
		1.0, (1.0 / LogBinary32(1.0+2.0)) / IdealDCG(1), 0.0,
		1.0, (1.0 + 1.0/LogBinary32(1.0+2.0)) / IdealDCG(3), 1.0,
		1.0,
	}).Equal(ReportNDCG(Y, 4, LabelVectors{
		LabelVector{0, 9, ^uint32(0), ^uint32(0)}, LabelVector{0, 9, ^uint32(0), ^uint32(0)}, LabelVector{0, 9, ^uint32(0), ^uint32(0)},
		LabelVector{9, 0, 3, ^uint32(0)}, LabelVector{9, 0, 3, ^uint32(0)}, LabelVector{9, 0, 3, ^uint32(0)},
		LabelVector{9, 0, 2, 3},
	}))
	goassert.New(t, []float32{
		1.0, (1.0 / LogBinary32(1.0+2.0)) / IdealDCG(1), 0.0,
		1.0, (1.0 + 1.0/LogBinary32(1.0+2.0)) / IdealDCG(3), 1.0,
		1.0,
	}).Equal(ReportNDCG(Y, 5, LabelVectors{
		LabelVector{0, 9, ^uint32(0), ^uint32(0), ^uint32(0)}, LabelVector{0, 9, ^uint32(0), ^uint32(0), ^uint32(0)}, LabelVector{0, 9, ^uint32(0), ^uint32(0), ^uint32(0)},
		LabelVector{9, 0, 3, ^uint32(0), ^uint32(0)}, LabelVector{9, 0, 3, ^uint32(0), ^uint32(0)}, LabelVector{9, 0, 3, ^uint32(0), ^uint32(0)},
		LabelVector{9, 0, 2, 3, ^uint32(0)},
	}))
	// Also with smaller K, it should work.
	goassert.New(t, []float32{
		1.0, (1.0 / LogBinary32(1.0+2.0)) / IdealDCG(1), 0.0,
		1.0, 1.0, 1.0,
		1.0,
	}).Equal(ReportNDCG(Y, 2, LabelVectors{
		LabelVector{0, 9, ^uint32(0), ^uint32(0), ^uint32(0)}, LabelVector{0, 9, ^uint32(0), ^uint32(0), ^uint32(0)}, LabelVector{0, 9, ^uint32(0), ^uint32(0), ^uint32(0)},
		LabelVector{9, 0, 3, ^uint32(0), ^uint32(0)}, LabelVector{9, 0, 3, ^uint32(0), ^uint32(0)}, LabelVector{9, 0, 3, ^uint32(0), ^uint32(0)},
		LabelVector{9, 0, 2, 3, ^uint32(0)},
	}))
}

func TestReportMaxPrecision(t *testing.T) {
	Y := LabelVectors{
		LabelVector{0, 9}, LabelVector{9}, LabelVector{8},
		LabelVector{0, 3, 9}, LabelVector{0, 8, 9}, LabelVector{0, 9},
		LabelVector{0, 2, 3, 9},
	}
	// Illegal case: K = 0
	goassert.New(t, "[]float32{NaN"+strings.Repeat(", NaN", len(Y)-1)+"}").Equal(fmt.Sprintf("%#v", ReportMaxPrecision(Y, 0)))
	goassert.New(t, []float32{
		1.0 / 1, 1.0 / 1, 1.0 / 1,
		1.0 / 1, 1.0 / 1, 1.0 / 1,
		1.0 / 1,
	}).Equal(ReportMaxPrecision(Y, 1))
	goassert.New(t, []float32{
		2.0 / 2, 1.0 / 2, 1.0 / 2,
		2.0 / 2, 2.0 / 2, 2.0 / 2,
		2.0 / 2,
	}).Equal(ReportMaxPrecision(Y, 2))
	goassert.New(t, []float32{
		2.0 / 3, 1.0 / 3, 1.0 / 3,
		3.0 / 3, 3.0 / 3, 2.0 / 3,
		3.0 / 3,
	}).Equal(ReportMaxPrecision(Y, 3))
	goassert.New(t, []float32{
		2.0 / 4, 1.0 / 4, 1.0 / 4,
		3.0 / 4, 3.0 / 4, 2.0 / 4,
		4.0 / 4,
	}).Equal(ReportMaxPrecision(Y, 4))
	goassert.New(t, []float32{
		2.0 / 5, 1.0 / 5, 1.0 / 5,
		3.0 / 5, 3.0 / 5, 2.0 / 5,
		4.0 / 5,
	}).Equal(ReportMaxPrecision(Y, 5))
}

func TestReportPrecision(t *testing.T) {
	ds := &Dataset{
		X: FeatureVectors{
			FeatureVector{KeyValue32{0, -1.0}}, FeatureVector{KeyValue32{0, -1.0}}, FeatureVector{KeyValue32{0, -1.0}},
			FeatureVector{KeyValue32{0, 1.0}}, FeatureVector{KeyValue32{0, 0.0}}, FeatureVector{KeyValue32{0, 0.0}},
			FeatureVector{KeyValue32{0, 1.0}},
		},
		Y: LabelVectors{
			LabelVector{0, 9}, LabelVector{9}, LabelVector{8},
			LabelVector{0, 3, 9}, LabelVector{0, 8, 9}, LabelVector{0, 9},
			LabelVector{0, 2, 3, 9},
		},
	}
	// Illegal case: K = 0
	goassert.New(t, "[]float32{NaN"+strings.Repeat(", NaN", len(ds.X)-1)+"}").Equal(fmt.Sprintf("%#v", ReportPrecision(ds.Y, 0, LabelVectors{
		LabelVector{}, LabelVector{}, LabelVector{},
		LabelVector{}, LabelVector{}, LabelVector{},
		LabelVector{},
	})))
	goassert.New(t, []float32{
		1.0 / 1, 0.0 / 1, 0.0 / 1,
		1.0 / 1, 1.0 / 1, 1.0 / 1,
		1.0 / 1,
	}).Equal(ReportPrecision(ds.Y, 1, LabelVectors{
		LabelVector{0}, LabelVector{0}, LabelVector{0},
		LabelVector{9}, LabelVector{9}, LabelVector{9},
		LabelVector{9},
	}))
	goassert.New(t, []float32{
		2.0 / 2, 1.0 / 2, 0.0 / 2,
		2.0 / 2, 2.0 / 2, 2.0 / 2,
		2.0 / 2,
	}).Equal(ReportPrecision(ds.Y, 2, LabelVectors{
		LabelVector{0, 9}, LabelVector{0, 9}, LabelVector{0, 9},
		LabelVector{9, 0}, LabelVector{9, 0}, LabelVector{9, 0},
		LabelVector{9, 0},
	}))
	goassert.New(t, []float32{
		2.0 / 3, 1.0 / 3, 0.0 / 3,
		3.0 / 3, 2.0 / 3, 2.0 / 3,
		3.0 / 3,
	}).Equal(ReportPrecision(ds.Y, 3, LabelVectors{
		LabelVector{0, 9, ^uint32(0)}, LabelVector{0, 9, ^uint32(0)}, LabelVector{0, 9, ^uint32(0)},
		LabelVector{9, 0, 3}, LabelVector{9, 0, 3}, LabelVector{9, 0, 3},
		LabelVector{9, 0, 2},
	}))
	goassert.New(t, []float32{
		2.0 / 4, 1.0 / 4, 0.0 / 4,
		3.0 / 4, 2.0 / 4, 2.0 / 4,
		4.0 / 4,
	}).Equal(ReportPrecision(ds.Y, 4, LabelVectors{
		LabelVector{0, 9, ^uint32(0), ^uint32(0)}, LabelVector{0, 9, ^uint32(0), ^uint32(0)}, LabelVector{0, 9, ^uint32(0), ^uint32(0)},
		LabelVector{9, 0, 3, ^uint32(0)}, LabelVector{9, 0, 3, ^uint32(0)}, LabelVector{9, 0, 3, ^uint32(0)},
		LabelVector{9, 0, 2, 3},
	}))
	goassert.New(t, []float32{
		2.0 / 5, 1.0 / 5, 0.0 / 5,
		3.0 / 5, 2.0 / 5, 2.0 / 5,
		4.0 / 5,
	}).Equal(ReportPrecision(ds.Y, 5, LabelVectors{
		LabelVector{0, 9, ^uint32(0), ^uint32(0), ^uint32(0)}, LabelVector{0, 9, ^uint32(0), ^uint32(0), ^uint32(0)}, LabelVector{0, 9, ^uint32(0), ^uint32(0), ^uint32(0)},
		LabelVector{9, 0, 3, ^uint32(0), ^uint32(0)}, LabelVector{9, 0, 3, ^uint32(0), ^uint32(0)}, LabelVector{9, 0, 3, ^uint32(0), ^uint32(0)},
		LabelVector{9, 0, 2, 3, ^uint32(0)},
	}))
	// Also with smaller K, it should work.
	goassert.New(t, []float32{
		2.0 / 2, 1.0 / 2, 0.0 / 2,
		2.0 / 2, 2.0 / 2, 2.0 / 2,
		2.0 / 2,
	}).Equal(ReportPrecision(ds.Y, 2, LabelVectors{
		LabelVector{0, 9, ^uint32(0), ^uint32(0), ^uint32(0)}, LabelVector{0, 9, ^uint32(0), ^uint32(0), ^uint32(0)}, LabelVector{0, 9, ^uint32(0), ^uint32(0), ^uint32(0)},
		LabelVector{9, 0, 3, ^uint32(0), ^uint32(0)}, LabelVector{9, 0, 3, ^uint32(0), ^uint32(0)}, LabelVector{9, 0, 3, ^uint32(0), ^uint32(0)},
		LabelVector{9, 0, 2, 3, ^uint32(0)},
	}))
}
