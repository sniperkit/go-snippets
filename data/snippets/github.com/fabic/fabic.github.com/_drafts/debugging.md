---
layout: post
title: "Draft - debugging"
tagline: "Supporting tagline"
category : notes
tags : [draft]
---

* [Capstone – The Ultimate Disassembler](http://www.capstone-engine.org/) [github.com/aquynh/capstone](https://github.com/aquynh/capstone)
* [Voltron](https://github.com/snare/voltron) – written in Python, uses Capstone.

## GDB

* <http://www.yolinux.com/TUTORIALS/GDB-Commands.html>

## LLDB

* <https://lldb.llvm.org/tutorial.html>

### Voltron

#### Install

    $ ./install.sh
    ...
    ==============================================================
    Installed for GDB (/home/fabi/dev/llvm-clang/local/bin/gdb):
      Python:             /usr/bin/python3.6
      Packages directory: /usr/lib/python3.6/site-packages
      Added voltron to:   /home/fabi/.gdbinit
      Entry point:        /usr/lib/python3.6/site-packages/voltron/entry.py
    Installed for LLDB (/usr/bin/lldb):
      Python:             /usr/bin/python2.7
      Packages directory: /usr/lib/python2.7/site-packages
      Added voltron to:   /home/fabi/.lldbinit
      Entry point:        /usr/lib/python2.7/site-packages/voltron/entry.py

It added itself to `~/.lldbinit` :

    $ cat /home/fabi/.lldbinit
    command script import /usr/lib/python2.7/site-packages/voltron/entry.py

As well as `~/.gdbinit` :

    $ cat /home/fabi/.gdbinit
    source /usr/lib/python3.6/site-packages/voltron/entry.py

##### Running Voltron

    $ gdb build/src/nstd/test_length

    GNU gdb (GDB) 7.11.50.20160803-git
    ...
    Voltron loaded.
    Reading symbols from build/src/nstd/test_length...(no debugging symbols found)...done.
    (gdb)

### Other / older stuff

* <https://code.google.com/archive/p/hyperdbg/>

__EOF__
