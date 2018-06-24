List
====

The list command lists files in a simple and consistent format. It works well
together with other commands such as grep and sort.

Usage
-----

	list [flag...] [file...]

	-a	Do not ignore files starting with a dot.
	-d	List directories themselves, not their contents.
	-f	Show file type.
	-r	List subdirectories recursively.
	-s	Show size.
	-t	Show time of last modification.

If no file is specified, the current directory is listed.

Sorting
-------

The output from the list command is always sorted by file name. To change the
sort order, pipe the output to the sort command.

Sort by file type:

	list -f | sort

Sort by size:

	list -s | sort -n

Sort by time of last modification:

	list -t | sort

Sort in reverse order:

	list | sort -r

Searching
---------

The -r flag lists directories recursively. This is useful when searching for
files. Pipe the output to the grep command to select files, or sort the list
with the sort command.

Find all PNG files:

	list -r | grep 'png$'

Find all symbolic links:

	list -r -f | grep '^link'

Find files from December 2017:

	list -r -t | grep '^2017-12'

Find the largest file:

	list -r -s | sort -n | tail -n 1

License
-------

The list command is distributed under the terms of the MIT license. See the
file LICENSE.md for more information.
