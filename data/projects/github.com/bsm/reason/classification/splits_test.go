package classification_test

import (
	"github.com/bsm/reason/classification"
	"github.com/bsm/reason/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SplitCriterion", func() {
	var pre *util.Vector
	var post1, post2, post3 *util.VectorDistribution

	BeforeEach(func() {
		pre = util.NewVectorFromSlice(9.0, 6.0)

		post1 = new(util.VectorDistribution)
		post1.Add(0, 0, 3.0)
		post1.Add(0, 1, 2.0)
		post1.Add(1, 0, 4.0)
		post1.Add(2, 0, 2.0)
		post1.Add(2, 1, 4.0)

		post2 = new(util.VectorDistribution)
		post2.Add(0, 0, 1.0)
		post2.Add(1, 0, 2.0)
		post2.Add(2, 0, 1.0)
		post2.Add(3, 0, 1.0)
		post2.Add(4, 0, 1.0)
		post2.Add(5, 0, 2.0)
		post2.Add(5, 1, 1.0)
		post2.Add(6, 0, 1.0)
		post2.Add(6, 1, 1.0)
		post2.Add(7, 1, 2.0)
		post2.Add(8, 1, 1.0)
		post2.Add(9, 1, 1.0)

		post3 = new(util.VectorDistribution)
		post3.Add(0, 0, 9.0)
		post3.Add(0, 1, 6.0)
	})

	Describe("GiniImpurity", func() {
		var subject = classification.GiniImpurity{}

		It("should have range", func() {
			Expect(subject.Range(nil)).To(Equal(1.0))
			Expect(subject.Range(pre)).To(Equal(1.0))
		})

		It("should evaluate split", func() {
			Expect(subject.Merit(nil, nil)).To(Equal(0.0))
			Expect(subject.Merit(pre, post1)).To(BeNumerically("~", 0.338, 0.001))
			Expect(subject.Merit(pre, post2)).To(BeNumerically("~", 0.155, 0.001))
			Expect(subject.Merit(pre, post3)).To(BeNumerically("~", 0.480, 0.001))
		})
	})

	Describe("InformationGain", func() {
		var subject = classification.InformationGain{MinBranchFraction: 0.1}

		It("should have range", func() {
			Expect(subject.Range(nil)).To(Equal(1.0))
			Expect(subject.Range(pre)).To(Equal(1.0))
			Expect(subject.Range(util.NewVectorFromSlice(1, 2, 3))).To(BeNumerically("~", 1.58, 0.01))
		})

		It("should evaluate split", func() {
			Expect(subject.Merit(nil, nil)).To(Equal(0.0))
			Expect(subject.Merit(pre, post1)).To(BeNumerically("~", 0.280, 0.001))
			Expect(subject.Merit(pre, post2)).To(BeNumerically("~", 0.654, 0.001))
			Expect(subject.Merit(pre, post3)).To(Equal(0.0))

		})

		It("should calculate merit with fraction limit", func() {
			x := classification.InformationGain{MinBranchFraction: 0.3}
			Expect(x.Merit(pre, post1)).To(BeNumerically("~", 0.280, 0.001))

			x = classification.InformationGain{MinBranchFraction: 0.35}
			Expect(x.Merit(pre, post1)).To(Equal(0.0))
		})
	})

	Describe("GainRatio", func() {
		var (
			base    = classification.DefaultSplitCriterion()
			subject = classification.GainRatio{SplitCriterion: base}
		)

		It("should reduce merit of 'super-attributes'", func() {
			Expect(base.Merit(pre, post1)).To(BeNumerically("~", 0.28, 0.01))
			Expect(subject.Merit(pre, post1)).To(BeNumerically("~", 0.18, 0.01))

			Expect(base.Merit(pre, post2)).To(BeNumerically("~", 0.65, 0.01))
			Expect(subject.Merit(pre, post2)).To(BeNumerically("~", 0.21, 0.01))
		})

	})
})
