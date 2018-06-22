#!/usr/bin/env bash
cd "$(dirname "$0")"
pwd
docker rm $(docker ps -a -q  --filter ancestor=registry:2) -f
docker run -d -p 5000:5000 --name registry -v `pwd`/auth:/auth -e "REGISTRY_AUTH=htpasswd" -e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" -e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd registry:2

i=0
until [ ${i} -ge 10 ] || $(curl --output /dev/null --silent --head --fail http://localhost:5000); do
    printf '.'
    sleep .2
    ((i++))
done

if (( i > 10 )); then exit 1; fi
exit 0