// Package fixtures is a collection of Go code samples.
package fixtures

// detached comment

/* detached comment 2 */

/*
	detached
	mutiline
*/

//go:generate echo directive

// var block comment
var (
	// a is a variable
	a int
	b int // inline comment
)

// foo is a sample function
func foo() {
	// function body comment
}

// data is a sample struct
//
// It's not very useful.
type data struct {
	a string // struct field comment
}
