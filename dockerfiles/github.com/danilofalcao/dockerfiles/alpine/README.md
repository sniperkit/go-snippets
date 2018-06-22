alpine-openrc
=============

Alpine Linux base container using OpenRC

Set up:

```
$ docker pull danilo/alpine
```

Run:

```
$ docker run --name alpine -d -it danilo/alpine 
e15f68438e3ad8b2eb1f82b783b35b513fb829d1cee6de18cb121f248ac28ebb
$ docker exec -it alpine sh
/ # apk update
fetch http://dl-cdn.alpinelinux.org/alpine/edge/main/x86_64/APKINDEX.tar.gz
fetch http://dl-cdn.alpinelinux.org/alpine/edge/community/x86_64/APKINDEX.tar.gz
fetch http://dl-cdn.alpinelinux.org/alpine/edge/testing/x86_64/APKINDEX.tar.gz
v3.4.0-727-g7baefde [http://dl-cdn.alpinelinux.org/alpine/edge/main]
v3.4.0-686-g229be22 [http://dl-cdn.alpinelinux.org/alpine/edge/community]
v3.4.0-726-g849e69a [http://dl-cdn.alpinelinux.org/alpine/edge/testing]
OK: 9558 distinct packages available
/ # apk add nginx
(1/3) Installing nginx-common (1.10.1-r2)
Executing nginx-common-1.10.1-r2.pre-install
(2/3) Installing pcre (8.39-r0)
(3/3) Installing nginx (1.10.1-r2)
Executing busybox-1.24.2-r9.trigger
OK: 8 MiB in 15 packages
/ # rc-update add nginx
 * service nginx added to runlevel default
/ # service nginx start
 * Caching service dependencies ...                                                                                                                                            [ ok ]
 * /run/nginx: creating directory
 * /run/nginx: correcting owner                                                                                                                                                [ ok ]
 * Starting nginx ...                                                                                                                                                          [ ok ]
/ #
```

Voilà!
