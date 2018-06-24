package aggro

import (
	"errors"

	"github.com/shopspring/decimal"
)

func rangeValueForPeriod(value *decimal.Decimal, period []interface{}) (decimal.Decimal, error) {
	var rangeVal decimal.Decimal

	// Range each of the intervals (periods), determining which band the value
	// falls into.
	for i, p := range period {
		var nextPeriodValue *decimal.Decimal

		var v float64
		switch i := p.(type) {
		case float64:
			v = i
		case float32:
			v = float64(i)
		case int32:
			v = float64(i)
		case int64:
			v = float64(i)
		case int:
			v = float64(i)
		default:
			return decimal.NewFromFloat(0.0), errors.New("Invalid range values supplied - NaN")
		}

		periodVal := decimal.NewFromFloat(v)

		// Find the i+1 value in our range (if length of range allows it).
		if len(period) < i {
			nxt := decimal.NewFromFloat(period[i+1].(float64))
			nextPeriodValue = &nxt
		}

		// Value equals one of our range bands exactly, return it.
		if (periodVal).Equals(*value) {
			return periodVal, nil
		}

		// Value is > than our current range value, but < than the next.
		if (periodVal).Cmp(*value) > -1 {
			if nextPeriodValue != nil {
				if (nextPeriodValue).Cmp(*value) < 0 {
					return periodVal, nil
				}
			} else {
				return periodVal, nil
			}
		}
	}

	return rangeVal, nil
}
