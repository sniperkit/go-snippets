// Copyright 2014-2015 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

type ArrayExpr struct {
	ExprName string     `json:"expression_name"`
	Type     *ArrayType `json:"type"`
}

type ArrayLit struct {
	ExprName string     `json:"expression_name"`
	Type     *ArrayType `json:"type"`
	Elts     []Expr     `json:"elements"`
}

type ArrayType struct {
	// Dimensions
	Dims []int64 `json:"dimensions"`

	Elt Expr `json:"element_type,omitempty"` // element type
}

type AssignStmt struct {
	StmtName string `json:"statement_name"`
	LHS      []Expr `json:"left_hand_side"`
	RHS      []Expr `json:"right_hand_side"`
	Line     int64  `json:"line"`
}

type Attr struct {
	Var
	Constant bool `json:"constant"`
	Static   bool `json:"static"`
}

type AttrRef struct {
	ExprName string `json:"expression_name"`
	Name     *Ident `json:"name"`
}

type BasicLit struct {
	ExprName string `json:"expression_name"`
	Kind     string `json:"kind"`
	Value    string `json:"value"`
}

type BinaryExpr struct {
	ExprName  string `json:"expression_name"`
	LeftExpr  Expr   `json:"left_expression,omitempty"`  // left operand
	Op        string `json:"operator"`                   // operator
	RightExpr Expr   `json:"right_expression,omitempty"` // right operand
}

type CallExpr struct {
	ExprName string   `json:"expression_name"`
	Fun      *FuncRef `json:"function"`  // Reference to the function
	Args     []Expr   `json:"arguments"` // function arguments
	Line     int64    `json:"line"`      // line number
}

type ClassDecl struct {
	Doc                   []string           `json:"doc,omitempty"`
	Name                  string             `json:"name"`
	Visibility            string             `json:"visibility"`
	ExtendedClasses       []*ClassRef        `json:"extended_classes,omitempty"`
	ImplementedInterfaces []*InterfaceRef    `json:"implemented_interfaces,omitempty"`
	Attrs                 []*Attr            `json:"attributes,omitempty"`
	Constructors          []*ConstructorDecl `json:"constructors,omitempty"`
	Destructors           []*DestructorDecl  `json:"destructors,omitempty"`
	Methods               []*MethodDecl      `json:"methods,omitempty"`
	NestedClasses         []*ClassDecl       `json:"nested_classes,omitempty"`
	Mixins                []*TraitRef        `json:"mixins,omitempty"`
}

type ClassLit struct {
	ExprName              string             `json:"expression_name"`
	ExtendedClasses       []*ClassRef        `json:"extended_classes,omitempty"`
	ImplementedInterfaces []*InterfaceRef    `json:"implemented_interfaces,omitempty"`
	Attrs                 []*Attr            `json:"attributes,omitempty"`
	Constructors          []*ConstructorDecl `json:"constructors,omitempty"`
	Destructors           []*DestructorDecl  `json:"destructors,omitempty"`
	Methods               []*MethodDecl      `json:"methods,omitempty"`
}

type ClassRef struct {
	Namespace string `json:"namespace"`
	ClassName string `json:"class_name"`
}

type Constant struct {
	Doc        []string `json:"doc"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`  // TODO rename into TypeName or use a type Type
	Value      string   `json:"value"` // TODO use an Expr instead of string value
	IsPointer  bool     `json:"is_pointer"`
	Visibility string   `json:"visibility,omitempty"`
}

type ConstructorCallExpr struct {
	CallExpr
}

type ConstructorDecl struct {
	Doc        []string `json:"doc,omitempty"`
	Name       string   `json:"name"`
	Params     []*Field `json:"parameters,omitempty"`
	Body       []Stmt   `json:"body,omitempty"`
	Visibility string   `json:"visibility"`
	LoC        int64    `json:"loc"`
}

type DeclStmt struct {
	AssignStmt
	Kind string `json:"kind"`
}

type DestructorDecl struct {
	ConstructorDecl
}

type EnumDecl struct {
	Doc                   []string           `json:"doc,omitempty"`
	Name                  string             `json:"name"`
	Visibility            string             `json:"visibility"`
	ImplementedInterfaces []*InterfaceRef    `json:"implemented_interfaces,omitempty"`
	EnumConstants         []*Ident           `json:"enum_constants,omitempty"`
	Attrs                 []*Attr            `json:"attributes,omitempty"`
	Constructors          []*ConstructorDecl `json:"constructors,omitempty"`
	Destructors           []*DestructorDecl  `json:"destructors,omitempty"`
	Methods               []*MethodDecl      `json:"methods,omitempty"`
}

type Expr interface{}

type ExprStmt struct {
	StmtName string `json:"statement_name"`
	X        Expr   `json:"expression"` // expression
}

type FuncDecl struct {
	Doc        []string  `json:"doc,omitempty"`
	Name       string    `json:"name"`
	Type       *FuncType `json:"type"`
	Body       []Stmt    `json:"body,omitempty"`
	Visibility string    `json:"visibility"`
	LoC        int64     `json:"loc"` // Lines of Code
}

type FuncLit struct {
	ExprName string    `json:"expression_name"`
	Type     *FuncType `json:"type"`
	Body     []Stmt    `json:"body,omitempty"`
	LoC      int64     `json:"loc"` // Lines of Code
}

type FuncRef struct {
	Namespace string `json:"namespace"`
	FuncName  string `json:"function_name"`
}

type FuncType struct {
	Params  []*Field `json:"parameters,omitempty"`
	Results []*Field `json:"results,omitempty"`
}

// GlobalDecl represents any declaration (var, const, type) declared outside of
// a function, class, trait, etc.
type GlobalDecl struct {
	Doc        []string `json:"doc,omitempty"`   // associated documentation; or nil
	Name       *Ident   `json:"name"`            // name of the var, const, or type
	Value      Expr     `json:"value,omitempty"` // default value; or nil
	Type       *Ident   `json:"type,omitempty"`  // type identifier; or nil
	Visibility string   `json:"visibility"`      // visibility (see the constants for the list of supported visibilities)
}

type Ident struct {
	ExprName string `json:"expression_name"`
	Name     string `json:"name"`
}

type IfStmt struct {
	StmtName string `json:"statement_name"`
	Init     Stmt   `json:"initialization,omitempty"`
	Cond     Expr   `json:"condition"`
	Body     []Stmt `json:"body"`
	Else     []Stmt `json:"else,omitempty"`
	Line     int64  `json:"line"` // Line number of the statement relatively to the function.
}

type IncDecExpr struct {
	ExprName string `json:"expression_name"`
	X        Expr   `json:"operand"`
	Op       string `json:"operator"` // INC or DEC
	IsPre    bool   `json:"is_pre"`   // pre = ++i, not pre = i++
}

type IndexExpr struct {
	ExprName string `json:"expression_name"`
	X        Expr   `json:"expression,omitempty"` // expression
	Index    Expr   `json:"index,omitempty"`      // index expression
}

type Interface struct {
	Doc                   []string        `json:"doc,omitempty"`
	Name                  string          `json:"name"`
	ImplementedInterfaces []*InterfaceRef `json:"implemented_interfaces,omitempty"`
	Protos                []*ProtoDecl    `json:"prototypes"`
	Visibility            string          `json:"visibility"`
}

type InterfaceRef struct {
	Namespace     string `json:"namespace"`
	InterfaceName string `json:"interface_name"`
}

type ListLit struct {
	Type *ListType `json:"type"`
	Elts []Expr    `json:"elements"`
}

type ListType struct {
	Len int64 `json:"length,omitempty"`
	Max int64 `json:"capacity,omitempty"` // maximum capacity
	Elt Expr  `json:"element_type"`
}

type LoopStmt struct {
	StmtName   string `json:"statement_name"`
	Init       []Stmt `json:"initialization,omitempty"`
	Cond       Expr   `json:"condition,omitempty"`
	Post       []Stmt `json:"post_iteration_statement,omitempty"`
	Body       []Stmt `json:"body"`
	Else       []Stmt `json:"else,omitempty"`
	IsPostEval bool   `json:"is_post_evaluated"`
	Line       int64  `json:"line"` // Line number of the statement relatively to the function.
}

type MapLit struct {
	Type *MapType        `json:"type"`
	Elts []*KeyValuePair `json:"elements"`
}

type KeyValuePair struct {
	Key   Expr `json:"key"`
	Value Expr `json:"value"`
}

type MapType struct {
	KeyType   Expr `json:"key_type"`
	ValueType Expr `json:"value_type"`
}

type MethodDecl struct {
	FuncDecl
	Override bool `json:"override"`
}

type OtherStmt struct {
	StmtName string `json:"statement_name"`
	Body     []Stmt `json:"body,omitempty"`
	Line     int64  `json:"line"` // Line number of the statement relatively to the function.
}

// Method/Function prototype declaration
type ProtoDecl struct {
	Doc        []string  `json:"doc"`
	Name       *Ident    `json:"name"`
	Type       *FuncType `json:"type"`
	Visibility string    `json:"visibility"`
}

type RangeLoopStmt struct {
	StmtName string `json:"statement_name"`
	Vars     []Expr `json:"variables,omitempty"`
	Iterable Expr   `json:"iterable,omitempty"`
	Body     []Stmt `json:"body"`
	Line     int64  `json:"line"` // Line number of the statement relatively to the function.
}

// A ReturnStmt represents a return statement.
type ReturnStmt struct {
	StmtName string `json:"statement_name"`
	Results  []Expr `json:"results,omitempty"` // result expressions; or nil
	Line     int64  `json:"line"`
}

type Stmt interface{}

// StructType represents a structured type. Most of the Object Oriented
// languages use a Class or a Trait instead.
//
// In Go, a StructType would be something of the form:
//    struct {
//       Bar string
//    }
type StructType struct {
	// This field is only used by the unmarshaller to "guess" the type while it
	// is unmarshalling a generic type. Since the StructType is considered as
	// an expression (which is represented by an interface{}), this is the only
	// way for the unmarshaller to know what type is it.
	//
	// The value of the ExprName for a StructType must always be "STRUCT", as
	// defined by the constant src.StructTypeName.
	ExprName string `json:"expression_name"`

	Doc    []string `json:"doc"`              // associated documentation; or nil
	Name   *Ident   `json:"name,omitempty"`   // name of the struct; or nil
	Fields []*Field `json:"fields,omitempty"` // the fields of the struct; or nil
}

// Field represents a pair name/type.
type Field struct {
	Doc  []string `json:"doc,omitempty"`  // associated documentation; or nil
	Name string   `json:"name,omitempty"` // name of the field; or nil
	Type string   `json:"type,omitempty"` // type of the field; or nil
}

type SwitchStmt struct {
	StmtName    string        `json:"statement_name"`
	Init        Stmt          `json:"initialization,omitempty"`
	Cond        Expr          `json:"condition,omitempty"` // TODO rename with a more appropriate name
	CaseClauses []*CaseClause `json:"case_clauses,omitempty"`
	Default     []Stmt        `json:"default,omitempty"`
}

type CaseClause struct {
	Conds []Expr `json:"conditions,omitempty"`
	Body  []Stmt `json:"body,omitempty"`
}

type TernaryExpr struct {
	ExprName string `json:"expression_name"`
	Cond     Expr   `json:"condition"`
	Then     Expr   `json:"then"`
	Else     Expr   `json:"else"`
}

type ThrowStmt struct {
	StmtName string `json:"statement_name"`
	X        Expr   `json:"expression"`
}

type Trait struct {
	Name    string        `json:"name"`
	Attrs   []*Attr       `json:"attributes"`
	Methods []*MethodDecl `json:"methods"`
	Classes []*ClassDecl  `json:"classes"`
	Traits  []*Trait      `json:"traits"`
}

type TraitRef struct {
	Namespace string `json:"namespace"`
	TraitName string `json:"trait_name"`
}

type TryStmt struct {
	StmtName     string         `json:"statement_name"`
	Body         []Stmt         `json:"body"`
	CatchClauses []*CatchClause `json:"catch_clauses,omitempty"`
	Finally      []Stmt         `json:"finally,omitempty"`
}

type CatchClause struct {
	Params []*Field `json:"parameters,omitempty"`
	Body   []Stmt   `json:"body,omitempty"`
}

// TypeSpec represents a type declaration. Most of the object oriented languages
// does not have such a node, they use classes and traits instead.
//
// In Go, a TypeSpec would be something of the form:
//    type Foo struct {
//       Bar string
//    }
type TypeSpec struct {
	Doc  []string `json:"doc,omitempty"`  // associated documentation; or nil
	Name *Ident   `json:"name"`           // type name (in the exemple, the name is "Foo")
	Type Expr     `json:"type,omitempty"` // *Ident or any of the *XxxType; or nil
}

type UnaryExpr struct {
	ExprName string `json:"expression_name"`
	Op       string `json:"operator"`          // operator
	X        Expr   `json:"operand,omitempty"` // operand (XXX investigate the omitempty)
}

type ValueSpec struct {
	ExprName string `json:"expression_name"`
	Name     *Ident `json:"name"`
	Type     *Ident `json:"type"`
}

type Var struct {
	Doc        []string `json:"doc,omitempty"`
	Name       string   `json:"name"`
	Type       string   `json:"type,omitempty"`
	Value      string   `json:"value,omitempty"`
	IsPointer  bool     `json:"is_pointer"`
	Visibility string   `json:"visibility,omitempty"`
}
