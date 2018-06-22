package fixtures

func chans() {
	c1 := make(chan string)
	c2 := make(chan int, 1)
	select {
	case s := <-c1:
		_ = s
	case v, ok := <-c2:
		_, _ = v, ok
	case c2 <- len(c2):
	default:
		close(c2)
	}
	go func() {
		// routine
	}()
}
