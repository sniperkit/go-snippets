# The DevMine project

Initially, DevMine is a project about creating a database of developers and
establishing developers profiles based on their open source contributions.

Since its inception, its goals have evolved into a broader vision. The DevMine
project crawls metadata from open source projects and gathers source code.

Along with crawling and fetching tools, software to analyze this data are
provided as well.

DevMine consist of a complete toolchain from which several tools can be combined
together to do analysis. The output of one tool is typically piped to the input
of another to produce the final data and all of the serialization works thanks
to JSON.

DevMine has its own AST that can be seen as a generic AST to provide an
abstraction over the various programming languages over which the analysis can
be performed.

# About this repository

This repository serves as a generic repository to report issues that are *not*
specific to a DevMine subproject.

# About the authors

We are a group of open source enthusiasts with various motivations about the
project.

# Thanks

DevMine is made possible thanks to the [EPFL][epfl] [DATA lab][datalab] of
[Christoph Koch][ck] which provides us with the necessary resources. Think
about the fact that DevMine has to deal with millions of developers metadata and
millions of source code repositories as well that occupy terabytes of data
storage. Performing analysis requires a lot of computing resources too and
all of this would be impossible at large scale without the lab.

[epfl]: http://epfl.ch/
[datalab]: http://data.epfl.ch/
[ck]: http://people.epfl.ch/christoph.koch
