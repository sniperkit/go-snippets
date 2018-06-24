package eval_test

import (
	"github.com/bsm/reason/classification/eval"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogLoss", func() {

	It("should calculate value", func() {
		e1 := eval.NewLogLoss()
		e1.Record(0.8, 1)
		e1.Record(0.9, 1)
		e1.Record(0.1, 1)
		e1.Record(0.6, 1)
		Expect(e1.Value()).To(BeNumerically("~", 0.785, 0.001))

		e2 := eval.NewLogLoss()
		e2.Record(0.4, 1)
		e2.Record(0.8, 1)
		e2.Record(0.7, 1)
		e2.Record(0.15, 1)
		Expect(e2.Value()).To(BeNumerically("~", 0.848, 0.001))
		e2.Record(1.0, 1)
		Expect(e2.Value()).To(BeNumerically("~", 0.679, 0.001))
		e2.Record(0.0, 1)
		Expect(e2.Value()).To(BeNumerically("~", 6.322, 0.001))

		e3 := eval.NewLogLoss()
		e3.Record(0.0, 1)
		Expect(e3.Value()).To(BeNumerically("~", 34.5, 0.1))
	})

})
