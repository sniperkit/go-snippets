package composer

import "fmt"
import "regexp/syntax"
import "strings"

// A regexp sub-expression
type RxElement string

func rxf(format string, a ...interface{}) RxElement {
	return RxElement(fmt.Sprintf(format, a...))
}

// Wrap rx in a non-capturing group
func Group(elt RxElement) RxElement {
	return rxf("(?:%s)", elt)
}

// Wrap rx in a named capturing group
func Capture(name string, elt RxElement) RxElement {
	return rxf("(?P<%s>%s)", name, elt)
}

// Convert any value to RxElement
func Element(v interface{}) RxElement {
	switch v.(type) {
	case string:
		return Group(RxElement(v.(string)))
	case fmt.Stringer:
		return Group(RxElement(v.(fmt.Stringer).String()))
	default:
		return Group(rxf("%v", v))
	}
}

// Wrap all elts in capturing groups, join with glue (only "|" or ""
// seem to make much sense), wrap whole thing with another
// non-capturing group
func Join(glue string, elts []RxElement) RxElement {
	wrapped := make([]string, len(elts))
	for i, elt := range elts {
		wrapped[i] = string(Group(elt))
	}
	return rxf("(?:%s)", strings.Join(wrapped, glue))
}

// Regexp alternation
func Alternation(elts ...RxElement) RxElement {
	return Group(Join("|", elts))
}

// Sequence of regexps
func Sequence(elts ...RxElement) RxElement {
	return Group(Join("", elts))
}

// Return a sequence, anchored at the beginning and end
func Anchor(elts ...RxElement) RxElement {
	return Sequence(Beginning, Sequence(elts...), End)
}

// Return a sequence, anchored at the beginning and end, trimming
// (i.e. matching and not capturing) leading and trailing whitespace
func TrimAnchor(elts ...RxElement) RxElement {
	return Sequence(Beginning, AnyWhitespace, Sequence(elts...), AnyWhitespace, End)
}

// Modify (with "*", "*?", "+", "?" etc)
func Mod(elt RxElement, modifier string) RxElement {
	return Group(elt) + RxElement(modifier)
}

// Add "?" modifier to sequence of elts
func Optional(elts ...RxElement) RxElement {
	return Mod(Sequence(elts...), "?")
}

// Add "*" modifier to sequence of elts
func Any(elts ...RxElement) RxElement {
	return Mod(Sequence(elts...), "*")
}

// Add "+" modifier to sequence of elts
func Some(elts ...RxElement) RxElement {
	return Mod(Sequence(elts...), "+")
}

// Quotes str as literal substring
func Literal(str string) RxElement {
	rx, err := syntax.Parse(str, syntax.Literal)
	if err != nil {
		panic(err)
	}
	return RxElement(rx.Simplify().String())
}

// Just some elements, this should be more complete
var (
	HexDigit      = RxElement(`[0-9a-f]`)
	HexNumber     = Some(HexDigit)
	DecimalDigit  = RxElement(`\d`)
	DecimalNumber = Some(DecimalDigit)
	WordChar      = RxElement(`\w`)
	Word          = Some(WordChar)
	Base64Char    = RxElement(`[./a-zA-Z0-9_-]`)
	Base64        = Sequence(Some(Base64Char), Any(`=`))
	Beginning     = RxElement(`^`)
	End           = RxElement(`$`)
	AnyWhitespace = RxElement(`\s*`)
)
