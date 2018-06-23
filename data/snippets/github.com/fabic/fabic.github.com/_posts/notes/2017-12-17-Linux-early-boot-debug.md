---
layout: note
title: "Linux debug crash problem at boot with netconsole (failed)"
tagline: "Failed report of me trying to get the netconsole kernel module to work..."
published: true
category : notes
tags : [Linux, crash, boot]
---

__TL;DR:__ Couln't get the `netconsole` kernel module to effectivelly send logs to my other computer -\_- &ndash; __EDIT:__ Turns out it was my graphics card (!) o\_O\` which btw. was also causing some strange failed SATA commands sent to one of my hard-disk &ndash; no more system hangs, and the good thing is I didn't loose a costly hard-disk ;-.

__Story:__ -\_- Some strange crash problem occured today where my screen would
go blank early on during boot &ndash; did reboot in recovery mode and investigated
&ndash; `journalctl -b1` revealed a few big kernel __Oops__ with long stacktraces,
but I couldn't infer if and which module (if any) was causing it.

Hence I spent some time fiddling around, trying to get `netconsole` to work, and
ultimately quit (for I had plans early this morning upon wakeup, to code or whatever).

It turned out that the problem _disappear(s/ed)_ when I boot with the new Grub
menu entry below, maybe this is just chance (or bad luck for non-reproducibility
of things).

Had to generate another __initrd__ image for the Arch-provided one didn't have
the `netconsole.ko` kernel module in it :

    # vim /etc/mkinitcpio.conf
    MODULES="configfs atl1c netconsole"

    # mkinitcpio -g /boot/initramfs-linux-fabic.img

Edited `/boot/grub/grub.cfg` and duplicated one of the entries there
(note that __grub.cfg__ should not be edited by hand for it may be regenerated) :

    linux  /boot/vmlinuz-linux \
      root=UUID=09022016-0924-1909-2411-2709cade0da3 rw \
      resume=UUID=b09ea959-d8f7-400c-a829-91445640746c  \
      debug pause_on_oops=3 earlymodules=configfs,atl1c,netconsole break  \
      netconsole=7777@192.168.2.2/eth0,7777@192.168.2.1/e8:9d:87:53:61:24 \

    initrd  /boot/initramfs-linux-fabic.img

* `debug` will have the kernel not filter its messages (ignoring `loglevel` (?)).
* `pause_on_oops=3` is supposed to have the kernel pause upon emitting a __Oops__
so that one may read something.
* `break` have the __initrd__ drop us with a BusyBox shell (ArchLinux only?).
* `netconsole=...` is supposed to configure the corresponding module.

Also tried to configure this module from the __BusyBox__ shell that `break` left
me with :

    # modprobe netconsole \
        netconsole=4678@192.168.2.2/enp4s0,7777@192.168.2.1/e8:9d:87:53:61:24

On the other, target machine :

    $ netcat -u -l -p 7777

But nothing was received, briefly checked my recent __iptables__ rules, seemed okay.

* Note: Specifying the _early name_ __eth0__ instead of __enp4s0__ had the effect of
preventing the interface from, well, being renamed.  So there's a chance that __netconsole__
grabbed the interface for doing its job.
* [Netconsole @ ArchLinux wiki](https://wiki.archlinux.org/index.php/Netconsole)
* [kernel.org/doc/Documentation/networking/__netconsole.txt__](https://www.kernel.org/doc/Documentation/networking/netconsole.txt)
* [Mkinitcpio @ Arch wiki](https://wiki.archlinux.org/index.php/Mkinitcpio)

__EOF__
