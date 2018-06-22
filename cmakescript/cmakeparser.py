#!/usr/bin/env python
"""
Module for parsing a CMake source file

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

class UnclosedChildBlockError(Exception):
	pass

class InputExhaustedError(Exception):
	pass

def parse_string(instr):
	parser = CMakeParser(ParseInput(instr))
	parser.parse()
	return parser

def parse_file(filename):
	cmakefile = open(filename, 'r')
	instr = cmakefile.read()
	cmakefile.close()
	return parse_string(instr)

class ParseInput():
	"""Class providing an iterable interface to the parser's input"""

	def __init__(self, strdata):
		self._data = strdata.splitlines()
		# Add a None to the end as a sentry.
		self._data.append(None)
		self._lineindex = 0
		self.alldone = False
		self.gotline = False
		self.alreadyseen = False

	def __iter__(self):
		"""Be usable as an iterator."""
		return self

	def next(self):
		"""Return the current line each time we are iterated.
		We don't go to the next line unless we've been accepted."""

		# Will always hit this condition if all works well
		if self._lineindex == len(self._data):
			self.alldone = True
			raise StopIteration

		# If we break, break big
		assert self._lineindex < len(self._data)

		# OK, we can actually return the data now
		self.alreadyseen = self.gotline
		self.gotline = True
		return self._data[self._lineindex]

	def accept(self):
		"""Signal that we've processed this line and should go to the next"""

		# We shouldn't accept a line we haven't even seen
		assert self.gotline

		self._lineindex = self._lineindex + 1
		self.gotline = False

	def merge(self):
		# Don't want to merge current line with the next if we don't
		# know what the current line is
		assert self.gotline

		ln = self._lineindex
		if len(self._data) <= ln + 2:
			# we want 2 - one more text line, and the EOF sentinel
			raise InputExhaustedError
		newdata = self._data[:ln] + [self._data[ln] + "\n" + self._data[ln+1]]
		newdata.extend(self._data[ln+2:])
		# Make sure we didn't change anything - leave off end of file sentinel
		old = "\n".join(self._data[:-1])
		thenew = "\n".join(newdata[:-1])
		assert old == thenew

		self._data = newdata
		return self._data[ln]

class CMakeParser():

	def __init__(self, parseinput):
		self.input = parseinput
		self.parsetree = None

	def parse(self):
		self.parsetree = self.parse_block_children(None)
		if self.parsetree is None:
			self.parsetree = []

	def parse_block_children(self, startTag):
		if startTag is None:
			# AKA, the block is the entire file
			isEnder = lambda x: (x is None)

		elif grammar.reBlockBeginnings.match(startTag):
			# can have children
			endblock = grammar.dReBlockTagsDict[startTag.lower()]
			isEnder = endblock.match

		else:
			# This function can have no children
			return None

		block = []
		for line in self.input:
			# Keep merging lines until we either:
			# - successfully parse
			# - run out of input (that is, input contains incomplete statement)
			while 1:
				try:
					func, args, comment = grammar.parse_line(line)
				except grammar.IncompleteStatementError:
					try:
						line = self.input.merge()
					except InputExhaustedError:
						raise grammar.IncompleteStatementError
				else:
					# If no exception in parsing, break out
					# of this while loop: we can keep going now
					break

			if startTag == None and isEnder(func):
				return block
			elif func is not None and isEnder(func) and not self.input.alreadyseen:
				return block

			# Not an ender, so we accept this child.
			self.input.accept()
			children = self.parse_block_children(func)
			statement = ( func, args, comment, children)
			block.append( statement )

		# If we make it this far, we never found our Ender.
		raise UnclosedChildBlockError


#if __name__ == "__main__":
#	pass
