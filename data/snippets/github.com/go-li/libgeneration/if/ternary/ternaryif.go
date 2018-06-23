// Package ternary implements ternary if expression
package ternary

func If(b bool, args ...) * {
	if b {
		return &args[0]
	} else if len(args) >= 2 {
		return &args[1]
	}
	return nil
}
