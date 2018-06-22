### Docker-Alpine-Shadowsocks

#### Get Docker Container

```
$ docker pull ncc0706/alpine-shadowsocks
```

#### Run Default Container

Run with default `METHOD` of `aes-256-cfb` and default `PASSWORD` of `shadowsocks` :

-k password
-m method

```
$ docker run -d -p 8388:8388 --name ss ncc0706/alpine-shadowsocks  -k node001 -m aes-256-cfb
```
