package hoeffding

import (
	"fmt"

	"github.com/bsm/reason/core"
)

const numPivotBuckets = 11

// PivotPoints determines the optimum split points between min and max
// for a given number of buckets.
func PivotPoints(min, max float64) []float64 {
	inc := (max - min) / float64(numPivotBuckets+1)
	if inc <= 0 {
		return nil
	}

	pp := make([]float64, numPivotBuckets)
	for i := 0; i < numPivotBuckets; i++ {
		pp[i] = min + inc*float64(i+1)
	}
	return pp
}

// FormatNodeCondition returns the node condition description.
func FormatNodeCondition(feat *core.Feature, pos int, pivot float64) string {
	if feat.Kind.IsNumerical() {
		if pos == 0 {
			return fmt.Sprintf(`%s <= %.2f`, feat.Name, pivot)
		} else {
			return fmt.Sprintf(`%s > %.2f`, feat.Name, pivot)
		}
	}

	cat := core.Category(pos)
	return fmt.Sprintf(`%s = %s`, feat.Name, feat.ValueOf(cat))
}
