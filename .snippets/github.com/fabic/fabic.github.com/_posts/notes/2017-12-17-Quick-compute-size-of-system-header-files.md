---
layout: note
title: "Quick compute size of system header files"
tagline: ""
published: true
category : notes
tags : [note, shell]
---

Using __awk__ to sum the output of __du__ ([SO](https://stackoverflow.com/a/450821/643087)):

    $ find /usr/include/ -type f \
        -iregex '.+\.\(c\|h\|hh\|cpp\|hpp\|cxx\|hxx\|h\.inc\|s\)$' \
          | xargs -r du -b \
          | awk '{s+=$1} END {print s}'

One-liner, using Bash's `$(( ... ))` to compute the size in kilo-bytes:

    $ echo "Size of your system header files: $(( $(find /usr/include/ -type f -iregex '.+\.\(c\|h\|hh\|cpp\|hpp\|cxx\|hxx\|h\.inc\|s\)$' | xar
    gs -r du -b | awk '{s+=$1} END {print s}') / 1024 )) kBytes."

    Size of your system header files: 270115 kBytes.
