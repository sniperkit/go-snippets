package eval_test

import (
	"github.com/bsm/reason/classification/eval"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kappa", func() {

	It("should calculate score", func() {
		k1 := eval.NewKappa()
		k1.Record(0, 0, 22)
		k1.Record(0, 1, 7)
		k1.Record(1, 0, 9)
		k1.Record(1, 1, 13)
		Expect(k1.Score()).To(BeNumerically("~", 0.353, 0.001))

		k2 := eval.NewKappa()
		k2.Record(1, 3, 1.1)
		Expect(k2.Score()).To(BeNumerically("~", 0.0, 0.001))

		k3 := eval.NewKappa()
		k3.Record(1, 1, 1.0)
		k3.Record(1, 1, 1.0)
		k3.Record(1, 0, 1.0)
		k3.Record(0, 0, 1.0)
		k3.Record(0, 0, 1.0)
		k3.Record(0, 1, 1.0)
		k3.Record(1, 1, 1.0)
		k3.Record(1, 1, 1.0)
		k3.Record(1, 1, 2.0)
		k3.Record(0, 0, 1.0)
		k3.Record(0, 1, 1.0)
		Expect(k3.Score()).To(BeNumerically("~", 0.471, 0.001))

		k4 := eval.NewKappa()
		k4.Record(0, 0, 1.0)
		k4.Record(0, 0, 1.0)
		Expect(k4.Score()).To(BeNumerically("~", 1.0, 0.001))
	})

})
