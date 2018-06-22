package normalizer

import (
	"fmt"

	"gopkg.in/bblfsh/sdk.v1/uast"
)

// ToNode is an instance of `uast.ObjectToNode`, defining how to transform an
// into a UAST (`uast.Node`).
//
// https://godoc.org/gopkg.in/bblfsh/sdk.v1/uast#ObjectToNode
var ToNode = &uast.ObjectToNode{
	TopLevelIsRootNode: true,
	InternalTypeKey:    "kind",
	OffsetKey:          "pos",
	EndOffsetKey:       "end",
	LineKey:            "line",
	ColumnKey:          "col",
	TokenKeys: map[string]bool{
		"text": true,
	},
	IsNode: func(v map[string]interface{}) bool {
		_, ok := v["kind"].(string)
		return ok
	},
	// Current sdk doesn't accept arrays in AST
	// https://github.com/bblfsh/sdk/pull/213
	Modifier: func(n map[string]interface{}) error {
		if flags, ok := n["flags"].([]interface{}); ok {
			var newFlags map[string]bool
			if len(flags) > 0 {
				newFlags = make(map[string]bool, len(flags))
			}
			for _, f := range flags {
				flagStr, ok := f.(string)
				if !ok {
					return fmt.Errorf("flag %+v isn't string", f)
				}
				newFlags[flagStr] = true
			}
			n["flags"] = newFlags
		}
		return nil
	},
}
