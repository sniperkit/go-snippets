znc
===

ZNC running on Alpine Linux container

Set up:

```
$ docker pull danilo/znc

$ docker run -p 6697:6697 --name znc -d -it danilo/znc
```

Configuration:

```
$ docker exec -it znc su - irc -c 'znc --makeconf'
```

Run:

```
$ docker stop znc

$ docker start znc
```

Voilà!
