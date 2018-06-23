---
layout: post
title: "LLVM Clang: Building it from scratch - take 1"
description: ""
category: diary
tags: [C++, LLVM, Clang, build, compile, old]
---

_Been a while since my last C++ line of code -- spent a few days during my vacations reading about the recently released C++14 specs -- went through the C++11 language features -- wandering around the web about all this -- got no project at hand, so I'm just goin' to kill some time today building Clang from its sources, with some vague purposes..._

* <http://llvm.org/docs/CMake.html>
* <http://clang.llvm.org/get_started.html>
* <http://libcxx.llvm.org/> &
  <http://libcxxabi.llvm.org/>
* Using Github SVN repositories mirror: <https://github.com/llvm-mirror>
  - LLVM: <git@github.com:llvm-mirror/llvm.git>
  - Clang: <git@github.com:llvm-mirror/clang.git>
  - Clang tools: <git@github.com:llvm-mirror/clang-tools-extra.git>
  - compiler-rt: <git@github.com:llvm-mirror/compiler-rt.git>
  - libcxx: <git@github.com:llvm-mirror/libcxx.git>
  - libcxxabi: <git@github.com:llvm-mirror/libcxxabi.git>

**NOTE:** `libcxx` is flagged as experimental on Linux, `libcxxabi` is clearly unsupported, goin' down this road anyway...

**NOTE:** Using version 3.5.0 released September 4th 2014.

**NOTE:** Using trunk of `libcxxabi` as it doesn't have releases.

**NOTE:** Importing code as git subtrees turns out to be effectively a not so good idea.

_**DIDN'T SUCCEED:** a few projects fail at link-time due to "undefined symbols" -- I suspect CMake to have produced a few possibly incorrect Makefiles for these for I do find those symbols here a there in static/dynamic libs -- retrying with the `configure` script -- compilation takes hell of a long time -\_- -- may retry tomorrow as I found more materials about building this :_

* [LLVM doc. - GettingStarted: Requirements](http://llvm.org/docs/GettingStarted.html#host-c-toolchain-both-compiler-and-standard-library)
* [LLVM doc. index](http://llvm.org/docs/index.html)
* [musl-libc wiki: Building LLVM](http://wiki.musl-libc.org/wiki/BuildingLLVM)
  (note however that building LLVM/Clang against `musl` is [flagged as "Build fails"](http://wiki.musl-libc.org/wiki/Pkgsrc_results))


## Initializing a new Git repository :

	[fabi@wall] ~/dev  mkdir llvm
	[fabi@wall] ~/dev  cd llvm

	[fabi@wall] ~/dev/llvm  git init .
	Initialized empty Git repository in /home/fabi/dev/llvm/.git/

Searching for a stable release :

	[fabi@wall] ~/dev/llvm_/llvm (live34 ✗)  glsr git@github.com:llvm-mirror/llvm.git
	...
	...
	a93239b7c6f0d78cb8836768c3ffbc39fb15b79f	refs/heads/release_33
	5c6aa738fb3325ae499454877f1e2926d2368135	refs/heads/release_34
	df820bfd87d7d70d8f1a739333496e191631a3e2	refs/heads/release_35
	...
	...


## Checking out stuff

This may not be a good nor simple idea, but let's fetch all this by means of **git subtree** into a brand new project._

**NOTE:** I'm mainly subtree-importing code for I want the various LLVM projects' source-trees to be "at the root" here ; I'll then drop symlinks where applicable under `llvm/`.

**NOTE:** The whole thing ends up about 285 MB.

### llvm

	[fabi@wall] ~/dev/llvm (master ✓)  git subtree add --prefix=llvm --squash git@github.com:llvm-mirror/llvm.git refs/heads/release_35
	git fetch git@github.com:llvm-mirror/llvm.git refs/heads/release_35
	warning: no common commits
	remote: Counting objects: 883756, done.
	remote: Compressing objects: 100% (165444/165444), done.
	remote: Total 883756 (delta 714239), reused 883744 (delta 714230)
	Receiving objects: 100% (883756/883756), 229.28 MiB | 813.00 KiB/s, done.
	Resolving deltas: 100% (714239/714239), done.
	From github.com:llvm-mirror/llvm
	 * branch            release_35 -> FETCH_HEAD
	Added dir 'llvm'

### clang

	[fabi@wall] ~/dev/llvm (master ✓)  git subtree add --prefix=clang --squash git@github.com:llvm-mirror/clang.git refs/heads/release_35
	git fetch git@github.com:llvm-mirror/clang.git refs/heads/release_35
	warning: no common commits
	remote: Counting objects: 483866, done.
	remote: Compressing objects: 100% (76403/76403), done.
	remote: Total 483866 (delta 401905), reused 483864 (delta 401904)
	Receiving objects: 100% (483866/483866), 142.42 MiB | 805.00 KiB/s, done.
	Resolving deltas: 100% (401905/401905), done.
	From github.com:llvm-mirror/clang
	 * branch            release_35 -> FETCH_HEAD
	Added dir 'clang'

### clang-tools-extra

	[fabi@wall] ~/dev/llvm (master ✓)  git subtree add --prefix=clang-tools-extra --squash git@github.com:llvm-mirror/clang-tools-extra.git refs/heads/release_35
	git fetch git@github.com:llvm-mirror/clang-tools-extra.git refs/heads/release_35
	warning: no common commits
	remote: Counting objects: 5386, done.
	remote: Compressing objects: 100% (1861/1861), done.
	remote: Total 5386 (delta 3431), reused 5374 (delta 3424)
	Receiving objects: 100% (5386/5386), 1.38 MiB | 179.00 KiB/s, done.
	Resolving deltas: 100% (3431/3431), done.
	From github.com:llvm-mirror/clang-tools-extra
	 * branch            release_35 -> FETCH_HEAD
	Added dir 'clang-tools-extra'

### compiler-rt

	[fabi@wall] ~/dev/llvm (master ✓)  git subtree add --prefix=compiler-rt --squash git@github.com:llvm-mirror/compiler-rt.git refs/heads/release_35
	git fetch git@github.com:llvm-mirror/compiler-rt.git refs/heads/release_35
	warning: no common commits
	remote: Counting objects: 37479, done.
	remote: Compressing objects: 100% (10836/10836), done.
	remote: Total 37479 (delta 26439), reused 37461 (delta 26428)
	Receiving objects: 100% (37479/37479), 14.73 MiB | 796.00 KiB/s, done.
	Resolving deltas: 100% (26439/26439), done.
	From github.com:llvm-mirror/compiler-rt
	 * branch            release_35 -> FETCH_HEAD
	Added dir 'compiler-rt'

### libcxx

	[fabi@wall] ~/dev/llvm (master ✓)  git subtree add --prefix=libcxx --squash git@github.com:llvm-mirror/libcxx.git refs/heads/release_35
	git fetch git@github.com:llvm-mirror/libcxx.git refs/heads/release_35
	warning: no common commits
	remote: Counting objects: 34096, done.
	remote: Compressing objects: 100% (14444/14444), done.
	remote: Total 34096 (delta 18945), reused 34094 (delta 18945)
	Receiving objects: 100% (34096/34096), 12.09 MiB | 597.00 KiB/s, done.
	Resolving deltas: 100% (18945/18945), done.
	From github.com:llvm-mirror/libcxx
	 * branch            release_35 -> FETCH_HEAD
	Added dir 'libcxx'


### libcxx-abi

There seems to be no 3.5 release tag for this one :

	[fabi@wall] ~/dev/llvm (master ✓)  glsr git@github.com:llvm-mirror/libcxxabi.git
	68fcfa16650b79a767953a9a4b684d228bd5920c	HEAD
	68fcfa16650b79a767953a9a4b684d228bd5920c	refs/heads/master
	7e5a90713832cd48cc6e15b9d7a7a94e49ac4f37	refs/heads/release_32
	...

Checking out the trunk of it :

	[fabi@wall] ~/dev/llvm (master ✓)  git subtree add --prefix=libcxxabi --squash git@github.com:llvm-mirror/libcxxabi.git refs/heads/master
	git fetch git@github.com:llvm-mirror/libcxxabi.git refs/heads/master
	warning: no common commits
	remote: Counting objects: 1719, done.
	remote: Compressing objects: 100% (603/603), done.
	remote: Total 1719 (delta 1091), reused 1719 (delta 1091)
	Receiving objects: 100% (1719/1719), 1.27 MiB | 85.00 KiB/s, done.
	Resolving deltas: 100% (1091/1091), done.
	From github.com:llvm-mirror/libcxxabi
	 * branch            master     -> FETCH_HEAD
	Added dir 'libcxxabi'


## Symlinking stuff under `llvm/`

	[fabi@wall] ~/dev/llvm (master ✓)  cd llvm/tools/

	[fabi@wall] ~/dev/llvm/llvm/tools (master ✓)  ln -s ../../clang .

	[fabi@wall] ~/dev/llvm/llvm/tools (master ✓)  cd -P clang/tools/
	[fabi@wall] ~/dev/llvm/clang/tools (master ✓)  ln -s ../../clang-tools-extra/ extra

	[fabi@wall] ~/dev/llvm/clang/tools (master ✓)  cd -

	[fabi@wall] ~/dev/llvm/llvm/tools (master ✓)  cd ../projects/
	ls
	[fabi@wall] ~/dev/llvm/llvm/projects (master ✓)  ln -s ../../compiler-rt .

	[fabi@wall] ~/dev/llvm/llvm/projects (master ✓)  ln -siv ../../libcxx* .
	‘./libcxx’ -> ‘../../libcxx’
	‘./libcxxabi’ -> ‘../../libcxxabi’


## Make preps

### Check host Clang version

	[fabi@wall] ~/dev/llvm_/libcxxabi (master ✓)  clang --version
	clang version 3.3 (tags/RELEASE_33/final)
	Target: x86_64-pc-linux-gnu
	Thread model: posix

	[fabi@wall] ~/dev/llvm_/libcxxabi (master ✓)  clang++ --version
	clang version 3.3 (tags/RELEASE_33/final)
	Target: x86_64-pc-linux-gnu
	Thread model: posix
	[fabi@wall] ~/dev/llvm_/libcxxabi (master ✓)

### CMake invocation

_lunch break... ;-_

**NOTE:**I initially used llvm's `../configure --prefix=/home/fabi/dev/llvm/local` but it said :

	configure: WARNING: Unknown project (libcxx) won't be configured automatically
	configure: WARNING: Unknown project (libcxxabi) won't be configured automatically


Out-of-source building under `build/`

	[fabi@wall] ~/dev/llvm/build (master ✗)  cd ../..
	[fabi@wall] ~/dev/llvm/build (master ✗)  mkdir build
	[fabi@wall] ~/dev/llvm/build (master ✗)  cd build/

	[fabi@wall] ~/dev/llvm/build (master ✗)  CC=clang CXX=clang++
	[fabi@wall] ~/dev/llvm/build (master ✗)  export CC CXX

	[fabi@wall] ~/dev/llvm/build (master ✗)  cmake ../llvm/

#### Some CMake output that seems ok :

It seems to have found our host's Clang, and `libcxxabi` may well have been found :

	-- The C compiler identification is Clang 3.3.0
	-- The CXX compiler identification is Clang 3.3.0
	-- Check for working C compiler: /usr/bin/clang
	-- Check for working C compiler: /usr/bin/clang -- works
	-- Detecting C compiler ABI info
	-- Detecting C compiler ABI info - done
	-- Check for working CXX compiler: /usr/bin/clang++
	-- Check for working CXX compiler: /usr/bin/clang++ -- works
	-- Detecting CXX compiler ABI info
	-- Detecting CXX compiler ABI info - done
	-- Looking for C++ include cxxabi.h
	-- Looking for C++ include cxxabi.h - found
	...
	...

#### Some tests failed however :

Can't tell if this is bad or not :

	-- Performing Test LIBCXXABI_HAS_WX_FLAG - Failed
	-- Performing Test LIBCXXABI_HAS_NO_WX_FLAG - Failed
	-- Performing Test LIBCXXABI_HAS_EHSC_FLAG - Failed
	-- Performing Test LIBCXXABI_HAS_NO_EHS_FLAG - Failed
	-- Performing Test LIBCXXABI_HAS_NO_EHA_FLAG - Failed
	-- Performing Test LIBCXXABI_HAS_NO_GR_FLAG - Failed
	...
	...
	-- Performing Test COMPILER_RT_HAS_FFREESTANDING_FLAG - Failed
	-- Performing Test COMPILER_RT_HAS_FNO_FUNCTION_SECTIONS_FLAG - Failed
	-- Performing Test COMPILER_RT_HAS_GR_FLAG - Failed
	-- Performing Test COMPILER_RT_HAS_GS_FLAG - Failed
	-- Performing Test COMPILER_RT_HAS_MT_FLAG - Failed
	-- Performing Test COMPILER_RT_HAS_Oy_FLAG - Failed
	-- Performing Test COMPILER_RT_HAS_Zi_FLAG - Failed
	-- Performing Test COMPILER_RT_HAS_W3_FLAG - Failed
	-- Performing Test COMPILER_RT_HAS_WX_FLAG - Failed
	-- Performing Test COMPILER_RT_HAS_WD4722_FLAG - Failed
	...
	...

## Make

We're in for a long while...

	[fabi@wall] ~/dev/llvm/build (master ✗)  time make -j2

_Goin' for a coffee break, build is 29% done..._

_Grrr! Laptop wasn't plugged to powa -\_- and went off brutally, resuming, 35% done..._


### Failed -\_-

Link-time errors for a few targets, including Clang itself : It appears a library do not get passed as argument to the linker since the same symbols are said to be _undefined_ ; went through the few generated libs with `nm -C ...|less`, notably `lib/libclangParse.a` & `lib/libclang.so.3.5` where I do find everything ; briefly investigated some CMake-generated Makefiles, but couldn't quickly spot "_where the magic happen_".

## Re-trying with `../configure`

Using the provided `configure` script, compilation took about 6 hours to complete (with `-j1`) ! A few binaries from `clang-tools-extra` didn't make it through (compilation errors) :

	tools/clang/tools/extra/clang-apply-replacements
	tools/clang/tools/extra/remove-cstr-calls
	tools/clang/tools/extra/tool-template
	tools/clang/tools/extra/modularize
	tools/clang/tools/extra/module-map-checker
	tools/clang/tools/extra/pp-trace

The usually repported error is something related to a `DIAG` macro :

	...: error: too few arguments provided to function-like macro invocation

	/usr/include/clang/Basic/DiagnosticIDs.h:50:9: note: macro 'DIAG' defined here
	#define DIAG(ENUM,FLAGS,DEFAULT_MAPPING,DESC,GROUP,\

### Clang got built however ;-]

And got linked against `libgcc` & `libstdc++` as expected :

	[fabi@wall] ~/dev/llvm/build (master ✗)  ldd Release+Asserts/bin/clang
		linux-vdso.so.1 (0x00007ffffc1ec000)
		libz.so.1 => /lib64/libz.so.1 (0x00007fc56bd34000)
		libpthread.so.0 => /lib64/libpthread.so.0 (0x00007fc56bb17000)
		libedit.so.0 => /lib64/libedit.so.0 (0x00007fc56b8da000)
		libncurses.so.5 => /lib64/libncurses.so.5 (0x00007fc56b684000)
		libdl.so.2 => /lib64/libdl.so.2 (0x00007fc56b480000)
		libm.so.6 => /lib64/libm.so.6 (0x00007fc56b182000)
		libstdc++.so.6 => /usr/lib/gcc/x86_64-pc-linux-gnu/4.7.3/libstdc++.so.6 (0x00007fc56ae60000)
		libgcc_s.so.1 => /usr/lib/gcc/x86_64-pc-linux-gnu/4.7.3/libgcc_s.so.1 (0x00007fc56ac49000)
		libc.so.6 => /lib64/libc.so.6 (0x00007fc56a899000)
		/lib64/ld-linux-x86-64.so.2 (0x00007fc56bf4a000)

### Pondering...

Since I have a possibly working Clang, I wonder if I should use it to bootstrap a new build of it, and try to get it linked against `libcxx` - `libcxxabi` / `libcxxrt`, hum...

_End of transmission_
