---
#layout: default
layout: page
title: "Networks: Linux IP stack, routing, iptables firewall, etc (notes)"
tagline: "Personal notes about Linux networking."
category : notes
tags : [draft, network, wide]
published: true
---

???

`host fabic.net`
`dig`
`curl`
`netcat` vs `socat`

Query the RIPE database with `whois <ip>`

```bash
$ ip rule list
```

## Monitoring

### VnStat

* [VnStat @ Arch Linux Wiki](https://wiki.archlinux.org/index.php/VnStat)
* [humdi.net/vnstat/](http://humdi.net/vnstat/) (official web page).

```bash
# pacman -S vnstat
```

Edit `/etc/vnstat.conf` and set the default interface (_optional_) :

```bash
# vim /etc/vnstat.conf  # EDIT: `Interface = eth0`
```

Start the service :

```bash
# systemctl start vnstat
```

Create the databases, one for each NIC :

```bash
# vnstat -u -i wlp4s0
# vnstat -u -i enp0s25
```

Query details ab
```
# vnstat -q  [-i wlp4s0]
```

Monitor live traffic :

```bash
# vnstat -l -i wlp4s0
```


## Multi-IPs hosts

### Source routing

At [SO](https://unix.stackexchange.com/a/111423) for a solution where we partition
the internet address space, with a specific local source address for each :

```bash
function setup_basic_source_routing() {
  ip route add  32.0.0.0/3 via 91.134.136.1 dev ens3 src 178.32.40.98
  ip route add  64.0.0.0/3 via 91.134.136.1 dev ens3 src 92.222.48.73
  ip route add 128.0.0.0/3 via 91.134.136.1 dev ens3 src 91.134.136.248
}
```

#### Multi-IPs on one NIC and source address selection

* __Solution:__ Had to resort to __Netfilter__ SNAT.
* __Script:__ [server-setup-multi-ips-source-based-routing.sh](https://github.com/fabic/bash-it/blob/master/fabic/bin/server-setup-multi-ips-source-based-routing.sh)
* Packet marking using __iptables/Netfilter__ didn't work (but should have),
  because it turned out (march 2018) that the source address for locally generated
  new outbound connection... was somehow already set upon entering the Mangle/OUTPUT
  chain.
    - Got help investigating the situation from pekster @ irc.freenode.net/#Netfilter (2018-03-31)
* [Chapter 11. Netfilter & iproute - marking packets @ tldp.org “Linux Advanced Routing & Traffic Control HOWTO”](http://www.tldp.org/HOWTO/Adv-Routing-HOWTO/lartc.netfilter.html)
* Closest solution from [SO](https://unix.stackexchange.com/a/201752), which didn't
  work for our specific case (of multiple IPs on one NIC).


Match packets starting a new connection, that have not yet been marked

```bash
iptables -t mangle -A OUTPUT -m conntrack --ctstate NEW \
  -m connmark --mark 0x00/0x30 \
  -m statistic --mode random --probability 0.333 \
  -j CONNMARK --set-mark 0x10/0x30  &&

iptables -t mangle -A OUTPUT -m conntrack --ctstate NEW \
  -m connmark --mark 0x00/0x30 \
  -m statistic --mode random --probability 0.500 \
  -j CONNMARK --set-mark 0x20/0x30  &&
```

Finally we may transfer (or not) the __connmark__ bits to the netfilter __fwmark__
with this rule :

```bash
# As per http://ipset.netfilter.org/iptables-extensions.man.html#lbCS :
# "Copy the ctmark to the nfmark. If a mask is specified, only those bits are copied."
iptables -t mangle -A OUTPUT -m conntrack --ctstate NEW \
  -m connmark \! --mark 0x00/0x30 \
  -j CONNMARK --restore-mark --ctmask 0x30
```

^ this is optional though, it depends if we'll be using the __fwmark__ later on
(see [MARK target](http://ipset.netfilter.org/iptables-extensions.man.html#lbDE)),
like for ex. with `ip rule from all fwmark 0x20 priority <P> table <N>`.
See [CONNMARK](http://ipset.netfilter.org/iptables-extensions.man.html#lbCS) about
the `--restore-mark --nfmask 0x... --ctmask 0x...` and the computation of the resulting
__nfmask = (nfmark & ~nfmask) ^ (ctmark & ctmask)__.

__TODO:__ About the 2nd `0.500` probability? somewhat related to the
[Monty Hall problem](https://en.wikipedia.org/wiki/Monty_Hall_problem) and
statistical packets classification (with e.g. __fwmark, connmark__) :

> Vos Savant's response was that the contestant should switch to the other door
> (vos Savant 1990a). Under the standard assumptions, contestants who switch
> have a 2/3 chance of winning the car, while contestants who stick to their
> initial choice have only a 1/3 chance.

## Pointers, references

*  [A Deep Dive into Iptables and Netfilter Architecture (2015, digitalocean.com)](https://www.digitalocean.com/community/tutorials/a-deep-dive-into-iptables-and-netfilter-architecture) __must read__
* [netfilter / iptables project documentation](https://netfilter.org/documentation/index.html)
* [man iptables-extensions](http://ipset.netfilter.org/iptables-extensions.man.html#index)

* [Linux Advanced Routing & Traffic Control HOWTO @ lartc.org](http://lartc.org/howto/)

* [ip-sysctl.txt @ kernel.org](https://www.kernel.org/doc/Documentation/networking/ip-sysctl.txt)
  about those `sysctl -a | grep '^net\.'` IPv4 and IPv6 kernel variables, e.g. :

    - ip\_forward
    - fwmark\_reflect
    - tcp\_syncookies
    - ip\_local\_port\_range
    - ip\_local\_reserved\_ports
    - ip\_unprivileged\_port\_start
    - forwarding
    - secure\_redirects
    - rp\_filter

* [#Netfilter @ Freenode.net “Best practices”](https://gist.github.com/Thermi/70c9d77dc96523885e81e3f86f59f587)
  (several pointers about debugging Linux nf/iptables problems, before asking questions).

* Linux' Netfilter packet flow graph :

    ![Netfilter-packet-flow.SVG](https://upload.wikimedia.org/wikipedia/commons/3/37/Netfilter-packet-flow.svg)
    ([Netfilter-packet-flow.PNG](http://inai.de/images/nf-packet-flow.png))

* About [route selection](http://linux-ip.net/html/routing-selection.html)
 [Guide to IP Layer Network Administration with Linux @ linux-ip-net ~ 2007](http://linux-ip.net/html/index.html),
  and "policy-based routing" :
  The Routing Policy Database (RPDB) (`ip rule ...`),
  the Forward Information Base (FIB) (`ip route show cache`),
  multiple routing tables ability (`ip route ... table <NNN|alias>`)
* About [source address selection](http://linux-ip.net/html/routing-saddr-selection.html)
  `ip route ... via ... src 12.13.14.1 dev ens3`

Finding out the source address that would be chosen for a given destination address :

```bash
ip route get 216.58.213.142 # google.com
```

## EOF
