---
layout: note
title: "OpenVPN Server setup"
tagline: ""
published: false
category : notes
tags : [note]
---

```bash
# pacman -S easy-rsa

# cd /etc/easy-rsa && 
  export EASYRSA=$(pwd) &&
  easyrsa init-pki &&
  easyrsa build-ca ; \
  echo RETVAL=$?
```

```bash
Note: using Easy-RSA configuration from: /etc/easy-rsa/vars

Your newly created PKI dir is: /etc/easy-rsa/pki

Note: using Easy-RSA configuration from: /etc/easy-rsa/vars
Generating a 2048 bit RSA private key
writing new private key to '/etc/easy-rsa/pki/private/ca.key.Rg3dMdgRYm'
Enter PEM pass phrase: XXXXX
Verifying - Enter PEM pass phrase: XXXXX
Common Name (eg: your user, host, or server name) [Easy-RSA CA]: fabic.net

CA creation complete and you may now import and sign cert requests.
Your new CA certificate file for publishing is at:
/etc/easy-rsa/pki/ca.crt
```

```bash
scp pki/ca.crt fabi@vps.fabic.net:~
## move that into /etc/openvpn/server + chown root
```
