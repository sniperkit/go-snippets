package internal_test

import (
	"bytes"
	"testing"

	"github.com/bsm/reason/classification/ftrl/internal"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Optimizer", func() {
	var subject *internal.Optimizer

	model := testdata.RegressionModel()

	BeforeEach(func() {
		subject = internal.NewOptimizer(model, "hours", 10)
	})

	It("should init", func() {
		Expect(subject).To(Equal(&internal.Optimizer{
			Model:   model,
			Target:  "hours",
			Sums:    make([]float64, 10),
			Weights: make([]float64, 10),
		}))
	})

	It("should write and read", func() {
		buf := new(bytes.Buffer)
		Expect(subject.WriteTo(buf)).To(Equal(int64(332)))

		dup := new(internal.Optimizer)
		Expect(dup.ReadFrom(buf)).To(Equal(int64(332)))
		Expect(dup).To(Equal(subject))
	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "classification/ftrl/internal")
}
