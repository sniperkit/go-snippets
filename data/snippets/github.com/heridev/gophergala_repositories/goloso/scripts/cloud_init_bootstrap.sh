#!/bin/bash 

curl="curl -s --connect-timeout 5"

nsqd_address='127.0.0.1:4151'
uuid=`uuidgen`
# ip_address=`$curl http://169.254.169.254/latest/meta-data/local-ipv4`
# instance_id=`$curl http://169.254.169.254/latest/meta-data/instance-id`

ip_address="127.0.0.10"
instance_id="i-806086"

if [ -e /etc/system-release ] ; then
  os=ami
elif [ -e /etc/debian_version ] ; then
  os=ubuntu
else
  os=unknown
fi

# payload="{\"event\":\"bootstrap\",\"uuid\":\"${uuid}\",\"instance\":{\"id\":\"${instance_id}\",\"ipaddress\":\"${ip_address}\",\"os\":\"${os}\"}}"
payload="{ \"event\":\"bootstrap\",\"uuid\":\"${uuid}\",\"instanceid\":\"${instance_id}\",\"ipaddress\":\"${ip_address}\",\"os\":\"${os}\"}"


ok=""
for count in 1 2 3 4 5 6 7 8 9 10 ; do
  ok=`$curl -d "$payload" -X POST http://${nsqd_address}/pub?topic=ec2`
  if [ "$ok" != "OK" ] ; then
    sleep 1;
  else
    break
  fi
done

if [ "$ok" != "OK" ] ; then
  touch /tmp/bootstrap-failed
fi
