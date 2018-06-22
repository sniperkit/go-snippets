package tools

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/bblfsh/sdk.v1/uast"
)

func TestCyclomaticComplexity(t *testing.T) {
	require := require.New(t)
	n := &uast.Node{InternalType: "module",
		Children: []*uast.Node{
			{InternalType: "root"}, // 1 (initial)
			// Prefix is the default so it doesnt need any role
			{InternalType: "if1", Roles: []uast.Role{uast.Statement, uast.If}, Children: []*uast.Node{ // 2 (If, Statement)
				{InternalType: "if1else1", Roles: []uast.Role{uast.If, uast.Then}, Children: []*uast.Node{ // 0
					{InternalType: "if1else1foreach", Roles: []uast.Role{uast.Statement, uast.For, uast.Iterator}, Children: []*uast.Node{ // 3 (For, Statement)
						{InternalType: "foreach_child1"},                                                             // 0
						{InternalType: "foreach_child2_continue", Roles: []uast.Role{uast.Statement, uast.Continue}}, // 4 (Statement, Continue)
					}},
					{InternalType: "if1else1if", Roles: []uast.Role{uast.Statement, uast.If}, Children: []*uast.Node{ // 5 (Statement, If)
						{InternalType: "elseif_child1"},                                                    // 0
						{InternalType: "opAnd", Roles: []uast.Role{uast.Operator, uast.Boolean, uast.And}}, // 6 (Operator, Boolean)
						{InternalType: "elseif_child2"},                                                    // 0
					}},
				}},
				{InternalType: "break", Roles: []uast.Role{uast.Statement, uast.Break}},
			},
			}}}
	require.Equal(cyclomaticComplexity(n), 6)
}
