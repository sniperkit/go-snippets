package regression_test

import (
	"github.com/bsm/reason/regression"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Evaluator", func() {
	var subject *regression.Evaluator

	BeforeEach(func() {
		subject = regression.NewEvaluator()
		subject.Record(25, 26, 1.0)
		subject.Record(25, 20, 1.0)
		subject.Record(22, 24, 1.0)
		subject.Record(23, 21, 1.0)
		subject.Record(24, 23, 1.0)
		subject.Record(29, 25, 1.0)
		subject.Record(28, 27, 1.0)
		subject.Record(26, 28, 2.0)
		subject.Record(30, 29, 1.0)
		subject.Record(18, 22, 1.0)
	})

	It("should calculate stats", func() {
		Expect(subject.Total()).To(Equal(11.0))
		Expect(subject.Mean()).To(BeNumerically("~", 24.8, 0.1))
		Expect(subject.MAE()).To(BeNumerically("~", 2.27, 0.01))
		Expect(subject.MSE()).To(BeNumerically("~", 7.00, 0.01))
		Expect(subject.RMSE()).To(BeNumerically("~", 2.64, 0.01))
		Expect(subject.R2()).To(BeNumerically("~", 0.39, 0.01))

		subject.Record(28, 28, 2.0)
		Expect(subject.R2()).To(BeNumerically("~", 0.47, 0.01))
	})

})
