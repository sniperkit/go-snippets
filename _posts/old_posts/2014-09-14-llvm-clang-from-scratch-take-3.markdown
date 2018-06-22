---
layout: post
title: "LLVM Clang: Building it from scratch - take 3 (libunwind vs libgcc_s)"
description: "Compiling LLVM Clang from sources, 2nd attempt."
category: diary
tags: [C++, LLVM, Clang, build, compile, old]
---

_Investigating the possibility of replacing usage of `libgcc_s` with `libunwind`. **Edit/sept. 19th:** even though this "worked", this is not that much a correct way of doing things._

* Found a nice 2012 email-thread with loads of explanations about all this `libcxx`, `libcxxabi`, `compiler-rt` & `libunwind` _versus_ `libgcc*` _work-in-progress_ things (at that time) :

  [[Soc-coordination] Report 5: Provide an alternative to libstdc++	with libc++ ((by Andrej Belym))](http://lists.alioth.debian.org/pipermail/soc-coordination/2012-July/001316.html)

  Also: [libstdc++'s ABI Policy and Guidelines](http://gcc.gnu.org/onlinedocs/libstdc++/manual/abi.html)

  And: [Mini FAQ about the misc libc/gcc crt files.](http://dev.gentoo.org/~vapier/crt.txt)

  And: [wiki.osdev.org's Creating a C Library -- Program initialization](http://wiki.osdev.org/Creating_a_C_Library#Program_Initialization) ;
  found through [forum thread](http://forum.osdev.org/viewtopic.php?p=217724#p217724).

  And: [Libgcc @ wiki.osdev.org](http://wiki.osdev.org/Libgcc)

* [Building LLVM @ musl-libc wiki](http://wiki.musl-libc.org/wiki/BuildingLLVM#Build_libunwind)
  - **NOTE:** ^ it seems here that a fork is being used: [GH: pathscale/libunwind](https://github.com/pathscale/libunwind).
  - Last common commit with the original project dates back [August 4th 2010](https://github.com/pathscale/libunwind/commit/982d590ddb778f0d301fe2e5647abd9135a2f9bc) ;
  - and both projects have a recent history this year -\_-
  - _where am I goin' with this ?_


* **[libunwind project page @ nongnu.org](http://www.nongnu.org/libunwind/)**

* **[Git repo.](http://git.savannah.gnu.org/gitweb/?p=libunwind.git):** `git://git.sv.gnu.org/libunwind.git`


## Libunwind

### Importing `libunwind` from its Git repo.

	[fabi@wall] ~/dev/llvm (master ✗)  git submodule add -b master git://git.sv.gnu.org/libunwind.git libunwind

	Cloning into 'libunwind'...
	remote: Counting objects: 14396, done.
	remote: Compressing objects: 100% (3806/3806), done.
	remote: Total 14396 (delta 10522), reused 14388 (delta 10518)
	Receiving objects: 100% (14396/14396), 3.32 MiB | 630.00 KiB/s, done.
	Resolving deltas: 100% (10522/10522), done.
	Checking connectivity... done.

### Building it

**NOTE:** There's no `./configure` script when retrieving project from its repo. **and** `autoreconf` won't complete (errors).

**Hence** I'm just getting around this problem by checking out v1.1 and extracting the released tarball over it :

	[fabi@wall] ~/dev/llvm/libunwind (master)  gco -b v1.1 v1.1

	[fabi@wall] ~/dev/llvm/libunwind (v1.1)  wget -O - 'http://download.savannah.gnu.org/releases/libunwind/libunwind-1.1.tar.gz' |
		tar -zxf - --strip-components=1

^ **Ok**, empty git diff ;-

	# Using my local build of Clang :
	[fabi@wall] ~/dev/llvm/libunwind (v1.1 ✗)  source ../environment.sh

	[fabi@wall] ~/dev/llvm/libunwind (v1.1 ✗)  ./configure --prefix=/home/fabi/dev/llvm/local | tee -a clog.txt

	[fabi@wall] ~/dev/llvm/libunwind (v1.1 ✗)  time make -j2

	[fabi@wall] ~/dev/llvm/libunwind (v1.1 ✗)  time make install


### Testing if `libgcc_s` gets out of the way

Tried the symlinks quickfix below (also rebuilt `test.cpp`), but it didn't work (`ldd a.out` lists the host's `libgcc_s`) :

	[fabi@wall] ~/dev/llvm/local/lib (master ✗)  ln -siv libunwind.so libgcc_s.so
	‘libgcc_s.so’ -> ‘libunwind.so’
	[fabi@wall] ~/dev/llvm/local/lib (master ✗)  ln -siv libunwind.a libgcc.a
	‘libgcc.a’ -> ‘libunwind.a’
	[fabi@wall] ~/dev/llvm/local/lib (master ✗)  ln -siv libunwind.a libgcc_eh.a

Also tried to pass `-lunwind` hoping resolution of symbols through `libunwind` would have made `libgcc_s` disappear, it did not :

	[fabi@wall] ~/dev/llvm (master ✗)  clang++ -std=c++11 -stdlib=libc++          -nodefaultlibs -lc++ -lc++abi -lm -lc -lunwind -Wl,-rpath=$PWD/local/lib test.cpp

^ `a.out` did end up linked against `luunwind` however, but `libgcc_s` remained.


## Patching Clang

Mainly found out what to do from "[Building LLVM @ musl-libc wiki](http://wiki.musl-libc.org/wiki/BuildingLLVM)" (thanks to [Jeremy Huntwork](https://github.com/jhuntwork) probably) -- with miscellaneous investigations and trial'n'errors along the way.

Searching...

	[fabi@wall] ~/dev/llvm (master ✗)  time find -L llvm/ -type f \( -name '*.cpp' -o -name '*.c' \) \
	    -exec grep -EHn '\-lgcc' {} \;
	...
	... output is mainly some c/o lines ...
	...
	llvm/tools/clang/lib/Driver/ToolChains.cpp
	llvm/tools/clang/lib/Driver/Tools.cpp

Patching...

	[fabi@wall] ~/dev/llvm (master ✗)  sed -i -e 's@"-lgcc[^"]*"@"-lunwind"@g' \
	    llvm/tools/clang/lib/Driver/{Tools,TooChains}.cpp

### Replacing the hard-coded dynamic loader :

While at it I'm replacing the hard-coded dynamic-loader "name" -- I'll be dropping a symlink there later on that points do `musl-libc` ;- or not :

	[fabi@wall] ~/dev/llvm (master ✗)  time find -L llvm/ -type f \( -name '*.cpp' -o -name '*.c' \) \
	    -exec grep -EHn 'ld-linux-x86' {} \;
	...
	... ok, one meaningful file only ...
	llvm/tools/clang/lib/Driver/Tools.cpp:7231:    return "/lib64/ld-linux-x86-64.so.2";

**Ok, do it :**

	[fabi@wall] ~/dev/llvm (master ✗)  sed -i -e 's@/lib64/ld-linux-x86-64.so.2@/lib/ld-hack-x86_64.so@' \
	    llvm/tools/clang/lib/Driver/Tools.cpp

**Symlinking** under `/lib` :

	[fabi@wall] ~/dev/llvm (master ✗)  su -
	[root@wall] ~ # cd /lib

	[root@wall] lib # ls ld* -l
	-rwxr-xr-x 1 root root 144920 May 29 23:29 ld-2.17.so
	lrwxrwxrwx 1 root root     22 Aug 24 19:33 ld-linux.so.2 -> ../lib32/ld-linux.so.2
	lrwxrwxrwx 1 root root     10 Aug 24 19:33 ld-linux-x86-64.so.2 -> ld-2.17.so

	[root@wall] lib # ln -s ld-2.17.so ld-hack-x86_64.so
	[root@wall] lib # ln -s ../lib32/ld-linux.so.2 ld-hack.so


### Setting `libcxx` as the default standard lib. :

	[fabi@wall] ~/dev/llvm (master ✗)  sed -i '/^  return ToolChain::CST_Libstdcxx/s@stdcxx@cxx@' \
	    llvm/tools/clang/lib/Driver/ToolChain.cpp

## Re-bootstraping LLVM Clang!

**NOTE:** You did try to have it built with `-stdlib=libc++` but `cmake` conf. tests revealed a few unexpected failed tests, like the one complaining about `-fPIC` being not supported by the CXX ABI.

	[fabi@wall] ~/dev/llvm(master ✗)  screen
	[fabi@wall] ~  cd dev/llvm

	[fabi@wall] ~/dev/llvm(master ✗)  mkdir build
	[fabi@wall] ~/dev/llvm(master ✗)  cd build/

	[fabi@wall] ~/dev/llvm/build (master ✗)  cmake -DCMAKE_INSTALL_PREFIX=/home/fabi/dev/llvm/local ../llvm/

	[fabi@wall] ~/dev/llvm/build (master ✗)  time make -j3 -l3

_It's late -- tired -- spent too much time doin' these kind of "computer" things today._

_**End of transmission** -- Sept. 14th_
