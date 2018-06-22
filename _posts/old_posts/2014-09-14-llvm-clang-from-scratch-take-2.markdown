---
layout: post
title: "LLVM Clang: Building it from scratch - take 2 (ok)"
description: "Compiling LLVM Clang from sources, 2nd attempt."
category: diary
tags: [C++, LLVM, Clang, build, compile, old]
---

Trying to proceed the CMake-based build that [failed yesterday](/diary/2014/09/13/llvm-clang-from-scratch) -- Huh! 2 GB of resident memory while linking `clang` target, and with only 4 GB of total RAM, my laptop ends up using all the 4 GB swap -- killing some processes -- possibly too late: link failled with the same "undefined symbols" reason.

Rebooting laptop -- trying again from console mode -- got loads of free mem. -- linker quickly jumps up using 2 GB of memory, not swaping this time and 2 GB seems to be the peak mem. needed for this target -- 3 minutes of work and... **failed** with the same errors -- ok that wasn't due to some obscur "bug" possibly related to memory exhaustion.

## Re-trying, bootstrapping Clang with yesterdays build

Ok, restarting from scratch, -- from console-mode --, but this time I'm using my yesterday's build of Clang (which is linked against `libgcc` & `libstdc++`) to bootstrap this one :

    # Mostly successful build using ../llvm/configure :
    mv build build__conf
    mkdir build
    cd build/

    export PATH=$PWD/build__conf/Release+Asserts/bin:$PATH

    CC=clang CXX=clang++ ; export CC CXX ; env|grep clang ; which clang

    cmake ../llvm/

    time make -j3 -l4

_lunch break_

_48% built after 1 hour or so, 1.4 GB free mem. + 1.8 GB cached & no swap used, 15 min. load average is 2.75, ok._

_Wandering 'round -- reading `man vfork` -- I'm a bit dizzy -- something like 2h20 to go -- nap break!_

_Back from nap ~1hr -- 98% done, fingers-crossed -- `libclang.so` is being linked, target `clang` just started! -- free mem. is dropping fast: can now see the linker process peaking a 2 GB memory -- little swapping occurs (other jobs running too) -- `ld` process peaks at 2.2 GB, and disappeared! Great!: `[ 98%] Built target clang` ;-)_

_Damn :/ can't find `clang` binary under `bin/` -- worried -- make is at 100% performing stuff about `libcxx_tsan/msan` -- `find -type f -perm -u+x |sort`: no `clang` anywhere :/ I hope this is due to some yet uncompleted build steps -- wating..._

_Argh! I was searching for it in yesterday's failled build dir. -\_-_

_Build completed, **ok**._

## libcxx any ?

Nope, it linked itself against the host's `libgcc` & `libstdc++` libs. ; since `libc++.so` was built I was expecting that Clang would build itself against it magically, but it did not :

    [fabi@wall] ~/dev/llvm/build2 (master ■)  ldd bin/clang
            linux-vdso.so.1 (0x00007fffcdbfe000)
            librt.so.1 => /lib64/librt.so.1 (0x00007ffb15492000)
            libdl.so.2 => /lib64/libdl.so.2 (0x00007ffb1528e000)
            libncurses.so.5 => /lib64/libncurses.so.5 (0x00007ffb15038000)
            libpthread.so.0 => /lib64/libpthread.so.0 (0x00007ffb14e1b000)
            libz.so.1 => /lib64/libz.so.1 (0x00007ffb14c05000)
            libstdc++.so.6 => /usr/lib/gcc/x86_64-pc-linux-gnu/4.7.3/libstdc++.so.6 (0x00007ffb148e3000)
            libm.so.6 => /lib64/libm.so.6 (0x00007ffb145e5000)
            libgcc_s.so.1 => /usr/lib/gcc/x86_64-pc-linux-gnu/4.7.3/libgcc_s.so.1 (0x00007ffb143ce000)
            libc.so.6 => /lib64/libc.so.6 (0x00007ffb1401e000)
            /lib64/ld-linux-x86-64.so.2 (0x00007ffb1569a000)

This is fine though since I now know that I can have a full `make all` run that completes w/o errors, and I have a clean Clang 3.5.0 toolchain for proceeding with further experiments.

## Hence...

* Build LLVM/Clang with `../llvm/configure && make -j3 -l3` so as to obtain a working toolchain (with possibly a few minor targets that fail in `clang-tools-extra`) ;
* Bootstrap-build again using this fresh `Clang`.

## Some brief testing

### hello world++

    [fabi@wall] ~/dev/llvm (master ✗)  cat > test.cpp <<EOF
    /** File: test.cpp
     */

    #include <iostream>

    using std::cout;

    int main(int argc, char *argv[])
    {
        int arr[3] = {1,2,3};
        for(auto i : arr)
            cout << i << ": Hello dude!\n";
        return 127;
    }
    EOF

Builds, ok :

    [fabi@wall] ~/dev/llvm (master ✗)  clang++ -std=c++11 test.cpp

    [fabi@wall] ~/dev/llvm (master ✗)  ./a.out ; echo $?
    1: Hello dude!
    2: Hello dude!
    3: Hello dude!
    127

It got linked with `libstdc++` & `libgcc_s` as expected :

    [fabi@wall] ~/dev/llvm (master ✗)  ldd a.out
        linux-vdso.so.1 (0x00007fffff1fe000)
        libstdc++.so.6 => /usr/lib/gcc/x86_64-pc-linux-gnu/4.7.3/libstdc++.so.6 (0x00007f4741996000)
        libm.so.6 => /lib64/libm.so.6 (0x00007f4741698000)
        libgcc_s.so.1 => /usr/lib/gcc/x86_64-pc-linux-gnu/4.7.3/libgcc_s.so.1 (0x00007f4741481000)
        libc.so.6 => /lib64/libc.so.6 (0x00007f47410d1000)
        /lib64/ld-linux-x86-64.so.2 (0x00007f4741cb8000)

#### With `-stdlib=c++`

    [fabi@wall] ~/dev/llvm (master ✗)  clang++ -std=c++11 -stdlib=libc++ test.cpp
    test.cpp:9:10: fatal error: 'iostream' file not found
    #include <iostream>
             ^
    1 error generated.

Ok with `-lc++abi` and specifying `libcxx`'s include path :

    [fabi@wall] ~/dev/llvm (master ✗)  clang++ -std=c++11 -stdlib=libc++ \
        -Illvm/projects/libcxx/include/ \
        -lc++abi \
        -v test.cpp

Ok too as per <http://libcxx.llvm.org/> :

    [fabi@wall] ~/dev/llvm (master ✗)  clang++ -std=c++11 -stdlib=libc++ \
        -Illvm/projects/libcxx/include/ \
        -nodefaultlibs -lc++ -lc++abi -lm -lc -lgcc_s -lgcc \
        test.cpp

    [fabi@wall] ~/dev/llvm (master ✗)  ldd ./a.out
        linux-vdso.so.1 (0x00007fffcdd69000)
        libc++.so.1 => not found
        libc++abi.so.1 => not found
        libm.so.6 => /lib64/libm.so.6 (0x00007f643ff30000)
        libc.so.6 => /lib64/libc.so.6 (0x00007f643fb80000)
        libgcc_s.so.1 => /usr/lib/gcc/x86_64-pc-linux-gnu/4.7.3/libgcc_s.so.1 (0x00007f643f969000)
        /lib64/ld-linux-x86-64.so.2 (0x00007f644022e000)

Ok, it runs :

    [fabi@wall] ~/dev/llvm (master ✗)  export LD_LIBRARY_PATH=$PWD/build/lib

    [fabi@wall] ~/dev/llvm (master ✗)  ldd ./a.out
        linux-vdso.so.1 (0x00007fff34fc8000)
        libc++.so.1 => /home/fabi/dev/llvm/build/lib/libc++.so.1 (0x00007f2f94f31000)
        libc++abi.so.1 => /home/fabi/dev/llvm/build/lib/libc++abi.so.1 (0x00007f2f94cae000)
        libm.so.6 => /lib64/libm.so.6 (0x00007f2f949b0000)
        libc.so.6 => /lib64/libc.so.6 (0x00007f2f94600000)
        libgcc_s.so.1 => /usr/lib/gcc/x86_64-pc-linux-gnu/4.7.3/libgcc_s.so.1 (0x00007f2f943e9000)
        libpthread.so.0 => /lib64/libpthread.so.0 (0x00007f2f941cc000)
        librt.so.1 => /lib64/librt.so.1 (0x00007f2f93fc4000)
        /lib64/ld-linux-x86-64.so.2 (0x00007f2f952cf000)

    [fabi@wall] ~/dev/llvm (master ✗)  ./a.out
    1: Hello dude!
    2: Hello dude!
    3: Hello dude!

Also **ok** if hard-coding library path with `-Wl,-rpath=$PWD/local/lib test.cpp` :

    [fabi@wall] ~/dev/llvm (master ✗)  clang++ -std=c++11 -stdlib=libc++ \
        -nodefaultlibs -lc++ -lc++abi -lm -lc -lgcc_s -lgcc \
        -Wl,-rpath=$PWD/local/lib test.cpp

    [fabi@wall] ~/dev/llvm (master ✗)  ldd ./a.out
        linux-vdso.so.1 (0x00007fff4d597000)
        libc++.so.1 => /home/fabi/dev/llvm/local/lib/libc++.so.1 (0x00007f78864de000)
        libc++abi.so.1 => /home/fabi/dev/llvm/local/lib/libc++abi.so.1 (0x00007f788625b000)
        libm.so.6 => /lib64/libm.so.6 (0x00007f7885f5d000)
        libc.so.6 => /lib64/libc.so.6 (0x00007f7885bad000)
        libgcc_s.so.1 => /usr/lib/gcc/x86_64-pc-linux-gnu/4.7.3/libgcc_s.so.1 (0x00007f7885996000)
        libpthread.so.0 => /lib64/libpthread.so.0 (0x00007f7885779000)
        librt.so.1 => /lib64/librt.so.1 (0x00007f7885571000)
        /lib64/ld-linux-x86-64.so.2 (0x00007f788687c000)


## Epilogue

It seems I have a good starting point for playing around with compiling stuff against `libcxx/abi` ;- under Linux...

_End of transmission._
