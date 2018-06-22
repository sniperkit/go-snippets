---
layout: note
title: "Check all hard-drives SMART info."
tagline: ""
published: true
category : notes
tags : [note, hardware, HDD, cli]
---

Query S.M.A.R.T. information from all hard-disks :

    # for dev in /dev/sd? ; do
        echo -e "\n\n### DEV: $dev ###\n\n" ;
        smartctl -a $dev || break ;
        read -p "-- ^ $dev -- Continue ? --" ;
      done

* Use `smartctl -t short /dev/sdX` to initiate a short self-test
(`long` for extented self-tests).

__smartd__ may be configured (`/etc/smartd.conf`) to schedule automatic run of self-tests,
albeit resorting to some arcane configuration, here `(S/../.././12|L/../../6/19)`,
where (_I guess_) '__S__' stands for __Short__, and '__L__' for __Long__ :

    # /etc/smartd.conf

    # Short self-test everyday between 12:00 - 13:00 ;
    # and extended self-test on Saturdays btw. 19:00 20:00 :
    DEVICESCAN -a -o on -S on -n standby,q \
      -W 4,35,40                           \
      -s (S/../.././12|L/../../6/19)       \
      -m cadet.fabien+sys@gmail.com

* [Arch wiki: SMART](https://wiki.archlinux.org/index.php/S.M.A.R.T.)
* Note that :
    - The first `DEVICESCAN ...` encountered will stop __smartd__ from
      reading the rest of its configuration file.
    - Configuration can be specified _per_ disk, ex. `/dev/sda -a -W 2,35,40`.
    - And that disks may be specified by their UUID instead of their block device
      path `/dev/sdX` so that configuration is not affected by e.g. some internal
      re-wiring of the drives on the motherboard/controller(s),
      see `ls -l /dev/disk/by-uuid/` :

             /dev/disk/by-uuid/09022016-0924-1909-2411-2709cade0da3 -a -o on -S on
