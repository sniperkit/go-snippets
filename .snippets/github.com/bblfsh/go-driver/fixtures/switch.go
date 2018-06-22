package fixtures

func switches() {
	switch a {
	case 1:
		break
	case 2, 3:
		fallthrough
	default:
		// body
	}

	switch c := a; c {
	case a + 1:
	}

	switch {
	case a == b:
	default:
	}

	switch a := a.(type) {
	case int:
		// int
	case string:
		// string
	}

}
