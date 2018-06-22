---
layout: page
title: "Algorithms – notes"
tagline: "Random personal notes about algorithms."
category : notes
tags : [draft, algorithms]
published: true
---


* <http://raghibm.com/blog/checklist-for-a-softwareEng-interview/>

### Lists

* One-way linear list “is a sequence of cells each of which (except the last)
  points to its successor.”

### Heap, priority queues

<https://en.wikipedia.org/wiki/Heap_(data_structure)>

### Trees

* [Van Emde Boas tree / priority queue](https://en.wikipedia.org/wiki/Van_Emde_Boas_tree)
  – “also known as a vEB tree, is a tree data structure which implements an
  associative array with m-bit integer keys. It performs all operations in
  O(log m) time, or equivalently in O(log log M) time, where M = 2m is the
  maximum number of elements that can be stored in the tree. [...]”.
* [Fractal tree index](https://en.wikipedia.org/wiki/Fractal_tree_index)
* [B-trees, Shadowing, and Clones – Ohad Rodeh – 2007](http://liw.fi/larch/ohad-btrees-shadowing-clones.pdf)
  – publication, file-systems related, found through
  [larch: Python B-tree library](http://liw.fi/larch/).

## Hash/map

* [New Concurrent Hash Maps for C++ | 2016, Jeff Preshing](http://preshing.com/20160201/new-concurrent-hash-maps-for-cpp/)
  (as well as [Using Quiescent States to Reclaim Memory | 2016, Jeff Preshing](http://preshing.com/20160726/using-quiescent-states-to-reclaim-memory/)).

## Other concepts

* [What is RCU, Fundamentally ? | 2007, lwn.net](https://lwn.net/Articles/262464/)

## Disk-based data structures

* __Blog post:__ [Datastructures for external memory – 2016 by Max Bolingbroke](http://blog.omega-prime.co.uk/?p=197)
* __Online course:__ [Disk-based data structures (Lecture #16, by Dr. Naveen Garg, ~2009) – NPTEL.ac.in](http://nptel.ac.in/courses/106102064/16) (vidéo course)


## Books

* [Algorithms, Etc (lecture notes by Jeff Erickson)](http://jeffe.cs.illinois.edu/teaching/algorithms/)
  ([PDF with all lecture notes, 1200+ pages](http://jeffe.cs.illinois.edu/teaching/algorithms/everything.pdf))
* [Elementary Algorithms (online free e-book, by github.com/liuxinyu95)](https://github.com/liuxinyu95/AlgoXY)
  ([PDF v0.618033](https://github.com/liuxinyu95/AlgoXY/releases/download/v0.618033/elementary-algorithms.pdf)).

### Other

* [github.com/haseebr/competitive-programming/](https://github.com/haseebr/competitive-programming/tree/master/Materials)
  stores PDFs of various e-books.
* [github.com/waydecs/books](https://github.com/waydecs/books) _likewise_.

## Misc.

* [msgpack (for C/C++)](https://github.com/msgpack/msgpack-c)
  – binary serialization library.
* [OpenDHT](https://github.com/savoirfairelinux/opendht)
  – A C++11 Distributed Hash Table implementation.

__EOF__
