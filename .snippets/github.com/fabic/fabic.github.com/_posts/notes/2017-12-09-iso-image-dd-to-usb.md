---
layout: note
title: "Write ISO image to USB stick (dd)"
tagline: ""
published: true
category : notes
tags : [note, data transfer, cli, image, USB]
---

    $ time \
        dd bs=4M status=progress \
           if=archlinux-2017.12.01-x86_64.iso \
           of=/dev/sdb && sync
