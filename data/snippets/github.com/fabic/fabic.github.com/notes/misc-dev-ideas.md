---
layout: page
title: "Miscellaneous notes & dev. ideas"
tagline: "Big ideas, don't get any..."
category : notes
tags : [draft, wide]
published: false
---

_Buddy: you note down here, whichever idea that pops up in your mind..._

### Quick serve directories with .html

For serving e.g. Doxygen-generated documentation. I'd like to have a command that would basically serve a folder with HTML in it, like :

```bash
$ quickserve ~/dev/some-project/doc/Doxygen/html/ ~/dev/other-project/doc/html
```

As an added bonus it may be able to process lightweight markup language like Markdown or reStructuredText or whatever.

* [pandoc.org](http://pandoc.org/)

* `[ ]` Write a simple Node.js app for this ?  Possibly with Electron maybe ?

### CLI: How to test if a TCP port is available for binding or not ?

Use case: I have that `jekyll-server.sh` script which defaults to have Jekyll
bind on port 3838 &ndash; and I run several Jekyll instances : I need a simple
way of checking that a TCP port is available for binding or not.

__EOF__
