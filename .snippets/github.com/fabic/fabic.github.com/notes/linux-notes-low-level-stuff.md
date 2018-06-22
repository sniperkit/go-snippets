---
layout: page
title: "Linux – notes – low-level stuff"
tagline: "Personal notes and pointers while crawling in the deepest..."
category : notes
tags : [draft, Linux, low-level, wide]
published: true
---

## TL;DR

&hellip; personal notes and pointers about the low-level internals of Linux/kernel/Glibc/etc &hellip;


## Links

* <https://filippo.io/linux-syscall-table/>
* <http://syscalls.kernelgrok.com/>
* <http://blog.rchapman.org/posts/Linux_System_Call_Table_for_x86_64/> (2012)
* <http://www.digilife.be/quickreferences/qrc/linux%20system%20call%20quick%20reference.pdf> (Jialong He)
* <https://blog.packagecloud.io/eng/2016/04/05/the-definitive-guide-to-linux-system-calls/> (2016, article)
* <http://cs.lmu.edu/~ray/notes/linuxsyscalls/>
* <http://pubs.opengroup.org/onlinepubs/9699919799/idx/headers.html>
OpenGroup POSIX.1-2008 [toc](http://pubs.opengroup.org/onlinepubs/9699919799/basedefs/contents.html)


## The Linux Kernel

* [LWN.net](https://lwn.net/)
* [Linux kernel documentation @ kernel.org](https://www.kernel.org/doc/html/latest/)
* [linux-stable source tree (repository)](https://git.kernel.org/pub/scm/linux/kernel/git/stable/linux-stable.git/)
    - [github.com/**torvalds/linux** : The Linux kernel source tree, mirror repository @ Github](https://github.com/torvalds/linux)

<p>
  <a href="https://commons.wikimedia.org/wiki/File:Simplified_Structure_of_the_Linux_Kernel.svg#/media/File:Simplified_Structure_of_the_Linux_Kernel.svg">
    <img src="https://upload.wikimedia.org/wikipedia/commons/thumb/6/65/Simplified_Structure_of_the_Linux_Kernel.svg/1200px-Simplified_Structure_of_the_Linux_Kernel.svg.png" alt="Simplified Structure of the Linux Kernel.svg"></a>
    <br>By <a href="//commons.wikimedia.org/wiki/User:ScotXW" title="User:ScotXW">ScotXW</a>
      - <span class="int-own-work" lang="en">Own work</span>,
      <a href="https://creativecommons.org/licenses/by-sa/4.0"
         title="Creative Commons Attribution-Share Alike 4.0">CC BY-SA 4.0</a>,
      <a href="https://commons.wikimedia.org/w/index.php?curid=47075153">Link</a>
</p>


## Shared libraries

* [SO: shared lib. initialization](https://stackoverflow.com/a/9759936/643087) :
    - Linker flags : `-Wl,-init,<function>` & `-Wl,-fini,<function>`
    - `__attribute__((constructor)) void foo(void) { printf("library loaded!\n"); }`
    - The constructor attr. may have a priority specification : `void before_main(void) __attribute__((constructor (101)))` ([SO](https://stackoverflow.com/a/32701238/643087))
    - Or ([via](https://stackoverflow.com/a/1681655/643087)):

        ```cpp
        namespace {
          struct initializer {
            initializer() {
              std::cout << "Loading the library" << std::endl;
            }

            ~initializer() {
              std::cout << "Unloading the library" << std::endl;
            }
          };
          static initializer i;
        }
        ```


## Memory allocation

* <https://en.wikipedia.org/wiki/Memory_management>
* <https://en.wikipedia.org/wiki/C_dynamic_memory_allocation>
  (and specifically the [References](https://en.wikipedia.org/wiki/C_dynamic_memory_allocation#References)
  section).
* [Process memory concepts, chap. 3.1](https://www.gnu.org/software/libc/manual/html_node/Memory-Concepts.html)
  from the [GNU C library](https://www.gnu.org/software/libc/manual/html_node/index.html#SEC_Contents)
  documentation.
* [Memory Management for System Programmers – 2005](http://www.enderunix.org/docs/memory.pdf)
  by [Baris Simsek](http://www.enderunix.org/simsek/)
(* [Buddy memory allocation](https://en.wikipedia.org/wiki/Buddy_memory_allocation)
* [Slab allocation](https://en.wikipedia.org/wiki/Slab_allocation)
* [SLOB _Simple List Of Blocks_](https://en.wikipedia.org/wiki/SLOB)
* <https://en.wikipedia.org/wiki/Sbrk>
* <https://en.wikipedia.org/wiki/Mmap>
* <http://stackoverflow.com/a/3479496/643087> & <http://stackoverflow.com/a/3479570/643087>
* <http://www.flounder.com/inside_storage_allocation.htm>
* <http://www.memorymanagement.org/> and specifically their huge
  [bibliography](http://www.memorymanagement.org/bib.html#bibliography).
* [CSc 553 — Principles of Compilation – Christian Collberg – 2011 (cs.arizona.edu)](https://www2.cs.arizona.edu/~collberg/Teaching/553/2011/)


### Implementations

* Heap-based ;
* [Doug Lea's __dlmalloc__](http://g.oswego.edu/dl/html/malloc.html)
  ([__malloc.c__](ftp://g.oswego.edu/pub/misc/malloc.c)) ;
  [Musl-libc's impl.](http://git.musl-libc.org/cgit/musl/tree/src/malloc)
* Musl-libc also has a `static void *__simple_malloc(size_t n)`
  ([lite\_malloc.c](http://git.musl-libc.org/cgit/musl/tree/src/malloc/lite_malloc.c)),
  most probably used internally.
* Glibc's ptmalloc (variant of dlmalloc).
* phkmalloc, jemalloc
* __tcmalloc :__ “Thread-caching malloc”
  [github.com/gperftools/gperftools](https://github.com/gperftools/gperftools)
* [Hoard's malloc](https://en.wikipedia.org/wiki/Hoard_memory_allocator)
  ([official web site](http://www.hoard.org/about/),
  [emeryberger.github.io/Hoard/](http://emeryberger.github.io/Hoard/)).
* [OpenSIPS has 2/3 (maybe 4 ?) custom implementation for fast memory allocation](https://github.com/OpenSIPS/opensips/tree/master/mem)

__Note:__ Mac OS X do not support `brk()` since it is based on FreeBSD ;
instead it provides a libc wrapper [brk.c](https://opensource.apple.com/source/Libc/Libc-763.12/emulated/brk.c) .

* [Course slides : Inner Workings of Malloc and Free &ndash; Prof. David August &ndash; COS 217](https://www.cs.princeton.edu/courses/archive/fall07/cos217/lectures/14Memory-2x2.pdf), and specifically the slides about
[Malloc Algorithms](https://www2.cs.arizona.edu/~collberg/Teaching/553/2011/Handouts/Handout-6.pdf).

#### Other impl.

* [memory-mgr.sourceforge.net](http://memory-mgr.sourceforge.net/),
  old thing ~2012, sort of custom memory manager for C++.


## Signals

* [§24. Signal Handling – GNU libc library doc.](https://www.gnu.org/software/libc/manual/html_node/Signal-Handling.html)
* <http://www.enderunix.org/simsek/articles/signals.pdf>


## Binaries

* <https://en.wikipedia.org/wiki/Executable_and_Linkable_Format>


## Linux kernel

* [0xax.gitbooks.io/linux-insides/](https://0xax.gitbooks.io/linux-insides/content/index.html)
  __must-read__


## Graphics

### DRM & KMS

* [Linux DRM Mode-Setting API | Sept. 2012, by David Herrmann](https://dvdhrm.wordpress.com/2012/09/13/linux-drm-mode-setting-api/)
  (see also the comments; found when reading about [(the now defunct) KMSCon / 2012](https://dvdhrm.wordpress.com/2012/12/10/kmscon-introduction/); note that there was also a subsequent [(now defunct `libsystemd-terminal` and `systemd-consoled`)](https://github.com/systemd/systemd/pull/747)).
* [David Herrmann's KMS howto nicely documented `modeset.c`](https://github.com/dvdhrm/docs/blob/master/drm-howto/modeset.c)
  – see also [a couple of other .c – `/master/drm-howto/`](https://github.com/dvdhrm/docs/tree/master/drm-howto)
* __Man pages :__
  - [drm](https://www.commandlinux.com/man-page/man7/drm.7.html)
  - [drm-memory, drm-mm, drm-gem, drm-ttm](https://www.commandlinux.com/man-page/man7/drm-gem.7.html)
  - [drm-kms – Kernel Mode-Setting](https://www.commandlinux.com/man-page/man7/drm-kms.7.html)
* [__libdrm__ Git repository](https://cgit.freedesktop.org/mesa/drm/)


## Other scaffholdings

* [Detailed description of C & C++ program startup – “Executable and Linkable Format (ELF)” [archive.org copy of ~charngda @ www.acsu.buffalo.edu – dated somewhere back 2010]](http://www.cs.stevens.edu/~jschauma/810/elf.html)
  ([copy](http://www.cs.stevens.edu/~jschauma/810/elf.html))
* [Mini FAQ about the misc libc/gcc crt files – crt.txt by ~vapier @ dev.gentoo.org](http://dev.gentoo.org/~vapier/crt.txt)
* [How initialization functions are handled – §17.20](https://gcc.gnu.org/onlinedocs/gccint/Initialization.html) from [Target description macros and functions – §17](https://gcc.gnu.org/onlinedocs/gccint/Target-Macros.html) of [the GNU compilers internals manual](https://gcc.gnu.org/onlinedocs/gccint/index.html).
* [Old paper (1999) about the ELF file format: ELF.txt – “Portable Formats Specification, version 1.1 – Tool Interface Standard (TIS)”](http://www.muppetlabs.com/~breadbox/software/ELF.txt)
* [The Cerberus ELF Interface (by mayhem@devhell.org) – phrack.org – Volume 0x0b, Issue 0x3d, Phile #0x08 of 0x0f](http://phrack.org/issues/61/8.html)
* [2011: Load-time relocation of shared libraries (Eli Bendersky's)](https://eli.thegreenplace.net/2011/08/25/load-time-relocation-of-shared-libraries/)
* [2016: Shared Libraries: Understanding Dynamic Loading - http://amir.rachum.com](http://amir.rachum.com/blog/2016/09/17/shared-libraries/#direct-dependencies)


## TODO :

* Have an online viewer for the Linux man pages _versus_ the BSD man pages,
  and eventually the related POSIX specifications (+wikipedia etc).


## Other pointers

* [“I Contribute to the **Windows Kernel**. We Are Slower Than Other Operating Systems. Here Is Why.” (by __Marc Bevand__, 10 May 2013)](http://blog.zorinaq.com/i-contribute-to-the-windows-kernel-we-are-slower-than-other-oper/)




__EOF__
