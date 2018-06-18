---
title: "Using Bettercap (v2.0) Proxy"
date: 2018-03-24T02:10:11+08:00
tags: ["intro", "networking"]
draft: false
disableToc: true
---

## What is Bettercap?

![bettercap logo](/images/bettercap/bettercap-logo-sm.png#featured)

If you are into networking tools for ARP spoof, sniffing, and other MITM-like utilities, you've probably heard about ettercap[^1]. Bettercap, which *sounds* similar to ettercap, is a networking tool like ettercap. Some people like to relate the name "Bettercap" as _better ettercap_.

<!--more-->

Quoting from its github repository at https://github.com/bettercap/bettercap, it's _"The state of the art network attack and monitoring framework"_.

## Why Bettercap?

This is where I find it interesting. The first time I knew about Bettercap was when it still built using ruby (v1.0). In Feb 2018, Bettercap author, Simon Margeritelly ([@evilsocket](https://twitter.com/evilsocket)), reimplemented Bettercap (v2.0) using Go[^2]. With this reimplementation, Bettercap is now a cross-platform single binary program. Unlike its previous version where you've to install ruby to run it, in Bettercap v2.0, you just need to download the binary and you're all set! 😎

## Using Bettercap: HTTP Proxy

Now that we know how easy is it to install Bettercap, let's go ahead and try one of its feature: **HTTP proxy**.

Bettercap has built-in feature to create HTTP proxy and use javascript file to interact with the connection. To demonstrate this functionality, we'll:

1. Run Bettercap on local machine
2. Set proxy and js file to change server response
3. Test if the proxy is working using curl

*The first step*, running Bettercap in local machine (I'm using macOS in this case), is pretty straightforward. Download the binary for your OS[^3], and run it from your command line: `path/to/bettercap`. An interactive Bettercap session will appear, something like this:

![bettercap init](/images/bettercap/bettercap-init.jpg) 
<p class="figure-text">Figure 1: initializing bettercap (command line)</p>

*Going to step 2*, we'll need to set a proxy and create a js file to interact with proxied connection. Before activating the proxy, let's create a `proxy-script.js` file in our desktop:

```js
// file: ~/desktop/proxy-script.js

// called when script is loaded
function onLoad() {
    console.log( "PROXY SCRIPT LOADED" );
}

// called after a request is proxied and there's a response
function onResponse(req, res) {
    res.ContentType = "rogue-content";
    res.Body = "me was here";
}
```
<p class="figure-text">Figure 2: js code inside desktop/proxy-script.js</p>

Now that the proxy-script is in place, we can bring up the HTTP proxy by running this command (in bettercap session):

![bettercap proxy on](/images/bettercap/bettercap-proxy-on.jpg) 
<p class="figure-text">Figure 3: activating bettercap http.proxy module</p>

In above image, you should've noticed that the interpreter displayed `PROXY SCRIPT LOADED` after loading the proxy script. Yep, that specific string `PROXY SCRIPT LOADED` is the one we set inside our `console.log` of `onLoad()` function in `proxy-script.js` (see: figure 2).

The `set http.proxy.script` tells Bettercap where to look for our js file. `set http.proxy.address` sets which address we'd like to bind the Bettercap proxy. Lastly, when all the variables for http.proxy are set, we activate the module using `http.proxy on`.

<p class="text-center">***</p>

*We're reaching our last step, step 3!* 

We've successfully activated `http.proxy` module. If you've issue activating it, try run the bettercap command as administrator (or via `sudo`). In this step, we'll test if our proxy is actually working by sending a request using curl command. If you're on windows and you don't have `curl` installed, you can download it from here: https://curl.haxx.se/download.html.

Based on our code in `proxy-script.js` file, the respond we should get when sending request using proxy is `me was here` (see `onResponse` function, figure 2). To test that, we'll first send a direct request to google to see what's the original respond, and then we'll send another request to google using proxy. The GIF in Figure 4 below shows how it works:

![bettercap proxy gif](/images/bettercap/bettercap-proxy-google.gif) 
<p class="figure-text">Figure 4: Bettercap proxy in action</p>

*Voila!* The proxy is working, and sending request to google.com (and to any other domain) via our proxy will get the respond modified to `me was here`. Notice that when we send the proxied request, our bettercap session displays new *event stream* `http.proxy.spoofed-response`. To sum up above operation, here's the logical flow: 

```
→ user sends request with bettercap proxy 
  → bettercap proxy handle the request and modify it 
    → user gets modified respond 
```

## Closing

That's it, a short demonstration on how we utilize Bettercap HTTP proxy module to interact with proxied connection and modify the response. In real world scenario, above operation can be used by any attacker to *smuggle* malicious script in HTTP response –– and the user will still think that nothing has changed, instead of changing the whole response to `me was here`. 

Knowing how easy to do this should give us more awareness when using random proxy in internet. Just like what we've demonstrated here, yeah, the proxy providers can easily do funky stuffs as we browse the internet. *Unless you really trust, don't use them to browse your confidential informations.*

I'd like to write one or two more posts to demonstrate other Bettercap features. But in the meantime, take your time exploring [Bettercap wiki](https://github.com/bettercap/bettercap/wiki) and see for yourself what it's capable of.

***Till next. See ya!***

[^1]:https://en.wikipedia.org/wiki/Ettercap_(software)
[^2]:https://www.evilsocket.net/2018/02/27/All-hail-bettercap-2-0-one-tool-to-rule-them-all/
[^3]:https://github.com/bettercap/bettercap/releases
