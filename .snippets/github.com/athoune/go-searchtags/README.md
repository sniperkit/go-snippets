Search tags
===========

Experimentation with golang.

All memory index for finding similarity within documents.

Less than 100k documents, document is qualified by 1k tags.

It's REST.

Make it
-------

Install go.

    go get
    ./build.sh

Test it
-------

In a terminal:

    GOMAXPROCS=3 ./searchtags

n-1 procs is a good number.

When application starts, 50k documents is randomly generated.

In another terminal:

    curl http://localhost:8000/similar/42

Looking for document like document 42.

Model
-----

A _document_ got tags. It's represented by a bitset.

Tags can have weight.

Document can have weight.
