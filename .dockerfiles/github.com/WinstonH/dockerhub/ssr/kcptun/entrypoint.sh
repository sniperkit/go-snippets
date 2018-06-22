#!/bin/sh
# launch ssr
python /shadowsocksr/shadowsocks/server.py -p $PORT -k $PASSWD -m $METHOD -O $PROTOCOL -o $OBFS -d start
# config kcptun
sed -i "s/listen_port/${listen_port}/g" /kcptun/${MODE}-config.json
sed -i "s/remote_server:remote_port/${remote_server}:${remote_port}/g" /kcptun/${MODE}-config.json
# launch kcptun server or client
/kcptun/${MODE}_linux_amd64 -c /kcptun/${MODE}-config.json
