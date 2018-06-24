#!/bin/sh
if [ -n "$UUID" ];then
    sed -i "s/UUID/$UUID/g" /etc/v2ray/config.json
    echo "You've set an UUID"
    echo "The UUID is: $UUID"
else
    export UUID=$(uuidgen)
    sed -i "s/UUID/$UUID/g" /etc/v2ray/config.json
    echo "You didn't set an UUID"
    echo "The auto-generated UUID is: $UUID"
fi

sed -i "s/PORT/$PORT/g" /etc/v2ray/config.json
v2ray -config=/etc/v2ray/config.json
