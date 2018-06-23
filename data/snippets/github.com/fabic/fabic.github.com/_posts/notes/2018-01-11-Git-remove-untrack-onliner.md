---
layout: note
title: "Git: remove untracked files onliner"
tagline: ""
published: true
category : notes
tags : [note, Git]
---

Yet another onliner for cleaning up untracked files from a checked-out Git tree.

```bash
$ git status -uall --porcelain \
    | grep -E '^\?\? +'        \
    | awk '{ print $2 }'       \
    | xargs -rt -d\\n rm -v
```

Or move to the trash instead with `gio trash ...` :

```bash
$ git status -uall --porcelain \
    | grep -E '^\?\? +'        \
    | awk '{ print $2 }'       \
    | xargs -rt -d\\n gio trash
```

__EOF__
