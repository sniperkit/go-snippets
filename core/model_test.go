package core_test

import (
	"github.com/bsm/reason/core"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Model", func() {

	model := core.NewModel(
		core.NewNumericalFeature("num"),
		core.NewCategoricalFeature("cat", []string{"a", "b"}),
	)

	It("should marshal/unmarshal", func() {
		bin, err := proto.Marshal(model)
		Expect(err).NotTo(HaveOccurred())

		other := new(core.Model)
		Expect(proto.Unmarshal(bin, other)).NotTo(HaveOccurred())
		Expect(other).To(Equal(model))
	})

})
