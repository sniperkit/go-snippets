---
title: "Ban Network Users Using Bettercap"
date: 2018-03-25T23:40:17+08:00
tags: ["networking"]
draft: true
---

![bettercap logo](/images/bettercap/bettercap-logo-sm.png#featured)

`TODO: write content` 






<p class="text-center">***</p>

*Outline:*

1. run bettercap
2. activate net.probe module to start probing for hosts in network: `net.probe on`
3. show hosts in network: `net.show`
4. start banning all hosts (except you): `arp.ban on`
    - run curl from seagate: timeout
    - run curl from current machine: pass
5. ban specific host:`set arp.spoof.targets 192.168.0.102; arp.ban on`
6. arp.spoof docs https://github.com/bettercap/bettercap/wiki/arp.spoof
7. think of arp.ban  as netcut (http://www.arcai.com/what-is-netcut/)