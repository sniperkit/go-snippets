---
layout: note
title: "Build SOCI C++ database abstraction library"
tagline: "Giving a try to SOCI, a database abstraction layer for C++ which supports various backends, of which MySQL."
published: true
category : notes
tags : [note, C++, database, library]
---

Build [SOCI](http://soci.sourceforge.net/)
([github.com/SOCI/soci](https://github.com/SOCI/soci)).

    $ git clone git@github.com:SOCI/soci.git vendor/soci &&
        cd vendor/soci/src/ &&
          ( [ -d build ] && rm -r build ) &&
          mkcd build &&
          cmake -G Ninja \
            -DCMAKE_BUILD_TYPE=RelWithDebInfo \
            -DSOCI_STATIC=ON \
            -DWITH_BOOST=ON \
            -DWITH_MYSQL=ON \
              ../src/ &&
          ninja -j2 ;
          echo "RETV=$? " ;
          pwd

__EDIT:__ -\_- dropped it, caused some seg. fault that I couldn't solve.
