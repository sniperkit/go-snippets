// Package FTRL performs Follow-The-Regularized-Leader adaptive learning.
package ftrl

import (
	"math"
	"sort"

	"github.com/bsm/reason/core"
)

func featureBV(feat *core.Feature, x core.Example, offset int) (int, float64) {
	switch feat.Kind {
	case core.Feature_CATEGORICAL:
		if cat := feat.Category(x); cat != core.NoCategory {
			return offset + int(cat), 1.0
		}
	case core.Feature_NUMERICAL:
		if num := feat.Number(x); !math.IsNaN(num) {
			return offset, num
		}
	}
	return -1, 0.0
}

func parseFeatures(features map[string]*core.Feature, target string) (predictors []string, offsets []int, size int) {
	predictors = make([]string, 0, len(features)-1)
	for _, feat := range features {
		if feat.Name != target {
			predictors = append(predictors, feat.Name)
		}
	}
	sort.Strings(predictors)

	offsets = make([]int, len(predictors))
	for i, name := range predictors {
		offsets[i] = size

		feat := features[name]
		switch feat.Kind {
		case core.Feature_CATEGORICAL:
			size += feat.NumCategories()
		case core.Feature_NUMERICAL:
			size += 1
		}
	}
	return
}
