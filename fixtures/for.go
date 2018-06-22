package fixtures

func fors() {
	for i := 0; i < 5; i++ {
		continue
	}
	for {
		break
	}
	for true {
		break
	}
	for i := range c {
		_ = i
	}
	for _, v := range c {
		_ = v
	}
	for range c {
		// body
	}
}
