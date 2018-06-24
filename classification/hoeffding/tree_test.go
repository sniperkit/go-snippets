package hoeffding_test

import (
	"bytes"

	"github.com/bsm/reason/classification/eval"
	"github.com/bsm/reason/classification/hoeffding"
	common "github.com/bsm/reason/common/hoeffding"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree", func() {

	var train = func(n int) (*hoeffding.Tree, *core.Model, []core.Example) {
		stream, model, err := testdata.OpenClassification("../../testdata")
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
		Expect(t1.Info()).To(Equal(&common.TreeInfo{NumNodes: 11, NumLearning: 9, MaxDepth: 3}))
		Expect(t1.Predict(nil, examples[4001]).Best().P(0)).To(BeNumerically("~", 0.273, 0.001))

		b1 := new(bytes.Buffer)
		Expect(t1.WriteTo(b1)).To(Equal(int64(b1.Len())))

		t2, err := hoeffding.Load(b1, c)
		Expect(err).NotTo(HaveOccurred())
		Expect(t2.Info()).To(Equal(&common.TreeInfo{NumNodes: 11, NumLearning: 9, MaxDepth: 3}))
		Expect(t2.Predict(nil, examples[4001]).Best().P(0)).To(BeNumerically("~", 0.273, 0.001))
	})

	It("should prune", func() {
		t, _, _ := train(3000)
		Expect(t.Info()).To(Equal(&common.TreeInfo{
			NumNodes:    11,
			NumLearning: 9,
			NumDisabled: 0,
			MaxDepth:    3,
		}))

		t.Prune(5)
		Expect(t.Info()).To(Equal(&common.TreeInfo{
			NumNodes:    11,
			NumLearning: 5,
			NumDisabled: 4,
			MaxDepth:    3,
		}))
	})

	It("should write TXT", func() {
		t, _, _ := train(3000)

		b := new(bytes.Buffer)
		Expect(t.WriteText(b)).To(Equal(int64(b.Len())))

		s := b.String()
		Expect(s).To(ContainSubstring(`ROOT [weight:600]`))
		Expect(s).To(ContainSubstring("\tc5 = v1 [weight:644]"))
	})

	It("should write DOT", func() {
		t, _, _ := train(3000)

		b := new(bytes.Buffer)
		Expect(t.WriteDOT(b)).To(Equal(int64(b.Len())))

		s := b.String()
		Expect(s).To(ContainSubstring(`N [label="weight: 600"];`))
		Expect(s).To(ContainSubstring(`N_0 [label="c5 = v1\nweight: 644"];`))
	})

	DescribeTable("should train & predict",
		func(n int, expInfo *common.TreeInfo, exp *testdata.ClassificationScore) {
			tree, model, examples := train(n)
			Expect(tree.Info()).To(Equal(expInfo))

			accuracy := eval.NewAccuracy()
			kappa := eval.NewKappa()
			logLoss := eval.NewLogLoss()

			for _, x := range examples[n:] {
				predicted, probability := tree.Predict(nil, x).Best().Top()
				actual := model.Feature("target").Category(x)

				accuracy.Record(predicted, actual, 1.0)
				kappa.Record(predicted, actual, 1.0)
				logLoss.Record(probability, 1.0)
			}

			Expect(accuracy.Accuracy() * 100).To(BeNumerically("~", exp.Accuracy, 0.1))
			Expect(kappa.Score()).To(BeNumerically("~", exp.Kappa, 0.001))
			Expect(logLoss.Value()).To(BeNumerically("~", exp.LogLoss, 0.001))
		},

		Entry("1,000", 1000, &common.TreeInfo{
			NumNodes:    6,
			NumLearning: 5,
			MaxDepth:    2,
		}, &testdata.ClassificationScore{
			Accuracy: 71.1,
			Kappa:    0.348,
			LogLoss:  0.349,
		}),
		Entry("10,000", 10000, &common.TreeInfo{
			NumNodes:    38,
			NumLearning: 30,
			MaxDepth:    4,
		}, &testdata.ClassificationScore{
			Accuracy: 80.3,
			Kappa:    0.594,
			LogLoss:  0.230,
		}),
		Entry("20,000", 20000, &common.TreeInfo{
			NumNodes:    65,
			NumLearning: 48,
			MaxDepth:    4,
		}, &testdata.ClassificationScore{
			Accuracy: 85.0,
			Kappa:    0.690,
			LogLoss:  0.183,
		}),
	)

})
