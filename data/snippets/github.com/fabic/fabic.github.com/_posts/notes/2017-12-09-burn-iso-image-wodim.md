---
layout: note
title: "Burn ISO image to CD / DVD"
tagline: ""
published: true
category : notes
tags : [note, data transfer, cli, image]
---

Fast blank the media (optional, possibly superfluous) :

    $ time \
        wodim -v speed=1 -sao dev=/dev/sr0 blank=fast

Burn the ISO image file to disk :

    $ time \
        wodim -v speed=1 -sao dev=/dev/sr0 \
          archlinux-2017.12.01-x86_64.iso

Verify data by computing a SHA-1 fingerprint.  Arg. `count=$(( ... ))` evaluates
the amount of 2048-bytes blocks that the source ISO file have.  This is _recommanded_
(as per Arch wiki) because the optical drive _may (or may not)_ read some
additional garbage.

    $ time \
        dd if=/dev/sr0 \
        bs=2048        \
        count=$(( $(du -b archlinux-2017.12.01-x86_64.iso | awk '{print $1}') / 2048 )) \
          | sha1sum -b


* `-V` (capital 'V') will display very verbose SCSI layer debugging details.
* `-scanbus` or `--devices` is supposed to find out usable devices, but didn't
  work -\_-
*
* [ArchLinux: Optical_disc_drive](https://wiki.archlinux.org/index.php/Optical_disc_drive#Erasing_CD-RW_and_DVD-RW)
