package pieces

// GoBlue
const (
	GoBlue = iota
	GoRed
	GoGreen
	GoPurple
	GoYellow
)

// Gopher ...
type Gopher struct {
	Type     byte
	Location int
}

// GopherFromInt ...
func GopherFromInt(i int) byte {
	switch i {
	case 1:
		return GoRed
	case 2:
		return GoGreen
	case 3:
		return GoPurple
	case 4:
		return GoYellow
	default:
		return GoBlue
	}
}
