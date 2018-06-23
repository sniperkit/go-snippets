package fixtures

func casts() {
	a := int(float64(5))
	var c interface{} = a
	_ = c.(int)
	_, ok := c.(int)
	_ = ok
}
