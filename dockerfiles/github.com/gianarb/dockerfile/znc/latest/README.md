```
docker run -it --name znc -p 6060:6060 gianarb/znc
```
Default username `testuser` password `testuser`

## override configuration
```
docker run -it --rm gianarb/znc znc --makeconf
```
Print new configuration on stdout, copy in a new file and share this file on
this path `/home/user/.znc/configs/znc.conf`
If you change the web admin panel port remember to bind another port and not the `6060`.

see you on IRC.
