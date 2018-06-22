---
layout: post
title: "Draft - Mac OS X - Mount NFS share + insecure"
tagline: "Supporting tagline"
category : notes
tags : [draft]
---


    sudo systemctl restart nfs-server
    sudo systemctl restart rpcbind

    sudo systemctl status nfs-server

    showmount -e winterfell.local

    sudo systemctl enable nfs-server
    sudo systemctl enable rpcbind

    pc40:Volumes root# mount -t nfs winterfell.local:/srv/nfs/tous_les_disques_durs /Volumes/ordinateur-fabien/
    mount_nfs: can't mount /srv/nfs/tous_les_disques_durs from winterfell.local onto /Volumes/ordinateur-fabien: Operation not permitted

    (09:11:59) [fabi@winterfell] 0 ~ $ vi /etc/exports

    man exports

    (09:15:19) [fabi@winterfell] 0 ~ $ sudo exportfs -rav
    exporting 192.168.1.0/24:/srv/nfs/tous_les_disques_durs

    (09:15:34) [fabi@winterfell] 0 ~ $ cat /etc/exports
    /srv/nfs/tous_les_disques_durs  192.168.1.0/24(rw,no_subtree_check,nohide,no_root_squash,sync,insecure)


    (09:19:24) [fabi@winterfell] 0 ~ $ showmount -e winterfell.local
    Export list for winterfell.local:
    /srv/nfs/tous_les_disques_durs 192.168.1.0/24

* <https://thornelabs.net/2013/10/15/operation-not-permitted-mounting-nfs-share-on-os-x-mountain-lion.html>
  alternate solution is to pass `-o resvport` from MacOS when mounting NFS shares.


__EOF__
