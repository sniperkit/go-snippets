package lazy_test

import (
	"testing"

	"github.com/cristaloleg/funf/lazy"
)

func Test(t *testing.T) {
	fn := func(a int) int { return a * 10 }
	value := lazy.New(fn, 20)
	if value == nil {
		t.Error("cannot instantiate")
	}

	if got := value.Value().(int); got != 200 {
		t.Error("expected 200")
	}
}
