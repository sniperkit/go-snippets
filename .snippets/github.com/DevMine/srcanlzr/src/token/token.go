// Copyright 2014-2015 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

// Literals
const (
	IntLit    = "INT"    // integer (1, 2, 3, ...)
	FloatLit  = "FLOAT"  // floatint point number (1.2)
	ImagLit   = "IMAG"   // imaginary number (complex number)
	CharLit   = "CHAR"   // character ('a', 'b', 'c', ...)
	StringLit = "STRING" // string ("foo")
	BoolLit   = "BOOl"   // boolean (true/false)
	NilLit    = "NIL"    // null value (nil, null, NULL, ...)
)

// Binary operators
const (
	// numerical operators
	ADD = "ADD" // addition
	SUB = "SUB" // subtraction
	MUL = "MUL" // multiplication
	QUO = "QUO" // division
	MOD = "MOD" // modulo

	// logical operators
	AND     = "AND"         // binary and (&)
	OR      = "OR"          // binary or (|)
	XOR     = "XOR"         // binary xor (^)
	SHL     = "SHIFT_LEFT"  // binary left shift <<
	SHR     = "SHIFT_RIGHT" // binary right shift >>
	AND_NOT = "AND_NOT"     // binary and not (&^)

	// comparison operators
	NEQ  = "NEQ"  // not equal
	LEQ  = "LEQ"  // less or equal
	GEQ  = "GEQ"  // greater or equal
	EQ   = "EQ"   // equal
	LSS  = "LSS"  // less
	GTR  = "GTR"  // greater
	LAND = "LAND" // and (&&)
	LOR  = "LOR"  // or (||)
)

// Kind of declarations
const (
	ConstDecl = "CONSTANT" // constant
	VarDecl   = "VAR"      // variable
)

// Increment/Decrement operators
const (
	INC = "INC" // increment operator (i++)
	DEC = "DEC" // decreement operator (i--)
)

// Unary operators
const (
	NOT  = "NOT"  // negation operator
	ADDR = "ADDR" // memory address (&foo)
	STAR = "STAR" // dereference operator (*foo)

	NEG = "NEG" // negative sign (-)
	POS = "POS" // positive sign (+)
)

// Statement names
const (
	IfStmtName        = "IF"         // if statement
	SwitchStmtName    = "SWITCH"     // switch statement
	LoopStmtName      = "LOOP"       // loop statement (for, while, do...while, ...)
	RangeLoopStmtName = "RANGE_LOOP" // range loops( for foo := range bar, for foo in bar, ...)
	AssignStmtName    = "ASSIGN"     // assign statement (foo = bar)
	DeclStmtName      = "DECL"       // declaration statement (int foo = 42, foo := 42, ...)
	ReturnStmtName    = "RETURN"     // return statement (return 42)
	ExprStmtName      = "EXPR"       // expression statement
	TryStmtName       = "TRY"        // try...catch statement
	ThrowStmtName     = "THROW"      // throw exception
	OtherStmtName     = "OTHER"      // any other not supported statement
)

// Expression names
const (
	// Expression names
	UnaryExprName           = "UNARY"            // unary expression (!foo)
	BinaryExprName          = "BINARY"           // binary expression (foo && bar)
	TernaryExprName         = "TERNARY"          // ternary expression (foo ? 42 : 0)
	IncDecExprName          = "INC_DEC"          // increment/decrement expression (i++/i--)
	CallExprName            = "CALL"             // call expression (foo(42))
	ConstructorCallExprName = "CONSTRUCTOR_CALL" // constructor call expression (foo = new Foo())
	ArrayExprName           = "ARRAY"            // array expression
	IndexExprName           = "INDEX"            // index expression (foo[i])

	// Literal names
	BasicLitName = "BASIC_LIT" // basic literal ("foo", 42, 12.34, ...)
	FuncLitName  = "FUNC_LIT"  // function literal (foo := func() { ... })
	ClassLitName = "CLASS_LIT" // class literal
	ArrayLitName = "ARRAY_LIT" // array literal (foo = [1,2,3])

	// Type names
	StructTypeName = "STRUCT_TYPE" // struct type

	// Other
	AttrRefName   = "ATTR_REF"   // attribute reference (this.foo)
	ValueSpecName = "VALUE_SPEC" // value specifier
	IdentName     = "IDENT"      // identifier
)

// Supported visiblities
const (
	PublicVisibility    = "public"    // public (visible from everywhere)
	PackageVisibility   = "package"   // package (visible from every entity inside the package)
	ProtectedVisibility = "protected" // protected (only visible from subclasses)
	PrivateVisibility   = "private"   // private
)

// list of all supported visibilities
var suppVisibility = []string{
	PublicVisibility,
	PackageVisibility,
	ProtectedVisibility,
	PrivateVisibility,
}

// Type names
const (
	TypeMapName         = "MAP"         // hash map
	TypeStructName      = "STRUCT"      // structure
	TypeArrayName       = "ARRAY"       // array
	TypeFuncName        = "FUNC"        // function
	TypeInterfaceName   = "INTERFACE"   // interface
	TypeUnsupportedName = "UNSUPPORTED" // unsupported type (for error)
)
