---
layout: page
title: "Concurrency, memory management"
tagline: ""
category : notes
tags : [drafts, memory allocation, concurrency]
published: true
---

## Misc.

* <https://en.wikipedia.org/wiki/Non-blocking_algorithm> &ndash; In computer science, an algorithm is called non-blocking if failure or suspension of any thread cannot cause failure or suspension of another thread;[1] for some operations, these algorithms provide a useful alternative to traditional blocking implementations. A non-blocking algorithm is lock-free if there is guaranteed system-wide progress, and wait-free if there is also guaranteed per-thread progress. __[...]__ With few exceptions, non-blocking algorithms use atomic read-modify-write primitives that the hardware must provide, the most notable of which is compare and swap (CAS). Critical sections are almost always implemented using standard interfaces over these primitives (in the general case, critical sections will be blocking, even when implemented with these primitives). Until recently, all non-blocking algorithms had to be written "natively" with the underlying primitives to achieve acceptable performance. However, the emerging field of software transactional memory promises standard abstractions for writing efficient non-blocking code. __[...]__ Much research has also been done in providing basic data structures such as stacks, queues, sets, and hash tables. These allow programs to easily exchange data between threads asynchronously.

* <https://en.wikipedia.org/wiki/Read-modify-write> &ndash; In computer science, read-modify-write is a class of atomic operations (such as test-and-set, fetch-and-add, and compare-and-swap) that both read a memory location and write a new value into it simultaneously, either with a completely new value or some function of the previous value. These operations prevent race conditions in multi-threaded applications. Typically they are used to implement mutexes or semaphores. These atomic operations are also heavily used in non-blocking synchronization.

* <https://en.wikipedia.org/wiki/Compare-and-swap> &ndash; In computer science, compare-and-swap (CAS) is an atomic instruction used in multithreading to achieve synchronization. It compares the contents of a memory location to a given value and, only if they are the same, modifies the contents of that memory location to a given new value. This is done as a single atomic operation. The atomicity guarantees that the new value is calculated based on up-to-date information; if the value had been updated by another thread in the meantime, the write would fail. The result of the operation must indicate whether it performed the substitution; this can be done either with a simple boolean response (this variant is often called compare-and-set), or by returning the value read from the memory location (not the value written to it).

* <https://en.wikipedia.org/wiki/Critical_section> &ndash; In concurrent programming, concurrent accesses to shared resources can lead to unexpected or erroneous behavior, so parts of the program where the shared resource is accessed is protected. This protected section is the critical section or critical region. It cannot be executed by more than one process. Typically, the critical section accesses a shared resource, such as a data structure, a peripheral device, or a network connection, that would not operate correctly in the context of multiple concurrent accesses.

* <https://en.wikipedia.org/wiki/Read%E2%80%93write_conflict> & <https://en.wikipedia.org/wiki/Write%E2%80%93read_conflict>

* <https://en.wikipedia.org/wiki/Timestamp-based_concurrency_control> - In computer science, a timestamp-based concurrency control algorithm is a non-lock concurrency control method. It is used in some databases to safely handle transactions, using timestamps.



## Software transactional memory (STM)

* <https://en.wikipedia.org/wiki/Software_transactional_memory>
  > In computer science, software transactional memory (STM) is a concurrency control mechanism analogous to database transactions for controlling access to shared memory in concurrent computing. It is an alternative to lock-based synchronization.
* [Transactional Memory in GCC](https://gcc.gnu.org/wiki/TransactionalMemory)
* [2012 - N3341 - Transactional Language Constructs for C++](http://www.open-std.org/jtc1/sc22/wg21/docs/papers/2012/n3341.pdf)
* [2012 - Draft Specification of Transactional Language Constructs for C++](https://3f993110-a-62cb3a1a-s-sites.googlegroups.com/site/tmforcplusplus/C%2B%2BTransactionalConstructs-1.1.pdf)


## Read-Copy-Update (RCU)

* <https://en.wikipedia.org/wiki/Read-copy-update>
  > In computer operating systems, read-copy-update (RCU) is a synchronization mechanism implementing a kind of mutual exclusion[note 1] that can sometimes be used as an alternative to a readers-writer lock. It allows extremely low overhead, wait-free reads. However, RCU updates can be expensive, as they must leave the old versions of the data structure in place to accommodate pre-existing readers. These old versions are reclaimed after all pre-existing readers finish their accesses.


__EOF__
