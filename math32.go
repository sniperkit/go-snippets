package sticker

import "math"

// Abs32 is the float32-version of math.Abs.
func Abs32(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

// Ceil32 is the float32-version of math.Ceil.
func Ceil32(x float32) float32 {
	return float32(math.Ceil(float64(x)))
}

// Exp32 is the float32-version of math.Exp.
func Exp32(x float32) float32 {
	return float32(math.Exp(float64(x)))
}

// Floor32 is the float32-version of math.Floor.
func Floor32(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

// Inf32 is the float32-version of math.Inf.
func Inf32(sign int) float32 {
	return float32(math.Inf(sign))
}

// IsInf32 is the float32-version of math.Inf.
func IsInf32(f float32, sign int) bool {
	return math.IsInf(float64(f), sign)
}

// IsNaN32 is the float32-version of math.IsNaN.
func IsNaN32(f float32) (is bool) {
	return math.IsNaN(float64(f))
}

// Log32 is the float32-version of math.Log.
func Log32(x float32) float32 {
	return float32(math.Log(float64(x)))
}

// LogBinary32 is the float32-version of math.Log2.
func LogBinary32(x float32) float32 {
	return float32(math.Log2(float64(x)))
}

// Modf32 is the float32-version of math.Modf.
func Modf32(x float32) (i, f float32) {
	i64, f64 := math.Modf(float64(x))
	return float32(i64), float32(f64)
}

// NaN32 is the float32-version of math.NaN.
func NaN32() float32 {
	return float32(math.NaN())
}

// Pow32 is the float32-version of math.Pow.
func Pow32(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

// Sqrt32 is the float32-version of math.Sqrt32.
func Sqrt32(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}
