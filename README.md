# Parsers

DevMine parsers are meant to parse source code repositories and produce an
extended custom abstract syntax tree as JSON, as defined by
[srcanlzr](http://devmine.ch/doc/srcanlzr/).

Each parser is required to be able to read either directly from a folder
containing the source code as text files or directly from a tar archive.

The output of the languages parsers can be redirected to `srcanlzr` or
[repotool](http://devmine.ch/doc/repotool/) for further processing.

## Installation

Each of the language parser can be installed using
[srctool](http://devmine.ch/doc/srctool/) with the `install` command.
See `srctool help` for more details.

## How to write a language parser

Documentation on how to write a parser can be found
[here](http://godoc.org/github.com/DevMine/srcanlzr/src).

An article explaining the steps to follow in order to write a parser can be
found [here](http://devmine.ch/news/2015/05/31/how-to-write-a-parser).
