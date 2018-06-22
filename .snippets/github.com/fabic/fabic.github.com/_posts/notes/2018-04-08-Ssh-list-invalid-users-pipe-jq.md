---
layout: post
title: "SSH: List failed password-based connections attempts (using jq)"
tagline: ""
category : notes
published: true
tags : [note, cli, security]
---

Basically one may extract the syslog message using __jq__ in this way :

```bash
$ journalctl -u ssh -o json-pretty -b | jq '.MESSAGE'
```

And we may list “failed password authentication” of valid users in this way :

```bash
$ journalctl -u ssh -o json -b \
    | jq '.MESSAGE | capture("^Failed password for (?<invalid>invalid user )?(?<user>\\w+)") | [.user, .invalid] | @sh' \
    | sed -e 's/"\(.*\)"/\1/' -e "s/'//g" \
    | sort -u -k1,1
    | grep -v 'invalid user *$'
```

__EOF__
