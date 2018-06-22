#!/usr/bin/env python
"""
Module for formatting a parsed CMake source file

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



class CMakeFormatter():

	def __init__(self, parsetree):
		self.parsetree = parsetree

	def output_as_cmake(self):
		return "\n".join(self.output_block(self.parsetree, 0))

	def output_block(self, block, level):
		lines = []
		if block is None:
			return lines
		for statement in block:
			lines.extend(self.output_statement(statement, level))
		return lines

	def output_statement(self, statement, level):

		func, args, comment, children = statement
		if func == None:
			# Allow visitors to replace one line with multiple lines
			# by making a None function have children.
			return self.output_block(children, level)

		if func == "" and comment is None:
			thisline = ""
		else:
			thisline = self.create_indent(level)

		thisline = self.output_function(statement, level, thisline)
		thisline = self.output_args(statement, level, thisline)
		thisline = self.output_comment(statement, level, thisline)
		output = [thisline]
		# Recurse into children
		output.extend(self.output_block(children, level + 1))
		return output

	def create_indent(self, level):
		return "\t" * level

	def output_function(self, statement, level, line):
		func, args, comment, children = statement
		return line + func

	def output_args(self, statement, level, line):
		func, args, comment, children = statement
		if args is None:
			args = ""
		newline = line
		if func is not None and func != "":
			newline = newline + "(" + args + ")"
		return newline

	def output_comment(self, statement, level, line):
		func, args, comment, children = statement
		newline = line
		if comment is not None:
			if newline != "" and newline != self.create_indent(level):
				newline = newline + "\t"
			newline = newline + comment
		return newline


class NiceFormatter(CMakeFormatter):
	def output_function(self, statement, level, line):
		func, args, comment, children = statement
		return line + func.lower()

	def output_args(self, statement, level, line):
		func, args, comment, children = statement
		if args is None:
			args = ""
		newline = line

		if func is not None and func != "":
			arglist = grammar.split_args(args)
			if len(newline + "(" + args + ")") > 72:
				newline = (newline + "(" +
							("\n"+self.create_indent(level+1)).join(arglist)
							+ ")")
			else:
				newline = newline + "(" + " ".join(arglist) + ")"
		return newline

#if __name__ == "__main__":
#	pass
