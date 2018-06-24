package internal_test

import (
	"github.com/bsm/reason/classification"
	"github.com/bsm/reason/classification/hoeffding/internal"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LeafNode", func() {
	model := testdata.ClassificationModel()
	examples := testdata.ClassificationData()

	var wrapper *internal.Node
	var subject *internal.LeafNode

	BeforeEach(func() {
		subject = new(internal.LeafNode)
		wrapper = internal.NewNode(&internal.Node_Leaf{Leaf: subject}, nil)

		target := model.Feature("play")
		for _, x := range examples {
			subject.Observe(model, target, x, 1.0, wrapper)
		}
	})

	It("should observe", func() {
		Expect(wrapper.Weight()).To(Equal(14.0))
		Expect(wrapper.Stats.Sparse).To(Equal(map[int64]float64{0: 9, 1: 5}))

		Expect(subject.FeatureStats).To(HaveLen(4))
		Expect(subject.FeatureStats).To(HaveKey("temp"))
		Expect(subject.FeatureStats["temp"].GetNumerical()).To(BeNil())
		Expect(subject.FeatureStats["temp"].GetCategorical().Len()).To(Equal(3))
		Expect(subject.WeightAtLastEval).To(Equal(0.0))
	})

	It("should evaluate splits", func() {
		crit := classification.DefaultSplitCriterion()
		Expect(subject.EvaluateSplit("unknown", crit, wrapper)).To(BeNil())

		cat := subject.EvaluateSplit("outlook", crit, wrapper)
		Expect(cat.Feature).To(Equal("outlook"))
		Expect(cat.Merit).To(BeNumerically("~", 0.247, 0.001))
		Expect(cat.Range).To(Equal(1.0))
		Expect(cat.Pivot).To(Equal(0.0))
		Expect(cat.PreSplit.Weight()).To(Equal(14.0))
		Expect(cat.PostSplit.Len()).To(Equal(3))
	})

	It("should allow to disable/enable", func() {
		Expect(subject.FeatureStats).To(HaveLen(4))
		Expect(subject.IsDisabled).To(BeFalse())
		subject.Enable()
		Expect(subject.FeatureStats).To(HaveLen(4))
		Expect(subject.IsDisabled).To(BeFalse())

		subject.Disable()
		Expect(subject.FeatureStats).To(BeNil())
		Expect(subject.IsDisabled).To(BeTrue())

		subject.Enable()
		Expect(subject.FeatureStats).NotTo(BeNil())
		Expect(subject.FeatureStats).To(HaveLen(0))
		Expect(subject.IsDisabled).To(BeFalse())
	})

})
