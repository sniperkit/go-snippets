---
layout: post
title: "Mac OS X : how to forward X11 through SSH"
#tagline: "Supporting tagline"
category : notes
tags : [howto, macosx, x11, ssh, tunnel]
---

Ok, the X server is no longer shipped with Mac OS X;
at <https://support.apple.com/en-us/HT201341> (Oct 31, 2016) they say :

> X11 is no longer included with Mac, but X11 server and client libraries are available from the XQuartz project.

Head to [xquartz.org](https://www.xquartz.org/),
got the `XQuartz-2.7.11.dmg`,
mount it, run the XQuartz.pkg installer;
it mentions that it is meant for _Mac OS-X **10.6.3** or later_, okey;
it says I've got to log out and in again &ndash; indeed `$DISPLAy` is not defined
when opening new shells; did that, ok :

    [fabi@pc40.home] ~ $ echo $DISPLAY
    /private/tmp/com.apple.launchd.K9dCZkZZUW/org.macosforge.xquartz:0

## SSH tunneling is not working

But `ssh -X ...` is not able to setup X11 forwarding for some obscur reason
(and `-vv` gave no obvious clue about what's gone wrong).

## Disabling security with `xauth +` &ndash; ok

http://oroborosx.sourceforge.net/remotex.html#setdisplay
