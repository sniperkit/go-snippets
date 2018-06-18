func LazyInt(f func() int) func() int {
	var lv *int
	return func() {
		if lv == nil {
			lv = new(int)
			*lv = f()
		}
		return *lv
	}
}
