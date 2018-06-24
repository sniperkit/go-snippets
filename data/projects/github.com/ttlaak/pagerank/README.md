pagerank
========

A naive implementation of Google's PageRank.

You should be able to get and install:

    go get github.com/ttlaak/pagerank
    go install github.com/ttlaak/pagerank

Input to the program consists of multiple files in the current directory.
Say you have 3 files 'a', 'b' and 'c' then each line in one of those files
consists of the filename of one of them. This constitutes a link so if the
file 'a' looks like:

    a
    c

then node 'a' has links to itself and 'c'. Files can be empty when there are
no outgoing links.

The output file (of which the filename is the mandatory argument) has floating
point numbers in rows. Each number corresponds to a node (or file in the
current directory). The order that is used (directory order) can be printed in
the orderFile. Use the -o option to specify the filename of the orderFile.

So to run this program setup the directory, say 'linkgraphdir', as described.
Then run the following:

    cd path/to/linkgraphdir
    pagerank -o ../lg-order ../lg-out

This will create two files in the parent directory the orderFile 'lg-order' and
the outFile in 'lg-out'.
