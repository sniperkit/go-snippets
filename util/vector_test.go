package util_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/bsm/reason/util"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vector", func() {
	var sparse, dense *util.Vector
	var blank = new(util.Vector)

	BeforeEach(func() {
		sparse = new(util.Vector)
		sparse.Set(0, 2.0)
		sparse.Set(3, 7.0)
		sparse.Set(4, 9.0)

		dense = &util.Vector{Dense: []float64{}}
		dense.Set(0, 2.0)
		dense.Set(3, 7.0)
		dense.Set(4, 9.0)
	})

	It("should set", func() {
		Expect(sparse).To(Equal(&util.Vector{Sparse: map[int64]float64{0: 2, 3: 7, 4: 9}, SparseCap: 5}))
		Expect(dense).To(Equal(&util.Vector{Dense: []float64{2, 0, 0, 7, 9}}))

		sparse.Set(6, 8.0)
		sparse.Set(-2, 4.0)
		Expect(sparse).To(Equal(&util.Vector{Sparse: map[int64]float64{0: 2, 3: 7, 4: 9, 6: 8}, SparseCap: 7}))
	})

	It("should get", func() {
		Expect(sparse.Get(0)).To(Equal(2.0))
		Expect(sparse.Get(1)).To(Equal(0.0))
		Expect(sparse.Get(-1)).To(Equal(0.0))
		Expect(dense.Get(0)).To(Equal(2.0))
		Expect(dense.Get(1)).To(Equal(0.0))
		Expect(dense.Get(-1)).To(Equal(0.0))
	})

	It("should add", func() {
		sparse.Add(0, 1.0)
		sparse.Add(1, 7.0)
		sparse.Add(-1, 6.0)
		Expect(sparse).To(Equal(&util.Vector{Sparse: map[int64]float64{0: 3, 1: 7, 3: 7, 4: 9}, SparseCap: 5}))

		dense.Add(0, 1.0)
		dense.Add(1, 7.0)
		dense.Add(-1, 6.0)
		Expect(dense).To(Equal(&util.Vector{Dense: []float64{3, 7, 0, 7, 9}}))
	})

	It("should normalize", func() {
		sparse.Normalize()
		Expect(sparse).To(Equal(&util.Vector{Sparse: map[int64]float64{0: (1.0 / 9.0), 3: (7.0 / 18.0), 4: 0.5}, SparseCap: 5}))

		dense.Normalize()
		Expect(dense).To(Equal(&util.Vector{Dense: []float64{(1.0 / 9.0), 0, 0, (7.0 / 18.0), 0.5}}))
	})

	It("should count", func() {
		Expect(sparse.Len()).To(Equal(3))
		Expect(dense.Len()).To(Equal(3))
		Expect(blank.Len()).To(Equal(0))
	})

	It("should clear", func() {
		sparse.Clear()
		Expect(sparse).To(Equal(&util.Vector{Sparse: map[int64]float64{}}))

		dense.Clear()
		Expect(dense).To(Equal(&util.Vector{Dense: []float64{}}))
	})

	It("should calculate total weight", func() {
		Expect(sparse.Weight()).To(Equal(18.0))
		Expect(dense.Weight()).To(Equal(18.0))
		Expect(blank.Weight()).To(Equal(0.0))
	})

	It("should calculate min", func() {
		i, w := sparse.Min()
		Expect(i).To(Equal(0))
		Expect(w).To(Equal(2.0))

		i, w = dense.Min()
		Expect(i).To(Equal(0))
		Expect(w).To(Equal(2.0))

		i, w = blank.Min()
		Expect(i).To(Equal(-1))
		Expect(w).To(Equal(0.0))
	})

	It("should calculate max", func() {
		i, w := sparse.Max()
		Expect(i).To(Equal(4))
		Expect(w).To(Equal(9.0))

		i, w = dense.Max()
		Expect(i).To(Equal(4))
		Expect(w).To(Equal(9.0))

		i, w = blank.Max()
		Expect(i).To(Equal(-1))
		Expect(w).To(Equal(0.0))
	})

	It("should calculate mean", func() {
		Expect(sparse.Mean()).To(Equal(6.0))
		Expect(dense.Mean()).To(Equal(6.0))
		Expect(blank.Mean()).To(Equal(0.0))
	})

	It("should iterate (sparse)", func() {
		indicies := make([]int, 0)
		values := make([]float64, 0)
		sparse.ForEach(func(i int, v float64) bool {
			indicies = append(indicies, i)
			values = append(values, v)
			return true
		})
		Expect(indicies).To(ConsistOf(0, 3, 4))
		Expect(values).To(ConsistOf(2.0, 7.0, 9.0))

		values = values[:0]
		sparse.ForEachValue(func(v float64) bool {
			values = append(values, v)
			return true
		})
		Expect(values).To(ConsistOf(2.0, 7.0, 9.0))
	})

	It("should iterate (dense)", func() {
		indicies := make([]int, 0)
		values := make([]float64, 0)
		dense.ForEach(func(i int, v float64) bool {
			indicies = append(indicies, i)
			values = append(values, v)
			return true
		})
		Expect(indicies).To(ConsistOf(0, 3, 4))
		Expect(values).To(ConsistOf(2.0, 7.0, 9.0))

		values = values[:0]
		dense.ForEachValue(func(v float64) bool {
			values = append(values, v)
			return true
		})
		Expect(values).To(ConsistOf(2.0, 7.0, 9.0))
	})

	It("should calculate variances", func() {
		Expect(sparse.Variance()).To(BeNumerically("~", 13.0, 0.01))
		Expect(dense.Variance()).To(BeNumerically("~", 13.0, 0.01))
		Expect(math.IsNaN(blank.Variance())).To(BeTrue())
	})

	It("should calculate standard deviations", func() {
		Expect(sparse.StdDev()).To(BeNumerically("~", 3.61, 0.01))
		Expect(dense.StdDev()).To(BeNumerically("~", 3.61, 0.01))
		Expect(math.IsNaN(blank.StdDev())).To(BeTrue())
	})

	It("should calculate entropy", func() {
		Expect(sparse.Entropy()).To(BeNumerically("~", 1.38, 0.01))
		Expect(dense.Entropy()).To(BeNumerically("~", 1.38, 0.01))
	})

})

var _ = Describe("VectorDistribution", func() {
	var sparse, dense *util.VectorDistribution

	BeforeEach(func() {
		sparse = new(util.VectorDistribution)
		dense = &util.VectorDistribution{Dense: []util.VectorDistribution_Dense{}}
		for _, n := range []int{1, 2, 3, 4} {
			sparse.Add(0, n, 1.0)
			dense.Add(0, n, 1.0)
		}
		for _, n := range []int{5, 6, 7, 8, 9} {
			sparse.Add(1, n, 1.0)
			dense.Add(1, n, 1.0)
		}
	})

	It("should add", func() {
		sparse.Add(7, 12, 1)
		dense.Add(7, 12, 1)
		Expect(sparse.Get(7)).NotTo(BeNil())
		Expect(dense.Get(7)).NotTo(BeNil())
	})

	It("should have len", func() {
		Expect(sparse.Len()).To(Equal(2))
		Expect(dense.Len()).To(Equal(2))
		sparse.Add(7, 12, 1)
		dense.Add(7, 12, 1)
		Expect(sparse.Len()).To(Equal(3))
		Expect(dense.Len()).To(Equal(3))
	})

	It("should get", func() {
		Expect(sparse.Get(0)).To(Equal(&util.Vector{Sparse: map[int64]float64{1: 1, 2: 1, 3: 1, 4: 1}, SparseCap: 5}))
		Expect(sparse.Get(2)).To(BeNil())
		Expect(sparse.Get(-1)).To(BeNil())

		Expect(dense.Get(0)).To(Equal(sparse.Get(0)))
		Expect(dense.Get(2)).To(BeNil())
		Expect(dense.Get(-1)).To(BeNil())
	})

	It("should convert to dense", func() {
		Expect(sparse.Dense).To(BeNil())
		Expect(sparse.Sparse).To(HaveLen(2))
		Expect(sparse.SparseCap).To(Equal(int64(2)))

		for i := 1999; i >= 1950; i-- {
			sparse.Add(i, 2, 1)
		}
		Expect(sparse.Len()).To(Equal(52))
		Expect(sparse.Dense).To(BeNil())
		Expect(len(sparse.Sparse)).To(Equal(52))
		Expect(sparse.SparseCap).To(Equal(int64(2000)))

		for i := 999; i >= 800; i-- {
			sparse.Add(i, 2, 1)
		}
		Expect(sparse.Len()).To(Equal(252))
		Expect(len(sparse.Dense)).To(Equal(2000))
		Expect(sparse.Sparse).To(BeNil())
		Expect(sparse.SparseCap).To(Equal(int64(0)))

		_, err := proto.Marshal(sparse)
		Expect(err).NotTo(HaveOccurred())
	})

})

func BenchmarkVector(b *testing.B) {
	for _, t := range []struct {
		N, MaxIndex int
	}{
		{10, 100},
		{100, 100},

		{10, 1000},
		{100, 1000},
		{200, 1000},
		{1000, 1000},

		{10, 10000},
		{100, 10000},
		{1000, 10000},
		{2000, 10000},
	} {
		b.Run(fmt.Sprintf("%d/%d", t.N, t.MaxIndex), func(b *testing.B) {
			s := new(util.Vector)
			for i := 0; i < t.N; i++ {
				offset := int(float64(i) * float64(t.MaxIndex) / float64(t.N))
				s.Set(t.MaxIndex-offset, 1.0)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				n, m := 0, 0
				s.ForEach(func(i int, w float64) bool {
					n++
					if i > m {
						m = i
					}
					return true
				})
				if n != t.N {
					b.Fatalf("expected N to be %d, but was %d", t.N, n)
				}
				if m != t.MaxIndex {
					b.Fatalf("expected MaxIndex to be %d, but was %d", t.MaxIndex, m)
				}
			}
		})
	}
}

func BenchmarkVectorDistribution(b *testing.B) {
	for _, t := range []struct {
		N, MaxIndex int
	}{
		{10, 100},
		{100, 100},

		{10, 1000},
		{100, 1000},
		{200, 1000},
		{1000, 1000},

		{10, 10000},
		{100, 10000},
		{1000, 10000},
		{2000, 10000},
	} {
		b.Run(fmt.Sprintf("%d/%d", t.N, t.MaxIndex), func(b *testing.B) {
			s := new(util.VectorDistribution)
			for i := 0; i < t.N; i++ {
				offset := int(float64(i) * float64(t.MaxIndex) / float64(t.N))
				s.Add(t.MaxIndex-offset, 2, 1.0)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				n, m := 0, 0
				s.ForEach(func(i int, s *util.Vector) bool {
					n++
					if i > m {
						m = i
					}
					return true
				})
				if n != t.N {
					b.Fatalf("expected N to be %d, but was %d", t.N, n)
				}
				if m != t.MaxIndex {
					b.Fatalf("expected MaxIndex to be %d, but was %d", t.MaxIndex, m)
				}
			}
		})
	}
}
