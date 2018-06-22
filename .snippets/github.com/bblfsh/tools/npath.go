package tools

import (
	"fmt"

	"gopkg.in/bblfsh/sdk.v1/uast"
)

type NPath struct{}

type NPathData struct {
	Name       string
	Complexity int
}

func (np NPath) Exec(n *uast.Node) error {
	result := NPathComplexity(n)
	fmt.Println(result)
	return nil
}

func (nd *NPathData) String() string {
	return fmt.Sprintf("FuncName:%s, Complexity:%d\n", nd.Name, nd.Complexity)
}

//Npath computes the NPath of functions in a *uast.Node.
//
//PMD is considered the reference implementation to assert correctness.
//See: https://pmd.github.io/pmd-5.7.0/pmd-java/xref/net/sourceforge/pmd/lang/java/rule/codesize/NPathComplexityRule.html
func NPathComplexity(n *uast.Node) []*NPathData {
	var result []*NPathData
	var funcs []*uast.Node
	var names []string

	if containsRoles(n, []uast.Role{uast.Function, uast.Body}, nil) {
		funcs = append(funcs, n)
		names = append(names, "NoName")
	} else {
		funcDecs := deepChildrenOfRoles(n, []uast.Role{uast.Function, uast.Declaration}, []uast.Role{uast.Argument})
		for _, funcDec := range funcDecs {
			if containsRoles(funcDec, []uast.Role{uast.Function, uast.Name}, nil) {
				names = append(names, funcDec.Token)
			}
			childNames := childrenOfRoles(funcDec, []uast.Role{uast.Function, uast.Name}, nil)
			if len(childNames) > 0 {
				names = append(names, childNames[0].Token)
			}
			childFuncs := childrenOfRoles(funcDec, []uast.Role{uast.Function, uast.Body}, nil)
			if len(childFuncs) > 0 {
				funcs = append(funcs, childFuncs[0])
			}
		}
	}
	for i, function := range funcs {
		npath := visitFunctionBody(function)
		result = append(result, &NPathData{Name: names[i], Complexity: npath})
	}

	return result
}

func visitorSelector(n *uast.Node) int {
	if containsRoles(n, []uast.Role{uast.Statement, uast.If}, []uast.Role{uast.Then, uast.Else}) {
		return visitIf(n)
	}
	if containsRoles(n, []uast.Role{uast.Statement, uast.While}, nil) {
		return visitWhile(n)
	}
	if containsRoles(n, []uast.Role{uast.Statement, uast.Switch}, nil) {
		return visitSwitch(n)
	}
	if containsRoles(n, []uast.Role{uast.Statement, uast.DoWhile}, nil) {
		return visitDoWhile(n)
	}
	if containsRoles(n, []uast.Role{uast.Statement, uast.For}, nil) {
		return visitFor(n)
	}
	if containsRoles(n, []uast.Role{uast.Statement, uast.Return}, nil) {
		return visitReturn(n)
	}
	if containsRoles(n, []uast.Role{uast.Statement, uast.Try}, nil) {
		return visitTry(n)
	}
	return visitNotCompNode(n)
}

func complexityMultOf(n *uast.Node) int {
	npath := 1
	for _, child := range n.Children {
		npath *= visitorSelector(child)
	}
	return npath
}

func visitFunctionBody(n *uast.Node) int {
	return complexityMultOf(n)
}

func visitNotCompNode(n *uast.Node) int {
	return complexityMultOf(n)
}

func visitIf(n *uast.Node) int {
	// (npath of if + npath of else (or 1) + bool_comp of if) * npath of next
	npath := 0
	ifThen := childrenOfRoles(n, []uast.Role{uast.If, uast.Then}, nil)
	ifCondition := childrenOfRoles(n, []uast.Role{uast.If, uast.Condition}, nil)
	ifElse := childrenOfRoles(n, []uast.Role{uast.If, uast.Else}, nil)

	if len(ifElse) > 0 {
		npath += complexityMultOf(ifElse[0])
	} else {
		npath++
	}
	npath *= complexityMultOf(ifThen[0])
	npath += expressionComp(ifCondition[0])

	return npath
}

func visitWhile(n *uast.Node) int {
	// (npath of while + bool_comp of while + npath of else (or 1)) * npath of next
	npath := 0
	whileCondition := childrenOfRoles(n, []uast.Role{uast.While, uast.Condition}, nil)
	whileBody := childrenOfRoles(n, []uast.Role{uast.While, uast.Body}, nil)
	whileElse := childrenOfRoles(n, []uast.Role{uast.While, uast.Else}, nil)
	// Some languages like python can have an else in a while loop
	if len(whileElse) > 0 {
		npath += complexityMultOf(whileElse[0])
	} else {
		npath++
	}

	npath *= complexityMultOf(whileBody[0])
	npath += expressionComp(whileCondition[0])

	return npath
}

func visitDoWhile(n *uast.Node) int {
	// (npath of do + bool_comp of do + 1) * npath of next
	npath := 1
	doWhileCondition := childrenOfRoles(n, []uast.Role{uast.DoWhile, uast.Condition}, nil)
	doWhileBody := childrenOfRoles(n, []uast.Role{uast.DoWhile, uast.Body}, nil)

	npath *= complexityMultOf(doWhileBody[0])
	npath += expressionComp(doWhileCondition[0])

	return npath
}

func visitFor(n *uast.Node) int {
	// (npath of for + bool_comp of for + 1) * npath of next
	npath := 1
	forBody := childrenOfRoles(n, []uast.Role{uast.For, uast.Body}, nil)
	if len(forBody) > 0 {
		npath *= complexityMultOf(forBody[0])
	}
	npath++
	return npath
}

func visitReturn(n *uast.Node) int {
	if aux := expressionComp(n); aux != 1 {
		return aux - 1
	}
	return 1
}

func visitSwitch(n *uast.Node) int {
	caseDefault := childrenOfRoles(n, []uast.Role{uast.Switch, uast.Default}, nil)
	switchCases := childrenOfRoles(n, []uast.Role{uast.Statement, uast.Switch, uast.Case}, []uast.Role{uast.Body})
	npath := 0

	if len(caseDefault) > 0 {
		npath += complexityMultOf(caseDefault[0])
	} else {
		npath++
	}
	for _, switchCase := range switchCases {
		npath += complexityMultOf(switchCase)
	}
	return npath
}

func visitTry(n *uast.Node) int {
	/*
		In pmd they decided the complexity of a try is the summatory of the complexity
		of the try body, catch body and finally body.I don't think this is the most acurate way
		of doing this.
	*/

	tryBody := childrenOfRoles(n, []uast.Role{uast.Try, uast.Body}, nil)
	tryCatch := childrenOfRoles(n, []uast.Role{uast.Try, uast.Catch}, nil)
	tryFinaly := childrenOfRoles(n, []uast.Role{uast.Try, uast.Finally}, nil)

	catchComp := 0
	if len(tryCatch) > 0 {
		for _, catch := range tryCatch {
			catchComp += complexityMultOf(catch)
		}
	}
	finallyComp := 0
	if len(tryFinaly) > 0 {
		finallyComp = complexityMultOf(tryFinaly[0])
	}
	npath := complexityMultOf(tryBody[0]) + catchComp + finallyComp

	return npath
}

func visitConditionalExpr(n *uast.Node) {
	// TODO ternary operators are not defined on the UAST yet
}

func expressionComp(n *uast.Node) int {
	orCount := deepCountChildrenOfRoles(n, []uast.Role{uast.Operator, uast.Boolean, uast.And}, nil)
	andCount := deepCountChildrenOfRoles(n, []uast.Role{uast.Operator, uast.Boolean, uast.Or}, nil)

	return orCount + andCount + 1
}

func containsRoles(n *uast.Node, andRoles []uast.Role, notRoles []uast.Role) bool {
	roleMap := make(map[uast.Role]bool)
	for _, r := range n.Roles {
		roleMap[r] = true
	}
	for _, r := range andRoles {
		if !roleMap[r] {
			return false
		}
	}
	if notRoles != nil {
		for _, r := range notRoles {
			if roleMap[r] {
				return false
			}
		}
	}
	return true
}

func childrenOfRoles(n *uast.Node, andRoles []uast.Role, notRoles []uast.Role) []*uast.Node {
	var children []*uast.Node
	for _, child := range n.Children {
		if containsRoles(child, andRoles, notRoles) {
			children = append(children, child)
		}
	}
	return children
}

func deepChildrenOfRoles(n *uast.Node, andRoles []uast.Role, notRoles []uast.Role) []*uast.Node {
	var childList []*uast.Node
	for _, child := range n.Children {
		if containsRoles(child, andRoles, notRoles) {
			childList = append(childList, child)
		}
		childList = append(childList, deepChildrenOfRoles(child, andRoles, notRoles)...)
	}
	return childList
}

func countChildrenOfRoles(n *uast.Node, andRoles []uast.Role, notRoles []uast.Role) int {
	count := 0
	for _, child := range n.Children {
		if containsRoles(child, andRoles, notRoles) {
			count++
		}
	}
	return count
}

func deepCountChildrenOfRoles(n *uast.Node, andRoles []uast.Role, notRoles []uast.Role) int {
	count := 0
	for _, child := range n.Children {
		if containsRoles(child, andRoles, notRoles) {
			count++
		}
		count += deepCountChildrenOfRoles(child, andRoles, notRoles)
	}
	return count
}
