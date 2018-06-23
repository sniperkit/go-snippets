---
layout: note
title: "Linux limit [sub]process(s) maximum memory"
tagline: ""
published: true
category : notes
tags : [note]
---

Bash provides the shell builtin `ulimit`, and specifically `ulimit -v <kBytes>`
which sets the maximum amount of _virtual_ memory (i.e. heap) that it can allocate,
and it so happens that sub-processes will inherit this value as well.

    $ ulimit -v $((512*1024))

Quick test :

    cat <<EOF | clang++ -x c++ - && ( echo HEY; ulimit -v $((512*1024)) ; ./a.out )
    #include <iostream>
    using namespace std;

    int main() {
      cerr << "HEY";
      unsigned long amount = 0;
      const unsigned increment = 1024*1024;
      while(true) {
        char *buf = new char[ increment ];
        amount += increment;
        cerr << "\n" << amount << "        ";
        cerr << flush;
      }
      return -1;
    }
    EOF

* __TODO:__ Find out how to do this by resorting to __cgroups__ in an easy manner.
* Note that sub-processes will inherit this same max. virtual memory limit, so
we cannot restrict the maximum amount of memory of a whole process tree in this
way.
* Story: while doing some Jekyll dev. I ran into __{&#123; site | jsonify &#125;}__
freeze my whole system by consuming memory infinitely, certainly due to the fact
that `jsonify` does not handle cycles. _And I did it again when writing this line -\_-_.
