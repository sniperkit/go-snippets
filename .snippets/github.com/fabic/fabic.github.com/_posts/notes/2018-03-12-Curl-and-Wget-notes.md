---
layout: note
title: "Curl and Wget notes"
tagline: ""
published: false
category : notes
tags : [note, cli, http]
---

Fetch headers only :

```bash
$ curl -I http://oiswww.eumetsat.org/IPPS/html/latestImages/EUMETSAT_MSGIODC_WV062_EastIndianOcean.jpg
$ wget -S -O /dev/null http://oiswww.eumetsat.org/IPPS/html/latestImages/EUMETSAT_MSGIODC_WV062_EastIndianOcean.jpg
```

Custom GET with an authorization header (and dump response headers) :

```bash
curl -D /dev/stderr -X GET \
  -H"Authorization: sso-key XXXXX:ZZZZZ" \
  "https://api.godaddy.com/v1/domains/"
```

See also :

* [Bash-it's Curl aliases](https://github.com/fabic/bash-it/blob/master/aliases/available/curl.aliases.bash)

__EOF__
