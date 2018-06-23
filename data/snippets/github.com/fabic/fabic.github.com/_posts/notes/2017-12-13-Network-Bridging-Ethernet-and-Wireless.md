---
layout: note
title: "Network bridge ethernet and wireless"
tagline: ""
published: false
category : notes
tags : [note, network]
---

* [`nmcli` doc. @ gnome.org](https://developer.gnome.org/NetworkManager/stable/nmcli.html)

    # nmcli connection add type bridge ifname br0 \
        con-name "bridge-br0" \
        autoconnect yes       \
        save yes              \
        stp no

    Connection 'bridge-br0' (1c1e0df3-770b-49dd-a74c-4b1431b90934)
      successfully added.

    # nmcli connection add type bridge-slave ifname wlp4s0 master br0 \
        con-name "bridge-br0-wifi"

    Connection 'bridge-br0-wifi' (8b5c1c18-d6f1-4683-9fb4-5c591d26bb43)
      successfully added.

    # nmcli connection add type bridge-slave ifname enp0s25 master br0 \
        con-name "bridge-br0-eth0"

    Connection 'bridge-br0-eth0' (410534a4-9d03-43eb-978a-fa4e20dfba26)
      successfully added.

    # nmcli conn up bridge-br0

    Connection successfully activated (master waiting for slaves)
      (D-Bus active path:
        /org/freedesktop/NetworkManager/ActiveConnection/93)

    # systemctl restart NetworkManager

* Note: [__brigde-utils__ was deprecated a while ago](https://lwn.net/Articles/435845/).

* [2016: Create network bridge with nmcli for libvirt &ndash; http://blog.leifmadsen.com](http://blog.leifmadsen.com/blog/2016/12/01/create-network-bridge-with-nmcli-for-libvirt/)
* [2016: Build a network bridge with Fedora &ndash; fedoramagazine.org](https://fedoramagazine.org/build-network-bridge-fedora/)

## OpenVPN

* `import type openvpn ...` solution found at <https://unix.stackexchange.com/a/359923>.

```bash
$ journalctl -f

$ sudo pacman -S networkmanager-openvpn

$ nmcli connection import type openvpn file client.ovpn
Connection 'client' (585587fd-1947-46dd-8810-d66175b6a917) successfully added.
```

### Use the GUI

```bash
$ nm-connection-editor
```

### Or from the command line :

```bash
$ nmcli connection show

NAME                           UUID                                  TYPE      DEVICE
Livebox-9e74                   c3e93cfa-e864-492d-9ce1-377bbf9730bd  wifi      wlp4s0
client                         585587fd-1947-46dd-8810-d66175b6a917  vpn       --
```

Connection was named `client` (inferred from the `client.ovpn` file name) :

```bash
$ nmcli connection show client
```

Rename connection :

```bash
$ nmcli conn modify client connection.id matynwavevpn
```

```bash
$ nmcli conn modify matynwavevpn \
    connection.autoconnect no    \
    ipv6.method ignore           \
    ipv4.dns "192.9.200.188,192.9.200.251" \
    ipv4.dns-search nwavedigital.local     \
    ipv4.ignore-auto-routes yes            \
    ipv4.never-default yes
```

* `ipv4.never-default yes` will prevent that NM adds a "default route" entry
* and `ipv4.ignore-auto-routes yes` is questionable here.

Activate connection :

```bash
$ nmcli connection up matynwavevpn
Connection successfully activated (D-Bus active path: /org/freedesktop/NetworkManager/ActiveConnection/2)
```

```bash
$ ip a

$ ip route
```

```bash
$ cat /etc/nsswitch.conf
hosts: files mymachines mdns4_minimal [NOTFOUND=return] resolve [!UNAVAIL=return] dns myhostname
```

```bash
$ dig @192.9.200.188 scripts5.nwavedigital.local
$ ping scripts5
$ systemctl stop avahi-daemon
```
