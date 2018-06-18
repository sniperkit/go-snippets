---
title: "OSX, Time Machine, and SSHFS"
date: 2018-03-20T21:52:39+08:00
tags: [""]
draft: true
---

![time machine logo](/images/time-machine-mac-icon-sm.png#featured)
<p class="text-center"><i>*Time Machine Logo</i></p>

`TODO: write content` 



> The secret to creativity is knowing how to hide your sources. 
> -- <cite>adsf asdf [^1]</cite>

[^1]:http://www.quotedb.com/quotes/2112



<p class="text-center">***</p>

*Outline:*

1. no issue using external drive (partitioned)
2. keep resetting backup with NAS → changed 8mb bands to 128mb https://apple.stackexchange.com/questions/93215/how-can-i-avoid-repeated-time-machine-must-create-a-new-backup-errors-when-bac
3. no limit for backup with NAS and never finish first backup → create image, mount using sshfs, set mounted disk as tm destination 
4. easiest path could be using apple time capsule
5. seagate central 2TB https://www.cnet.com/products/seagate-central-series/review/
6. tplink AC1750 https://www.tp-link.com/us/products/details/cat-5508_RE450.html

limit size on nas: https://www.youtube.com/watch?v=Nq7mSizqUSI


limit using command but doesn't work: https://www.defaults-write.com/time-machine-setup-a-size-limit-for-backup-volumes/


sshfs https://www.digitalocean.com/community/tutorials/how-to-use-sshfs-to-mount-remote-file-systems-over-ssh

https://linode.com/docs/networking/ssh/using-sshfs-on-linux/


https://askubuntu.com/questions/777116/sshfs-giving-remote-host-has-disconnected

https://superuser.com/questions/709820/sshfs-is-failing-with-remote-host-has-disconnected

"(sudo) sshfs (-o allow_other) -o sshfs_debug user@computer:/mountpoint /localmountpoint"
