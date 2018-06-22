package human

import "errors"

// ErrInvalidTagName indicates that no tag name was specified.
var ErrInvalidTagName = errors.New("invalid tag name")

// ErrListSymbolsEmpty indicates that no list symbols were provided.
var ErrListSymbolsEmpty = errors.New("no list symbols provided")
