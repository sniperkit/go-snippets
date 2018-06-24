#!/usr/bin/with-contenv ash


if [ -z "$DNS_INIT_TIMEOUT" ]
then
    DNS_INIT_TIMEOUT=45
fi

( sleep $DNS_INIT_TIMEOUT ; [ -f /var/run/dns.init ] || ( echo "DNS : Timed out setting up DNS." && s6-nuke && kill -2 1 ) ) &

echo "DNS : Initial Setup"
cp /etc/hosts /etc/hosts.orig
cp /etc/hosts /tmp/hosts

echo "DNS STEP 1 : Creating the dnsmasq-resolv.conf"
if ! ( cat /etc/resolv.conf | grep "nameserver 127.0.0.1" &> /dev/null )
then
    cp -f /etc/resolv.conf /etc/dnsmasq-resolv.conf
    echo "nameserver 127.0.0.1" > /etc/resolv.conf
    echo "nameserver 8.8.8.8" >> /etc/resolv.conf
fi

echo "DNS : Contents of dnsmasq-resolv.conf"
echo "-------------------"
cat /etc/dnsmasq-resolv.conf
echo
echo


echo "DNS STEP 2 : Adding the linked services"
. /bin/dns_hack.sh

sort -u < /tmp/hosts > /etc/hosts

echo "DNS : Initial /etc/hosts calculated"
echo "-------------------"
cat /etc/hosts
echo
echo

echo "DNS : initial work complete"
touch /var/run/dns.init

