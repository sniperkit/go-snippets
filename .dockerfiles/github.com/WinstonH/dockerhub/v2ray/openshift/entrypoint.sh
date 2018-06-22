#!/bin/sh
if [ -n "$UUID" ];then
    sed -i "s/UUID/$UUID/g" /etc/v2ray/config.json
    echo "You've set an UUID"
    echo "The UUID is: $UUID"
else
    export UUID=$(v2ctl uuid)
    sed -i "s/UUID/$UUID/g" /etc/v2ray/config.json
    echo "You didn't set an UUID"
    echo "The auto-generated UUID is: $UUID"
fi
VERSION=$(v2ray --version |grep v |awk '{print $2}')
BUILDDATE=$(v2ray --version |grep v |awk '{print $5}')
REBOOTDATE=$(date)
sed -i "s/VERSION/$VERSION/g" /var/lib/nginx/html/index.html
sed -i "s/BUILDDATE/$BUILDDATE/g" /var/lib/nginx/html/index.html
sed -i "s/REBOOTDATE/$REBOOTDATE/g" /var/lib/nginx/html/index.html
nginx
v2ray -config=/etc/v2ray/config.json
