Introduction
------------
The GoMP project is an effort to add a simple way to
parallelize loops in Go programs, similar to that of OpenMP.
The project was designed as an entry for Gopher Gala 2015 hackathon
that was limited to 48 hours, therefore it lacks some of
the features that OpenMP has. However, it is capable of
parallelizing simple 'for' loops which may be of some use to a
curious developer. Besides, it may work as a proof-of-concept
for we know there was some interest from the community in
a similar project.
The developers of GoMP tried to adhere to the OpenMP
standard when making most of the important design decisions.
These decisions are explained in the Limitations section.
That the word 'parallel' is used throughout this text
and the word 'concurrent' is barely mentioned is also a
consequence of this fact.


Installation
------------
To install GoMP, make sure that you have Go 1.4 installed.
Then type

	$ go install github.com/gophergala/gomp
    
This will add the gomp binary to your system. If
you want more information on the installation process,
see the documentation for the go tool by typing
   
	$ go help install

or by reading https://golang.org/cmd/go/.


Usage and Examples
------------------
If GoMP has been installed correctly you can now use
the 'gomp' program which reads an annotated Go source file
from stdin and writes this file with all its annotated
loops parallelized to stdout. The exact format of the
annotations follows.

A comment of the form "//gomp" should be inserted before 
every 'for' loop that is to be parallelized. Although this 
tool manipulation through textual comment-like pragmas is similar
in spirit to that of go generate (and to any macro processor,
for that matter), the actual rules are rather liberate:
  * Blank lines between "//gomp" and the loop are allowed.
  * "//gomp" must be in a comment line by itself but it
     is not restricted to be in the beginning of the line.
  * Comments supporting the 'for' loop should not be removed
    as long as the last comment line before the loop is of
    the specified format.
  * No spaces between the double slash and the word "gomp"
    are allowed.
  * The exact form of the comment is "//gomp". For example, "// //gomp"
    will not work.


Here is an example of how it may look like:

	...
	func main() {
	...
	//gomp
	for i := 0; i < 100; i++ {
		// Do some work
		...
	}
	...
	

The file should be used as standard input to 'gomp':
	
	$ gomp <main.go >par.go


More examples are at http://github.com/gophergala/gomp/examples.


Limitations
-----------
Only 'for' loops that have non-empty InitStmt, Condition and
PostStmt as defined in the language specification are supported.

In case the total range spanned by the loop variable does
not fit in the 64-bit integer type the resulting loop
is not guaranteed to work as expected (in particular it
is not even guaranteed to be finite even if it was before
parallelization). We are trying to mimick what OpenMP does in
a similar situation.

The loops whose correctness depends on integer overflow
(as in "for i := ^uint(0); i < 10; i++") are not guaranteed
to work either.

Current implementation of GoMP strips the file off all its comments.
The reason for this is the garbling of the line numeration because
of the changes that we make to the loops. The project authors
did not have enough time to correct the positions of all the
comments. Despite that, the file that was processed by GoMP
is not intended to be read by a human but instead it should
be compiled and run by a machine so its unreadable state should
not be that big of a problem (compare this with go generate, again).

Nested loops are not supported: they are scanned and the outermost
loop that has the "//gomp" comment preceding it is parallelized
with all its descendants being ignored. If you still want to
parallelize both the outer and the inner loops, try abstracting
the inner one into a function.

TODOs
-----

This section describes the most interesting immediate
directions of the development.

1. Add OMP private.
2. Add OMP reductions.
3. Do not erase the comments.
4. Relieve the user from manually running gomp on every file.
