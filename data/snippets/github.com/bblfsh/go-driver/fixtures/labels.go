package fixtures

func labels() {
loop:
	for {
		for {
			break loop
		}
	}

start:
	if false {
		goto start
	}
}
