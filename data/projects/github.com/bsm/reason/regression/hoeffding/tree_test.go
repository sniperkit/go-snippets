package hoeffding_test

import (
	"bytes"

	common "github.com/bsm/reason/common/hoeffding"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/regression"
	"github.com/bsm/reason/regression/hoeffding"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree", func() {

	var train = func(n int) (*hoeffding.Tree, *core.Model, []core.Example) {
		stream, model, err := testdata.OpenRegression("../../testdata")
		Expect(err).NotTo(HaveOccurred())
		defer stream.Close()

		examples, err := stream.ReadN(n * 2)
		Expect(err).NotTo(HaveOccurred())

		tree, err := hoeffding.New(model, "target", nil)
		Expect(err).NotTo(HaveOccurred())

		for _, x := range examples[:n] {
			tree.Train(x, 1.0)
		}
		return tree, model, examples
	}

	It("should dump/load", func() {
		c := &hoeffding.Config{
			Config: common.Config{GracePeriod: 50},
		}

		t1, _, examples := train(3000)
		Expect(t1.Info()).To(Equal(&common.TreeInfo{NumNodes: 626, NumLearning: 625, MaxDepth: 2}))
		Expect(t1.Predict(nil, examples[4001]).Best().Mean()).To(BeNumerically("~", 0.260, 0.001))

		b1 := new(bytes.Buffer)
		Expect(t1.WriteTo(b1)).To(Equal(int64(b1.Len())))

		t2, err := hoeffding.Load(b1, c)
		Expect(err).NotTo(HaveOccurred())
		Expect(t2.Info()).To(Equal(&common.TreeInfo{NumNodes: 626, NumLearning: 625, MaxDepth: 2}))
		Expect(t2.Predict(nil, examples[4001]).Best().Mean()).To(BeNumerically("~", 0.260, 0.001))
	})

	It("should prune", func() {
		t, _, _ := train(3000)
		Expect(t.Info()).To(Equal(&common.TreeInfo{
			NumNodes:    626,
			NumLearning: 625,
			NumDisabled: 0,
			MaxDepth:    2,
		}))

		t.Prune(10)
		Expect(t.Info()).To(Equal(&common.TreeInfo{
			NumNodes:    626,
			NumLearning: 10,
			NumDisabled: 615,
			MaxDepth:    2,
		}))
	})

	It("should write TXT", func() {
		t, _, _ := train(3000)

		b := new(bytes.Buffer)
		Expect(t.WriteText(b)).To(Equal(int64(b.Len())))

		s := b.String()
		Expect(s).To(ContainSubstring(`ROOT [weight:3000 mean:0.5 variance:0.6]`))
		Expect(s).To(ContainSubstring("\tc1 = #4 [weight:4 mean:0.6 variance:0.0]"))
	})

	It("should write DOT", func() {
		t, _, _ := train(3000)

		b := new(bytes.Buffer)
		Expect(t.WriteDOT(b)).To(Equal(int64(b.Len())))

		s := b.String()
		Expect(s).To(ContainSubstring(`N [label="weight: 3000"];`))
		Expect(s).To(ContainSubstring(`N_4 [label="c1 = #4\nweight: 4"];`))
	})

	DescribeTable("should train & predict",
		func(n int, expInfo *common.TreeInfo, exp *testdata.RegressionScore) {
			tree, model, examples := train(n)
			Expect(tree.Info()).To(Equal(expInfo))

			eval := regression.NewEvaluator()
			for _, x := range examples[n:] {
				prediction := tree.Predict(nil, x).Best().Mean()
				actual := model.Feature("target").Number(x)
				eval.Record(prediction, actual, 1.0)
			}
			Expect(eval.R2()).To(BeNumerically("~", exp.R2, 0.001))
			Expect(eval.RMSE()).To(BeNumerically("~", exp.RMSE, 0.001))
		},

		Entry("1,000", 1000, &common.TreeInfo{
			NumNodes:    1,
			NumLearning: 1,
			MaxDepth:    1,
		}, &testdata.RegressionScore{
			R2:   0.002,
			RMSE: 0.854,
		}),
		Entry("5,000", 5000, &common.TreeInfo{
			NumNodes:    764,
			NumLearning: 763,
			MaxDepth:    2,
		}, &testdata.RegressionScore{
			R2:   0.186,
			RMSE: 0.960,
		}),
		Entry("10,000", 10000, &common.TreeInfo{
			NumNodes:    905,
			NumLearning: 903,
			NumDisabled: 0,
			MaxDepth:    3,
		}, &testdata.RegressionScore{
			R2:   0.618,
			RMSE: 0.616,
		}),
	)

})
