package core

import "math"

// NoCategory indicates a bad or invalid category.
const NoCategory Category = -1

// Category is the index of a categorical value.
type Category int

// IsCat returns true if cat is a valid category.
func IsCat(cat Category) bool { return cat > NoCategory }

// IsNum returns true if num is a valid number.
func IsNum(num float64) bool { return !math.IsNaN(num) }
