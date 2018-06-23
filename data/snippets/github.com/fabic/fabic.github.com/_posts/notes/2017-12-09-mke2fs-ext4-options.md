---
layout: note
title: "Write ISO image to USB stick (dd)"
tagline: ""
published: true
category : notes
tags : [note, filesystem]
---

    ~#  mkfs.ext4 -L archlinux -m 1 -U time -v /dev/sda3

* Note that default filesystem features (`-O ...`) are read from `/etc/mke2fs.conf`.
* `-L ...` : filesystem label.
* `-m 1` : 1% reserved blocks pct.
* `-U time` : generate UUID from date-time.
