---
layout: post
title: "Draft - template"
tagline: "Supporting tagline"
category : notes
tags : [draft]
---

Canonical frame address, or control frame information, or call frame information :
_what on earth is this actually used for ?_

Stumbled upon <http://git.musl-libc.org/cgit/musl/tree/tools/add-cfi.x86_64.awk>

    CC_CMD = $(CC) $(CFLAGS_ALL) -c -o $@ $<

    # Choose invocation of assembler to be used
    ifeq ($(ADD_CFI),yes)
      AS_CMD = LC_ALL=C awk -f $(srcdir)/tools/add-cfi.common.awk \
                  -f $(srcdir)/tools/add-cfi.$(ARCH).awk $< \
                  | $(CC) $(CFLAGS_ALL) -x assembler -c -o $@ -
    else
      AS_CMD = $(CC_CMD)
    endif


* <http://stackoverflow.com/a/7535848/643087>
* <http://stackoverflow.com/a/23530662/643087>
* <https://en.wikipedia.org/wiki/Call_stack#Structure>
* Which is not to the same thing as “Control-Flow Integrity”, like in
  [Control-Flow Integrity – Principles, Implementations, and Applications – 2007](http://www.cse.usf.edu/~ligatti/papers/cfi-tissec.pdf)

__EOF__
