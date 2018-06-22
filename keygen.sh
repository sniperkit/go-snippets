#!/bin/sh

if [ ! -f "/etc/ssh/ssh_host_rsa_key" ]; then
  ssh-keygen -f /etc/ssh/ssh_host_rsa_key -N '' -t rsa > /dev/null
fi

if [ ! -f "/etc/ssh/ssh_host_dsa_key" ]; then
  ssh-keygen -f /etc/ssh/ssh_host_dsa_key -N '' -t dsa > /dev/null
fi

if [ ! -f "/etc/ssh/ssh_host_ecdsa_key" ]; then
  ssh-keygen -f /etc/ssh/ssh_host_ecdsa_key -N '' -t ecdsa > /dev/null
fi

if [ ! -f "/etc/ssh/ssh_host_ed25519_key" ]; then
  ssh-keygen -f /etc/ssh/ssh_host_ed25519_key -N '' -t ed25519 > /dev/null
fi

mkdir -m 700 /root/.ssh
echo "${AUTHORIZED_KEYS}" > /root/.ssh/authorized_keys
chmod 600 /root/.ssh/authorized_keys
