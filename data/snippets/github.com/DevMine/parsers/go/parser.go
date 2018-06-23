// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"archive/tar"
	"encoding/json"
	"flag"
	"fmt"
	stdAst "go/ast"
	"go/parser"
	stdToken "go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DevMine/srcanlzr/src"
	"github.com/DevMine/srcanlzr/src/ast"
	"github.com/DevMine/srcanlzr/src/token"
)

var (
	// List of all packages of the projects.
	pkgs = make(map[string]*src.Package)

	goLang = src.Language{
		Lang:      src.Go,
		Paradigms: []string{src.Compiled, src.Concurrent, src.Imperative, src.Structured},
	}

	goToastToken = map[stdToken.Token]string{
		stdToken.ADD: token.ADD,
		stdToken.SUB: token.SUB,
		stdToken.MUL: token.MUL,
		stdToken.QUO: token.QUO,
		stdToken.REM: token.MOD,

		stdToken.NEQ:  token.NEQ,
		stdToken.LEQ:  token.LEQ,
		stdToken.GEQ:  token.GEQ,
		stdToken.EQL:  token.EQ,
		stdToken.LSS:  token.LSS,
		stdToken.GTR:  token.GTR,
		stdToken.LAND: token.LAND,
		stdToken.LOR:  token.LOR,

		stdToken.NOT: token.NOT,

		stdToken.INT:    token.IntLit,
		stdToken.FLOAT:  token.FloatLit,
		stdToken.IMAG:   token.ImagLit,
		stdToken.CHAR:   token.CharLit,
		stdToken.STRING: token.StringLit,
	}
)

func parseFuncDecl(fd *stdAst.FuncDecl, s *src.SrcFile) int64 {
	var loc int64 = 1

	var params []*ast.Field
	if fd.Type.Params != nil {
		params = parseFieldList(fd.Type.Params.List)
	}

	var results []*ast.Field
	if fd.Type.Results != nil {
		results = parseFieldList(fd.Type.Results.List)
	}

	var stmts []ast.Stmt
	if fd.Body != nil {
		stmts = parseStmt(fd.Body.List, &loc)
	}

	var typ *ast.FuncType
	if (params != nil && len(params) > 0) || (results != nil && len(results) > 0) {
		typ = &ast.FuncType{
			Params:  params,
			Results: results,
		}
	}

	s.Funcs = append(s.Funcs, &ast.FuncDecl{
		Doc:  formatDoc(fd.Doc),
		Name: fd.Name.Name,
		Type: typ,
		Body: stmts,
	})

	return loc
}

func parseCallExpr(ce *stdAst.CallExpr) *ast.CallExpr {
	var funcName string

	switch ce.Fun.(type) {
	case *stdAst.Ident:
		funcName = ce.Fun.(*stdAst.Ident).Name
	}

	return &ast.CallExpr{
		ExprName: token.CallExprName,
		// XXX handle  FuncRef correctly
		Fun: &ast.FuncRef{
			Namespace: "",
			FuncName:  funcName,
		},
		Args: parseExprSlice(ce.Args),
		Line: 0,
	}
}

func parseStmt(s []stdAst.Stmt, loc *int64) []ast.Stmt {
	stmts := []ast.Stmt{}

	// FIXME suppport line properly
	var line int64 = 1

	for _, s := range s {
		switch s.(type) {
		case *stdAst.AssignStmt:
			as := s.(*stdAst.AssignStmt)

			stmts = append(stmts, ast.AssignStmt{
				StmtName: token.AssignStmtName,
				LHS:      parseExprSlice(as.Lhs),
				RHS:      parseExprSlice(as.Rhs),
				Line:     line,
			})
		case *stdAst.BlockStmt:
			stmts = append(stmts, parseStmt(s.(*stdAst.BlockStmt).List, loc))
		case *stdAst.BranchStmt:
			// this should never happen since the switch statement is parsed
			// in a separate function
		case *stdAst.DeclStmt:
			// TODO
		case *stdAst.DeferStmt:
			// Defer is considered as a function call
			stmts = append(stmts, &ast.ExprStmt{
				StmtName: token.ExprStmtName,
				X:        parseCallExpr(s.(*stdAst.DeferStmt).Call),
			})
		case *stdAst.EmptyStmt:
			// does not count as a line of code
			(*loc)--
		case *stdAst.ExprStmt:
			es := s.(*stdAst.ExprStmt)
			stmts = append(stmts, &ast.ExprStmt{
				StmtName: token.ExprStmtName,
				X:        parseExpr(es.X),
			})
		case *stdAst.GoStmt:
			// Go is considered as a function call
			// XXX try to keep the semantic of a concurrent call
			stmts = append(stmts, &ast.ExprStmt{
				StmtName: token.ExprStmtName,
				X:        parseCallExpr(s.(*stdAst.GoStmt).Call),
			})
		case *stdAst.IfStmt:
			is := s.(*stdAst.IfStmt)
			stmts = append(stmts, ast.IfStmt{
				StmtName: token.IfStmtName,
				Cond:     parseExpr(is.Cond),
				Body:     parseStmt(is.Body.List, loc),
				Line:     line,
			})
		case *stdAst.ForStmt:
			fs := s.(*stdAst.ForStmt)
			stmts = append(stmts, ast.LoopStmt{
				StmtName: token.LoopStmtName,
				Body:     parseStmt(fs.Body.List, loc),
				Line:     line,
			})
		case *stdAst.LabeledStmt:
			stmts = append(stmts, ast.OtherStmt{StmtName: token.OtherStmtName, Line: line})
		case *stdAst.RangeStmt:
			rs := s.(*stdAst.RangeStmt)

			vars := []ast.Expr{}

			if rs.Key != nil {
				vars = append(vars, parseExpr(rs.Key))
			}

			if rs.Value != nil {
				vars = append(vars, parseExpr(rs.Value))
			}

			stmts = append(stmts, &ast.RangeLoopStmt{
				StmtName: token.RangeLoopStmtName,
				Vars:     vars,
				Iterable: parseExpr(rs.X),
				Body:     parseStmt(rs.Body.List, loc),
				Line:     line,
			})
		case *stdAst.ReturnStmt:
			stmts = append(stmts, &ast.ReturnStmt{
				StmtName: token.ReturnStmtName,
				Results:  parseExprSlice(s.(*stdAst.ReturnStmt).Results),
				Line:     line,
			})
		case *stdAst.SelectStmt:
			// Select statement is not yet supported, however, we add the number
			// of statements in its Body. This is not accurate since the
			// statements contained in the Body may contain a list of statements
			// but for now this estimation is good enough.
			*loc += int64(len(s.(*stdAst.SelectStmt).Body.List))
		case *stdAst.SendStmt:
			stmts = append(stmts, ast.OtherStmt{StmtName: token.OtherStmtName, Line: line})
		case *stdAst.SwitchStmt:
			stmts = append(stmts, parseSwitchStmt(s.(*stdAst.SwitchStmt), loc))
		default:
			// this should never happen, except if a new statement is introduced and
			// not yet supported by this parser
			stmts = append(stmts, ast.OtherStmt{StmtName: token.OtherStmtName, Line: line})
		}

		line++
		(*loc)++
	}

	return stmts
}

func parseSwitchStmt(ss *stdAst.SwitchStmt, loc *int64) *ast.SwitchStmt {
	switchstmt := ast.SwitchStmt{StmtName: token.SwitchStmtName}

	if ss.Init != nil {
		initStmt := parseStmt([]stdAst.Stmt{ss.Init}, loc)
		if len(initStmt) > 0 {
			switchstmt.Init = initStmt[0]
		}
	}

	if ss.Tag != nil {
		switchstmt.Cond = parseExpr(ss.Tag)
	}

	body := ss.Body.List
	if body != nil && len(body) > 0 {
		cases, defaultCase := parseCaseClause(body, loc)
		switchstmt.CaseClauses = cases

		if defaultCase != nil {
			switchstmt.Default = defaultCase
		}
	}

	return &switchstmt
}

func parseCaseClause(es []stdAst.Stmt, loc *int64) (cases []*ast.CaseClause, defaultCase []ast.Stmt) {
	cases = []*ast.CaseClause{}

	if len(es) == 0 {
		return
	}

	for _, val := range es {
		// XXX: should never fail, but a check could be nice too
		caseClause := val.(*stdAst.CaseClause)
		if caseClause == nil {
			// default
			defaultCase = parseStmt(caseClause.Body, loc)
			continue
		}

		cases = append(cases, &ast.CaseClause{
			Conds: parseExprSlice(caseClause.List),
			Body:  parseStmt(caseClause.Body, loc),
		})
	}

	return
}

func parseExpr(e stdAst.Expr) ast.Expr {
	switch e.(type) {
	case *stdAst.BinaryExpr:
		be := e.(*stdAst.BinaryExpr)

		return &ast.BinaryExpr{
			ExprName:  token.BinaryExprName,
			LeftExpr:  parseExpr(be.X),
			Op:        goToastToken[be.Op],
			RightExpr: parseExpr(be.Y),
		}
	case *stdAst.CallExpr:
		ce := e.(*stdAst.CallExpr)
		namespace, funcName := parseFuncRef(ce.Fun)

		return &ast.CallExpr{
			ExprName: token.CallExprName,
			Fun: &ast.FuncRef{
				Namespace: namespace,
				FuncName:  funcName,
			},
			Args: parseExprSlice(ce.Args),
			Line: 0, // XXX: count
		}
	case *stdAst.IndexExpr:
		ie := e.(*stdAst.IndexExpr)

		return &ast.IndexExpr{
			ExprName: token.IndexExprName,
			X:        parseExpr(ie.X),
			Index:    parseExpr(ie.Index),
		}
	case *stdAst.KeyValueExpr:
		/*kv := e.(*stdAst.KeyValueExpr)

		return &ast.KeyValuePair{
			ExprName: ast.KeyValuePairName,
			Key:      parseExpr(kv.Key),
			Value:    parseExpr(kv.Value),
		}*/
		// TODO: handle this properly
		return nil
	case *stdAst.ParenExpr:
		return parseExpr(e.(*stdAst.ParenExpr).X)
	case *stdAst.SelectorExpr:
		// TODO: handle this case correctly
		se := e.(*stdAst.SelectorExpr)

		return &ast.Ident{ExprName: token.IdentName, Name: se.Sel.Name}
	case *stdAst.SliceExpr:
		// XXX not yet supported by the source API
		return nil
	case *stdAst.StarExpr:
		// XXX add support of the * operator in the ast API
		se := e.(*stdAst.StarExpr)

		return &ast.UnaryExpr{
			ExprName: token.UnaryExprName,
			Op:       token.ADDR,
			X:        parseExpr(se.X),
		}
	case *stdAst.TypeAssertExpr:
		// XXX not yet supported by the source API
		return nil
	case *stdAst.UnaryExpr:
		ue := e.(*stdAst.UnaryExpr)

		return &ast.UnaryExpr{
			ExprName: token.UnaryExprName,
			Op:       goToastToken[ue.Op],
			X:        parseExpr(ue.X),
		}
	case *stdAst.BasicLit:
		bl := e.(*stdAst.BasicLit)

		return &ast.BasicLit{
			ExprName: token.BasicLitName,
			Kind:     goToastToken[bl.Kind],
			Value:    formatStringLit(bl.Value),
		}
	case *stdAst.CompositeLit:
		// XXX check if really useful
	case *stdAst.FuncLit:
		fl := e.(*stdAst.FuncLit)

		// FIXME pass loc as arg everywhere
		var loc int64

		var params []*ast.Field
		if fl.Type.Params != nil {
			parseFieldList(fl.Type.Params.List)
		}

		var results []*ast.Field
		if fl.Type.Results != nil {
			parseFieldList(fl.Type.Results.List)
		}

		return &ast.FuncLit{
			ExprName: token.FuncLitName,
			Type: &ast.FuncType{
				Params:  params,
				Results: results,
			},
			Body: parseStmt(fl.Body.List, &loc),
		}
	case *stdAst.Ident:
		ident := e.(*stdAst.Ident)

		return &ast.Ident{ExprName: token.IdentName, Name: ident.Name}
		// FIXME move the specifiers
		/*case *stdAst.TypeSpec:
		ts := e.(*stdAst.TypeSpec)

		return &ast.TypeSpec{
			Doc: formatDoc(ts.Doc),
			Name: &ast.Ident{
				ExprName: ast.IdentName,
				Name:     ts.Name.Name,
			},
			Type: parseExpr(ts.Type),
		}*/
		/*case *stdAst.ValueSpec:
		// XXX not yet fully supported
		return nil*/
	case *stdAst.StructType:
		strct := e.(*stdAst.StructType)

		// XXX count LoC
		return &ast.StructType{
			ExprName: token.StructTypeName,
			Fields:   parseFieldList(strct.Fields.List),
		}
	}

	return nil
}

func parseFieldList(fs []*stdAst.Field) []*ast.Field {
	if fs == nil || len(fs) == 0 {
		return nil
	}

	fields := make([]*ast.Field, len(fs), len(fs))

	for ind, val := range fs {
		var name string

		if val.Names != nil {
			names := make([]string, len(val.Names), len(val.Names))
			for i, n := range val.Names {
				names[i] = n.Name
			}

			name = strings.Join(names, ",")
		}

		var typ string
		if ident := parseTypeAsIdent(val.Type); ident != nil {
			typ = ident.Name
		}

		fields[ind] = &ast.Field{
			Doc:  formatDoc(val.Doc),
			Name: name,
			Type: typ,
		}
	}

	return fields
}

func parseTypeAsIdent(e stdAst.Expr) *ast.Ident {
	if e == nil {
		return nil
	}

	var name string

	switch e.(type) {
	case *stdAst.ArrayType:
		// XXX we just use "ARRAY" for now, but it would be nice to have more
		// information about this type.
		name = token.TypeArrayName
	case *stdAst.ChanType:
		name = token.TypeUnsupportedName
	case *stdAst.FuncType:
		// XXX we just use "FUNC" for now, but it would be nice to have more
		// information about this type.
		name = token.TypeFuncName
	case *stdAst.InterfaceType:
		// XXX we just use "INTERFACE" for now, but it would be nice to have more
		// information about this type.
		name = token.TypeInterfaceName
	case *stdAst.MapType:
		// XXX we just use "MAP" for now, but it would be nice to have more
		// information about this type.
		name = token.TypeMapName
	case *stdAst.StructType:
		// XXX we just use "STRUCT" for now, but it would be nice to have more
		// information about this type.
		name = token.TypeStructName
	case *stdAst.Ident:
		// XXX it would be great to "normalize" the type name. For instance:
		// int, int32 and int64 should just be an "INTEGER".
		name = e.(*stdAst.Ident).Name
	default:
		return nil
	}

	return &ast.Ident{
		ExprName: token.IdentName,
		Name:     name,
	}
}

func parseExprSlice(es []stdAst.Expr) []ast.Expr {
	ses := make([]ast.Expr, 0)
	for _, val := range es {
		if expr := parseExpr(val); expr != nil {
			ses = append(ses, expr)
		}
	}

	return ses
}

func parseFuncRef(e stdAst.Expr) (namespace string, funcName string) {
	switch e.(type) {
	case *stdAst.SelectorExpr:
		se := e.(*stdAst.SelectorExpr)
		namespace = fmt.Sprintf("%v", se.X)
		funcName = se.Sel.Name
	case *stdAst.Ident:
		i := e.(*stdAst.Ident)
		funcName = i.Name
	}

	return
}

func parseSrcFileConst(cst *stdAst.ValueSpec, s *src.SrcFile, doc []string) int64 {
	for ind, val := range cst.Names {
		if len(cst.Values) == len(cst.Names) {
			s.Constants = append(s.Constants, &ast.GlobalDecl{
				Name:       &ast.Ident{ExprName: token.IdentName, Name: val.Name},
				Value:      parseExpr(cst.Values[ind]),
				Type:       parseTypeAsIdent(cst.Type),
				Visibility: visibility(val.Name),
			})

			continue
		}

		s.Constants = append(s.Constants, &ast.GlobalDecl{
			Name:       &ast.Ident{ExprName: token.IdentName, Name: val.Name},
			Type:       parseTypeAsIdent(cst.Type),
			Visibility: visibility(val.Name),
		})
	}

	return 1
}

func parseSrcFileVar(v *stdAst.ValueSpec, s *src.SrcFile, doc []string) int64 {
	for ind, val := range v.Names {
		if len(v.Values) == len(v.Names) {
			s.Vars = append(s.Vars, &ast.GlobalDecl{
				Doc:        doc,
				Name:       &ast.Ident{ExprName: token.IdentName, Name: val.Name},
				Value:      parseExpr(v.Values[ind]),
				Type:       parseTypeAsIdent(v.Type),
				Visibility: visibility(val.Name),
			})

			continue
		}

		s.Vars = append(s.Vars, &ast.GlobalDecl{
			Doc:        doc,
			Name:       &ast.Ident{ExprName: token.IdentName, Name: val.Name},
			Type:       parseTypeAsIdent(v.Type),
			Visibility: visibility(val.Name),
		})
	}

	return 1
}

func parseSrcFileType(ts *stdAst.TypeSpec, s *src.SrcFile, doc []string) int64 {
	s.TypeSpecs = append(s.TypeSpecs, &ast.TypeSpec{
		Doc: doc,
		Name: &ast.Ident{
			ExprName: token.IdentName,
			Name:     ts.Name.Name,
		},
		Type: parseExpr(ts.Type),
	})

	return 1
}

func parseGenDecl(genDecl *stdAst.GenDecl, s *src.SrcFile) int64 {
	var loc int64

	for _, spec := range genDecl.Specs {
		switch genDecl.Tok {
		case stdToken.CONST:
			loc += parseSrcFileConst(spec.(*stdAst.ValueSpec), s, formatDoc(genDecl.Doc))
		case stdToken.VAR:
			loc += parseSrcFileVar(spec.(*stdAst.ValueSpec), s, formatDoc(genDecl.Doc))
		case stdToken.TYPE:
			loc += parseSrcFileType(spec.(*stdAst.TypeSpec), s, formatDoc(genDecl.Doc))
		}
	}

	return loc
}

func parseDecl(decl interface{}, s *src.SrcFile) int64 {
	var loc int64

	// Field, XxxSpec, FuncDecl, LabeledStmt, AssignStmt, Scope;
	switch decl.(type) {
	case *stdAst.FuncDecl:
		loc += parseFuncDecl(decl.(*stdAst.FuncDecl), s)
	case *stdAst.GenDecl:
		loc += parseGenDecl(decl.(*stdAst.GenDecl), s)
	case *stdAst.LabeledStmt, *stdAst.AssignStmt:
		// TODO
		//loc += parseStmt([]stdAst.Stmt{decl.(stdAst.Stmt)}, &loc)
	case *stdAst.Scope:
		loc += parseScope(decl.(*stdAst.Scope), s)
	default:
		// XxxSpec
		// TODO
		//loc += parseExpr(decl.(stdAst.Expr), s)
	}

	return loc
}

func parseScope(scope *stdAst.Scope, s *src.SrcFile) int64 {
	var loc int64

	for _, v := range scope.Objects {
		loc += parseDecl(v.Decl, s)

		// XXX: necessary?
		/*switch v.Kind {
		case stdAst.Con:
			loc += parseConstants(v, &s)
		case stdAst.Var:
			loc += processVariables(v, &s)
		case stdAst.Fun:
			loc += processFunctions(k, v, &s, astCode)
		}*/
	}

	return loc
}

// formatStringLit formats a string literal, ie. it removes double-quotes.
func formatStringLit(str string) string {
	return strings.Trim(str, "\"")
}

func formatDoc(docGroup *stdAst.CommentGroup) []string {
	if docGroup == nil {
		return nil
	}

	list := docGroup.List
	doc := make([]string, len(list), len(list))

	for ind, val := range list {
		doc[ind] = trimComment(val.Text)
	}

	return doc
}

func trimComment(comm string) string {
	newComm := comm
	newComm = strings.TrimPrefix(newComm, "//")
	newComm = strings.Trim(newComm, " ")
	newComm = strings.TrimSuffix(newComm, "\n")
	return newComm
}

func isExported(name string) bool {
	if len(name) == 0 {
		return false
	}

	firstLetter := name[0:1]
	return firstLetter == strings.ToUpper(firstLetter)
}

func visibility(ident string) string {
	if isExported(ident) {
		return token.PublicVisibility
	}
	return token.PackageVisibility
}

func runParser(path string, data []byte) *src.SrcFile {
	var loc int64

	s := src.SrcFile{
		Path:      path,
		Lang:      &goLang,
		Imports:   []string{},
		Constants: []*ast.GlobalDecl{},
		Vars:      []*ast.GlobalDecl{},
		Funcs:     []*ast.FuncDecl{},
	}

	fset := stdToken.NewFileSet()

	f, err := parser.ParseFile(fset, path, data, parser.ParseComments)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		// skip parsing errors
		return nil
	}

	for _, i := range f.Imports {
		name := i.Path.Value
		s.Imports = append(s.Imports, name[1:len(name)-1])
	}

	for _, d := range f.Decls {
		loc += parseDecl(d, &s)
	}

	s.LoC = loc

	return &s
}

func Visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	var data []byte
	if !info.IsDir() {
		data, err = ioutil.ReadFile(path)
		if err != nil {
			return err
		}
	}
	return visitAux(path, info, data)
}

func visitAux(path string, info os.FileInfo, data []byte) error {
	// skip hidden files/directories
	// TODO: check in subdir
	if name := info.Name(); name[0] == '.' || name[0] == '~' || name[len(name)-1] == '~' {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	// skip non-go files
	if filepath.Ext(path) != ".go" {
		return nil
	}

	// skip test files
	re := regexp.MustCompile(".+\\_test\\.go")
	if re.MatchString(filepath.Base(path)) {
		return nil
	}

	if info.IsDir() {
		addPkg(path)
		return nil
	}

	// ensure that the package has already been created
	var pkg *src.Package
	var ok bool
	if pkg, ok = pkgs[filepath.Dir(path)]; !ok {
		pkg = addPkg(filepath.Dir(path))
	}

	if sf := runParser(path, data); sf != nil {
		pkg.SrcFiles = append(pkg.SrcFiles, sf)
		pkg.LoC += sf.LoC
	}

	return nil
}

func addPkg(path string) *src.Package {
	pkg := &src.Package{
		Name:     filepath.Base(path),
		Path:     path,
		Doc:      []string{},
		SrcFiles: make([]*src.SrcFile, 0),
	}

	pkgs[path] = pkg

	return pkg
}

func readTarArchive(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	tr := tar.NewReader(f)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		info := hdr.FileInfo()

		var data []byte
		if !info.IsDir() {
			data, err = ioutil.ReadAll(tr)
			if err != nil {
				return err
			}
		}
		adjustedPath := filepath.Join(strings.TrimSuffix(path, ".tar"), hdr.Name)
		if err := visitAux(adjustedPath, info, data); err != nil && err != filepath.SkipDir {
			return err
		}
	}

	return nil
}

func main() {
	flag.Parse()

	root := flag.Arg(0)

	pkgs[root] = &src.Package{
		Path:     root,
		SrcFiles: make([]*src.SrcFile, 0),
	}

	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(root); err != nil {
		log.Fatal(err)
	}

	// TODO: sexify
	var prjName string
	if fi.IsDir() {
		prjName = filepath.Base(root)
		err = filepath.Walk(root, Visit)
	} else {
		switch filepath.Ext(root) {
		case ".tar":
			prjName = strings.TrimSuffix(filepath.Base(root), ".tar")
			err = readTarArchive(root)
		}
	}
	if err != nil {
		log.Fatal(err)
	}

	pkgs[root].Name = prjName

	prj := src.Project{
		Name:     prjName,
		Packages: make([]*src.Package, 0),
		Langs:    []*src.Language{&goLang},
	}

	var loc int64
	for _, p := range pkgs {
		if len(p.SrcFiles) == 0 {
			continue
		}
		prj.Packages = append(prj.Packages, p)
		loc += p.LoC
	}

	prj.LoC = loc

	bs, err := json.Marshal(prj)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bs))
}
