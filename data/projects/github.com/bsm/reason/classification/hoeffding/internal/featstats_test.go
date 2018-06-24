package internal_test

import (
	"github.com/bsm/reason/classification/hoeffding/internal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FeatureStats_Numerical", func() {
	var subject *internal.FeatureStats_Numerical

	BeforeEach(func() {
		subject = new(internal.FeatureStats_Numerical)
		subject.Add(1.4, 0, 1.0)
		subject.Add(1.3, 0, 1.0)
		subject.Add(1.5, 0, 1.0)
		subject.Add(4.1, 1, 1.0)
		subject.Add(3.7, 1, 1.0)
		subject.Add(4.9, 1, 1.0)
		subject.Add(4.0, 1, 1.0)
		subject.Add(3.3, 1, 1.0)
		subject.Add(6.3, 2, 1.0)
		subject.Add(5.8, 2, 1.0)
		subject.Add(5.1, 2, 1.0)
		subject.Add(5.3, 2, 1.0)
	})

	It("should add", func() {
		Expect(subject.Min.Sparse).To(Equal(map[int64]float64{0: 1.3, 1: 3.3, 2: 5.1}))
		Expect(subject.Max.Sparse).To(Equal(map[int64]float64{0: 1.5, 1: 4.9, 2: 6.3}))
		Expect(subject.Stats.Len()).To(Equal(3))
	})

	It("should calculate pivot points", func() {
		pp := subject.PivotPoints()
		Expect(pp).To(HaveLen(11))
		Expect(pp[0]).To(BeNumerically("~", 1.72, 0.01))
		Expect(pp[1]).To(BeNumerically("~", 2.13, 0.01))
		Expect(pp[9]).To(BeNumerically("~", 5.47, 0.01))
		Expect(pp[10]).To(BeNumerically("~", 5.88, 0.01))

		Expect(new(internal.FeatureStats_Numerical).PivotPoints()).To(BeEmpty())
	})

	It("should calculate post-splits", func() {
		s1 := subject.PostSplit(2.4)
		Expect(s1.Len()).To(Equal(2))
		Expect(s1.Get(0).Sparse).To(Equal(map[int64]float64{0: 3}))
		Expect(s1.Get(1).Sparse).To(Equal(map[int64]float64{1: 5, 2: 4}))

		s2 := subject.PostSplit(4.8)
		Expect(s2.Len()).To(Equal(2))
		Expect(s2.Get(0).Sparse).To(Equal(map[int64]float64{0: 3, 1: 4.55925906389872}))
		Expect(s2.Get(1).Sparse).To(Equal(map[int64]float64{1: 0.44074093610127996, 2: 4}))
	})
})

var _ = Describe("FeatureStats_Categorical", func() {
	var subject *internal.FeatureStats_Categorical

	BeforeEach(func() {
		subject = new(internal.FeatureStats_Categorical)

		// outlook: rainy=0 overcast=1 sunny=2
		// play: 		yes=0 no=1
		subject.Add(0, 1, 1.0) // rainy -> no
		subject.Add(0, 1, 1.0) // rainy -> no
		subject.Add(1, 0, 1.0) // overcast -> yes
		subject.Add(2, 0, 1.0) // sunny -> yes
		subject.Add(2, 0, 1.0) // sunny -> yes
		subject.Add(2, 1, 1.0) // sunny -> no
		subject.Add(1, 0, 1.0) // overcast -> yes
		subject.Add(0, 1, 1.0) // rainy -> no
		subject.Add(0, 0, 1.0) // rainy -> yes
		subject.Add(2, 0, 1.0) // sunny -> yes
		subject.Add(0, 0, 1.0) // rainy -> yes
		subject.Add(1, 0, 1.0) // overcast -> yes
		subject.Add(1, 0, 1.0) // overcast -> yes
		subject.Add(2, 1, 1.0) // sunny -> no
	})

	It("should add", func() {
		Expect(subject.Len()).To(Equal(3))
	})

	It("should calculate post-splits", func() {
		s := subject.PostSplit()
		Expect(s.Len()).To(Equal(3))
		Expect(s.Get(0).Sparse).To(Equal(map[int64]float64{0: 2, 1: 3}))
		Expect(s.Get(1).Sparse).To(Equal(map[int64]float64{0: 4}))
		Expect(s.Get(2).Sparse).To(Equal(map[int64]float64{0: 3, 1: 2}))
	})

})
