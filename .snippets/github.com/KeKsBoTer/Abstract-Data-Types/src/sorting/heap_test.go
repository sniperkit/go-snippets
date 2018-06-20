package sorting

import (
	"testing"
	"sorting/test"
)

func TestSort(t *testing.T) {
	test.Run(50, Sort, t)
}
