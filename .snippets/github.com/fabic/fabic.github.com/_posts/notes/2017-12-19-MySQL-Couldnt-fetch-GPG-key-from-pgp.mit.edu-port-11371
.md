---
layout: note
title: "MySQL: Couldnt fetch GPG key from pgp.mit.edu port 11371"
tagline: "While trying to build __mysql-connector-cpp__"
published: true
category : notes
tags : [note, C++, MySQL, library]
---

This one command will attempt connection on port __11371__ at __pgp.mit.edu__ :

    gpg --keyserver hkp://pgp.mit.edu \
        --recv-keys A4A9406876FCBD3C456770C88C718D3B5072E1F5

Protocol scheme `hkp://` is for port 11371 :

    grep hkp /etc/services

Tought it might be my firewall, but looking at the output of `sudo watch -d -n1 iptables -nvL`
showed no packets being dropped at all.  Didn't go look at my home router, might
be that it filters out traffic (?).

Anyway, found out that one may request keys through port __80__ :

    gpg --keyserver hkp://pgp.mit.edu:80 \
        --recv-keys A4A9406876FCBD3C456770C88C718D3B5072E1F5

__Story:__ Trying to build [mysql-connector-cpp](https://dev.mysql.com/doc/connector-cpp/en/connector-cpp-installation-source-distribution.html)
on ArchLinux ([AUR](https://aur.archlinux.org/packages/mysql-connector-c%2B%2B/))
and `makepkg -si` fails with:

    ==> Verifying source file signatures with gpg...
        mysql-connector-c++-1.1.9.tar.gz ... FAILED (unknown public key 8C718D3B5072E1F5)
        ==> ERROR: One or more PGP signatures could not be verified!

