package hoeffding_test

import (
	"testing"

	"github.com/bsm/reason/regression/hoeffding"
	"github.com/bsm/reason/testdata"
)

func BenchmarkTree_Train(b *testing.B) {
	const N = 1000

	stream, model, err := testdata.OpenRegression("../../testdata")
	if err != nil {
		b.Fatal(err)
	}
	defer stream.Close()

	examples, err := stream.ReadN(N)
	if err != nil {
		b.Fatal(err)
	}

	tree, err := hoeffding.New(model, "target", nil)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		tree.Train(examples[i%N], 1.0)
	}
}

func BenchmarkTree_Train_parallel(b *testing.B) {
	const N = 1000

	stream, model, err := testdata.OpenRegression("../../testdata")
	if err != nil {
		b.Fatal(err)
	}
	defer stream.Close()

	examples, err := stream.ReadN(N)
	if err != nil {
		b.Fatal(err)
	}

	tree, err := hoeffding.New(model, "target", nil)
	if err != nil {
		b.Fatal(err)
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			tree.Train(examples[i%N], 1.0)
		}
	})
}
