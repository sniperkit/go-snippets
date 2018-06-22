#!/usr/bin/with-contenv sh

# activate ssl

if [[ ! -f /etc/nginx/keys/dh.pem && $NGINX_ENABLE_SSL = 1 ]]; then
    apk --no-cache add openssl
    echo "Generating dh.pem"
    openssl dhparam -out /etc/nginx/keys/dh.pem "$NGINX_DH_SIZE"
    echo "Generating insecure https certificate"
    openssl req -newkey rsa:4096 -x509 -days 365 -nodes -sha256 \
        -subj "$NGINX_SSL_SUBJECT" \
        -keyout "/etc/nginx/keys/web.key" \
        -out "/etc/nginx/keys/web.crt"
    cat /etc/nginx/keys/web.crt /etc/nginx/keys/web.key > /etc/nginx/keys/web.pem
    cat /etc/nginx/ssl/ssl.conf > /etc/nginx/conf.d/14-ssl.conf
    cat /etc/nginx/ssl/default.conf > /etc/nginx/hosts.d/default.conf
    # # fix permissions on SSL certificates
    chmod 0400 /etc/nginx/keys/*
    chmod 0644 /etc/nginx/keys/web.crt
fi