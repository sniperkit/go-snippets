---
layout: note
title: "Discover network MTU with ping"
tagline: "Find out the maximum tranfer unit to a given host, using ping with fragmentation disabled."
published: true
category : notes
tags : [note, networking]
---

Decrease the value of `-s 1500` :

```bash
ping -M do -s 1500 -c 1 vps.fabic.net
PING vps.fabic.net (91.134.136.248) 1500(1528) bytes of data.
ping: local error: Message too long, mtu=1500

--- vps.fabic.net ping statistics ---
1 packets transmitted, 0 received, +1 errors, 100% packet loss, time 0ms
```

Until we get a ping reply :

```bash
$ ping -M do -s 1464 -c 1 vps.fabic.net
PING vps.fabic.net (91.134.136.248) 1464(1492) bytes of data.
1472 bytes from 248.ip-91-134-136.eu (91.134.136.248): icmp_seq=1 ttl=50 time=239 ms

--- vps.fabic.net ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 239.746/239.746/239.746/0.000 ms
```

Where `man ping` for `-M pmtudisc_option` :

> Select Path MTU Discovery strategy.  `pmtudisc_option` may be either :
>
>    * `do` (prohibit fragmentation, even local one) ;
>    * `want` (do PMTU discovery, fragment locally when packet size is large) ;
>    * or `dont` (do not set DF flag).

Found at [ArchLinux: OpenVPN](https://wiki.archlinux.org/index.php/OpenVPN#Configure_the_MTU_with_Fragment_and_MSS).

__EOF__
