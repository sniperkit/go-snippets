---
title: "MITM à la Bettercap"
date: 2018-03-26T01:53:42+08:00
tags: ["networking"]
draft: true
---

`TODO: write content` 

set net.sniff.verbose false
net.sniff on

# http proxy
set http.proxy.script caplets/proxy-script-test.js
set http.proxy.sslstrip true
set http.proxy.port 80
http.proxy on

# https proxy
set https.proxy.script caplets/proxy-script-test.js
set https.proxy.sslstrip true
set https.proxy.port 443
https.proxy on

# redirect connections to our machine
arp.spoof on

<p class="text-center">***</p>

*Outline:*

1. 