---
layout: note
title: "Netfilter/iptables: watch command"
tagline: "A simple command for debugging firewall problems, or just monitor network traffic."
published: true
category : notes
tags : [note, cli, networking]
---

Dude, u do this for monitoring traffic as you fiddle with Linux' Netfilter/iptables :

```bash
~# watch -d -n1 'echo "--- MANGLE ---\n" ; iptables -t mangle -nvL ; echo "\n--- FILTER ---\n"; iptables -t filter -nvL ; echo "\n--- NAT ---\n"; iptables -t nat -nvL'
```
