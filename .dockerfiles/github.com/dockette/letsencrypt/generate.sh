#!/bin/bash

if [ -z "$DOMAINS" ] ; then
  echo "No domains set, please fill -e DOMAINS='example.com www.example.com'"
  exit 1
fi

if [ -z "$EMAIL" ] ; then
  echo "No email set, please fill -e EMAIL='your@email.tld'"
  exit 1
fi

# Prepare args
DARRAYS=(${DOMAINS})
NGINX_DOMAINS=${DOMAINS}
EMAIL_ADDRESS=${EMAIL}
LE_DOMAINS=("${DARRAYS[*]/#/-d }")
CERTIFICATES=/etc/letsencrypt/live

# Inform
echo "Creating certificates for: $NGINX_DOMAINS"

# Replace domains
sed -i "s/\$DOMAINS/$NGINX_DOMAINS/g" /etc/nginx/nginx.conf

# Start nginx
echo "- start nginx"
service nginx start

# Run LetsEncrypt
echo "- start letsencrypt"
${LE_BIN} certonly --webroot -w /var/www/acme-certs ${LE_DOMAINS} --email ${EMAIL_ADDRESS} --agree-tos

# Copy created certs
if [ -d ${CERTIFICATES} ] ; then
  echo "- copy certificates to /var/www/certs"
  cp ${CERTIFICATES}/. /var/www/certs/ -R -L
else
  echo "- certificates folder $CERTIFICATES not found"
fi

# Stop nginx
echo "- stop nginx"
service nginx stop
