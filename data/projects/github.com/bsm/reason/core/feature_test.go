package core_test

import (
	"math"

	"github.com/bsm/reason/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Feature", func() {

	DescribeTable("Number",
		func(f *core.Feature, x core.Example, exp float64) {
			if math.IsNaN(exp) {
				Expect(core.IsNum(f.Number(x))).To(BeFalse())
			} else {
				Expect(f.Number(x)).To(Equal(exp))
			}
		},

		Entry("categorical",
			core.NewCategoricalFeature("cat", nil), core.MapExample{}, math.NaN()),
		Entry("numerical (no value)",
			core.NewNumericalFeature("num"), core.MapExample{}, math.NaN()),
		Entry("numerical (from float)",
			core.NewNumericalFeature("num"), core.MapExample{"num": 1.23}, 1.23),
		Entry("numerical (from int)",
			core.NewNumericalFeature("num"), core.MapExample{"num": 101}, 101.0),
		Entry("numerical (from string)",
			core.NewNumericalFeature("num"), core.MapExample{"num": "12.3"}, 12.3),
	)

	DescribeTable("Category",
		func(f *core.Feature, x core.Example, exp core.Category) {
			Expect(f.Category(x)).To(Equal(exp))
		},

		Entry("numerical",
			core.NewNumericalFeature("num"), core.MapExample{}, core.NoCategory),

		Entry("categorical (no value)",
			core.NewCategoricalFeature("cat", []string{"x"}), core.MapExample{}, core.NoCategory),
		Entry("categorical (unknown value)",
			core.NewCategoricalFeature("cat", []string{"x"}), core.MapExample{"cat": "y"}, core.NoCategory),
		Entry("categorical (string value)",
			core.NewCategoricalFeature("cat", []string{"x"}), core.MapExample{"cat": "x"}, core.Category(0)),
		Entry("categorical (byte value)",
			core.NewCategoricalFeature("cat", []string{"x"}), core.MapExample{"cat": []byte("x")}, core.Category(0)),
		Entry("categorical (int value)",
			core.NewCategoricalFeature("cat", []string{"2"}), core.MapExample{"cat": 2}, core.Category(0)),

		Entry("categorical, no vocabulary, no hash buckets",
			core.NewCategoricalFeature("cat", nil),
			core.MapExample{"cat": "x"}, core.NoCategory),
		Entry("categorical, vocabulary, no hash buckets (known value)",
			core.NewCategoricalFeature("cat", []string{"a", "b"}),
			core.MapExample{"cat": "b"}, core.Category(1)),
		Entry("categorical, vocabulary, no hash buckets (unknown value)",
			core.NewCategoricalFeature("cat", []string{"a"}),
			core.MapExample{"cat": "x"}, core.NoCategory),

		Entry("categorical, hash buckets",
			core.NewCategoricalFeatureHashBuckets("cat", 10),
			core.MapExample{"cat": "z"}, core.Category(6)),

		Entry("categorical, vocabulary, hash buckets (known value)",
			core.NewCategoricalFeatureVocabularyHashBuckets("cat", []string{"a", "b"}, 10),
			core.MapExample{"cat": "a"}, core.Category(0)),
		Entry("categorical, vocabulary, hash buckets (unknown value)",
			core.NewCategoricalFeatureVocabularyHashBuckets("cat", []string{"a", "b"}, 10),
			core.MapExample{"cat": "z"}, core.Category(8)),

		Entry("categorical, vocabulary, expandable (known value)",
			core.NewCategoricalFeatureExpandable("cat", []string{"a", "b"}),
			core.MapExample{"cat": "b"}, core.Category(1)),
		Entry("categorical, vocabulary, expandable (unknown value)",
			core.NewCategoricalFeatureExpandable("cat", []string{"a", "b"}),
			core.MapExample{"cat": "x"}, core.Category(2)),

		Entry("categorical, identity (no value)",
			core.NewCategoricalFeatureIdentity("cat"), core.MapExample{}, core.NoCategory),
		Entry("categorical, identity (bad value)",
			core.NewCategoricalFeatureIdentity("cat"), core.MapExample{"cat": "x"}, core.NoCategory),
		Entry("categorical, identity (string value)",
			core.NewCategoricalFeatureIdentity("cat"), core.MapExample{"cat": "2"}, core.Category(2)),
		Entry("categorical, identity (byte value)",
			core.NewCategoricalFeatureIdentity("cat"), core.MapExample{"cat": []byte("1")}, core.Category(1)),
		Entry("categorical, identity (int value)",
			core.NewCategoricalFeatureIdentity("cat"), core.MapExample{"cat": 5}, core.Category(5)),
		Entry("categorical, identity (bool value)",
			core.NewCategoricalFeatureIdentity("cat"), core.MapExample{"cat": true}, core.Category(1)),
		Entry("categorical, identity (numeric value)",
			core.NewCategoricalFeatureIdentity("cat"), core.MapExample{"cat": customNumeric(8)}, core.Category(8)),
		Entry("categorical, identity (stringish value)",
			core.NewCategoricalFeatureIdentity("cat"), core.MapExample{"cat": customString("7")}, core.Category(7)),
		Entry("categorical, identity (int ptr)",
			core.NewCategoricalFeatureIdentity("cat"), core.MapExample{"cat": intPtr(6)}, core.Category(6)),
	)

})
