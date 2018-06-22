// Copyright 2014-2015 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package src provides a set of structures for representing a project with its
	related source code independently of the language. In other words, it
	provides a generic representation (abstraction) of a source code.


	Goal

	The goal of this package is to provide a generic representation of a project
	that can be analyzed by the anlzr package as well as an API for
	encoding/decoding it to/from JSON.

	A presentation video is available on the DevMine website:

	   http://devmine.ch/news/2015/06/08/srcanlzr-presentation/

	Usage

	There are two kinds of program that interact with a src.Project: language
	parsers and VCS support tools. The former visits all source files inside
	the project folder and parse every source file in order to fill the
	src.Project.Packages field (and few others). The latter read the VCS folder
	that contains VCS data and fill the src.Project.Repo structure. The next
	two chapters treat about them more in details.


	Language parsers

	Language parsers must output the same structure as defined by the
	src.Project type. They have to first parse a project in order to get the
	specific AST for that project. Then, they have to make that AST match with
	our generic AST defined in the package:

	   http://godoc.org/github.com/DevMine/srcanlzr/src/ast

	To get more detail about how to write a language parser for srcanlzr, refer
	to that tutorial:

	   http://devmine.ch/news/2015/05/31/how-to-write-a-parser/


	VCS support tools

	Language parsers must not provide any information related to
	Version Control Systems (VCS). VCS metadata is the job of
	repotool:

	   http://devmine.ch/news/2015/06/01/repotool-presentation/
	   http://devmine.ch/doc/repotool/


	Example

	For the following Go source file (greet/main.go):

		package main

		import (
		        "fmt"
		)

		func greet(name string) {
		        fmt.Printf("Hello, %s!\n", name)
		}

		func main() {
		        name := "World"
		        greet(name)
		}

	The language parser must produce the following JSON output:

		{
		   "name": "greet",
		   "loc": 5,
		   "languages": [
		      {
				"language": "go",
		         "paradigms": [
		            "compiled",
		            "concurrent",
		            "imperative",
		            "structured"
		         ]
		      }
		   ],
		   "packages": [
		      {
		         "loc": 5,
		         "name": "greet",
		         "path": "/home/revan/go/src/foo/greet",
		         "source_files": [
		            {
		               "functions": [
		                  {
		                     "body": [
		                        {
		                           "expression": {
		                              "arguments": [
		                                 {
		                                    "expression_name": "BASIC_LIT",
		                                    "kind": "STRING",
		                                    "value": "Hello, %s!\\n"
		                                 },
		                                 {
		                                    "expression_name": "IDENT",
		                                    "name": "name"
		                                 }
		                              ],
		                              "expression_name": "CALL",
		                              "function": {
		                                 "function_name": "Printf",
		                                 "namespace": "fmt"
		                              },
		                              "line": 0
		                           },
		                           "statement_name": "EXPR"
		                        }
		                     ],
		                     "loc": 0,
		                     "name": "greet",
		                     "type": {
		                        "parameters": [
		                           {
		                              "name": "name",
		                              "type": "string"
		                           }
		                        ]
		                     },
		                     "visibility": ""
		                  },
		                  {
		                     "body": [
		                        {
		                           "left_hand_side": [
		                              {
		                                 "expression_name": "IDENT",
		                                 "name": "name"
		                              }
		                           ],
		                           "line": 1,
		                           "right_hand_side": [
		                              {
		                                 "expression_name": "BASIC_LIT",
		                                 "kind": "STRING",
		                                 "value": "World"
		                              }
		                           ],
		                           "statement_name": "ASSIGN"
		                        },
		                        {
		                           "expression": {
		                              "arguments": [
		                                 {
		                                    "expression_name": "IDENT",
		                                    "name": "name"
		                                 }
		                              ],
		                              "expression_name": "CALL",
		                              "function": {
		                                 "function_name": "greet",
		                                 "namespace": ""
		                              },
		                              "line": 0
		                           },
		                           "statement_name": "EXPR"
		                        }
		                     ],
		                     "loc": 0,
		                     "name": "main",
		                     "type": null,
		                     "visibility": ""
		                  }
		               ],
		               "imports": [
		                  "fmt"
		               ],
		               "language": {
		                  "language": "go",
		                  "paradigms": [
		                     "compiled",
		                     "concurrent",
		                     "imperative",
		                     "structured"
		                  ]
		               },
		               "loc": 5,
		               "path": "/home/revan/go/src/foo/greet/main.go"
		            }
		         ]
		      }
		   ]
		}

	Lines of Code counting

	The number of real lines of code must be precomputed by the language
	parsers. This is the only "feature" that must be precomputed because it may
	have multiple usages:

	1. Eliminate empty projects

	2. Evalutate project size

	3. Verify that the decoding is correct

	4. Normalize various counts

	5. ...

	Therefore, this count must be accurate and strictly follow the following
	rules:

	We only count statements and declarations as a line of code. Comments,
	package declaration, imports, expression, etc. must not be taken into
	account. Since an exemple is worth more than a thousand words, let's
	consider the following snippet:

	   // Package doc (does not count as a line of code)
	   package main // does not count as a line of code

	   import "fmt" // does not count as a line of code

	   func main() { // count as 1 line of code
	     fmt.Println(
	        "Hello, World!
	     ) // count as 1 line of code
	   }

	The expected number of lines of code is 2: The main function declaration
	and the call to fmt.Println function.


	Performance

	DevMine project is dealing with Terabytes of source code, therefore the
	JSON decoding must be efficient. That is why we implemented our own JSON
	decoder that focuses on performance. To do so, we had to make some
	choices and add some constraints for language parsers in order to make this
	process as fast as possible.

	JSON is usually unpredicatable which forces JSON parsers to be generic to
	deal with every possible kind of input. In DevMine, we have a well defined
	structure, thus instead of writting a generic JSON decoder we wrote one that
	decodes only src.Project objects. This really improves the performances
	since we don't need to use reflextion, generic types (interface{}) and type
	assertion. The drawback of this choice is that we have to update the decoder
	everytime we modify our structures.

	Most JSON parsers assume that the JSON input is potentially invalid
	(ie. malformed). We don't. Unlike json.Unmarshal, we don't Check for
	well-formedness.

	We also force the language parsers to put the "expression_name" and
	"statement_name" fields at the beginning of the JSON object. We use that
	convention to decode generic ast.Expr and ast.Stmt without reading the whole
	JSON object.

	Besides, we restrict the supported JSON types to:
	   string
	   int64
	   float64
	   bool
	   object
	   array
	All objects used (even inside an array) must absolutely be a pointer. This
	is required by the decoder generator.

	The only officially supported encoding is UTF-8.
*/
package src
