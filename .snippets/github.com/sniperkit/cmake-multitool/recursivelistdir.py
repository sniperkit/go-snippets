#!/usr/bin/env python
"""
Functions to recurse down directory trees

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
import os.path

###
# third-party packages
# - none

###
# internal packages
# - none

def recursive_listdir(dir):
	"""
	Return a list of all regular, non-hidden files, recursing directories

	Warning: if you run this on a directory with a lot of files, it may
	end up producing some enormous lists!  Use the iterator version instead
	in those cases, as the input to a "filter" expression.

	This returns the full absolute path to the files, and in no particular
	order (whatever order they are received from os.listdir, when recursing
	depth-first into sub-directories)
	"""
	files = []
	for item in os.listdir(dir):
		# Skip hidden files/directories
		if item[0] == ".":
			continue

		current = os.path.join(dir, item)

		if os.path.isfile(current):
			files.append(current)
		else:
			files.extend(recursive_listdir(current))

	return files

def recursive_listdir_iter(dir):
	"""
	An iterator over all regular, non-hidden files in dir and below.

	To best handle directories with many children, use a filter expression
	with recursive_listdir_iter(yourdir) as the "sequence" argument.

	This returns the full absolute path to the files, and in no particular
	order (whatever order they are received from os.listdir, when recursing
	depth-first into sub-directories)
	"""
	for item in os.listdir(dir):
		if item[0] == ".":
			continue

		current = os.path.join(dir, item)
		if os.path.isfile(current):
			yield current
		else:
			for recursive_item in recursive_listdir_iter(current):
				yield recursive_item
