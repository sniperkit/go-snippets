---
layout: note
title: "Bash: Test if variable is set (exists) or not"
tagline: "TL;DR: `[ -z ${EDITOR+x} ]`"
published: true
category : notes
tags : [note, shell, idioms, Bash]
---

Use Bash substitution `${var+x}` for testing whether or not a variable (incl.
environment one) is set or not ([via](https://stackoverflow.com/a/13864829/643087)).

    [ -z ${EDITOR+x} ] && echo "\$EDITOR IS NOT SET" || echo "SET: $EDITOR" ]

Example :

    [ -z ${CC+x}  ] && type -p clang >/dev/null
      && CC=clang
      && echo "| \$CC was set to $CC"

Likewise for `CXX` :

    [ -z ${CXX+x} ] && type -p clang++ >/dev/null
      && CXX=clang++
      && echo "| \$CXX was set to $CXX"

__EOF__
