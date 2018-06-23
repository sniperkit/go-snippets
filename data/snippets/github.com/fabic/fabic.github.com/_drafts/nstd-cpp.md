---
layout: page
title: "Draft - NSTD ~ non-standard STD -\_-"
tagline: "NSTD is not std::"
category : notes
tags : [draft]
---

* Compiler & linker flags.
* `_start()` -> ... -> `main()`
* `memset()`
    * <http://stackoverflow.com/a/35638043/643087>
    * [memset.c @ Apple's libc](https://opensource.apple.com/source/Libc/Libc-166/string.subproj/memset.c)
    * <http://www.xs-labs.com/en/blog/2013/08/06/optimising-memset/>
* `void* operator new(std::size_t sz)`
    * [Gotcha #62:  Replacing Global New and Delete (2002, Stephen C. Dewhurst)](http://www.informit.com/articles/article.aspx?p=30642&seqNum=3)
    * <http://blog.aaronballman.com/2011/08/the-placement-new-operator/>

## Requirements, constraints

* [LP64](http://en.cppreference.com/w/cpp/language/types)
  __4/8/8__ (`int` is 32-bit, `long` and pointer are 64-bit).
* `char` is 8-bits “As of C++14, char must represent 256 distinct values, bijectively convertible to the values 0..255 of unsigned char [...]”.
* __C++14/17__ only, so that we do not have to resort to `#define-s` for backward
  compatible declarations.
* Try best no to use `#define-s`, at least for logic/code implementations.

### Coding style

* _work in progress_
* Let's _not_ prefix anything with `_` (underscores), and specifically let's not
  follow the __std::__ naming convention which does cripples code readability.

## Pointers

* [SGI's STL](http://www.sgi.com/tech/stl/download.html)
* [STLport](https://sourceforge.net/projects/stlport/) (last updated ~2008).
* The [Itanium C++ ABI](https://mentorembedded.github.io/cxx-abi/abi.html).
* [JTC1/SC22/WG21 - The C++ Standards Committee - ISOCPP](http://www.open-std.org/JTC1/SC22/WG21/)

### libc

* <http://codingrelic.geekhold.com/2008/11/six-million-dollar-libc.html>
  about Google's [Bionic libc](https://github.com/android/platform_bionic) impl.

__EOF__
