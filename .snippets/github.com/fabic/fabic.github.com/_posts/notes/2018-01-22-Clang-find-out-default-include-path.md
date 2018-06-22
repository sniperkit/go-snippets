---
layout: note
title: "Clang/GCC: find out default include path"
tagline: "Using Sed to select compiler hard-coded include paths."
published: true
category : notes
tags : [note, cli, Clang, sed]
---

This will list the hard-coded compiler include paths :

```bash
clang++ -Wp,-v -x c++ - -fsyntax-only < /dev/null
```

Or :

```bash
clang++ -E -x c++ - -v < /dev/null
```

Using __Sed__ to select those include path ([SO](https://stackoverflow.com/a/17990436)):

```bash
$ clang++ -Wp,-v -x c++ - -fsyntax-only < /dev/null 2>&1 |
    sed -e '/^#include <...>/,/^End of search/{ //!b };d'
```

Where `//!b` :

> @Brendan the instruction __//!b__ reads if the current line is neither one of
> the lines that match the range, break and therefore print those lines
> otherwise all other lines are deleted. – potong Feb 17 '17 at 1:14

__EOF__
