#!/usr/bin/env python
"""
Module for modifying a parsed CMake source file

Original Author:
2010 Ryan Pavlik <rpavlik@iastate.edu> <abiryan@ryand.net>
http://academic.cleardefinition.com
Iowa State University HCI Graduate Program/VRAC

          Copyright Iowa State University 2010.
 Distributed under the Boost Software License, Version 1.0.
    (See accompanying file LICENSE_1_0.txt or copy at
          http://www.boost.org/LICENSE_1_0.txt)
"""

###
# standard packages
import re

###
# third-party packages
# - none

###
# internal packages
import cmakegrammar

grammar = cmakegrammar

class CMakeBlock():
	def __init__(self, block):
		self.data = [CMakeStatement(x) for x in block]

	def __repr__(self):
		return repr([repr(x) for x in self.data])

	def get(self):
		output = []
		for x in self.data:
			output.extend(x.get())
		return output

	def accept(self, visitor):
		visitor.visit_block(self)
		for x in self.data:
			x.accept(visitor)

class CMakeStatement():
	def __init__(self, statement):
		self.func, self.args, self.comment, children = statement
		if children is not None:
			self.children = CMakeBlock(children)
		else:
			self.children = None

	def __repr__(self):
		return repr( (self.func, self.args, self.comment, repr(self.children) ))

	def get(self):
		if self.func is None:
			return self.children.get()
		elif self.children is None:
			return [(self.func, self.args, self.comment, None )]
		else:
			return [(self.func, self.args, self.comment, self.children.get() )]

	def replace_with_statements(self, statements):
		if self.func is not None:
			if self.comment is not None:
				statements.insert(0, ("", None, self.comment, None))
			self.func = None
			self.args = None
			self.comment = None
		self.children = CMakeBlock(statements)

	def accept(self, visitor):
		visitor.visit_statement(self)
		if self.children is not None:
			self.children.accept(visitor)


class CMakeVisitor:

	def visit_block(self, block):
		pass

	def visit_statement(self, statement):
		pass


class VisitorRemoveRedundantConditions(CMakeVisitor):
	funcs = ["else", "endif", "endmacro", "endfunction", "endforeach", "endwhile"]
	_reFuncs = (r"(?ix)" +		# case-insensitive and verbose
						r"(?P<RedundantConditionFuncs>" +
						"|".join(funcs) +
						r")$")
	reFuncs = re.compile(_reFuncs)
	def visit_statement(self, statement):
		if self.reFuncs.match(statement.func):
			statement.args = None

class VisitorReplaceSubdirs(CMakeVisitor):
	funcs = ["subdirs"]
	_reFuncs = (r"(?ix)" +		# case-insensitive and verbose
						r"(?P<SubdirsFuncs>" +
						"|".join(funcs) +
						r")$")
	reFuncs = re.compile(_reFuncs)
	def visit_statement(self, statement):
		if self.reFuncs.match(statement.func):
			args = grammar.split_args(statement.args)
			if len(args) == 1:
				statement.func = "add_subdirectory"
			else:
				statement.replace_with_statements([("add_subdirectory", x, None, None) for x in args ])

class VisitorFindModuleDependencies(CMakeVisitor):
	def __init__(self):
		CMakeVisitor.__init__(self)
		self.findmodules = []
		self.modules = []
		self.optionalmodules = []
		self.files = []
		self.optionalfiles = []
		self.directories = []
	def visit_statement(self, statement):
		if re.match(r"(?i)find_package$", statement.func):
			args = grammar.split_args(statement.args)
			if args[0]:
				self.findmodules.append("Find"+args[0])
		elif re.match(r"(?i)include$", statement.func):
			args = grammar.split_args(statement.args)
			if args[0]:
				if re.search(r"(?i)[/.]", args[0]):
					file = args[0]
					module = None
				else:
					file = None
					module = args[0]
				if "OPTIONAL" in args:
					self.optionalmodules.append(module)
					self.optionalfiles.append(file)
				else:
					self.modules.append(module)
					self.files.append(file)
		elif re.match(r"(?i)add_subdirectory$", statement.func):
			args = grammar.split_args(statement.args)
			if args[0]:
				self.directories.append(args[0])
		self.findmodules = [x for x in self.findmodules if x is not None]
		self.modules = [x for x in self.modules if x is not None]
		self.optionalmodules = [x for x in self.optionalmodules if x is not None]
		self.files = [x for x in self.files if x is not None]
		self.optionalfiles = [x for x in self.optionalfiles if x is not None]
		self.directories = [x for x in self.directories if x is not None]


def apply_all_cleanup_visitors(tree):
	rootBlock = CMakeBlock(tree)
	rootBlock.accept(VisitorRemoveRedundantConditions())
	rootBlock.accept(VisitorReplaceSubdirs())
	return rootBlock.get()

#if __name__ == "__main__":
#	pass
