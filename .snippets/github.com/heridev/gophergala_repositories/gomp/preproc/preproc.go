package preproc

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strconv"
	"strings"

	"github.com/gophergala/gomp/gensym"
)

type Cond int

type Context struct {
	genSym        func() string
	runtimeCalled bool
	cmap          ast.CommentMap
}

// ok is set to true when loop init part looks like:
// for variable := begin ; ... {}
func parseForInit(stmt *ast.Stmt) (variable *ast.Ident, begin *ast.Expr, ok bool) {
	if stmt == nil {
		return
	}
	var assignStmt *ast.AssignStmt
	if assignStmt, ok = (*stmt).(*ast.AssignStmt); !ok {
		return
	}
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return
	}
	if variable, ok = assignStmt.Lhs[0].(*ast.Ident); !ok {
		return
	}
	begin = &assignStmt.Rhs[0]
	return
}

// ok is set to true when loop cond part looks like:
// for ... ; variable (< | <= | > | >= ) end ; ... {}
// In this case op is set to token.{LSS|LEQ|GTR|GEQ}.
func parseForCond(expr *ast.Expr) (variable *ast.Ident, op token.Token, end *ast.Expr, ok bool) {
	if expr == nil {
		return
	}
	binaryExpr, ok := (*expr).(*ast.BinaryExpr)
	if !ok {
		return
	}
	switch binaryExpr.Op {
	case token.LEQ, token.LSS, token.GTR, token.GEQ:
		op = binaryExpr.Op
	default:
		return
	}
	if variable, ok = binaryExpr.X.(*ast.Ident); !ok {
		return
	}
	end = &binaryExpr.Y
	return
}

// ok is set to true when loop post part looks like:
// for ... ; (variable++ | variable-- | variable += step | variable -= step)
// In this case op is set to token.{ADD_ASSIGN|SUB_ASSIGN}.
// Also, in case of ++ or -- operators step is set to 1.
func parseForPost(stmt *ast.Stmt) (variable *ast.Ident, op token.Token, step *ast.Expr, ok bool) {
	if stmt == nil {
		return
	}

	if incDecStmt, isIncDec := (*stmt).(*ast.IncDecStmt); isIncDec {
		variable, ok = incDecStmt.X.(*ast.Ident)
		if !ok {
			return
		}
		newStmt := &ast.AssignStmt{
			Lhs: []ast.Expr{variable},
			Rhs: []ast.Expr{mkIntLit(1)},
		}
		switch incDecStmt.Tok {
		case token.INC:
			newStmt.Tok = token.ADD_ASSIGN
		case token.DEC:
			newStmt.Tok = token.SUB_ASSIGN
		default:
			panic("Unknown op in IncDecStmt")
		}
		*stmt = newStmt
	}
	if assignStmt, isAssignStmt := (*stmt).(*ast.AssignStmt); isAssignStmt {
		if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
			return
		}
		if variable, ok = assignStmt.Lhs[0].(*ast.Ident); !ok {
			return
		}
		switch assignStmt.Tok {
		case token.ADD_ASSIGN, token.SUB_ASSIGN:
			op = assignStmt.Tok
		default:
			ok = false
			return
		}
		step = &assignStmt.Rhs[0]
	}
	return
}

func mkIntLit(n int) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(n)}
}

func mkSym(context *Context) *ast.Ident {
	return &ast.Ident{Name: context.genSym()}
}

func mkIdent(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func mkTypeConv(expr ast.Expr, t string) ast.Expr {
	return &ast.CallExpr{
		Fun:  mkIdent(t),
		Args: []ast.Expr{expr},
	}
}

func mkAssign(lhs, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{rhs},
	}
}

func mkNegate(expr ast.Expr) ast.Expr {
	return &ast.UnaryExpr{Op: token.SUB, X: expr}
}

func mkGoLambda(body *ast.BlockStmt, arg *ast.Ident) *ast.GoStmt {
	return &ast.GoStmt{
		Call: &ast.CallExpr{
			Fun: &ast.FuncLit{
				Type: &ast.FuncType{Params: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Names: []*ast.Ident{arg},
							Type:  mkIdent("int")}}}},
				Body: body,
			},
			Args: []ast.Expr{mkTypeConv(arg, "int")},
		},
	}
}

// Creates a channel for synchronization between goroutines. Size
// of channel equals to number of goroutines.
func mkSyncChannelDecl(chanType ast.Expr, numOfGoRoutines *ast.Ident) ast.Expr {
	return &ast.CallExpr{
		Fun: mkIdent("make"),
		Args: []ast.Expr{
			&ast.ChanType{
				Value: chanType,
				Dir:   ast.SEND | ast.RECV,
			},
			numOfGoRoutines,
		}}
}

// Creates an outer loop that schedules execution of disjoint parts of the range to goroutines.
// The loop looks like:
// for loopVar := 0; cond; loopVar++ {
//    go func(i int) {
//       body
//       channel <- ...
//    }(int(loopVar))
// }
func mkOuterLoop(loopVar *ast.Ident, cond, channel, channelType ast.Expr, body ast.Stmt) *ast.ForStmt {
	return &ast.ForStmt{
		Init: mkAssign(loopVar, mkIntLit(0)),
		Cond: cond,
		Post: &ast.IncDecStmt{
			X:   loopVar,
			Tok: token.INC,
		},
		Body: &ast.BlockStmt{List: []ast.Stmt{
			mkGoLambda(
				&ast.BlockStmt{
					List: []ast.Stmt{
						body,
						&ast.SendStmt{
							Chan: channel,
							Value: &ast.CompositeLit{
								Type: channelType,
							}}}},
				loopVar),
		},
		},
	}
}

// Creates an inner loop that is going to be executed on an individual goroutine.
// The loop looks like:
// for loopVar, counter := begin, 0;
//     loopVar <= end && counter < taskSize;
//     loopVar, counter = loopVar + step, counter + 1 {}
//
// Note that loop body is not set.
func mkInnerLoop(loopVar *ast.Ident, begin, end, step, taskSize ast.Expr, context *Context) *ast.ForStmt {
	counter := mkSym(context)
	return &ast.ForStmt{
		Init: &ast.AssignStmt{
			Lhs: []ast.Expr{loopVar, counter},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{begin, mkIntLit(0)},
		},
		Cond: &ast.BinaryExpr{
			X: &ast.BinaryExpr{
				X:  loopVar,
				Op: token.LEQ,
				Y:  end,
			},
			Op: token.LAND,
			Y: &ast.BinaryExpr{
				X:  counter,
				Op: token.LSS,
				Y:  taskSize,
			},
		},
		Post: &ast.AssignStmt{
			Lhs: []ast.Expr{loopVar, counter},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{&ast.BinaryExpr{
				X:  loopVar,
				Op: token.ADD,
				Y:  step},
				&ast.BinaryExpr{
					X:  counter,
					Op: token.ADD,
					Y:  mkIntLit(1)}},
		}}
}

// Creates a loop that reads from the channel numOfGoRoutines times.
func mkSyncChannelLoop(numOfGoRoutines, channel *ast.Ident, context *Context) ast.Stmt {
	loopVar := mkSym(context)
	return &ast.ForStmt{
		Init: mkAssign(loopVar, mkIntLit(0)),
		Cond: &ast.BinaryExpr{
			X:  loopVar,
			Op: token.LSS,
			Y:  numOfGoRoutines,
		},
		Post: &ast.IncDecStmt{
			X:   loopVar,
			Tok: token.INC,
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.UnaryExpr{
						X:  channel,
						Op: token.ARROW,
					}}}}}
}

// Calculates task size as:
// taskSize := (end - begin + numCPU * step) / (numCPU * step).
//
// Task size is an estimation of number of iterations per go routine.
// Returns new symbol + expression.
func calcTaskSize(begin, end, step *ast.Ident, context *Context) ast.Expr {
	denom := ast.BinaryExpr{
		X:  step,
		Op: token.MUL,
		Y: &ast.CallExpr{Fun: &ast.SelectorExpr{
			X:   mkIdent("runtime"),
			Sel: mkIdent("NumCPU")}},
	}
	num := ast.BinaryExpr{
		X: &ast.BinaryExpr{
			X:  end,
			Op: token.SUB,
			Y:  begin,
		},
		Op: token.ADD,
		Y:  &denom,
	}
	context.runtimeCalled = true
	return &ast.BinaryExpr{X: &num, Op: token.QUO, Y: &denom}
}

// Calculates number of goroutines as:
// numOfGoRoutines := (end - begin) / (taskSize * step) + 1
func calcNumOfGoRoutines(begin, end, step, taskSize *ast.Ident) ast.Expr {
	return &ast.BinaryExpr{
		X: &ast.BinaryExpr{
			X: &ast.BinaryExpr{
				X:  end,
				Op: token.SUB,
				Y:  begin,
			},
			Op: token.QUO,
			Y: &ast.BinaryExpr{
				X:  taskSize,
				Op: token.MUL,
				Y:  step,
			},
		},
		Op: token.ADD,
		Y:  mkIntLit(1),
	}
}

// Creates code for parallel loop.
func emitSchedulerLoop(loopVar, begin, end, step *ast.Ident, context *Context, loopBody *ast.BlockStmt) (code []ast.Stmt) {
	taskSize := mkSym(context)
	code = append(code, mkAssign(taskSize, calcTaskSize(begin, end, step, context)))

	routineId := mkSym(context)

	// routineBegin := begin + routineId * taskSize * step
	routineBegin := ast.BinaryExpr{
		X:  begin,
		Op: token.ADD,
		Y: &ast.BinaryExpr{
			X:  routineId,
			Op: token.MUL,
			Y: &ast.BinaryExpr{
				X:  taskSize,
				Op: token.MUL,
				Y:  step,
			},
		},
	}

	// begin + routineId * taskSize * step <= end
	routineBeginCheckExpr := ast.BinaryExpr{
		X:  &routineBegin,
		Op: token.LEQ,
		Y:  end,
	}

	numOfGoRoutines := mkSym(context)
	code = append(code, mkAssign(
		numOfGoRoutines,
		calcNumOfGoRoutines(begin, end, step, taskSize)))

	emptyStruct := ast.StructType{
		Fields: &ast.FieldList{},
	}

	channel := mkSym(context)
	code = append(code, mkAssign(
		channel,
		mkSyncChannelDecl(&emptyStruct, numOfGoRoutines)))

	{
		innerLoop := mkInnerLoop(loopVar, &routineBegin, end, step, taskSize, context)
		innerLoop.Body = loopBody

		outerLoop := mkOuterLoop(routineId, &routineBeginCheckExpr, channel, &emptyStruct, innerLoop)
		code = append(code, outerLoop)
	}
	code = append(code, mkSyncChannelLoop(numOfGoRoutines, channel, context))
	return code
}

func visitFor(stmt *ast.ForStmt, context *Context) *ast.BlockStmt {
	initVar, initExpr, initOk := parseForInit(&stmt.Init)
	condVar, condOp, condExpr, condOk := parseForCond(&stmt.Cond)
	postVar, postOp, postExpr, postOk := parseForPost(&stmt.Post)

	if !initOk || !condOk || !postOk {
		return nil
	}
	if initVar.Name != condVar.Name || initVar.Name != postVar.Name {
		return nil
	}

	// Following code transforms loops in the form:
	// for i := b ; i > e; i (+= | -=) d { ... }
	// to the form:
	// for i := b ; i >= (e + 1); i (+= | -=) d {...}
	if condOp == token.GTR {
		*condExpr = &ast.BinaryExpr{
			X:  *condExpr,
			Op: token.ADD,
			Y:  mkIntLit(1),
		}
		condOp = token.GEQ
	}

	// Following code transforms loops in the form:
	// for i := b ; i >= e; i -= d { ... }
	// to the form:
	// for i := b - d * ((b - e) / d) ; i <= b ; i -= d { ... }
	if condOp == token.GEQ {
		if postOp == token.ADD_ASSIGN {
			postOp = token.SUB_ASSIGN
			*postExpr = mkNegate(*postExpr)
		}
		newInitExpr := &ast.BinaryExpr{
			X:  *initExpr,
			Op: token.SUB,
			Y: &ast.BinaryExpr{
				X:  *postExpr,
				Op: token.MUL,
				Y: &ast.BinaryExpr{
					X: &ast.BinaryExpr{
						X:  *initExpr,
						Op: token.SUB,
						Y:  *condExpr,
					},
					Op: token.QUO,
					Y:  *postExpr,
				},
			},
		}
		*initExpr, *condExpr = newInitExpr, *initExpr
		condOp = token.LEQ
	}

	if condOp == token.LSS {
		if postOp == token.SUB_ASSIGN {
			postOp = token.ADD_ASSIGN
			*postExpr = mkNegate(*postExpr)
		}

		condOp = token.LEQ
		*condExpr = &ast.BinaryExpr{X: *condExpr, Op: token.SUB, Y: mkIntLit(1)}
	}

	block := new(ast.BlockStmt)
	block.List = []ast.Stmt{}
	initVarSym, condVarSym, incVarSym := mkSym(context), mkSym(context), mkSym(context)
	{
		boundsDecl := ast.AssignStmt{
			Lhs: []ast.Expr{initVarSym, condVarSym, incVarSym},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{*initExpr, *condExpr, *postExpr},
		}

		*initExpr, *condExpr = ast.Expr(initVarSym), ast.Expr(condVarSym)
		stmt.Post = &ast.AssignStmt{
			Lhs: []ast.Expr{initVar},
			Tok: postOp,
			Rhs: []ast.Expr{incVarSym},
		}

		block.List = append(block.List, &boundsDecl)
	}

	switch condOp {
	case token.LSS, token.LEQ:
		block.List = append(
			block.List,
			emitSchedulerLoop(initVar, initVarSym, condVarSym, incVarSym, context, stmt.Body)...)
	default:
		block.List = append(block.List, ast.Stmt(stmt))
	}
	return block
}

func visitExpr(e *ast.Expr, context *Context) {
	if e == nil {
		return
	}
	switch t := (*e).(type) {
	case *ast.FuncLit:
		if t.Body == nil {
			return
		}
		for _, s := range t.Body.List {
			visitStmt(&s, context)
		}
	}
}

func shouldParalellize(stmt *ast.Stmt, context *Context) bool {
	commentGroups := ((*context).cmap)[(*stmt).(ast.Node)]
	length := len(commentGroups)
	if length == 0 {
		return false
	}
	commentGroup := *commentGroups[length-1]
	length = len(commentGroup.List)
	return (length > 0) && strings.HasPrefix(commentGroup.List[length-1].Text, "//gomp")
}

func visitStmt(stmt *ast.Stmt, context *Context) {
	if stmt == nil {
		return
	}
	switch t := (*stmt).(type) {
	case *ast.AssignStmt:
		for _, e := range t.Rhs {
			visitExpr(&e, context)
		}
	case *ast.ForStmt:
		if shouldParalellize(stmt, context) {
			if block := visitFor(t, context); block != nil {
				*stmt = block
				//TODO: save old comments here
			}
		} else {
			visitBlock(t.Body, context)
		}
	case *ast.BlockStmt:
		visitBlock(t, context)
	case *ast.IfStmt:
		visitBlock(t.Body, context)
	case *ast.SwitchStmt:
		visitBlock(t.Body, context)
	case *ast.TypeSwitchStmt:
		visitBlock(t.Body, context)
	case *ast.CaseClause:
		for i, _ := range t.Body {
			visitStmt(&t.Body[i], context)
		}
	}
}

func visitBlock(stmt *ast.BlockStmt, context *Context) {
	if stmt != nil {
		for i, _ := range stmt.List {
			visitStmt(&stmt.List[i], context)
		}
	}
}

func visitFunction(f *ast.FuncDecl, context *Context) {
	if f.Body != nil {
		visitBlock(f.Body, context)
	}
}

// Run preprocessor on a source. filename is used for error reporting.
// This function is currently not implemented.
func PreprocFile(source, filename string) (result string, err error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, source,
		parser.ParseComments|parser.AllErrors)
	if err != nil {
		return
	}
	context := Context{gensym.MkGen(source), false, ast.NewCommentMap(fset, file, file.Comments)}

	for _, decl := range file.Decls {
		switch t := decl.(type) {
		case *ast.FuncDecl:
			visitFunction(t, &context)
		}
	}

	if context.runtimeCalled {
		const runtimePath = `"runtime"`
		runtimeImported := false
		for _, spec := range file.Imports {
			if spec.Path != nil && spec.Path.Value == runtimePath {
				runtimeImported = true
				break
			}
		}
		if !runtimeImported {
			runtimeImport := ast.ImportSpec{
				Path: &ast.BasicLit{Value: runtimePath, Kind: token.STRING}}
			runtimeDecl := ast.GenDecl{Tok: token.IMPORT, Specs: []ast.Spec{&runtimeImport}}
			file.Decls = append([]ast.Decl{&runtimeDecl}, file.Decls...)
			file.Imports = append(file.Imports, &runtimeImport)
		}
	}
	file.Imports = []*ast.ImportSpec{}

	//Delete all comments from file
	file.Comments = nil
	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), file)
	result = buf.String()
	return
}
