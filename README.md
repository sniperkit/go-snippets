# docker-opencv

This repository contains the x86_64, armhf, and arm64 architecture images of Ubuntu and Alpine distributions, including opencv3.4.1 and python3.6 .

It's very convenient to use raspberry pie, server or local to develop interesting programs with Python and OpenCV. You only need to take corresponding images.

## How to use it ？

```

root@ubuntu:~/docker-opencv# docker run --rm -it jermine/opencv:alpine-arm64-3.4.1 sh
Unable to find image 'jermine/opencv:alpine-arm64-3.4.1' locally
alpine-arm64-3.4.1: Pulling from jermine/opencv
e629e6244a09: Already exists 
aafc292a94ab: Already exists 
fce7233400d4: Pull complete 
66837379c661: Pull complete 
55195849748f: Pull complete 
72eb230e9210: Pull complete 
Digest: sha256:93910d3c5a05c57c21084161b821674d6d47fc7734f1550aa00ddd8482c1950f
Status: Downloaded newer image for jermine/opencv:alpine-arm64-3.4.1
/ # python3
Python 3.6.3 (default, Nov 22 2017, 21:32:12) 
[GCC 6.4.0] on linux
Type "help", "copyright", "credits" or "license" for more information.
imort>>> import cv2
>>> cv2.__version__
'3.4.1'
>>> exit()
/ # pip list
Package    Version
---------- -------
numpy      1.14.5 
pip        10.0.1 
setuptools 28.8.0 
six        1.11.0 
wheel      0.31.1 
/ # uname -a
Linux 8524802fbd9a 4.14.37-chainsx+stable #1 SMP PREEMPT Fri Apr 27 21:36:31 CST 2018 aarch64 Linux
/ # exit
root@ubuntu:~/docker-opencv# docker images
REPOSITORY                                  TAG                            IMAGE ID            CREATED             SIZE
jermine/opencv                              alpine-arm64-3.4.1             172e3578c9ee        6 minutes ago       310MB
nginx                                       latest                         9e0f046b6d6e        13 days ago         103MB
jermine/opencv                              alpine-armhf                   a80120167531        6 weeks ago         299MB
jermine/opencv                              alpine-armhf-3.4.1             a80120167531        6 weeks ago         299MB
weaveworks/weave-npc                        2.3.0                          982e879f62a9        2 months ago        49.4MB
weaveworks/weave-kube                       2.3.0                          75770647069a        2 months ago        97.8MB
weaveworks/weave-kube                       2weaveworks_weave-kube-arm64   75770647069a        2 months ago        97.8MB
weaveworks/weave-kube                       0weaveworks_weave-kube-amd64   f15514acce73        2 months ago        96.8MB
jermine/flannel                             v0.10.0-arm64                  122cdb7aa710        4 months ago        43.5MB
quay.io/coreos/flannel                      v0.10.0-arm64                  122cdb7aa710        4 months ago        43.5MB
mirrorgooglecontainers/kube-proxy-arm64     v1.9.2                         0f834c88e77d        5 months ago        112MB
gcr.io/google_containers/kube-proxy-amd64   v1.9.2                         0f834c88e77d        5 months ago        112MB
alpine                                      edge                           698b28d4f8d1        5 months ago        3.96MB
mirrorgooglecontainers/pause-arm64          3.1                            6cf7c80fe444        6 months ago        525kB
k8s.gcr.io/pause-arm64                      3.1                            6cf7c80fe444        6 months ago        525kB
```
