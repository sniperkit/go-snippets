package internal_test

import (
	"github.com/bsm/reason/classification/hoeffding/internal"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SplitNode_Children", func() {
	var sparse, dense *internal.SplitNode_Children

	BeforeEach(func() {
		sparse = new(internal.SplitNode_Children)
		sparse.SetRef(33, int64(101))
		sparse.SetRef(72, int64(102))

		dense = &internal.SplitNode_Children{Dense: []int64{}}
		dense.SetRef(33, int64(101))
		dense.SetRef(72, int64(102))
	})

	It("should iterate (sparse)", func() {
		nodes := make(map[int]int64)
		sparse.ForEach(func(i int, nodeRef int64) bool {
			nodes[i] = nodeRef
			return true
		})
		Expect(nodes).To(HaveLen(2))
		Expect(nodes).To(HaveKeyWithValue(33, int64(101)))
		Expect(nodes).To(HaveKeyWithValue(72, int64(102)))
	})

	It("should iterate (dense)", func() {
		nodes := make(map[int]int64)
		dense.ForEach(func(i int, nodeRef int64) bool {
			nodes[i] = nodeRef
			return true
		})
		Expect(nodes).To(HaveLen(2))
		Expect(nodes).To(HaveKeyWithValue(33, int64(101)))
		Expect(nodes).To(HaveKeyWithValue(72, int64(102)))
	})

	It("should cancel iterate", func() {
		nodes := make(map[int]int64)
		sparse.ForEach(func(i int, nodeRef int64) bool {
			nodes[i] = nodeRef
			return false
		})
		Expect(nodes).To(HaveLen(1))
	})

	It("should have len", func() {
		Expect(sparse.Len()).To(Equal(2))
		Expect(dense.Len()).To(Equal(2))
	})

	It("should set refs", func() {
		Expect(sparse.Sparse).To(HaveLen(2))
		Expect(sparse.SparseCap).To(Equal(int64(73)))
		Expect(sparse.Dense).To(BeNil())

		Expect(dense.Dense).To(HaveLen(73))
		Expect(dense.Dense).To(HaveCap(146))
		Expect(dense.Sparse).To(BeNil())
	})

	It("should get ref by index", func() {
		Expect(sparse.GetRef(-1)).To(BeZero())
		Expect(sparse.GetRef(1000)).To(BeZero())
		Expect(sparse.GetRef(32)).To(BeZero())
		Expect(sparse.GetRef(33)).To(Equal(int64(101)))

		Expect(dense.GetRef(-1)).To(BeZero())
		Expect(dense.GetRef(1000)).To(BeZero())
		Expect(dense.GetRef(32)).To(BeZero())
		Expect(dense.GetRef(33)).To(Equal(int64(101)))
	})

	It("should convert to dense", func() {
		Expect(sparse.Dense).To(BeNil())
		Expect(sparse.Sparse).To(HaveLen(2))
		Expect(sparse.SparseCap).To(Equal(int64(73)))

		for i := 1999; i >= 1950; i-- {
			sparse.SetRef(i, int64(i)*100)
		}
		Expect(sparse.Len()).To(Equal(52))
		Expect(sparse.Dense).To(BeNil())
		Expect(len(sparse.Sparse)).To(Equal(52))
		Expect(sparse.SparseCap).To(Equal(int64(2000)))

		for i := 999; i >= 800; i-- {
			sparse.SetRef(i, int64(i)*100)
		}
		Expect(sparse.Len()).To(Equal(252))
		Expect(len(sparse.Dense)).To(Equal(2000))
		Expect(sparse.Sparse).To(BeNil())
		Expect(sparse.SparseCap).To(Equal(int64(0)))

		_, err := proto.Marshal(sparse)
		Expect(err).NotTo(HaveOccurred())
	})

})
