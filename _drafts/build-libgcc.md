---
layout: post
title: "Draft - Build GCC's libgcc"
tagline: "Supporting tagline"
category : notes
tags : [draft, GCC, build]
---

Query the remote repository in search for the last version tags :

    ~/dev $  git ls-remote  git@github.com:gcc-mirror/gcc.git  |
                sed -e 's#refs/tags/gcc_5#refs/tags/gcc-5#'  |
                  sort -k2 -V  |
                    tail -n4

    9d0507742960aa9f2b99bc6e9584655ecc611792        refs/tags/gcc-5_4_0-release
    c441d9e8e0438dcf54274ec7a3539859e94ad201        refs/tags/gcc-6_1_0-release
    6ac74a62ba7258299cf85fbef9bf45333ddf10c0        refs/tags/gcc-6_2_0-release
    4b5e15daff8b54440e3fda451c318ad31e532fab        refs/tags/gcc-6_3_0-release

Git clone the last version found :

    $ git clone  -b gcc-6_3_0-release  --depth 1  git@github.com:gcc-mirror/gcc.git

Git leaves us on a detached HEAD, checkout a branch :

    ~/dev $  cd gcc/

    ~/dev/gcc $  git checkout  -b gcc-6_3_0-release  tags/gcc-6_3_0-release

    ~/dev/gcc $  mkdir build && cd build

    ~/dev/gcc/build $  ../configure --prefix=$HOME/.local \
                                    --enable-languages=c,c++ \
                                    --disable-multilib \
                                    --disable-bootstrap \
                                    --with-system-zlib \
                                        |& tee _confiture.log

    ~/dev/gcc/build $  time make all-target-libgcc
    ~/dev/gcc/build $  make install-target-libgcc

    ~/dev/gcc/build $  find  $HOME/.local/  -name libgcc\*

    /home/fabi/dev/llvm-clang/local/lib/libgcc_s.so.1
    /home/fabi/dev/llvm-clang/local/lib/libgcc_s.so
    /home/fabi/dev/llvm-clang/local/lib/gcc/x86_64-pc-linux-gnu/6.3.0/libgcc.a
    /home/fabi/dev/llvm-clang/local/lib/gcc/x86_64-pc-linux-gnu/6.3.0/libgcc_eh.a

* <http://stackoverflow.com/a/22707330/643087>
* <http://wiki.osdev.org/Libgcc>
* <http://www.linuxfromscratch.org/lfs/view/development/chapter06/gcc.html>
* <https://gcc.gnu.org/releases.html>
* <https://gcc.gnu.org/install/>

__EOF__
