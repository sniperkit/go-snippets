package sticker

import (
	"fmt"
	"testing"

	"github.com/hiro4bbh/go-assert"
)

func TestMath32Functions(t *testing.T) {
	goassert.New(t, float32(1.0)).Equal(Abs32(-1.0))
	goassert.New(t, float32(2.0)).Equal(Ceil32(1.5))
	goassert.New(t, float32(1.0)).Equal(Exp32(0.0))
	goassert.New(t, float32(1.0)).Equal(Floor32(1.5))
	goassert.New(t, true).Equal(IsInf32(Inf32(+1.0), +1))
	goassert.New(t, false).Equal(IsNaN32(-1.0))
	goassert.New(t, "+Inf(float32)").Equal(fmt.Sprintf("%g(%T)", Inf32(+1), Inf32(+1)))
	goassert.New(t, float32(0.0)).Equal(Log32(1.0))
	goassert.New(t, float32(2.0)).Equal(LogBinary32(4.0))
	goassert.New(t, float32(1.0), float32(0.5)).Equal(Modf32(1.5))
	goassert.New(t, "NaN(float32)").Equal(fmt.Sprintf("%g(%T)", NaN32(), NaN32()))
	goassert.New(t, float32(8.0)).Equal(Pow32(2.0, 3.0))
	goassert.New(t, float32(2.0)).Equal(Sqrt32(4.0))
}
