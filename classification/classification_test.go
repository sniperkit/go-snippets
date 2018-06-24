package classification_test

import (
	"testing"

	"github.com/bsm/reason/classification"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Prediction", func() {
	var subject *classification.Prediction

	BeforeEach(func() {
		v := new(util.Vector)
		v.Add(0, 12.0)
		v.Add(2, 72.0)
		v.Add(5, 36.0)

		subject = &classification.Prediction{Vector: *v}
	})

	It("should calculate top", func() {
		cat, prob := subject.Top()
		Expect(cat).To(Equal(core.Category(2)))
		Expect(prob).To(Equal(0.6))

		cat, weight := subject.TopW()
		Expect(cat).To(Equal(core.Category(2)))
		Expect(weight).To(Equal(72.0))
	})

	It("should calculate weight", func() {
		Expect(subject.W(0)).To(Equal(12.0))
		Expect(subject.W(1)).To(Equal(0.0))
		Expect(subject.W(2)).To(Equal(72.0))
		Expect(subject.W(5)).To(Equal(36.0))
	})

	It("should calculate probability", func() {
		Expect(subject.P(0)).To(Equal(0.1))
		Expect(subject.P(1)).To(Equal(0.0))
		Expect(subject.P(2)).To(Equal(0.6))
		Expect(subject.P(5)).To(Equal(0.3))
	})

})

// --------------------------------------------------------------------

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "classification")
}
