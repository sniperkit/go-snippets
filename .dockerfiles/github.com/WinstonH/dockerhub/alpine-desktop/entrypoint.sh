#!/bin/sh

#prepare sshd config
sudo mkdir -p /var/run/sshd
sudo sed -i "s/UsePrivilegeSeparation.*/UsePrivilegeSeparation no/g" /etc/ssh/sshd_config \
&& sudo sed -i "s/UsePAM.*/UsePAM no/g" /etc/ssh/sshd_config \
&& sudo sed -i "s/PermitRootLogin.*/PermitRootLogin yes/g" /etc/ssh/sshd_config \
&& sudo sed -i "s/#AuthorizedKeysFile/AuthorizedKeysFile/g" /etc/ssh/sshd_config \
&& sudo sed -i "s/#X11Forwarding no/X11Forwarding yes/g" /etc/ssh/sshd_config \
&& sudo sed -i "s/#PermitUserEnvironment no/PermitUserEnvironment yes/g" /etc/ssh/sshd_config \
&& sudo echo "ForwardX11Trusted yes" >> /etc/ssh/ssh_config
# generate host keys if not present
sudo ssh-keygen -A
#prepare xauth
sudo touch /root/.Xauthority

sudo /usr/bin/supervisord -c /etc/supervisord.conf
exec "$@"
