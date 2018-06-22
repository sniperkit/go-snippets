package normalizer

import (
	"gopkg.in/bblfsh/sdk.v2/uast"
	"gopkg.in/bblfsh/sdk.v2/uast/role"
	. "gopkg.in/bblfsh/sdk.v2/uast/transformer"
	"gopkg.in/bblfsh/sdk.v2/uast/transformer/positioner"
)

// Native is the of list `transformer.Transformer` to apply to a native AST.
// To learn more about the Transformers and the available ones take a look to:
// https://godoc.org/gopkg.in/bblfsh/sdk.v2/uast/transformer
var Native = Transformers([][]Transformer{
	{
		// ResponseMetadata is a transform that trims response metadata from AST.
		//
		// https://godoc.org/gopkg.in/bblfsh/sdk.v2/uast#ResponseMetadata
		ResponseMetadata{
			TopLevelIsRootNode: true,
		},
		Mappings(
			ASTMap("remove loc",
				Part("_", Obj{"loc": AnyNode(Obj{})}),
				Part("_", Obj{}),
			),
			ASTMap("remove extra",
				Part("_", Obj{"extra": AnyNode(Obj{})}),
				Part("_", Obj{}),
			),
		),
	},
	// The main block of transformation rules.
	{Mappings(Annotations...)},
	{
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

func mapASTCustom(typ string, ast, norm ObjectOp, rop ArrayOp, roles ...role.Role) Mapping {
	return ASTMap(typ,
		ASTObjectLeft(typ, ast),
		ASTObjectRight(typ, norm, rop, roles...),
	)
}

var (
	unaryRoles = StringToRolesMap(map[string][]role.Role{
		// Unary operators
		"-":      {role.Arithmetic, role.Negative},
		"+":      {role.Arithmetic, role.Positive},
		"!":      {role.Boolean, role.Not},
		"~":      {role.Bitwise, role.Not},
		"typeof": {role.Type},
		"void":   {role.Null},
		"delete": {role.Incomplete},
		"throw":  {role.Throw},

		// Update operators
		"++": {role.Arithmetic, role.Increment},
		"--": {role.Arithmetic, role.Decrement},
	})
	binaryRoles = StringToRolesMap(map[string][]role.Role{
		"in":         {},
		"==":         {role.Relational, role.Equal},
		"!=":         {role.Relational, role.Equal, role.Not},
		"===":        {role.Relational, role.Identical},
		"!==":        {role.Relational, role.Identical, role.Not},
		"<":          {role.Relational, role.LessThan},
		"<=":         {role.Relational, role.LessThanOrEqual},
		">":          {role.Relational, role.GreaterThan},
		">=":         {role.Relational, role.GreaterThanOrEqual},
		"<<":         {role.Bitwise, role.LeftShift},
		">>":         {role.Bitwise, role.RightShift},
		">>>":        {role.Bitwise, role.RightShift, role.Unsigned},
		"+":          {role.Arithmetic, role.Add},
		"-":          {role.Arithmetic, role.Substract},
		"*":          {role.Arithmetic, role.Multiply},
		"/":          {role.Arithmetic, role.Divide},
		"%":          {role.Arithmetic, role.Modulo},
		"|":          {role.Bitwise, role.Or},
		"^":          {role.Bitwise, role.Xor},
		"&":          {role.Bitwise, role.And},
		"instanceof": {role.Type},
		"|>":         {role.Incomplete},
	})
	assignRoles = StringToRolesMap(map[string][]role.Role{
		"=":    {},
		"+=":   {role.Arithmetic, role.Add},
		"-=":   {role.Arithmetic, role.Substract},
		"*=":   {role.Arithmetic, role.Multiply},
		"/=":   {role.Arithmetic, role.Divide},
		"%=":   {role.Arithmetic, role.Modulo},
		"<<=":  {role.Bitwise, role.LeftShift},
		">>=":  {role.Bitwise, role.RightShift},
		">>>=": {role.Bitwise, role.RightShift, role.Unsigned},
		"|=":   {role.Bitwise, role.Or},
		"^=":   {role.Bitwise, role.Xor},
		"&=":   {role.Bitwise, role.And},
	})
	logicalRoles = StringToRolesMap(map[string][]role.Role{
		"||": {role.Or},
		"&&": {role.And},
		"??": {role.Incomplete},
	})
)

func literal(typ string, roles ...role.Role) Mapping {
	return mapAST(typ, Obj{
		"value": Var("val"),
	}, Obj{
		uast.KeyToken: Var("val"),
	}, roles...)
}

func function(typ string) Mapping {
	return mapAST(typ, Obj{
		"id":     OptObjectRoles("id"),
		"body":   ObjectRoles("body"),
		"params": EachObjectRolesByType("param", nil),
	}, Obj{
		"id":   OptObjectRoles("id", role.Function, role.Name),
		"body": ObjectRoles("body", role.Function, role.Body),
		// TODO: AnnotateType("RestElement", nil, role.ArgsList),
		"params": EachObjectRolesByType("param", map[string][]role.Role{
			"RestElement": {role.ArgsList},
			"":            {},
		}, role.Function, role.Argument),
	}, role.Declaration, role.Function)
}

// Annotations is a list of individual transformations to annotate a native AST with roles.
var Annotations = []Mapping{
	// ObjectToNode defines how to normalize common fields of native AST
	// (like node type, token, positional information).
	//
	// https://godoc.org/gopkg.in/bblfsh/sdk.v2/uast#ObjectToNode
	ObjectToNode{
		InternalTypeKey: "type",
		OffsetKey:       "start",
		EndOffsetKey:    "end",
	}.Mapping(),

	AnnotateType("File", nil, role.File),
	AnnotateType("Program", nil, role.Module),

	// Comments
	mapAST("CommentLine", Obj{
		"value": UncommentCLike("text"),
	}, Obj{
		uast.KeyToken: Var("text"),
	}, role.Comment),

	mapAST("CommentBlock", Obj{
		"value": UncommentCLike("text"),
	}, Obj{
		uast.KeyToken: Var("text"),
	}, role.Comment, role.Block),

	// Identifiers
	AnnotateType("Identifier",
		FieldRoles{
			"name": {Rename: uast.KeyToken},
		},
		role.Expression, role.Identifier,
	),
	AnnotateType("PrivateName", nil, role.Expression, role.Identifier, role.Qualified, role.Visibility, role.Instance),

	// Literals
	AnnotateType("RegExpLiteral",
		FieldRoles{
			"pattern": {Rename: uast.KeyToken},
		},
		role.Expression, role.Literal, role.Regexp,
	),
	AnnotateType("NullLiteral", nil, role.Expression, role.Literal, role.Null),
	literal("StringLiteral", role.Expression, role.Literal, role.String),
	literal("BooleanLiteral", role.Expression, role.Literal, role.Boolean),
	literal("NumericLiteral", role.Expression, role.Literal, role.Number),

	// Functions
	function("FunctionDeclaration"),
	function("ArrowFunctionExpression"),
	function("FunctionExpression"),
	function("ObjectMethod"),
	function("ClassMethod"),
	function("ClassPrivateMethod"),
	function("OptFunctionDeclaration"),

	// Statements
	AnnotateType("ExpressionStatement", nil, role.Statement),
	AnnotateType("BlockStatement", nil, role.Statement, role.Block, role.Scope),
	AnnotateType("EmptyStatement", nil, role.Statement),
	AnnotateType("DebuggerStatement", nil, role.Statement, role.Incomplete),

	AnnotateType("WithStatement",
		ObjRoles{
			"object": {role.Incomplete},
		},
		role.Statement, role.Scope, role.Block, role.Incomplete,
	),

	// Control flow
	AnnotateType("ReturnStatement", nil, role.Statement, role.Return),
	AnnotateType("LabeledStatement", nil, role.Statement, role.Incomplete),
	AnnotateType("BreakStatement", nil, role.Statement, role.Break),
	AnnotateType("ContinueStatement", nil, role.Statement, role.Continue),

	// Choice
	AnnotateType("IfStatement",
		ObjRoles{
			"test":       {role.If, role.Condition},
			"consequent": {role.If, role.Then, role.Body},
			"alternate":  {role.If, role.Else, role.Body},
		},
		role.If,
	),
	AnnotateType("ConditionalExpression",
		ObjRoles{
			"test":       {role.If, role.Condition},
			"consequent": {role.If, role.Then, role.Body},
			"alternate":  {role.If, role.Else, role.Body},
		},
		role.If,
	),
	AnnotateType("IfStatement", nil, role.Statement),
	AnnotateType("SwitchStatement",
		ObjRoles{
			"discriminant": {role.Switch, role.Condition},
		},
		role.Statement, role.Switch,
	),
	AnnotateType("SwitchCase",
		ObjRoles{
			"test": {role.Case, role.Condition},
		},
		role.Switch, role.Case,
	),

	// Exception
	AnnotateType("ThrowStatement", nil, role.Statement, role.Throw),
	AnnotateType("TryStatement",
		ObjRoles{
			"finalizer": {role.Try, role.Finally},
		},
		role.Statement, role.Try,
	),
	AnnotateType("CatchClause", nil, role.Try, role.Catch),

	// Loops
	AnnotateType("WhileStatement",
		ObjRoles{
			"test": {role.While, role.Condition},
			"body": {role.While, role.Body},
		},
		role.Statement, role.While,
	),
	AnnotateType("DoWhileStatement",
		ObjRoles{
			"test": {role.DoWhile, role.Condition},
			"body": {role.DoWhile, role.Body},
		},
		role.Statement, role.DoWhile,
	),
	AnnotateType("ForStatement",
		ObjRoles{
			"init":   {role.For, role.Initialization},
			"test":   {role.For, role.Condition},
			"update": {role.For, role.Update},
		},
		role.Statement, role.For,
	),
	AnnotateType("ForInStatement",
		ObjRoles{
			"left":  {role.For, role.Iterator},
			"right": {role.For},
			"body":  {role.For, role.Body},
		},
		role.Statement, role.For, role.Iterator,
	),
	AnnotateType("ForOfStatement",
		ObjRoles{
			"left":  {role.For, role.Iterator},
			"right": {role.For},
			"body":  {role.For, role.Body},
		},
		role.Statement, role.For, role.Iterator,
	),

	// Declarations
	AnnotateType("FunctionDeclaration", nil, role.Statement),
	AnnotateType("VariableDeclaration", nil, role.Statement, role.Declaration, role.Variable),
	AnnotateType("VariableDeclarator",
		ObjRoles{
			"init": {role.Initialization},
		},
		role.Declaration, role.Variable,
	),

	// Misc
	AnnotateType("Decorator", nil, role.Incomplete),
	AnnotateType("Directive", nil, role.Incomplete),
	AnnotateType("DirectiveLiteral",
		FieldRoles{
			"value": {Rename: uast.KeyToken},
		},
		role.Expression, role.Literal, role.Incomplete,
	),

	// Expressions
	AnnotateType("Super", nil, role.Expression, role.Identifier, role.Base),
	AnnotateType("Import", nil, role.Expression, role.Import),
	AnnotateType("ThisExpression", nil, role.Expression, role.This),
	AnnotateType("ArrowFunctionExpression", nil, role.Expression),
	AnnotateType("YieldExpression", nil, role.Expression, role.Return, role.Incomplete),
	AnnotateType("AwaitExpression", nil, role.Expression, role.Incomplete),
	AnnotateType("ArrayExpression", nil, role.Expression, role.Initialization, role.List, role.Literal),
	AnnotateType("ObjectExpression", nil, role.Expression, role.Initialization, role.Map, role.Literal),
	AnnotateType("SpreadElement", nil, role.Incomplete),
	AnnotateType("MemberExpression", nil, role.Qualified, role.Expression, role.Identifier),
	AnnotateType("BindExpression", nil, role.Expression, role.Incomplete),
	AnnotateType("ConditionalExpression", nil, role.Expression),
	AnnotateType("NewExpression", nil, role.Expression, role.Incomplete),
	AnnotateType("SequenceExpression", nil, role.Expression, role.List),
	AnnotateType("DoExpression",
		ObjRoles{
			"body": {role.Body},
		},
		role.Expression, role.Incomplete,
	),

	// Object properties
	AnnotateType("ObjectMethod",
		ObjRoles{
			"key":  {role.Map, role.Key, role.Function, role.Name},
			"body": {role.Map, role.Value},
		},
		role.Map,
	),
	AnnotateType("ObjectProperty",
		ObjRoles{
			"key":   {role.Map, role.Key},
			"value": {role.Map, role.Value},
		},
		role.Map,
	),

	// Function expressions
	AnnotateType("FunctionExpression", nil, role.Expression),
	mapAST("CallExpression", Obj{
		"callee":    ObjectRoles("callee"),
		"arguments": EachObjectRolesByType("argument", nil),
	}, Obj{
		"callee": ObjectRoles("callee", role.Call, role.Callee),
		"arguments": EachObjectRolesByType("argument", map[string][]role.Role{
			"SpreadElement": {role.ArgsList},
			"":              {},
		}, role.Call, role.Argument),
	}, role.Expression, role.Call),

	// Unary operations
	mapASTCustom("UnaryExpression", Obj{
		"operator": Var("op"),
	}, Fields{ // ->
		//{Name:"prefix", Op:  String("true")},
		{Name: "operator", Op: Operator("op", unaryRoles, role.Unary)},
	}, LookupArrOpVar("op", unaryRoles), role.Expression, role.Unary, role.Operator),

	mapASTCustom("UpdateExpression", Obj{
		"operator": Var("op"),
	}, Fields{ // ->
		//{Name:"prefix", Op:  String("true")},
		{Name: "operator", Op: Operator("op", unaryRoles, role.Unary)},
	}, LookupArrOpVar("op", unaryRoles), role.Expression, role.Unary, role.Operator),

	AnnotateType("UnaryExpression",
		FieldRoles{
			"prefix": {Op: Bool(false)},
		},
		role.Postfix,
	),

	AnnotateType("UpdateExpression",
		FieldRoles{
			"prefix": {Op: Bool(false)},
		},
		role.Postfix,
	),

	// Binary operations
	mapASTCustom("BinaryExpression", Obj{
		"operator": Var("op"),
		"left":     ObjectRoles("left"),
		"right":    ObjectRoles("right"),
	}, Fields{ // ->
		{Name: "operator", Op: Operator("op", binaryRoles, role.Binary)},
		{Name: "left", Op: ObjectRoles("left", role.Binary, role.Left)},
		{Name: "right", Op: ObjectRoles("right", role.Binary, role.Right)},
	}, LookupArrOpVar("op", binaryRoles), role.Expression, role.Operator, role.Binary),

	mapASTCustom("AssignmentExpression", Obj{
		"operator": Var("op"),
		"left":     ObjectRoles("left"),
		"right":    ObjectRoles("right"),
	}, Fields{ // ->
		{Name: "operator", Op: Operator("op", assignRoles, role.Assignment, role.Binary)},
		{Name: "left", Op: ObjectRoles("left", role.Assignment, role.Binary, role.Left)},
		{Name: "right", Op: ObjectRoles("right", role.Assignment, role.Binary, role.Right)},
	}, LookupArrOpVar("op", assignRoles), role.Expression, role.Assignment, role.Operator, role.Binary),

	mapASTCustom("LogicalExpression", Obj{
		"operator": Var("op"),
		"left":     ObjectRoles("left"),
		"right":    ObjectRoles("right"),
	}, Fields{ // ->
		{Name: "operator", Op: Operator("op", logicalRoles, role.Boolean, role.Binary)},
		{Name: "left", Op: ObjectRoles("left", role.Boolean, role.Binary, role.Left)},
		{Name: "right", Op: ObjectRoles("right", role.Boolean, role.Binary, role.Right)},
	}, LookupArrOpVar("op", logicalRoles), role.Boolean, role.Expression, role.Operator, role.Binary),

	// Template literals
	AnnotateType("TemplateLiteral", nil, role.Expression, role.Literal, role.Incomplete),
	AnnotateType("TaggedTemplateExpression", nil, role.Expression, role.Literal, role.Incomplete),
	AnnotateType("TemplateElement",
		FieldRoles{
			"value": {Skip: true}, // drop value field
		},
		role.Literal, role.String, role.Incomplete,
	),

	// Patterns
	AnnotateType("ObjectPattern", nil, role.Incomplete),
	AnnotateType("ArrayPattern", nil, role.Incomplete),
	AnnotateType("RestElement", nil, role.Incomplete),
	AnnotateType("AssignmentPattern", nil, role.Assignment, role.Incomplete),

	// Classes
	AnnotateType("ClassBody", nil, role.Type, role.Body),
	AnnotateType("ClassDeclaration",
		ObjRoles{
			"id":         {role.Type, role.Name},
			"superClass": {role.Type, role.Base},
		},
		role.Declaration, role.Type,
	),
	AnnotateType("ClassExpression",
		ObjRoles{
			"id":         {role.Type, role.Name},
			"superClass": {role.Type, role.Base},
		},
		role.Declaration, role.Type,
	),
	AnnotateType("OptClasDeclaration",
		ObjRoles{
			"id":         {role.Type, role.Name},
			"superClass": {role.Type, role.Base},
		},
		role.Declaration, role.Type,
	),
	AnnotateType("ClassDeclaration", nil, role.Statement),
	AnnotateType("ClassExpression", nil, role.Expression),
	AnnotateType("ClassMethod",
		ObjRoles{
			"key":  {role.Key, role.Name},
			"body": {role.Value},
		},
		role.Statement,
	),
	AnnotateType("ClassPrivateMethod",
		ObjRoles{
			"key":  {role.Key, role.Name},
			"body": {role.Value},
		},
		role.Statement,
	),
	AnnotateType("ClassProperty",
		ObjRoles{
			"key":   {role.Key, role.Name},
			"value": {role.Value, role.Initialization},
		},
		role.Variable,
	),
	AnnotateType("ClassPrivateProperty",
		ObjRoles{
			"key":   {role.Key, role.Name},
			"value": {role.Value, role.Initialization},
		},
		role.Variable,
	),

	AnnotateType("MetaProperty", nil, role.Expression, role.Incomplete),

	// Modules
	AnnotateType("ImportDeclaration",
		FieldRoles{
			"specifiers": {Arr: true, Roles: role.Roles{role.Import}},
			"source":     {Roles: role.Roles{role.Import, role.Pathname}},
		},
		role.Statement, role.Declaration, role.Import,
	),
	AnnotateType("ImportSpecifier",
		ObjRoles{
			"local": {role.Import},
		},
		role.Import,
	),
	AnnotateType("ImportDefaultSpecifier",
		ObjRoles{
			"local": {role.Import},
		},
		role.Import,
	),
	AnnotateType("ImportNamespaceSpecifier",
		ObjRoles{
			"local": {role.Import},
		},
		role.Import,
	),
	AnnotateType("ImportSpecifier",
		ObjRoles{
			"imported": {role.Import},
		},
	),
	AnnotateType("ExportNamedDeclaration", nil, role.Statement, role.Declaration, role.Visibility, role.Module, role.Incomplete),
	AnnotateType("ExportDefaultDeclaration", nil, role.Statement, role.Declaration, role.Visibility, role.Module, role.Incomplete),
	AnnotateType("ExportAllDeclaration", nil, role.Statement, role.Declaration, role.Visibility, role.Module, role.Incomplete),
	AnnotateType("ExportNamedDeclaration",
		FieldRoles{
			"declaration": {Roles: role.Roles{role.Incomplete}},
			"specifiers":  {Arr: true, Roles: role.Roles{role.Incomplete}},
			"source":      {Opt: true, Roles: role.Roles{role.Pathname, role.Incomplete}},
		},
		role.Statement, role.Declaration, role.Visibility, role.Module, role.Incomplete,
	),
	AnnotateType("ExportSpecifier",
		ObjRoles{
			"local":    {role.Incomplete},
			"exported": {role.Incomplete},
		},
		role.Incomplete,
	),
	AnnotateType("OptFunctionDeclaration",
		ObjRoles{
			"id": {role.Name, role.Incomplete},
		},
		role.Statement, role.Incomplete,
	),
	AnnotateType("OptClasDeclaration",
		ObjRoles{
			"id": {role.Name, role.Incomplete},
		},
		role.Statement, role.Incomplete,
	),
}
