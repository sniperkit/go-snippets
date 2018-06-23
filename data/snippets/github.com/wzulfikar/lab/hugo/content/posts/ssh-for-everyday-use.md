---
title: "SSH for Everyday Use"
date: 2018-03-07T17:19:17+08:00
tags: [""]
draft: true
---

`TODO: write content` 


private key starts 
```
-----BEGIN RSA PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: AES-128-CBC,BC89BA4DC3C8888DC41FDE0A302F5197

ozWvbvi11uIMpjTkWZwuHflvNR3gm8+4zSEHdTaCHVOM/nS6DSoauR3DhNWOZ3e0
kQ2NLbtXMqlrdEU77S3Q3JJQNrb5CfRstDiU811Yhp9URwdbgy6xHnHdHDViCW+C
... redacted ...
kQ2NLbtXMqlrdEU77S3Q3JJQNrb5CfRstDiU811Yhp9URwdbgy6xHnHdHDViCW+C
ozWvbvi11uIMpjTkWZwuHflvNR3gm8+4zSEHdTaCHVOM/nS6DSoauR3DhNWOZ3e0
-----END RSA PRIVATE KEY-----
```

public key starts with (open ssh) "ssh-rsa" + key + comment
- portion of public key is a "comment". you can remove it, or use it to identify info related to the public key.
```
ssh-rsa AAAAB3NzaC1yc2EAAAA...(redacted)...Fi9wrf+M7Q== me@mycomputer.local
```

change ssh private password (rarely). can also use to add password to existing private key:
`ssh-keygen -p -f ~/.ssh/id_rsa`

regenerating public key from private key: `ssh-keygen -t dsa -y > ~/.ssh/id_dsa.pub`

common error: 
- `chmod 400 mykey.pem` (permission) (or chmod 600?)
- reset host (remove from known_host): `ssh-keygen -R {hostname}`, ie. `ssh-keygen -R server.mydomain.com`


https://blogs.oracle.com/pranav/how-to-send-message-to-users-logged-on-to-a-unix-terminal

<p class="text-center">***</p>

*Outline:*

1. difference public vs private key
    - how to easily tell if it's private or public
2. securing private key with password
3. adding public key to authorized_keys
4. revoking private key
5. 