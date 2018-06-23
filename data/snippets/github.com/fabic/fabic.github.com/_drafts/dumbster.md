---
layout: page
title: "Draft - dumbster – the dumb parser"
tagline: "Dumbster: very simple source code parser with no semantical analysis."
category : notes
tags : [draft, build, tool, parser]
---

It takes about 20 seconds to read all sources files of the LLVM/Clang project and
dump these into a single file :

    time \
      bin/findsources.sh \
        llvm-clang/{clang,llvm,lld,lldb,libcxx,libcxxabi,compiler-rt,clang-tools-extra,test-suite} -- -print0 \
          | xargs -0 cat > tutu.txt

And the result weights ~200 MB :

    $ ls -lh tutu.txt
    -rw-rw-r-- 1 fabi fabi 204M Mar  1 13:08 tutu.txt

So it's okey to proceed with this “simple build tool” idea, probably...

## Have to...

* Locate all source files from a given directory.
* For each file :
    - Read the whole file into memory (or `mmap()`) ;
    - Parse the code ⇒ have an abstract syntax tree (AST) ;
    - Walk the AST in search for comments left upon artifacts of interests :
        - top-level file comment (and footer [?])
        - comments at global “anonymous” scope
        - comments at namespace scopes
        - comments _of_ a given namespace
        - comments found on "free functions" (those that aren't attached to structs/classes),
        - comments found on classes/structs (but _not_ the inner classes,
          nor ad-hoc structs one would find into function definitions).
    - For each comment :
        - Extract the textual content (removing leadin `*` and adhering to `\`
          line continuations).
        - Find out if the comment provides `! ...` prefixed ~~compilation~~
          __build directives__, like ex. `! executable dumbster`

### Examples / use cases

{% highlight c++ %}
/**
! executable: dumbster PIE W:no-main
! may library: libdumbster
! needs header : boostcpp/..
! link against library/ies: libnstd libstdc++ (or libc++)
! link prepend object : crti.o  !first!
! link prepend object : crtbegin.o  !if c++ filetype!
! link prepend object : crtend.o    !if c++ filetype!
! link append object  : crti.o  !last!
*/
int main(int argc, const char *argv[], const char *envp[])
{
  // ...
}
{% endhighlight %}









__EOF__
