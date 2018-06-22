---
layout: note
title: "ArchLinux & ZFS: Upgrade system howto"
tagline: ""
published: true
category : notes
tags : [note, ZFS, cli, ArchLinux]
---

ZFS installed from the __AUR__ causes __Pacman__ to fail during dependency
resolution (due to the fact that the AUR packages are tied to a specific
(strict) Linux kernel version).  Bypass the problem by just ignoring the upgrade
of the Linux kernel :

    # pacman -Syu \
        --ignore spl-linux,spl-linux-headers,zfs-linux,zfs-linux-headers,linux,linux-headers,linux-firmware

* __TODO:__ Find out the correct way to upgrade the whole system, including ZFS
stuff, for ex. [DKMS](https://wiki.archlinux.org/index.php/Dynamic_Kernel_Module_Support)
* __And__ generally, find out how to _"correctly"_ manage your bunch of installed
AUR packages.
* [ZFS @ ArchLinux](https://wiki.archlinux.org/index.php/ZFS)
