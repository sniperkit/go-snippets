package normalizer

import (
	"go/token"

	"gopkg.in/bblfsh/sdk.v2/uast"
	"gopkg.in/bblfsh/sdk.v2/uast/role"
	. "gopkg.in/bblfsh/sdk.v2/uast/transformer"
	"gopkg.in/bblfsh/sdk.v2/uast/transformer/positioner"
)

// Native is the of list `transformer.Transformer` to apply to a native AST.
// To learn more about the Transformers and the available ones take a look to:
// https://godoc.org/gopkg.in/bblfsh/sdk.v2/uast/transformer
var Native = Transformers([][]Transformer{
	// The main block of transformation rules.
	{Mappings(Annotations...)},
	{
		Mappings(
			AnnotateIfNoRoles("FieldList", role.Incomplete),
		),
		// RolesDedup is used to remove duplicate roles assigned by multiple
		// transformation rules.
		RolesDedup(),
	},
}...)

// Code is a special block of transformations that are applied at the end
// and can access original source code file. It can be used to improve or
// fix positional information.
//
// https://godoc.org/gopkg.in/bblfsh/sdk.v2/uast/transformer/positioner
var Code = []CodeTransformer{
	positioner.NewFillLineColFromOffset(),
}

// mapAST is a helper for describing a single AST transformation for a given node type.
func mapAST(typ string, ast, norm ObjectOp, roles ...role.Role) Mapping {
	return mapASTCustom(typ, ast, norm, nil, roles...)
}

func rolesByType(typ string) role.Roles {
	return typeRoles[typ]
}

func mapASTCustom(typ string, ast, norm ObjectOp, rop ArrayOp, roles ...role.Role) Mapping {
	return MapASTCustomType(typ, ast, norm, rolesByType, rop, roles...)
}

func annotateType(typ string, fields ObjAnnotator, roles ...role.Role) Mapping {
	return annotateTypeCustom(typ, fields, nil, roles...)
}

func annotateTypeCustom(typ string, fields ObjAnnotator, rop ArrayOp, roles ...role.Role) Mapping {
	return AnnotateTypeCustom(mapASTCustom, typ, fields, rop, roles...)
}

func operator(field, vr string, lookup map[uast.Value]ArrayOp, roles ...role.Role) Field {
	return Field{Name: field, Op: Operator(vr, lookup, roles...)}
}

var _ ObjAnnotator = astField{}

type astField struct {
	Inherit bool
	Roles   []role.Role
}

func (f astField) left(pref string) ObjectOp {
	return ASTObjectLeft("Field", Obj{
		"Names": EachObjectRoles(pref + "field_name"),
		"Type":  ObjectRoles(pref + "field_type"),
	})
}

func (f astField) right(pref string) ObjectOp {
	nameRoles := []role.Role{role.Name}
	typeRoles := []role.Role{role.Type}
	if f.Inherit {
		n := len(f.Roles)
		nameRoles = append(f.Roles[:n:n], nameRoles...)
		typeRoles = append(f.Roles[:n:n], typeRoles...)
	}
	return ASTObjectRightCustom("Field", Obj{
		"Names": EachObjectRoles(pref+"field_name", nameRoles...),
		"Type":  ObjectRoles(pref+"field_type", typeRoles...),
	}, rolesByType, nil, f.Roles...)
}

func (f astField) MappingParts(pref string) (src, dst ObjectOp) {
	return f.left(pref), f.right(pref)
}

var (
	literalRoles = TokenToRolesMap(map[token.Token][]role.Role{
		token.STRING: {role.String},
		token.CHAR:   {role.Character},
		token.INT:    {role.Number},
		token.FLOAT:  {role.Number},
		token.IMAG:   {role.Incomplete}, // TODO: IMAG
	})
	opRolesBinary = TokenToRolesMap(map[token.Token][]role.Role{
		token.ADD: {role.Arithmetic, role.Add},
		token.SUB: {role.Arithmetic, role.Substract},
		token.MUL: {role.Arithmetic, role.Multiply},
		token.QUO: {role.Arithmetic, role.Divide},
		token.REM: {role.Arithmetic, role.Modulo},

		token.XOR:     {role.Bitwise, role.Xor},
		token.OR:      {role.Bitwise, role.Or},
		token.AND:     {role.Bitwise, role.And},
		token.SHL:     {role.Bitwise, role.LeftShift},
		token.SHR:     {role.Bitwise, role.RightShift},
		token.AND_NOT: {role.Bitwise, role.And, role.Negative},

		token.LOR:  {role.Boolean, role.Or},
		token.LAND: {role.Boolean, role.And},

		token.EQL: {role.Relational, role.Equal},
		token.NEQ: {role.Relational, role.Not, role.Equal},
		token.LSS: {role.Relational, role.LessThan},
		token.GTR: {role.Relational, role.GreaterThan},
		token.LEQ: {role.Relational, role.LessThanOrEqual},
		token.GEQ: {role.Relational, role.GreaterThanOrEqual},
	})
	opRolesUnary = TokenToRolesMap(map[token.Token][]role.Role{
		token.ADD: {role.Arithmetic, role.Positive},
		token.SUB: {role.Arithmetic, role.Negative},

		token.MUL: {role.Dereference},
		token.AND: {role.TakeAddress},

		token.XOR: {role.Bitwise, role.Negative},
		token.NOT: {role.Boolean, role.Negative},

		token.ARROW: {role.Incomplete},
	})
	opIncDec = TokenToRolesMap(map[token.Token][]role.Role{
		token.INC: {role.Increment},
		token.DEC: {role.Decrement},
	})
	opAssign = TokenToRolesMap(map[token.Token][]role.Role{
		token.ASSIGN: {},
		token.DEFINE: {role.Declaration},

		token.ADD_ASSIGN: {role.Operator, role.Arithmetic, role.Add},
		token.SUB_ASSIGN: {role.Operator, role.Arithmetic, role.Substract},
		token.MUL_ASSIGN: {role.Operator, role.Arithmetic, role.Multiply},
		token.QUO_ASSIGN: {role.Operator, role.Arithmetic, role.Divide},
		token.REM_ASSIGN: {role.Operator, role.Arithmetic, role.Modulo},

		token.OR_ASSIGN:      {role.Operator, role.Bitwise, role.Or},
		token.AND_ASSIGN:     {role.Operator, role.Bitwise, role.And},
		token.XOR_ASSIGN:     {role.Operator, role.Bitwise, role.Xor},
		token.SHL_ASSIGN:     {role.Operator, role.Bitwise, role.LeftShift},
		token.SHR_ASSIGN:     {role.Operator, role.Bitwise, role.RightShift},
		token.AND_NOT_ASSIGN: {role.Operator, role.Bitwise, role.And, role.Negative},
	})
	branchRoles = TokenToRolesMap(map[token.Token][]role.Role{
		token.CONTINUE:    {role.Continue},
		token.BREAK:       {role.Break},
		token.GOTO:        {role.Goto},
		token.FALLTHROUGH: {role.Incomplete}, // TODO: fallthrough
	})
)

func goTok(tok token.Token) uast.Value {
	return uast.String(tok.String())
}

func TokenToRolesMap(m map[token.Token][]role.Role) map[uast.Value]ArrayOp {
	out := make(map[uast.Value]ArrayOp, len(m))
	for tok, roles := range m {
		out[goTok(tok)] = Roles(roles...)
	}
	return out
}

// Annotations is a list of individual transformations to annotate a native AST with roles.
var Annotations = []Mapping{

	annotateType("File", nil, role.File),

	annotateType("CommentGroup", nil, role.Comment, role.List),

	mapAST("Comment", Obj{
		"Text": UncommentCLike("text"),
	}, Obj{ // ->
		uast.KeyToken: Var("text"),
	}, role.Comment),

	annotateType("BadExpr", nil, role.Incomplete),

	annotateType("Ident", FieldRoles{
		"Name": {Rename: uast.KeyToken},
	}, role.Identifier),

	mapASTCustom("BasicLit", Obj{
		"Value": Var("val"),
		"Kind":  Var("kind"),
	}, Fields{ // ->
		{Name: "Kind", Op: Var("kind")},
		{Name: uast.KeyToken, Op: Var("val")},
	}, LookupArrOpVar("kind", literalRoles),
		role.Literal, role.Primitive),

	mapASTCustom("BinaryExpr", Obj{
		"X":  ObjectRoles("left"),
		"Y":  ObjectRoles("right"),
		"Op": Var("op"),
	}, Pre(Fields{ // ->
		operator("Op", "op", opRolesBinary, role.Binary),
	}, Obj{
		"X": ObjectRoles("left", role.Binary, role.Left),
		"Y": ObjectRoles("right", role.Binary, role.Right),
	}), LookupArrOpVar("op", opRolesBinary), role.Binary),

	mapASTCustom("UnaryExpr", Obj{
		"X":  Var("x"),
		"Op": Var("op"),
	}, Fields{ // ->
		operator("Op", "op", opRolesUnary, role.Unary),
		{Name: "X", Op: Var("x")},
	}, LookupArrOpVar("op", opRolesUnary), role.Unary),

	mapASTCustom("IncDecStmt", Obj{
		"X":   Var("x"),
		"Tok": Var("op"),
	}, Fields{ // ->
		operator("Op", "op", opIncDec, role.Unary),
		{Name: "X", Op: Var("x")},
	}, LookupArrOpVar("op", opIncDec), role.Unary),

	annotateType("BlockStmt", nil, role.Block, role.Scope),

	mapASTCustom("AssignStmt", Obj{
		"Lhs": EachObjectRoles("left"),
		"Rhs": EachObjectRoles("right"),
		"Tok": Var("op"),
	}, Fields{ // ->
		operator("Op", "op", opAssign, role.Operator, role.Assignment, role.Binary),
		{Name: "Lhs", Op: EachObjectRoles("left", role.Assignment, role.Binary, role.Left)},
		{Name: "Rhs", Op: EachObjectRoles("right", role.Assignment, role.Binary, role.Right)},
	}, LookupArrOpVar("op", opAssign),
		role.Assignment, role.Binary),

	annotateType("IfStmt", FieldRoles{
		"Init": {Opt: true, Roles: role.Roles{role.If, role.Initialization}},
		"Cond": {Roles: role.Roles{role.If, role.Condition}},
		"Body": {Roles: role.Roles{role.Then, role.Body}},
		"Else": {Opt: true, Roles: role.Roles{role.Else}},
	}, role.If),

	annotateType("SwitchStmt", FieldRoles{
		"Init": {Opt: true, Roles: role.Roles{role.Switch, role.Initialization}},
		"Body": {Roles: role.Roles{role.Switch, role.Body}},
	}, role.Switch),

	annotateType("TypeSwitchStmt", FieldRoles{
		"Init": {Opt: true, Roles: role.Roles{role.Switch, role.Initialization}},
		"Body": {Roles: role.Roles{role.Switch, role.Body}},
	}, role.Switch, role.Incomplete),

	annotateType("SelectStmt", ObjRoles{
		"Body": {role.Switch, role.Body},
	}, role.Switch, role.Incomplete),

	annotateType("ForStmt", FieldRoles{
		"Init": {Opt: true, Roles: role.Roles{role.For, role.Initialization}},
		"Cond": {Opt: true, Roles: role.Roles{role.For, role.Condition}},
		"Body": {Roles: role.Roles{role.For, role.Body}},
		"Post": {Opt: true, Roles: role.Roles{role.For, role.Update}},
	}, role.For),

	annotateType("RangeStmt", FieldRoles{
		"Key":   {Opt: true, Roles: role.Roles{role.For, role.Iterator, role.Key}},
		"Value": {Opt: true, Roles: role.Roles{role.For, role.Iterator, role.Value}},
		"X":     {},
		"Body":  {Roles: role.Roles{role.For, role.Body}},
	}, role.For, role.Iterator),

	annotateTypeCustom("BranchStmt", FieldRoles{
		"Tok":   {Op: Var("tok")},
		"Label": {},
	}, LookupArrOpVar("tok", branchRoles)),

	annotateType("ImportSpec", FieldRoles{
		"Name": {Opt: true, Roles: role.Roles{role.Import, role.Alias}},
		"Path": {Roles: role.Roles{role.Import, role.Pathname}},
	}, role.Import, role.Declaration),

	annotateType("ValueSpec", ObjRoles{
		"Type": {role.Type},
	}, role.Declaration),

	annotateType("TypeSpec", ObjRoles{
		"Type": {role.Type},
	}, role.Declaration),

	annotateType("ArrayType", ObjRoles{
		"Elt": {role.Entry},
	}, role.Type, role.List),

	annotateType("MapType", ObjRoles{
		"Key":   {role.Key},
		"Value": {role.Entry},
	}, role.Type, role.Map),

	annotateType("FuncLit", ObjRoles{
		"Type": {role.Type},
		"Body": {role.Body},
	}),

	annotateType("StructType", FieldRoles{
		"Fields": {Sub: FieldRoles{
			"List": {Arr: true, Sub: astField{
				Roles: role.Roles{role.Entry},
			}},
		}},
	}, role.Type),

	annotateType("InterfaceType", FieldRoles{
		"Methods": {Sub: FieldRoles{
			"List": {Arr: true,
				Sub: astField{Roles: role.Roles{role.Entry}},
			},
		}, Roles: role.Roles{role.Function, role.List}},
	}, role.Type, role.Incomplete),

	annotateType("FuncType", FieldRoles{
		"Params": {Sub: FieldRoles{
			"List": {Arr: true,
				Sub: astField{Roles: role.Roles{role.Argument}},
			},
		}, Roles: role.Roles{role.ArgsList}},
		"Results": {Opt: true, Sub: FieldRoles{
			"List": {Arr: true,
				Sub: astField{Roles: role.Roles{role.Return, role.Argument}},
			},
		}, Roles: role.Roles{role.Return, role.ArgsList}},
	}, role.Function, role.Type),

	annotateType("FuncDecl", FieldRoles{
		"Recv": {Opt: true, Sub: FieldRoles{
			"List": {Arr: true, Sub: astField{
				Inherit: true, Roles: role.Roles{role.Function, role.Receiver},
			}},
		}, Roles: role.Roles{role.Function, role.Receiver, role.List}},
		"Name": {Roles: role.Roles{role.Function, role.Name}},
		"Type": {Roles: role.Roles{role.Function, role.Type}},
		"Body": {Roles: role.Roles{role.Function, role.Body}},
	}, role.Function, role.Declaration),

	annotateType("GenDecl", FieldRoles{
		"Tok": {Op: Is(goTok(token.VAR))},
		"Specs": {Arr: true, Sub: FieldRoles{
			"Names": {Arr: true, Roles: role.Roles{role.Variable, role.Name}},
		}},
	}, role.Variable, role.Declaration),

	annotateType("GenDecl", FieldRoles{
		"Tok": {Op: Is(goTok(token.CONST))},
		"Specs": {Arr: true, Sub: FieldRoles{
			"Names": {Arr: true, Roles: role.Roles{role.Name}},
		}},
	}, role.Incomplete, role.Declaration),

	annotateType("GenDecl", FieldRoles{
		"Tok": {Op: Is(goTok(token.TYPE))},
		"Specs": {Arr: true, Sub: FieldRoles{
			"Name": {Roles: role.Roles{role.Type, role.Name}},
		}},
	}, role.Type, role.Declaration),

	annotateType("GenDecl", FieldRoles{
		"Tok": {Op: Is(goTok(token.IMPORT))},
	}, role.Declaration),

	annotateType("CallExpr", FieldRoles{
		"Fun":  {Roles: role.Roles{role.Callee}},
		"Args": {Arr: true, Roles: role.Roles{role.Argument, role.Positional}},
	}, role.Call),

	annotateType("KeyValueExpr", ObjRoles{
		"Key":   {role.Key},
		"Value": {role.Value},
	}),

	annotateType("CaseClause", nil, role.Case),
	annotateType("CommClause", nil, role.Case, role.Incomplete),
	// TODO: default (length of List is zero)

	annotateType("ReturnStmt", nil, role.Return),
	annotateType("GoStmt", nil, role.Incomplete),       // TODO: Async role
	annotateType("SelectorExpr", nil, role.Incomplete), // TODO: new role

	annotateType("CompositeLit", nil, role.Literal),
	annotateType("ChanType", nil, role.Type, role.Incomplete),

	annotateType("ExprStmt", nil),
	annotateType("DeclStmt", nil),
	annotateType("DeferStmt", nil, role.Incomplete),
	annotateType("SendStmt", nil, role.Incomplete),
	annotateType("LabeledStmt", nil),
	annotateType("Ellipsis", nil),
	annotateType("SliceExpr", nil),
	annotateType("IndexExpr", nil),
	annotateType("StarExpr", nil),
	annotateType("TypeAssertExpr", nil, role.Incomplete),
	annotateType("ParenExpr", nil, role.Incomplete),
}
