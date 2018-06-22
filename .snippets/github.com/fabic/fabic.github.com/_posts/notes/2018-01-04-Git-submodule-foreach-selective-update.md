---
layout: note
title: "Git submodule foreach selective update"
tagline: "Using the foreach Git command for fetching submodules updates."
published: true
category : notes
tags : [note, Git, submodules]
---

It is possible to use Bash expressions as arguments to `git submodule foreach ...`
Hence given a project with loads of dependencies checked-out as submodules,
updating only a subset of these that reside in, say `vendor/extra/`, may be
achieved with :

```bash
cd vendor/extra/ &&
  git submodule foreach \
    '[[ $path =~ ^../ ]] || (time git remote update)'
```

Enter each submodule that live under the current directory, and checkout
branch `release_60` that each have (it's LLVM/Clang btw.) ; __upon error__
we're dropped into a Bash shell for fixing things :

```bash
git submodule foreach       \
  '[[ $path =~ ^../ ]] ||   \
                            \
      (git checkout -b release_60 origin/release_60 || bash)'
```
