# go-system/utilities

A collection of System's Utility software.

## ssh-suite

1. mc-ssh : A SSH-Signer-Key generator. Supports DSA, ECDSA and RSA.

### mc-ssh

getting it:
```sh
go get github.com/maxymania/go-system/utilities/ssh-suite/mc-ssh
```

usage:
```sh
# rsa-key:
mc-ssh -rsa > rsa-key.pem
# dsa-key:
mc-ssh -dsa > dsa-key.pem
# ecdsa-key:
mc-ssh -ecdsa > ecdsa-key.pem
```


