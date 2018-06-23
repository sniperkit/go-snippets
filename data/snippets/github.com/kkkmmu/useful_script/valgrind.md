#Valgrind's Tool Suite

The Valgrind distribution includes the following debugging and profiling tools:

##Memcheck Cachegrind Callgrind Massif Helgrind DRD Other Tools

# Memcheck
 Memcheck detects memory-management problems, and is aimed primarily at C and C++ programs. When a program is run under Memcheck's supervision, all reads and writes of memory are checked, and calls to malloc/new/free/delete are intercepted. As a result, Memcheck can detect if your program:

 Accesses memory it shouldn't (areas not yet allocated, areas that have been freed, areas past the end of heap blocks, inaccessible areas of the stack).
 Uses uninitialised values in dangerous ways.
 Leaks memory.
 Does bad frees of heap blocks (double frees, mismatched frees).
 Passes overlapping source and destination memory blocks to memcpy() and related functions.
 Memcheck reports these errors as soon as they occur, giving the source line number at which it occurred, and also a stack trace of the functions called to reach that line. Memcheck tracks addressability at the byte-level, and initialisation of values at the bit-level. As a result, it can detect the use of single uninitialised bits, and does not report spurious errors on bitfield operations. Memcheck runs programs about 10--30x slower than normal.

# Cachegrind
 Cachegrind is a cache profiler. It performs detailed simulation of the I1, D1 and L2 caches in your CPU and so can accurately pinpoint the sources of cache misses in your code. It identifies the number of cache misses, memory references and instructions executed for each line of source code, with per-function, per-module and whole-program summaries. It is useful with programs written in any language. Cachegrind runs programs about 20--100x slower than normal.

# Callgrind
 Callgrind, by Josef Weidendorfer, is an extension to Cachegrind. It provides all the information that Cachegrind does, plus extra information about callgraphs. It was folded into the main Valgrind distribution in version 3.2.0. Available separately is an amazing visualisation tool, KCachegrind, which gives a much better overview of the data that Callgrind collects; it can also be used to visualise Cachegrind's output.
# Massif
 Massif is a heap profiler. It performs detailed heap profiling by taking regular snapshots of a program's heap. It produces a graph showing heap usage over time, including information about which parts of the program are responsible for the most memory allocations. The graph is supplemented by a text or HTML file that includes more information for determining where the most memory is being allocated. Massif runs programs about 20x slower than normal.

# Helgrind
 Helgrind is a thread debugger which finds data races in multithreaded programs. It looks for memory locations which are accessed by more than one (POSIX p-)thread, but for which no consistently used (pthread\_mutex\_) lock can be found. Such locations are indicative of missing synchronisation between threads, and could cause hard-to-find timing-dependent problems. It is useful for any program that uses pthreads. It is a somewhat experimental tool, so your feedback is especially welcome here.

# DRD
 DRD is a tool for detecting errors in multithreaded C and C++ programs. The tool works for any program that uses the POSIX threading primitives or that uses threading concepts built on top of the POSIX threading primitives. While Helgrind can detect locking order violations, for most programs DRD needs less memory to perform its analysis.

# Lackey, Nulgrind
 Lackey and Nulgrind are also included in the Valgrind distribution. They don't do very much, and are there for testing and demonstrative purposes.

# Other Tools
 Several other Valgrind tools have been created: you can find a list of them on the Variants / Patches page
