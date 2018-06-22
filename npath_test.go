package tools

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/bblfsh/sdk.v1/protocol"
	"gopkg.in/bblfsh/sdk.v1/uast"
)

func TestCountChildrenOfRole(t *testing.T) {
	require := require.New(t)

	n1 := &uast.Node{InternalType: "module", Children: []*uast.Node{
		{InternalType: "Statement", Roles: []uast.Role{uast.Statement}},
		{InternalType: "Statement", Roles: []uast.Role{uast.Statement}},
		{InternalType: "If", Roles: []uast.Role{uast.Statement, uast.If}},
	}}
	n2 := &uast.Node{InternalType: "module", Children: []*uast.Node{
		{InternalType: "Statement", Roles: []uast.Role{uast.Statement}, Children: []*uast.Node{
			{InternalType: "Statement", Roles: []uast.Role{uast.Statement}, Children: []*uast.Node{
				{InternalType: "If", Roles: []uast.Role{uast.Statement, uast.If}},
				{InternalType: "Statemenet", Roles: []uast.Role{uast.Statement}},
			}},
		}},
	}}
	result := countChildrenOfRoles(n1, []uast.Role{uast.Statement}, nil)
	expect := 3
	require.Equal(expect, result)

	result = countChildrenOfRoles(n2, []uast.Role{uast.Statement}, nil)
	expect = 1
	require.Equal(expect, result)

	result = deepCountChildrenOfRoles(n1, []uast.Role{uast.Statement}, nil)
	expect = 3
	require.Equal(expect, result)

	result = deepCountChildrenOfRoles(n2, []uast.Role{uast.Statement}, nil)
	expect = 4
	require.Equal(expect, result)
}

func TestChildrenOfRole(t *testing.T) {
	require := require.New(t)

	n1 := &uast.Node{InternalType: "module", Children: []*uast.Node{
		{InternalType: "Statement", Roles: []uast.Role{uast.Statement}},
		{InternalType: "Statement", Roles: []uast.Role{uast.Statement}},
		{InternalType: "If", Roles: []uast.Role{uast.If}},
	}}
	n2 := &uast.Node{InternalType: "module", Children: []*uast.Node{
		{InternalType: "Statement", Roles: []uast.Role{uast.Statement}, Children: []*uast.Node{
			{InternalType: "Statement", Roles: []uast.Role{uast.Statement}, Children: []*uast.Node{
				{InternalType: "If", Roles: []uast.Role{uast.If}},
				{InternalType: "Statemenet", Roles: []uast.Role{uast.Statement}},
			}},
		}},
	}}

	result := childrenOfRoles(n1, []uast.Role{uast.Statement}, nil)
	expect := 2
	require.Equal(expect, len(result))

	result = childrenOfRoles(n2, []uast.Role{uast.Statement}, nil)
	expect = 1
	require.Equal(expect, len(result))

	result = deepChildrenOfRoles(n1, []uast.Role{uast.Statement}, nil)
	expect = 2
	require.Equal(expect, len(result))

	result = deepChildrenOfRoles(n2, []uast.Role{uast.Statement}, nil)
	expect = 3
	require.Equal(expect, len(result))
}

func TestContainsRole(t *testing.T) {
	require := require.New(t)
	n := &uast.Node{InternalType: "node", Roles: []uast.Role{uast.Statement, uast.If}}

	result := containsRoles(n, []uast.Role{uast.If}, nil)
	require.Equal(true, result)

	result = containsRoles(n, []uast.Role{uast.Switch}, nil)
	require.Equal(false, result)
}

func TestExpresionComplex(t *testing.T) {
	require := require.New(t)

	n := &uast.Node{InternalType: "ifCondition", Roles: []uast.Role{uast.If, uast.Condition}, Children: []*uast.Node{
		{InternalType: "bool_and", Roles: []uast.Role{uast.Operator, uast.Boolean, uast.And}},
		{InternalType: "bool_xor", Roles: []uast.Role{uast.Operator, uast.Boolean, uast.Xor}},
	}}
	n2 := &uast.Node{InternalType: "ifCondition", Roles: []uast.Role{uast.If, uast.Condition}, Children: []*uast.Node{
		{InternalType: "bool_and", Roles: []uast.Role{uast.Operator, uast.Boolean, uast.And}, Children: []*uast.Node{
			{InternalType: "bool_or", Roles: []uast.Role{uast.Operator, uast.Boolean, uast.Or}, Children: []*uast.Node{
				{InternalType: "bool_xor", Roles: []uast.Role{uast.Operator, uast.Boolean, uast.Xor}},
			}},
		}},
	}}

	result := expressionComp(n)
	expect := 2
	require.Equal(expect, result)

	result = expressionComp(n2)
	expect = 3
	require.Equal(expect, result)
}

func TestNPathComplexity(t *testing.T) {
	require := require.New(t)
	var result []int
	var expect []int

	andBool := &uast.Node{InternalType: "bool_and", Roles: []uast.Role{uast.Operator, uast.Boolean, uast.And}}
	orBool := &uast.Node{InternalType: "bool_or", Roles: []uast.Role{uast.Operator, uast.Boolean, uast.Or}}
	statement := &uast.Node{InternalType: "Statement", Roles: []uast.Role{uast.Statement}}

	n := &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		statement,
	}}

	npathData := NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 1)

	/*
			if(3conditions){
				Statement
				Statement
			}else if(3conditions){
				Statement
				Statement
		  }else{
				Statement
				Statement
		  } Npath = 7
	*/
	ifCondition := &uast.Node{InternalType: "Condition", Roles: []uast.Role{uast.If, uast.Condition}, Children: []*uast.Node{
		andBool,
		orBool,
	}}
	ifBody := &uast.Node{InternalType: "Body", Roles: []uast.Role{uast.If, uast.Then}, Children: []*uast.Node{
		statement,
		statement,
	}}
	elseIf := &uast.Node{InternalType: "elseIf", Roles: []uast.Role{uast.If, uast.Else}, Children: []*uast.Node{
		&uast.Node{InternalType: "If", Roles: []uast.Role{uast.Statement, uast.If}, Children: []*uast.Node{
			ifCondition,
			ifBody,
		}},
	}}
	ifElse := &uast.Node{InternalType: "else", Roles: []uast.Role{uast.If, uast.Else}, Children: []*uast.Node{
		ifBody,
	}}
	nIf := &uast.Node{InternalType: "if", Roles: []uast.Role{uast.Statement, uast.If}, Children: []*uast.Node{
		ifCondition,
		ifBody,
		elseIf,
		ifElse,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		nIf,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 7)

	/*
	  if(condition){
	    Statement
	    Statement
	  }Npath = 2
	*/
	nSimpleIf := &uast.Node{InternalType: "If", Roles: []uast.Role{uast.Statement, uast.If}, Children: []*uast.Node{
		{InternalType: "ifCondition", Roles: []uast.Role{uast.If, uast.Condition}, Children: []*uast.Node{}},
		ifBody,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		nSimpleIf,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 2)

	/*
		The same if structure of the example above
		but repeated three times, in sequencial way
		Npath = 343
	*/

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		nIf,
		nIf,
		nIf,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 343)

	/*
		if(3conditions){
			if(3conditions){
				if(3conditions){
					Statement
					Statemenet
				}else{
					Statement
					Statement
				}
			}
		} Npath = 10
	*/
	nestedIfBody := &uast.Node{InternalType: "body", Roles: []uast.Role{uast.If, uast.Then}, Children: []*uast.Node{
		{InternalType: "if2", Roles: []uast.Role{uast.Statement, uast.If}, Children: []*uast.Node{
			ifCondition,
			{InternalType: "body2", Roles: []uast.Role{uast.If, uast.Then}, Children: []*uast.Node{
				{InternalType: "if3", Roles: []uast.Role{uast.Statement, uast.If}, Children: []*uast.Node{
					ifCondition,
					ifBody,
					ifElse,
				}},
			}},
		}},
	}}
	nNestedIf := &uast.Node{InternalType: "if1", Roles: []uast.Role{uast.Statement, uast.If}, Children: []*uast.Node{
		ifCondition,
		nestedIfBody,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		nNestedIf,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 10)

	/*
		while(2condtions){
			Statement
			Statement
			Statement
		}else{
			Statement
			Statement
		} Npath = 3
	*/
	whileCondition := &uast.Node{InternalType: "WhileCondition", Roles: []uast.Role{uast.While, uast.Condition}, Children: []*uast.Node{
		andBool,
	}}
	whileBody := &uast.Node{InternalType: "WhileBody", Roles: []uast.Role{uast.While, uast.Body}, Children: []*uast.Node{
		statement,
		statement,
		statement,
	}}
	whileElse := &uast.Node{InternalType: "WhileElse", Roles: []uast.Role{uast.While, uast.Else}, Children: []*uast.Node{
		statement,
		statement,
	}}
	nWhile := &uast.Node{InternalType: "While", Roles: []uast.Role{uast.Statement, uast.While}, Children: []*uast.Node{
		whileCondition,
		whileBody,
		whileElse,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		nWhile,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 3)

	/*
		while(2conditions){
			while(2conditions){
				while(2conditions){
					Statement
					Statement
				}
			}
		} Npath = 7
	*/
	nestedWhileBody := &uast.Node{InternalType: "WhileBody1", Roles: []uast.Role{uast.While, uast.Body}, Children: []*uast.Node{
		{InternalType: "While2", Roles: []uast.Role{uast.Statement, uast.While}, Children: []*uast.Node{
			whileCondition,
			{InternalType: "WhileBody2", Roles: []uast.Role{uast.While, uast.Body}, Children: []*uast.Node{
				{InternalType: "While3", Roles: []uast.Role{uast.Statement, uast.While}, Children: []*uast.Node{
					whileCondition,
					whileBody,
				}},
			}},
		}},
	}}
	nNestedWhile := &uast.Node{InternalType: "While1", Roles: []uast.Role{uast.Statement, uast.While}, Children: []*uast.Node{
		whileCondition,
		nestedWhileBody,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		nNestedWhile,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 7)

	/*
			 for(init;2condition;update){
			 	Statement
				Statement
		 	 } Npath = 2
	*/
	forCondition := &uast.Node{InternalType: "forCondition", Roles: []uast.Role{uast.For, uast.Expression}, Children: []*uast.Node{
		orBool,
	}}
	forInit := &uast.Node{InternalType: "forInit", Roles: []uast.Role{uast.For, uast.Initialization}}
	forUpdate := &uast.Node{InternalType: "forUpdate", Roles: []uast.Role{uast.For, uast.Update}}
	forBody := &uast.Node{InternalType: "forBody", Roles: []uast.Role{uast.For, uast.Body}, Children: []*uast.Node{
		statement,
		statement,
	}}
	nFor := &uast.Node{InternalType: "for", Roles: []uast.Role{uast.Statement, uast.For}, Children: []*uast.Node{
		forInit,
		forCondition,
		forUpdate,
		forBody,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		nFor,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 2)

	/*
		for(init;2conditions;update){
			for(init;2conditions;update){
				for(init;2condtions;update){
					Statement
					Statement
				}
			}
		} Npath = 4
	*/
	nestedForBody := &uast.Node{InternalType: "forBody1", Roles: []uast.Role{uast.For, uast.Body}, Children: []*uast.Node{
		{InternalType: "for2", Roles: []uast.Role{uast.Statement, uast.For}, Children: []*uast.Node{
			forInit,
			forCondition,
			forUpdate,
			{InternalType: "forBody2", Roles: []uast.Role{uast.For, uast.Body}, Children: []*uast.Node{
				{InternalType: "for3", Roles: []uast.Role{uast.Statement, uast.For}, Children: []*uast.Node{
					forInit,
					forCondition,
					forUpdate,
					forBody,
				}},
			}},
		}},
	}}
	nNestedFor := &uast.Node{InternalType: "for1", Roles: []uast.Role{uast.Statement, uast.For}, Children: []*uast.Node{
		forInit,
		forCondition,
		forUpdate,
		nestedForBody,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		nNestedFor,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 4)

	/*
		do{
			Statement
			Statement
		}while(3conditions)
		Npath = 4
	*/
	doWhileCondition := &uast.Node{InternalType: "doWhileCondition", Roles: []uast.Role{uast.DoWhile, uast.Condition}, Children: []*uast.Node{
		orBool,
		orBool,
	}}
	doWhileBody := &uast.Node{InternalType: "doWhileBody", Roles: []uast.Role{uast.DoWhile, uast.Body}, Children: []*uast.Node{
		statement,
		statement,
	}}
	nDoWhile := &uast.Node{InternalType: "doWhile", Roles: []uast.Role{uast.Statement, uast.DoWhile}, Children: []*uast.Node{
		doWhileBody,
		doWhileCondition,
	}}
	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		nDoWhile,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 4)

	/*
		do{
			do{
				do{
					Statement
					Statement
				}while(3conditions)
			}while{3conditions}
		}while(3condtions)
		Npath = 10
	*/
	nestedDoWhileBody := &uast.Node{InternalType: "doWhileBody1", Roles: []uast.Role{uast.DoWhile, uast.Body}, Children: []*uast.Node{
		{InternalType: "doWhile2", Roles: []uast.Role{uast.Statement, uast.DoWhile}, Children: []*uast.Node{
			{InternalType: "doWhileBody2", Roles: []uast.Role{uast.DoWhile, uast.Body}, Children: []*uast.Node{
				{InternalType: "doWhile3", Roles: []uast.Role{uast.Statement, uast.DoWhile}, Children: []*uast.Node{
					doWhileBody,
					doWhileCondition,
				}},
			}},
			doWhileCondition,
		}},
	}}
	nNestedDoWhile := &uast.Node{InternalType: "doWhile1", Roles: []uast.Role{uast.Statement, uast.DoWhile}, Children: []*uast.Node{
		nestedDoWhileBody,
		doWhileCondition,
	}}
	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		nNestedDoWhile,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 10)

	/*
		switch(){
		case:
			Statement
			Statement
		case:
			Statement
			Statement
		default:
			Statement
			Statement
		} Npath = 3
	*/
	switchCondition := &uast.Node{InternalType: "switchCondition", Roles: []uast.Role{uast.Switch, uast.Case, uast.Condition}, Children: []*uast.Node{
		orBool,
		andBool,
	}}
	switchCaseBody := &uast.Node{InternalType: "switchCaseBody", Roles: []uast.Role{uast.Switch, uast.Case, uast.Body}, Children: []*uast.Node{
		statement,
		statement,
	}}
	switchCase := &uast.Node{InternalType: "switchCase", Roles: []uast.Role{uast.Statement, uast.Switch, uast.Case}, Children: []*uast.Node{
		switchCondition,
	}}
	defaultCase := &uast.Node{InternalType: "defaultCase", Roles: []uast.Role{uast.Switch, uast.Default}, Children: []*uast.Node{
		statement,
		statement,
	}}
	nSwitch := &uast.Node{InternalType: "switch", Roles: []uast.Role{uast.Statement, uast.Switch}, Children: []*uast.Node{
		switchCase,
		switchCaseBody,
		switchCase,
		switchCaseBody,
		defaultCase,
	}}
	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		nSwitch,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 3)

	/*
		switch(){
		case:
			Statement
			Statement
		case:
			Statement
			Statement
		default:
			switch(){
			case:
				Statement
				Statement
			case:
				Statement
				Statement
			default:
				Statement
				Statement
		} Npath = 9
	*/
	nestedDefaultCase := &uast.Node{InternalType: "defaultCase", Roles: []uast.Role{uast.Switch, uast.Default}, Children: []*uast.Node{
		nSwitch,
	}}
	nNestedSwitch := &uast.Node{InternalType: "switch", Roles: []uast.Role{uast.Statement, uast.Switch}, Children: []*uast.Node{
		switchCase,
		switchCaseBody,
		switchCase,
		switchCaseBody,
		nestedDefaultCase,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		nNestedSwitch,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 5)

	/*
		return
	*/
	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		{InternalType: "Return", Roles: []uast.Role{uast.Statement, uast.Return}},
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 1)

	/*
		statement
		statement
		return 3condition
	*/
	nReturn := &uast.Node{InternalType: "Return", Roles: []uast.Role{uast.Statement, uast.Return}, Children: []*uast.Node{
		orBool,
		andBool,
	}}
	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		statement,
		statement,
		nReturn,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 2)

	nForEach := &uast.Node{InternalType: "ForEach", Roles: []uast.Role{uast.Statement, uast.For, uast.Iterator}, Children: []*uast.Node{
		forInit,
		forCondition,
		forBody,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		nForEach,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 2)

	tryBody := &uast.Node{InternalType: "TryBody", Roles: []uast.Role{uast.Try, uast.Body}, Children: []*uast.Node{
		statement,
		statement,
	}}

	tryCatch := &uast.Node{InternalType: "TryCatch", Roles: []uast.Role{uast.Try, uast.Catch}, Children: []*uast.Node{
		statement,
		statement,
	}}

	nTry := &uast.Node{InternalType: "Try", Roles: []uast.Role{uast.Statement, uast.Try}, Children: []*uast.Node{
		tryBody,
		tryCatch,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		nTry,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 2)

	tryFinally := &uast.Node{InternalType: "TryFinally", Roles: []uast.Role{uast.Try, uast.Finally}, Children: []*uast.Node{
		nSimpleIf,
	}}
	nTryFinally := &uast.Node{InternalType: "Try", Roles: []uast.Role{uast.Statement, uast.Try}, Children: []*uast.Node{
		tryBody,
		tryCatch,
		tryCatch,
		tryCatch,
		tryFinally,
	}}

	n = &uast.Node{InternalType: "Function declaration body", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{
		nTryFinally,
	}}

	npathData = NPathComplexity(n)
	result = append(result, npathData[0].Complexity)
	expect = append(expect, 6)

	require.Equal(expect, result)
}

func TestNpathMultiFunc(t *testing.T) {
	require := require.New(t)
	var result []int
	expect := []int{7, 7, 7}
	andBool := &uast.Node{InternalType: "bool_and", Roles: []uast.Role{uast.Operator, uast.Boolean, uast.And}}
	orBool := &uast.Node{InternalType: "bool_or", Roles: []uast.Role{uast.Operator, uast.Boolean, uast.Or}}
	statement := &uast.Node{InternalType: "Statement", Roles: []uast.Role{uast.Statement}}

	ifCondition := &uast.Node{InternalType: "Condition", Roles: []uast.Role{uast.If, uast.Condition}, Children: []*uast.Node{
		andBool,
		orBool,
	}}
	ifBody := &uast.Node{InternalType: "Body", Roles: []uast.Role{uast.If, uast.Then}, Children: []*uast.Node{
		statement,
		statement,
	}}
	elseIf := &uast.Node{InternalType: "elseIf", Roles: []uast.Role{uast.If, uast.Else}, Children: []*uast.Node{
		&uast.Node{InternalType: "If", Roles: []uast.Role{uast.Statement, uast.If}, Children: []*uast.Node{
			ifCondition,
			ifBody,
		}},
	}}
	ifElse := &uast.Node{InternalType: "else", Roles: []uast.Role{uast.If, uast.Else}, Children: []*uast.Node{
		ifBody,
	}}
	nIf := &uast.Node{InternalType: "if", Roles: []uast.Role{uast.Statement, uast.If}, Children: []*uast.Node{
		ifCondition,
		ifBody,
		elseIf,
		ifElse,
	}}
	funcBody := &uast.Node{InternalType: "funcBody", Roles: []uast.Role{uast.Function, uast.Body}, Children: []*uast.Node{nIf}}

	func1 := &uast.Node{InternalType: "func1", Roles: []uast.Role{uast.Function, uast.Declaration}, Children: []*uast.Node{
		&uast.Node{InternalType: "funcName1", Roles: []uast.Role{uast.Function, uast.Name}, Children: []*uast.Node{}, Token: "Name1"},
		funcBody,
	}}
	func2 := &uast.Node{InternalType: "func2", Roles: []uast.Role{uast.Function, uast.Declaration}, Children: []*uast.Node{
		&uast.Node{InternalType: "funcName2", Roles: []uast.Role{uast.Function, uast.Name}, Children: []*uast.Node{}, Token: "Name2"},
		funcBody,
	}}
	func3 := &uast.Node{InternalType: "func3", Roles: []uast.Role{uast.Function, uast.Declaration}, Children: []*uast.Node{
		&uast.Node{InternalType: "funcName3", Roles: []uast.Role{uast.Function, uast.Name}, Children: []*uast.Node{}, Token: "Name3"},
		funcBody,
	}}

	n := &uast.Node{InternalType: "module", Children: []*uast.Node{
		func1,
		func2,
		func3,
	}}
	npathData := NPathComplexity(n)
	for _, v := range npathData {
		result = append(result, v.Complexity)
	}
	require.Equal(expect, result)
}
func TestZeroFunction(t *testing.T) {
	require := require.New(t)
	// Empty tree
	n := &uast.Node{InternalType: "module"}
	comp := NPathComplexity(n)
	require.Equal(0, len(comp))
}

func TestRealUAST(t *testing.T) {
	fileNames := []string{
		"fixtures/npath/ifelse.java.json",
		"fixtures/npath/do_while.java.json",
		"fixtures/npath/while.java.json",
		"fixtures/npath/for.java.json",
		"fixtures/npath/someFuncs.java.json",
		"fixtures/npath/switch.java.json",
	}

	require := require.New(t)
	var result []int
	for _, name := range fileNames {
		file, err := os.Open(name)
		require.NoError(err)

		reader := bufio.NewReader(file)
		dec := json.NewDecoder(reader)
		res := &protocol.ParseResponse{}
		err = dec.Decode(res)
		require.NoError(err)
		n := res.UAST
		npathData := NPathComplexity(n)
		for _, v := range npathData {
			result = append(result, v.Complexity)
		}
	}

	expect := []int{2, 2, 2, 2, 2, 6, 2, 6, 3, 5, 4}

	require.Equal(expect, result)

}
