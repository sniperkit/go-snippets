package eval_test

import (
	"github.com/bsm/reason/classification/eval"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Accuracy", func() {
	var subject *eval.Accuracy

	BeforeEach(func() {
		subject = eval.NewAccuracy()
		subject.Record(1, 1, 1.0)
		subject.Record(1, 1, 1.0)
		subject.Record(1, 0, 1.0)
		subject.Record(0, 0, 1.0)
		subject.Record(0, 0, 1.0)
		subject.Record(0, 1, 1.0)
		subject.Record(1, 1, 1.0)
		subject.Record(1, 1, 1.0)
		subject.Record(1, 1, 2.0)
		subject.Record(0, 0, 1.0)
		subject.Record(0, 1, 1.0)
	})

	It("should calculate stats", func() {
		Expect(subject.Correct()).To(Equal(9.0))
		Expect(subject.Total()).To(Equal(12.0))
		Expect(subject.Accuracy()).To(Equal(0.75))
	})

})
