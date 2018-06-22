#!/bin/sh
# launch website
nginx
# launch ssr
cd /shadowsocksr
chmod +x initcfg.sh
./initcfg.sh
# Warning! 
# The shell below is configured mannually! 
# Do not build this images by 'docker build' command
# Do not build this images by 'docker build' command
# Do not build this images by 'docker build' command
cat user-config.json
cd shadowsocks
./run.sh
# config kcptun
sed -i "s/listen_port/${listen_port}/g" /kcptun/${MODE}-config.json
sed -i "s/remote_server:remote_port/${remote_server}:${remote_port}/g" /kcptun/${MODE}-config.json
# launch kcptun server or client
/kcptun/${MODE}_linux_amd64 -c /kcptun/${MODE}-config.json
