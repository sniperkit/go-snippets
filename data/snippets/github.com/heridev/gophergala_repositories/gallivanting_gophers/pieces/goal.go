package pieces

// GlCarrot and friends are the goals generated for gophers to move towards.
const (
	GlCarrot = iota
	GlRadish
	GlPotatoes
	GlCabbage
	GlMushroom
	GlRaddish2
)

// Goal has a type (only one of each type is generated at a time) and a location.
type Goal struct {
	Type     byte
	Location int
}

// GoalFromInt ...
func GoalFromInt(i int) byte {
	switch i {
	case 1:
		return GlRadish
	case 2:
		return GlPotatoes
	case 3:
		return GlCabbage
	case 4:
		return GlMushroom
	case 5:
		return GlRaddish2
	default:
		return GlCarrot
	}
}
