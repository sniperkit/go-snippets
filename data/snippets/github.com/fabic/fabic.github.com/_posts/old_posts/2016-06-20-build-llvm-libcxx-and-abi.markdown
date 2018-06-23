---
layout: post
title: "LLVM Clang: Building `libcxx` & `libcxxabi`"
description: "How to build `libcxx` and `libcxxabi` under Linux"
category: diary
tags: [C++, LLVM, Clang, build, compile, libcxx, libcxxabi, libc++, old]
---

* <http://libcxx.llvm.org/docs/BuildingLibcxx.html>
* <http://libcxx.llvm.org/>
* <http://libcxxabi.llvm.org/>

## Why ?

* Got a problem with Boost v1.61.0 where component “boost `program_options`”
  is failing to achieve anything due to some strange character encoding exception.

* ^ though it would be good to extend the possibility for more problems by having
  Clang work with its `libcxx / libcxxabi` thing (which it appear isn't supported
  under Linux).

* and I have some spare time ahead of me this morning...

## Setup environment

	[fabi@winterfell] ~/dev/llvm/build $

	$ cmake --version
	cmake version 3.4.3

	$ uname -a

	Linux localhost.localdomain 4.2.0-sabayon #1 SMP
	  Fri Dec 25 19:51:26 UTC 2015
	  x86_64
	  Intel(R) Core(TM) i7-2620M CPU @ 2.70GHz GenuineIntel
	  GNU/Linux

	$ source ./environment-clang.sh


## Build `libcxx` & `libcxxabi` as part of LLVM

This is by far the easiest / most reliable way I found to do this.

### Symlink `libcxx` and `libcxxabi` under `llvm/projects/`

	ln -sivn ../../libcxx llvm/projects/
	ln -sivn ../../libcxxabi llvm/projects/

### Build the whole thing

* This will build the whole LLVM thing, and the too actual things we want (libcxx & abi) ;
* About ~1800 build targets will get built so this will take a while...

	mkdir build && cd build/

	cmake -G Ninja ../llvm

_patience..._


## Building `libcxxabi` then `libcxx` separately

### libcxxabi

	cd libcxxabi/ &&
	mkdir build &&
	cd build/

	$ cmake .. \
		-DCMAKE_INSTALL_PREFIX=$(cd ../../local && pwd) \
		-DLLVM_PATH=../../llvm \
		-DLIBCXXABI_LIBCXX_PATH=../../libcxx \
		-DLIBCXXABI_LIBCXX_INCLUDES=../../libcxx/include

		... ... ...
		-- Configuring done
		-- Generating done
		-- Build files have been written to: /media/g/home/fabi/dev/llvm/libcxxabi/build

	$ make

		Scanning dependencies of target cxxabi_objects
		[  5%] Building CXX object src/CMakeFiles/cxxabi_objects.dir/abort_message.cpp.o
		[ 10%] Building CXX object src/CMakeFiles/cxxabi_objects.dir/cxa_aux_runtime.cpp.o
		...
		[ 90%] Building CXX object src/CMakeFiles/cxxabi_objects.dir/cxa_thread_atexit.cpp.o
		[ 90%] Built target cxxabi_objects

		Scanning dependencies of target cxxabi_shared
		[ 95%] Linking CXX shared library ../lib/libc++abi.so
		[ 95%] Built target cxxabi_shared

		Scanning dependencies of target cxxabi_static
		[100%] Linking CXX static library ../lib/libc++abi.a
		[100%] Built target cxxabi_static

#### check

	$ ldd lib/libc++abi.so
	    linux-vdso.so.1 (0x00007ffc4c32a000)
	    libpthread.so.0 => /lib64/libpthread.so.0 (0x00007f8d44f9c000)
	    libc.so.6 => /lib64/libc.so.6 (0x00007f8d44beb000)
	    libgcc_s.so.1 => /usr/lib/gcc/x86_64-pc-linux-gnu/4.9.3/libgcc_s.so.1 (0x00007f8d449d4000)
	    /lib64/ld-linux-x86-64.so.2 (0x000055b9d688c000)

#### install

	$ make install
	...
	...
	[ 90%] Built target cxxabi_objects
	[ 95%] Built target cxxabi_shared
	[100%] Built target cxxabi_static
	Install the project...
	-- Install configuration: ""
	-- Installing: /media/g/home/fabi/dev/llvm/local/lib/libc++abi.so.1.0
	-- Installing: /media/g/home/fabi/dev/llvm/local/lib/libc++abi.so.1
	-- Installing: /media/g/home/fabi/dev/llvm/local/lib/libc++abi.so
	-- Installing: /media/g/home/fabi/dev/llvm/local/lib/libc++abi.a

__NOTE:__ Headers aren't installed here, but will be upon install of `libcxx`.

### libcxx

* Had to provide `-DLIBCXX_CXX_ABI_INCLUDE_PATHS=../../local/include`
  explicitly for those headers to be located ;
* as well as `-DLIBCXX_CXX_ABI_LIBRARY_PATH=../../local/lib`
  so that it finds and links against `libc++.so`.

	$ cd libcxx/
	$ mkdir build && cd build/

	$ cmake .. \
			-DCMAKE_INSTALL_PREFIX=$(cd ../../local && pwd)
			-DLLVM_PATH=../../llvm \
	        -DLIBCXX_CXX_ABI=libcxxabi \
	        -DLIBCXX_CXX_ABI_INCLUDE_PATHS=../../libcxxabi/include \
	        -DLIBCXX_CXX_ABI_LIBRARY_PATH=../../local/lib

		...
		...
		-- Configuring done
		-- Generating done
		-- Build files have been written to: /media/g/home/fabi/dev/llvm/build

	$ make

		Scanning dependencies of target cxx
		...
		...
		[100%] Linking CXX shared library libc++.so
		GENERATING SCRIPT: 'INPUT(libc++.so.1 -lc++abi)' as file /media/g/home/fabi/dev/llvm/libcxx/build/lib/libc++.so
		[100%] Built target cxx

#### make install

	$ make install

	[100%] Built target cxx
	Install the project...
	-- Install configuration: ""

	-- Installing: /media/g/home/fabi/dev/llvm/local/include/c++/v1/cxxabi.h
	-- Installing: /media/g/home/fabi/dev/llvm/local/include/c++/v1/__cxxabi_config.h
	...
	-- Installing: /media/g/home/fabi/dev/llvm/local/include/c++/v1/./ostream
	-- Installing: /media/g/home/fabi/dev/llvm/local/include/c++/v1/./type_traits
	-- Installing: /media/g/home/fabi/dev/llvm/local/include/c++/v1/./cstring
	-- Installing: /media/g/home/fabi/dev/llvm/local/include/c++/v1/./stdlib.h
	...
	-- Installing: /media/g/home/fabi/dev/llvm/local/include/c++/v1/./cstdlib
	-- Installing: /media/g/home/fabi/dev/llvm/local/include/c++/v1/./iostream
	-- Installing: /media/g/home/fabi/dev/llvm/local/include/c++/v1/./vector
	-- Installing: /media/g/home/fabi/dev/llvm/local/include/c++/v1/./list
	-- Installing: /media/g/home/fabi/dev/llvm/local/include/c++/v1/./cstdint
	-- Installing: /media/g/home/fabi/dev/llvm/local/include/c++/v1/./forward_list
	...
	-- Installing: /media/g/home/fabi/dev/llvm/local/lib/libc++.so.1.0
	-- Installing: /media/g/home/fabi/dev/llvm/local/lib/libc++.so.1
	-- Installing: /media/g/home/fabi/dev/llvm/local/lib/libc++.so
	-- Up-to-date: /media/g/home/fabi/dev/llvm/local/lib/libc++.so

## re-setup environment

* Set `CXXFLAGS=-stdlib=libc++`
* and `CPLUS_INCLUDE_PATH=local/include/c++/v1`

	$ source ./environment-clang.sh

	$ ./show-environment.sh

## test

	cat <<EOF | clang++ -std=c++1y -stdlib=libc++ -x c++ -Wall -o test -
	#include <iostream>
	int main(int argc, char *argv[])
	{
	    int arr[3] = {1,2,3};
	    for(auto i : arr)
	        std::cout << i << ": Hello dude!" << std::endl;
	    return 127;
	}
	EOF

### It did link against `libc++.so.1` and `libc++abi.so.1`

	$ ldd test
        linux-vdso.so.1 (0x00007fff02beb000)
        libc++.so.1 => /media/g/home/fabi/dev/llvm/local/lib/libc++.so.1 (0x00007faa8b20b000)
        libc++abi.so.1 => /media/g/home/fabi/dev/llvm/local/lib/libc++abi.so.1 (0x00007faa8af8f000)
        libm.so.6 => /lib64/libm.so.6 (0x00007faa8ac8b000)
        libgcc_s.so.1 => /usr/lib/gcc/x86_64-pc-linux-gnu/4.9.3/libgcc_s.so.1 (0x00007faa8aa74000)
        libc.so.6 => /lib64/libc.so.6 (0x00007faa8a6c4000)
        libpthread.so.0 => /lib64/libpthread.so.0 (0x00007faa8a4a7000)
        librt.so.1 => /lib64/librt.so.1 (0x00007faa8a29f000)
        /lib64/ld-linux-x86-64.so.2 (0x00007faa8b58c000)

### and it runs :

	$ ./test

	1: Hello dude!
	2: Hello dude!
	3: Hello dude!

### Re-building Boost C++ v1.61.0

	./install_boost.sh

It appears some libraries have linked against `libc++.so.12` as desired :

	$ ldd local/boost-1.61.0-clang/lib/libboost_filesystem-mt-d.so
        linux-vdso.so.1 (0x00007ffc9b4d5000)
        libboost_system-mt-d.so.1.61.0 => not found
        librt.so.1 => /lib64/librt.so.1 (0x00007fcbf916b000)
        libc++.so.1 => /media/g/home/fabi/dev/llvm/local/lib/libc++.so.1 (0x00007fcbf8dea000)
        libc++abi.so.1 => /media/g/home/fabi/dev/llvm/local/lib/libc++abi.so.1 (0x00007fcbf8b6e000)
        libm.so.6 => /lib64/libm.so.6 (0x00007fcbf8869000)
        libgcc_s.so.1 => /usr/lib/gcc/x86_64-pc-linux-gnu/4.9.3/libgcc_s.so.1 (0x00007fcbf8652000)
        libpthread.so.0 => /lib64/libpthread.so.0 (0x00007fcbf8435000)
        libc.so.6 => /lib64/libc.so.6 (0x00007fcbf8084000)
        /lib64/ld-linux-x86-64.so.2 (0x0000559cd8128000)

Including the one `boost_program_options` I wanted :

	$ ldd local/boost-1.61.0-clang/lib/libboost_program_options-mt.so
        linux-vdso.so.1 (0x00007ffcecce3000)
        librt.so.1 => /lib64/librt.so.1 (0x00007fe8ef260000)
        libc++.so.1 => /media/g/home/fabi/dev/llvm/local/lib/libc++.so.1 (0x00007fe8eeede000)
        libc++abi.so.1 => /media/g/home/fabi/dev/llvm/local/lib/libc++abi.so.1 (0x00007fe8eec62000)
        libm.so.6 => /lib64/libm.so.6 (0x00007fe8ee95e000)
        libgcc_s.so.1 => /usr/lib/gcc/x86_64-pc-linux-gnu/4.9.3/libgcc_s.so.1 (0x00007fe8ee746000)
        libpthread.so.0 => /lib64/libpthread.so.0 (0x00007fe8ee529000)
        libc.so.6 => /lib64/libc.so.6 (0x00007fe8ee179000)
        /lib64/ld-linux-x86-64.so.2 (0x0000564c1c0b3000)

However its “debug” counterpart didn't, which is puzzling me :

	ldd local/boost-1.61.0-clang/lib/libboost_program_options-mt-d.so
        linux-vdso.so.1 (0x00007ffe175fb000)
        librt.so.1 => /lib64/librt.so.1 (0x00007fa11feea000)
        libstdc++.so.6 => /usr/lib/gcc/x86_64-pc-linux-gnu/4.9.3/libstdc++.so.6 (0x00007fa11fbbb000)
        libm.so.6 => /lib64/libm.so.6 (0x00007fa11f8b7000)
        libgcc_s.so.1 => /usr/lib/gcc/x86_64-pc-linux-gnu/4.9.3/libgcc_s.so.1 (0x00007fa11f6a0000)
        libpthread.so.0 => /lib64/libpthread.so.0 (0x00007fa11f483000)
        libc.so.6 => /lib64/libc.so.6 (0x00007fa11f0d2000)
        /lib64/ld-linux-x86-64.so.2 (0x0000560800deb000)

#### Investigating :

	for lib in local/boost-1.61.0-clang/lib/libboost_*.so.*; do
		echo "$lib :"
		ldd $lib |
			grep stdc;
	done

Unfortunately the problem doesn't just occur on the `-mt-d` debug variants :
a couple of shared libraries ended up linked against
`/usr/lib/gcc/x86_64-pc-linux-gnu/4.9.3/libstdc++.so.6` :-| :

* `local/boost-1.61.0-clang/lib/libboost_graph-mt.so.1.61.0`
* `local/boost-1.61.0-clang/lib/libboost_locale-mt.so.1.61.0`
* `local/boost-1.61.0-clang/lib/libboost_log-mt.so.1.61.0`
* `local/boost-1.61.0-clang/lib/libboost_log_setup-mt.so.1.61.0`
* `local/boost-1.61.0-clang/lib/libboost_regex-mt.so.1.61.0`

_**Argh! :** it didn't link against `libc++.so.1` :-/ fail._

_End of transmission_
