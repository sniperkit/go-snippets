---
layout: note
title: "Quick count project lines count"
tagline: ""
published: true
category : notes
tags : [note, cli, metrics]
---

    $ find -type f -name \*.rs -print0 | wc -l --files0-from=-
