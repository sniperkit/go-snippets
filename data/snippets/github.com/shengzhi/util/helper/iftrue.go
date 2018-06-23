package helper

// IfTrueStr 二元运算
func IfTrueStr(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}

// IfTrueInt  二元运算
func IfTrueInt(b bool, t, f int) int {
	if b {
		return t
	}
	return f
}

// IfTrueInt64  二元运算
func IfTrueInt64(b bool, t, f int64) int64 {
	if b {
		return t
	}
	return f
}

// IfTrueFloat  二元运算
func IfTrueFloat(b bool, t, f float64) float64 {
	if b {
		return t
	}
	return f
}

// IfTrue  二元运算
func IfTrue(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}
