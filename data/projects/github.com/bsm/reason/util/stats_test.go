package util_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/bsm/reason/util"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("StreamStats", func() {
	var subject *util.StreamStats
	var blank = new(util.StreamStats)
	var weight1 = &util.StreamStats{Weight: 1.0}

	BeforeEach(func() {
		subject = new(util.StreamStats)
		for _, v := range []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9} {
			subject.Add(v, 1)
		}
	})

	It("should return total weight", func() {
		Expect(subject.Weight).To(Equal(9.0))
		subject.Add(2.2, 2)
		Expect(subject.Weight).To(Equal(11.0))
		Expect(blank.Weight).To(Equal(0.0))
	})

	It("should return value sum", func() {
		Expect(subject.Sum).To(Equal(49.5))
		Expect(blank.Sum).To(Equal(0.0))
	})

	It("should calc mean", func() {
		Expect(subject.Mean()).To(Equal(5.5))
		subject.Add(8.8, 8)
		Expect(subject.Mean()).To(BeNumerically("~", 7.05, 0.01))
		Expect(weight1.Mean()).To(Equal(0.0))
		Expect(math.IsNaN(blank.Mean())).To(BeTrue())
	})

	It("should calc variance", func() {
		Expect(subject.Variance()).To(BeNumerically("~", 9.07, 0.01))
		Expect(weight1.Variance()).To(Equal(0.0))
		Expect(blank.Variance()).To(Equal(0.0))
	})

	It("should calc std-dev", func() {
		Expect(subject.StdDev()).To(BeNumerically("~", 3.01, 0.01))
		Expect(weight1.StdDev()).To(Equal(0.0))
		Expect(blank.StdDev()).To(Equal(0.0))
	})

	It("should calculate probability density", func() {
		Expect(subject.Prob(1.2)).To(BeNumerically("~", 0.048, 0.001))
		Expect(subject.Prob(5.5)).To(BeNumerically("~", 0.132, 0.001))
		Expect(subject.Prob(13.3)).To(BeNumerically("~", 0.005, 0.001))
		Expect(subject.Prob(24.6)).To(BeNumerically("~", 0.000, 0.001))
		Expect(math.IsNaN(weight1.Prob(5.5))).To(BeTrue())
		Expect(math.IsNaN(blank.Prob(5.5))).To(BeTrue())
	})

	DescribeTable("should estimate",
		func(v, xlt, xeq, xgt float64) {
			lt, eq, gt := subject.Estimate(v)
			Expect(lt).To(BeNumerically("~", xlt, 0.01))
			Expect(eq).To(BeNumerically("~", xeq, 0.01))
			Expect(gt).To(BeNumerically("~", xgt, 0.01))
		},
		Entry("lower end", 1.2, 0.26, 0.43, 8.31),
		Entry("close to mean", 5.4, 3.19, 1.19, 4.62),
		Entry("top end", 9.1, 7.37, 0.58, 1.04),
	)

	It("should fail to estimate on insufficient weight", func() {
		lt, eq, gt := weight1.Estimate(1.2)
		Expect(math.IsNaN(lt)).To(BeTrue())
		Expect(math.IsNaN(eq)).To(BeTrue())
		Expect(math.IsNaN(gt)).To(BeTrue())

		lt, eq, gt = blank.Estimate(1.2)
		Expect(math.IsNaN(lt)).To(BeTrue())
		Expect(math.IsNaN(eq)).To(BeTrue())
		Expect(math.IsNaN(gt)).To(BeTrue())
	})

})

var _ = Describe("StreamStatsDistribution", func() {
	var sparse, dense *util.StreamStatsDistribution

	BeforeEach(func() {
		sparse = new(util.StreamStatsDistribution)
		dense = &util.StreamStatsDistribution{Dense: []util.StreamStatsDistribution_Dense{}}
		for _, v := range []float64{1.1, 2.2, 3.3, 4.4} {
			sparse.Add(0, v, 1)
			dense.Add(0, v, 1)
		}
		for _, v := range []float64{5.5, 6.6, 7.7, 8.8, 9.9} {
			sparse.Add(1, v, 1)
			dense.Add(1, v, 1)
		}
	})

	It("should add", func() {
		sparse.Add(7, 12.12, 1)
		Expect(sparse.Get(7)).NotTo(BeNil())

		dense.Add(7, 12.12, 1)
		Expect(dense.Get(7)).NotTo(BeNil())
	})

	It("should have len", func() {
		Expect(sparse.Len()).To(Equal(2))
		Expect(dense.Len()).To(Equal(2))
		sparse.Add(7, 12.12, 1)
		dense.Add(7, 12.12, 1)
		Expect(sparse.Len()).To(Equal(3))
		Expect(dense.Len()).To(Equal(3))
	})

	It("should get", func() {
		Expect(sparse.Get(0)).To(Equal(&util.StreamStats{Weight: 4, Sum: 11, SumSquares: 36.3}))
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
			sparse.Add(i, 1.1, 1)
		}
		Expect(sparse.Len()).To(Equal(52))
		Expect(sparse.Dense).To(BeNil())
		Expect(len(sparse.Sparse)).To(Equal(52))
		Expect(sparse.SparseCap).To(Equal(int64(2000)))

		for i := 999; i >= 800; i-- {
			sparse.Add(i, 1.1, 1)
		}
		Expect(sparse.Len()).To(Equal(252))
		Expect(len(sparse.Dense)).To(Equal(2000))
		Expect(sparse.Sparse).To(BeNil())
		Expect(sparse.SparseCap).To(Equal(int64(0)))

		_, err := proto.Marshal(sparse)
		Expect(err).NotTo(HaveOccurred())
	})

})

func BenchmarkStreamStatsDistribution(b *testing.B) {
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
			s := new(util.StreamStatsDistribution)
			for i := 0; i < t.N; i++ {
				offset := int(float64(i) * float64(t.MaxIndex) / float64(t.N))
				s.Add(t.MaxIndex-offset, 2.2, 1.0)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				n, m := 0, 0
				s.ForEach(func(i int, s *util.StreamStats) bool {
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
