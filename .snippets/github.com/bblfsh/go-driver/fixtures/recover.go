package fixtures

func fail() {
	defer func() {
		if r := recover(); r != nil {
			_ = r
		}
	}()
	panic("stop!")
}
