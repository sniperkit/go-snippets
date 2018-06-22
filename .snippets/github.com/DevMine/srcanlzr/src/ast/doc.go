// Copyright 2014-2015 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package ast represents a language agnostic Abstract Syntax Tree (AST).

	The aim of this package is to provide a generic syntactic representation.

	Serialization Constraints

	Every expression must have a field called "expression_name" that represents
	the type of the expression as defined by the package token. The same rules
	applies to statements, which must have a field "statement_name".

	IMPORTANT: That field MUST be the first one of the structure. This is
	imperative to make the JSON decoding work. This constraint is not really
	convenient but necessary in order to make performance optimizations in the
	JSON decoding process.

	For more information about how serialization work, refer to the
	documentation of the src package:

	   http://godoc.org/github.com/DevMine/srcanlzr/src
*/
package ast
