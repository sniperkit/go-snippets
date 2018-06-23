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

sed -i "s/PORT/$PORT/g" /etc/nginx/conf.d/default.conf
v2ray -config=/etc/v2ray/config.json
