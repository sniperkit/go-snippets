#!/bin/sh

sh keygen.sh
sh pwgen.sh

/usr/sbin/sshd -D
